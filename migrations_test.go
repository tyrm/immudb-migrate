package migrate

import (
	"fmt"
	"strings"
	"testing"
)

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

			switch {
			case err == nil && table.err != nil:
				// error was unexpectedly nil
				t.Errorf("[%d] expected error, got: '%v', want: '%v'", i, err, table.err)

				return
			case err != nil && table.err == nil:
				// expected nil error
				t.Errorf("[%d] unexpected nil error, got: '%v', want: '%v'", i, err, table.err)

				return
			case err != nil && table.err != nil && err.Error() != table.err.Error():
				// errors do not match
				t.Errorf("[%d] invalid error, got: '%v', want: '%v'", i, err, table.err)

				return
			}

			if migName != table.name {
				t.Errorf("[%d] invalid name, got: '%s', want: '%s'", i, migName, table.name)
			}
		})
	}
}
