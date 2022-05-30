package migrate

import "fmt"

// table migrations

const createTableMigrationsStatement = `
CREATE TABLE IF NOT EXISTS %s(
    id          INTEGER AUTO_INCREMENT,
    name        VARCHAR,
    group_id    INTEGER,
    migrated_at TIMESTAMP,
    PRIMARY KEY (id)
);`

func createTableMigrations(name string) string {
	return fmt.Sprintf(createTableMigrationsStatement, name)
}

const insertTableMigrationStatement = `
INSERT INTO %s (
    name,
    group_id,
    migrated_at
)
VALUES (
    '%s',
    %d,
    NOW()
);`

func insertTableMigration(tableName, name string, groupID int64) string {
	return fmt.Sprintf(insertTableMigrationStatement, tableName, name, groupID)
}

const selectAppliedTableMigrationsStatement = `
SELECT id, name, group_id, migrated_at FROM %s;`

func selectAppliedTableMigrations(tableName string) string {
	return fmt.Sprintf(selectAppliedTableMigrationsStatement, tableName)
}

// table migration lock

const createTableMigrationLocksStatement = `
CREATE TABLE IF NOT EXISTS %s(
    id          INTEGER AUTO_INCREMENT,
    table_name  VARCHAR[256],
    state       VARCHAR[256],
    PRIMARY KEY (id)
);`

func createTableMigrationLocks(name string) string {
	return fmt.Sprintf(createTableMigrationLocksStatement, name)
}

const createIndexMigrationLocksStatement = `
CREATE UNIQUE INDEX IF NOT EXISTS ON %s(table_name, state);`

func createIndexMigrationLocks(name string) string {
	return fmt.Sprintf(createIndexMigrationLocksStatement, name)
}

const insertTableMigrationLockStatement = `
INSERT INTO %s (
    table_name,
    state
)
VALUES (
    '%s',
    '%s'
);`

func insertTableMigrationLock(lockTableName, name string, locked bool) string {
	state := StateUnlocked
	if locked {
		state = StateLocked
	}
	return fmt.Sprintf(insertTableMigrationLockStatement, lockTableName, name, state)
}

const selectTableMigrationLockStatement = `
SELECT id, table_name, state FROM %s WHERE table_name = '%s';`

func selectTableMigrationLock(lockTableName string, name string) string {
	return fmt.Sprintf(selectTableMigrationLockStatement, lockTableName, name)
}

const upsertTableMigrationLockStatement = `
UPSERT INTO %s (
    id,
    table_name,
    state
)
VALUES (
    %d,
    '%s',
    '%s'
);`

func upsertTableMigrationLock(lockTableName string, id int64, name string, locked bool) string {
	state := StateUnlocked
	if locked {
		state = StateLocked
	}
	return fmt.Sprintf(upsertTableMigrationLockStatement, lockTableName, id, name, state)
}
