package migrate

import "fmt"

const (
	tableMigrationsColumnNameID         = "id"
	tableMigrationsColumnNameName       = "name"
	tableMigrationsColumnNameGroupID    = "group_id"
	tableMigrationsColumnNameMigratedAt = "migrated_at"
)

const (
	tableMigrationsColumnIndexID int64 = iota
	tableMigrationsColumnIndexName
	tableMigrationsColumnIndexGroupID
	tableMigrationsColumnIndexMigratedAt
)

const tableMigrationsAllColumns = tableMigrationsColumnNameID + ", " + // 0
	tableMigrationsColumnNameName + ", " + // 1
	tableMigrationsColumnNameGroupID + ", " + // 2
	tableMigrationsColumnNameMigratedAt // 3

const createTableMigrationsStatement = `
CREATE TABLE IF NOT EXISTS %[1]s (
    %[2]s INTEGER AUTO_INCREMENT,
    %[3]s VARCHAR,
    %[4]s INTEGER,
    %[5]s TIMESTAMP,
    PRIMARY KEY (%[2]s)
);`

func createTableMigrations(tableName string) string {
	return fmt.Sprintf(
		createTableMigrationsStatement,
		tableName,
		tableMigrationsColumnNameID,
		tableMigrationsColumnNameName,
		tableMigrationsColumnNameGroupID,
		tableMigrationsColumnNameMigratedAt,
	)
}

const insertTableMigrationStatement = `
INSERT INTO %[1]s (
    %[2]s,
    %[3]s,
    %[4]s
)
VALUES (
    @%[2]s,
    @%[3]s,
    NOW()
);`

func insertTableMigration(tableName string) string {
	return fmt.Sprintf(
		insertTableMigrationStatement,
		tableName,
		tableMigrationsColumnNameName,
		tableMigrationsColumnNameGroupID,
		tableMigrationsColumnNameMigratedAt,
	)
}

const selectAppliedTableMigrationsStatement = `
SELECT %[2]s FROM %[1]s;`

func selectAppliedTableMigrations(tableName string) string {
	return fmt.Sprintf(
		selectAppliedTableMigrationsStatement,
		tableName,
		tableMigrationsAllColumns,
	)
}
