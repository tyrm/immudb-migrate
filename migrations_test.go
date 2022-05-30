package migrate

import (
	"fmt"
	"strings"
	"testing"
)

func TestMigrationFile(t *testing.T) {
	if fpath := migrationFile(); !strings.Contains(fpath, "/go/src/testing/testing.go") {
		t.Errorf("NewMigrations returned invalid path, got: '%s', should contain: '%s'", fpath, "/go/src/testing/testing.go")
	}
}

func TestExtractMigrationName(t *testing.T) {
	//
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

			name, err := extractMigrationName(table.fpath)

			if err == nil && table.err != nil {
				// error was unexpectedly nil
				t.Errorf("[%d] expected error, got: '%v', want: '%v'", i, err, table.err)

				return
			} else if err != nil && table.err == nil {
				// expected nil error
				t.Errorf("[%d] unexpected nil error, got: '%v', want: '%v'", i, err, table.err)

				return
			} else if err != nil && table.err != nil && err.Error() != table.err.Error() {
				// errors do not match
				t.Errorf("[%d] invalid error, got: '%v', want: '%v'", i, err, table.err)

				return
			}

			if name != table.name {
				t.Errorf("[%d] invalid name, got: '%s', want: '%s'", i, name, table.name)
			}
		})
	}
}
