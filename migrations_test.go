package migrate

import (
	"strings"
	"testing"
)

func TestMigrationFile(t *testing.T) {
	if fpath := migrationFile(); !strings.Contains(fpath, "/go/src/testing/testing.go") {
		t.Errorf("NewMigrations returned invalid path, got: '%s', should contain: '%s'", fpath, "/go/src/testing/testing.go")
	}
}
