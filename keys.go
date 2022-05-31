package migrate

import "fmt"

func keyTableMigrationLock(lockKeyName, name string) []byte {
	return []byte(fmt.Sprintf("%s:%s", lockKeyName, name))
}
