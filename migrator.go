package migrate

import (
	"context"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
)

type MigratorOption func(m *Migrator)

func WithTableName(table string) MigratorOption {
	return func(m *Migrator) {
		m.table = table
	}
}

func WithLocksKeyName(table string) MigratorOption {
	return func(m *Migrator) {
		m.locksKey = table
	}
}

type Migrator struct {
	immudb     client.ImmuClient
	migrations *Migrations

	ms MigrationSlice

	table    string
	locksKey string
}

func NewMigrator(db client.ImmuClient, migrations *Migrations, opts ...MigratorOption) *Migrator {
	m := &Migrator{
		immudb:     db,
		migrations: migrations,

		ms: migrations.ms,

		table:    "immudb_migrations",
		locksKey: "immudb_migration_lock",
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

		group.Migrations = migrations[:i+1]

		if !cfg.nop && migration.Up != nil {
			// create transaction
			tx, err := m.immudb.NewTx(ctx)
			if err != nil {
				return group, err
			}

			// run migration
			if err := migration.Up(ctx, tx); err != nil {
				return group, Rollback(ctx, tx, err)
			}

			// commit
			_, err = tx.Commit(ctx)
			if err != nil {
				return group, err
			}
		}

		if err := m.MarkApplied(ctx, migration); err != nil {
			return group, err
		}
	}

	return group, nil
}

func (m *Migrator) Lock(ctx context.Context) error {
	// check for lock
	entry, err := m.immudb.Get(ctx, keyTableMigrationLock(m.locksKey, m.table))
	if err != nil && err.Error() != KeyNotFoundError {
		// error that isn't 'key not found'
		return NewLockError("migrations table is already locked", err)
	}
	if entry != nil && string(entry.Value) == StateLocked {
		// entry is locked
		return NewLockError("migrations table is already locked", nil)
	}

	// set lock
	var preconditions []*schema.Precondition
	if entry != nil {
		precondition := schema.PreconditionKeyNotModifiedAfterTX(
			keyTableMigrationLock(m.locksKey, m.table),
			entry.Tx,
		)

		preconditions = append(preconditions, precondition)
	}
	_, err = m.immudb.SetAll(ctx, &schema.SetRequest{
		KVs: []*schema.KeyValue{{
			Key:   keyTableMigrationLock(m.locksKey, m.table),
			Value: []byte(StateLocked),
		}},
		Preconditions: preconditions,
	})
	if err != nil {
		return NewLockError("migrations table is already locked", err)
	}

	return nil
}

func (m *Migrator) Unlock(ctx context.Context) error {
	// Without verification
	_, err := m.immudb.Set(ctx, keyTableMigrationLock(m.locksKey, m.table), []byte(StateUnlocked))
	if err != nil {
		return NewLockError("can't unlock migration table", err)
	}

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

	var ms MigrationSlice
	for _, row := range resp.GetRows() {
		migration := Migration{
			ID:         row.GetValues()[0].GetN(),
			Name:       row.GetValues()[1].GetS(),
			GroupID:    row.GetValues()[2].GetN(),
			MigratedAt: tsToTime(row.GetValues()[3].GetTs()),
		}

		ms = append(ms, migration)
	}

	return ms, nil
}

func (m *Migrator) validate() error {
	if len(m.ms) == 0 {
		return ErrNoNewMigrations
	}

	return nil
}
