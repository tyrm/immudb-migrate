package migrate

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

//revive:disable:add-constant

func TestNewMigrations(t *testing.T) {
	t.Parallel()

	tables := []struct {
		opts []MigrationsOption
	}{
		{[]MigrationsOption{}},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running extractMigrationName", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			migrations := NewMigrations(table.opts...)

			if !strings.Contains(migrations.implicitDirectory, "/go/src/testing") {
				t.Errorf("NewMigrations returned invalid implicitDirectory, got: '%s', should contain: '%s'", migrations.implicitDirectory, "/go/src/testing")
			}
			if len(migrations.ms) != 0 {
				t.Errorf("NewMigrations returned unexpected migrations, got: '%d', should contain: '%d'", len(migrations.ms), 0)
			}
		})
	}
}

func TestMigrations_Add(t *testing.T) {
	t.Parallel()

	tables := []struct {
		migrations []Migration
	}{
		{
			[]Migration{},
		},
		{
			[]Migration{
				{
					ID:         1,
					Name:       "20220506174129",
					GroupID:    1,
					MigratedAt: time.Date(2022, 05, 30, 16, 04, 27, 0, time.UTC),
				},
			},
		},
		{
			[]Migration{
				{
					ID:         1,
					Name:       "20220506174129",
					GroupID:    1,
					MigratedAt: time.Date(2022, 05, 30, 16, 04, 27, 0, time.UTC),
				},
				{
					ID:         2,
					Name:       "20220508174129",
					GroupID:    2,
					MigratedAt: time.Date(2022, 05, 30, 16, 04, 28, 0, time.UTC),
				},
				{
					ID:         3,
					Name:       "20220510174129",
					GroupID:    2,
					MigratedAt: time.Date(2022, 05, 30, 16, 04, 29, 0, time.UTC),
				},
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running Migrations.Add", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			migrations := Migrations{}
			for _, m := range table.migrations {
				migrations.Add(m)
			}

			if len(migrations.ms) != len(table.migrations) {
				t.Errorf("[%d] invalid length, got: '%d', want: '%d'", i, len(migrations.ms), len(table.migrations))

				return
			}

			for j, eMig := range table.migrations {
				oMig := migrations.ms[j]

				if oMig.ID != eMig.ID {
					t.Errorf("[%d] invalid migration id at %d, got: '%d', want: '%d'", i, j, oMig.ID, eMig.ID)
				}
				if oMig.Name != eMig.Name {
					t.Errorf("[%d] invalid migration name at %d, got: '%+v', want: '%+v'", i, j, oMig.Name, eMig.Name)
				}
				if oMig.GroupID != eMig.GroupID {
					t.Errorf("[%d] invalid migration group id at %d, got: '%d', want: '%d'", i, j, oMig.GroupID, eMig.GroupID)
				}
			}
		})
	}
}

func TestMigrations_Add_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should have paniced")
		}
	}()

	migrations := Migrations{}
	migrations.Add(Migration{})
}

func TestMigrations_Sorted(t *testing.T) {
	t.Parallel()

	tables := []struct {
		migrations Migrations
		slices     MigrationSlice
	}{
		{
			Migrations{
				ms: MigrationSlice{
					testMigration2,
					testMigration1,
					testMigration4,
					testMigration3,
				},
			},
			MigrationSlice{
				testMigration1,
				testMigration2,
				testMigration3,
				testMigration4,
			},
		},
		{
			Migrations{
				ms: MigrationSlice{
					testMigration5,
					testMigration4,
					testMigration3,
					testMigration2,
					testMigration1,
				},
			},
			MigrationSlice{
				testMigration1,
				testMigration2,
				testMigration3,
				testMigration4,
				testMigration5,
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running Migrations.Sorted", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			migrations := table.migrations.Sorted()

			if len(migrations) != len(table.slices) {
				t.Errorf("[%d] invalid length, got: '%d', want: '%d'", i, len(migrations), len(table.slices))

				return
			}

			for j, eMig := range table.slices {
				oMig := migrations[j]

				if oMig.ID != eMig.ID {
					t.Errorf("[%d] invalid migration id at %d, got: '%d', want: '%d'", i, j, oMig.ID, eMig.ID)
				}
				if oMig.Name != eMig.Name {
					t.Errorf("[%d] invalid migration name at %d, got: '%+v', want: '%+v'", i, j, oMig.Name, eMig.Name)
				}
				if oMig.GroupID != eMig.GroupID {
					t.Errorf("[%d] invalid migration group id at %d, got: '%d', want: '%d'", i, j, oMig.GroupID, eMig.GroupID)
				}
			}
		})
	}
}

func TestMigrationFile(t *testing.T) {
	t.Parallel()

	if fpath := migrationFile(); !strings.Contains(fpath, "/go/src/testing/testing.go") {
		t.Errorf("NewMigrations returned invalid path, got: '%s', should contain: '%s'", fpath, "/go/src/testing/testing.go")
	}
}

func TestExtractMigrationName(t *testing.T) {
	t.Parallel()

	tables := []struct {
		fpath string

		name string
		err  error
	}{
		{"/usr/local/go/src/testing/testing.go", "", NewMigrationNameError("unsupported migrate name format", "testing.go")},
		{"/go/src/github.com/tyrm/immudb-migrate/20220506174129_init.go", "20220506174129", nil},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running extractMigrationName", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			migName, err := extractMigrationName(table.fpath)

			if !testCompareErrors(t, i, err, table.err) {
				return
			}

			if migName != table.name {
				t.Errorf("[%d] invalid name, got: '%s', want: '%s'", i, migName, table.name)
			}
		})
	}
}

//revive:enable:add-constant
