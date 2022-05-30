package migrate

import "fmt"

func keyTableMigrationLock(lockTableName, name string) []byte {
	return []byte(fmt.Sprintf("%s:%s", lockTableName, name))
}
