package migrate

import "fmt"

const createTableMigrationsStatement = `
CREATE TABLE IF NOT EXISTS %s (
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
