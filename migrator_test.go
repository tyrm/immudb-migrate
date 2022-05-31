package migrate

import (
	"fmt"
	"testing"
)

//revive:disable:add-constant

func TestWithTableName(t *testing.T) {
	t.Parallel()

	migOption := WithTableName("new_table_name")
	migrator := &Migrator{}

	migOption(migrator)

	if migrator.table != "new_table_name" {
		t.Errorf("WithTableName returned invalid table, got: '%s', want: '%s'", migrator.table, "new_table_name")
	}
}

func TestWithLocksKeyName(t *testing.T) {
	t.Parallel()

	migOption := WithLocksKeyName("new_table_lock_name")
	migrator := &Migrator{}

	migOption(migrator)

	if migrator.locksKey != "new_table_lock_name" {
		t.Errorf("WithTableName returned invalid table, got: '%s', want: '%s'", migrator.table, "new_table_lock_name")
	}
}

func TestNewMigrator(t *testing.T) {
	t.Parallel()

	tables := []struct {
		opts       []MigratorOption
		migrations *Migrations

		output Migrator
	}{
		{
			[]MigratorOption{},
			&Migrations{
				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},
			},

			Migrator{
				migrations: &Migrations{
					ms: MigrationSlice{
						testMigration1,
						testMigration2,
						testMigration3,
					},
				},

				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},

				table:    "immudb_migrations",
				locksKey: "immudb_migration_lock",
			},
		},
		{
			[]MigratorOption{
				WithTableName("new_table_name"),
			},
			&Migrations{
				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},
			},

			Migrator{
				migrations: &Migrations{
					ms: MigrationSlice{
						testMigration1,
						testMigration2,
						testMigration3,
					},
				},

				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},

				table:    "new_table_name",
				locksKey: "immudb_migration_lock",
			},
		},
		{
			[]MigratorOption{
				WithLocksKeyName("new_table_lock_name"),
			},
			&Migrations{
				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},
			},

			Migrator{
				migrations: &Migrations{
					ms: MigrationSlice{
						testMigration1,
						testMigration2,
						testMigration3,
					},
				},

				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},

				table:    "immudb_migrations",
				locksKey: "new_table_lock_name",
			},
		},
		{
			[]MigratorOption{
				WithTableName("new_table_name"),
				WithLocksKeyName("new_table_lock_name"),
			},
			&Migrations{
				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},
			},

			Migrator{
				migrations: &Migrations{
					ms: MigrationSlice{
						testMigration1,
						testMigration2,
						testMigration3,
					},
				},

				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},

				table:    "new_table_name",
				locksKey: "new_table_lock_name",
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running NewMigrator", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			migrator := NewMigrator(nil, table.migrations, table.opts...)

			if migrator == nil {
				t.Errorf("NewMigrator returned nil")

				return
			}

			testCompareMigrationSlice(t, i, migrator.migrations.ms, table.output.migrations.ms)
			testCompareMigrationSlice(t, i, migrator.ms, table.output.ms)

			if migrator.table != table.output.table {
				t.Errorf("NewMigrator returned invalid table, got: '%s', want: '%s'", migrator.table, table.output.table)
			}
			if migrator.locksKey != table.output.locksKey {
				t.Errorf("NewMigrator returned invalid locksKey, got: '%s', want: '%s'", migrator.locksKey, table.output.locksKey)
			}
		})
	}
}

func TestMigrator_Validate(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input Migrator
		err   error
	}{
		{
			Migrator{},
			ErrNoNewMigrations,
		},
		{
			Migrator{
				ms: MigrationSlice{
					testMigration1,
					testMigration2,
					testMigration3,
				},
			},
			nil,
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running NewMigrator", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := table.input.validate()

			if !testCompareErrors(t, i, err, table.err) {
				return
			}
		})
	}
}

//revive:enable:add-constant
