package migrate

import (
	"bytes"
	"fmt"
	"testing"
)

func TestKeyTableMigrationLock(t *testing.T) {
	t.Parallel()

	tables := []struct {
		lockKeyName string
		name        string

		output []byte
	}{
		{"immudb_migration_lock", "immudb_migrations", []byte("immudb_migration_lock:immudb_migrations")},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running keyTableMigrationLock", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newKey := keyTableMigrationLock(table.lockKeyName, table.name)

			if !bytes.Equal(newKey, table.output) {
				t.Errorf("[%d] invalid key, got: '%s', want: '%s'", i, string(newKey), string(table.output))
			}
		})
	}
}
