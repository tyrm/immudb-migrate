package migrate

import (
	"context"
	"errors"
	"fmt"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"log"
	"time"
)

type MigratorOption func(m *Migrator)

func WithTableName(table string) MigratorOption {
	return func(m *Migrator) {
		m.table = table
	}
}

func WithLocksTableName(table string) MigratorOption {
	return func(m *Migrator) {
		m.locksTable = table
	}
}

type Migrator struct {
	immudb     client.ImmuClient
	migrations *Migrations

	ms MigrationSlice

	table        string
	locksTable   string
	locksTableID int64
}

func NewMigrator(db client.ImmuClient, migrations *Migrations, opts ...MigratorOption) *Migrator {
	m := &Migrator{
		immudb:     db,
		migrations: migrations,

		ms: migrations.ms,

		table:      "immudb_migrations",
		locksTable: "immudb_migration_locks",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Migrator) Init(ctx context.Context) error {
	// prep statements
	statements := []string{
		createTableMigrations(m.table),
		createTableMigrationLocks(m.locksTable),
		createIndexMigrationLocks(m.locksTable),
	}

	// create transaction
	tx, err := m.immudb.NewTx(ctx)
	if err != nil {
		return err
	}

	// run statements
	for _, c := range statements {
		err = tx.SQLExec(ctx, c, nil)
		if err != nil {
			return Rollback(ctx, tx, err)
		}
	}

	// commit
	_, err = tx.Commit(ctx)

	return err
}

// MarkApplied marks the migration as applied (completed).
func (m *Migrator) MarkApplied(ctx context.Context, migration *Migration) error {
	_, err := m.immudb.SQLExec(
		ctx,
		insertTableMigration(m.table, migration.Name, migration.GroupID),
		nil,
	)

	return err
}

// Migrate runs unapplied migrations. If a migration fails, migrate immediately exits.
func (m *Migrator) Migrate(ctx context.Context, opts ...MigrationOption) (*MigrationGroup, error) {
	cfg := newMigrationConfig(opts)

	if err := m.validate(); err != nil {
		return nil, err
	}

	if err := m.Lock(ctx); err != nil {
		return nil, err
	}
	defer m.Unlock(ctx) //nolint:errcheck

	migrations, lastGroupID, err := m.migrationsWithStatus(ctx)
	if err != nil {
		return nil, err
	}
	migrations = migrations.Unapplied()

	group := new(MigrationGroup)
	if len(migrations) == 0 {
		return group, nil
	}
	group.ID = lastGroupID + 1

	for i := range migrations {
		migration := &migrations[i]
		migration.GroupID = group.ID

		// Always mark migration as applied so the rollback has a chance to fix the database.
		if err := m.MarkApplied(ctx, migration); err != nil {
			return group, err
		}

		group.Migrations = migrations[:i+1]

		if !cfg.nop && migration.Up != nil {
			if err := migration.Up(ctx, m.immudb); err != nil {
				return group, err
			}
		}
	}

	return group, nil
}

func (m *Migrator) Lock(ctx context.Context) error {
	// check for lock
	entry, err := m.immudb.Get(ctx, keyTableMigrationLock(m.locksTable, m.table))
	fmt.Printf("err fmt: %T", err)
	if err != nil && err.Error() != KeyNotFoundError {
		// error that isn't key not found
		return fmt.Errorf("migrate: migrations table is already locked (%w)", err)
	}
	if entry != nil && string(entry.Value) == StateLocked {
		// entry is locked
		return errors.New("migrate: migrations table is already locked")
	}

	// set lock
	var preconditions []*schema.Precondition
	if entry != nil {
		precondition := schema.PreconditionKeyNotModifiedAfterTX(
			keyTableMigrationLock(m.locksTable, m.table),
			entry.Tx,
		)

		preconditions = append(preconditions, precondition)
	}
	_, err = m.immudb.SetAll(ctx, &schema.SetRequest{
		KVs: []*schema.KeyValue{{
			Key:   keyTableMigrationLock(m.locksTable, m.table),
			Value: []byte(StateLocked),
		}},
		Preconditions: preconditions,
	})
	if err != nil {
		return fmt.Errorf("migrate: migrations table is already locked (%w)", err)
	}

	return nil
}

func (m *Migrator) Unlock(ctx context.Context) error {
	// Without verification
	tx, err := m.immudb.Set(ctx, keyTableMigrationLock(m.locksTable, m.table), []byte(StateUnlocked))
	if err != nil {
		return fmt.Errorf("migrate: unlocking (%w)", err)
	}
	log.Printf("Set: tx: %d", tx.Id)

	return nil
}

// privates

func (m *Migrator) migrationsWithStatus(ctx context.Context) (MigrationSlice, int64, error) {
	sorted := m.migrations.Sorted()

	applied, err := m.selectAppliedMigrations(ctx)
	if err != nil {
		return nil, 0, err
	}

	appliedMap := migrationMap(applied)
	for i := range sorted {
		m1 := &sorted[i]
		if m2, ok := appliedMap[m1.Name]; ok {
			m1.ID = m2.ID
			m1.GroupID = m2.GroupID
			m1.MigratedAt = m2.MigratedAt
		}
	}

	return sorted, applied.LastGroupID(), nil
}

// selectAppliedMigrations selects applied (applied) migrations in descending order.
func (m *Migrator) selectAppliedMigrations(ctx context.Context) (MigrationSlice, error) {
	resp, err := m.immudb.SQLQuery(
		ctx,
		selectAppliedTableMigrations(m.table),
		nil,
		true,
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("applied migrations: %+v\n\n", resp)

	var ms MigrationSlice
	for _, row := range resp.GetRows() {
		migratedAt := row.GetValues()[2].GetTs()

		migratedAtSec := migratedAt / 1000000
		migratedAtNSec := (migratedAt % 1000000) * 1000

		fmt.Printf("ts: %d, sec: %d nsec: %d\n\n", migratedAt, migratedAtSec, migratedAtNSec)
		migratedAtT := time.Unix(migratedAtSec, migratedAtNSec)

		migration := Migration{
			ID:         row.GetValues()[0].GetN(),
			Name:       row.GetValues()[1].GetS(),
			GroupID:    row.GetValues()[2].GetN(),
			MigratedAt: migratedAtT,
		}

		ms = append(ms, migration)
	}
	fmt.Printf("applied migration slices: %+v\n\n", resp)
	return ms, nil
}

func (m *Migrator) validate() error {
	if len(m.ms) == 0 {
		return errors.New("migrate: there are no any migrations")
	}
	return nil
}