package migrate

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const testCreateTableMigrations = `
CREATE TABLE IF NOT EXISTS immudb_migrations (
    id          INTEGER AUTO_INCREMENT,
    name        VARCHAR,
    group_id    INTEGER,
    migrated_at TIMESTAMP,
    PRIMARY KEY (id)
);`

func TestCreateTableMigrations(t *testing.T) {
	result := createTableMigrations("immudb_migrations")

	if result != testCreateTableMigrations {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(result, testCreateTableMigrations, false)

		t.Errorf("invalid statement:\n%s", dmp.DiffPrettyText(diffs))
	}
}

const testInsertTableMigration = `
INSERT INTO immudb_migrations (
    name,
    group_id,
    migrated_at
)
VALUES (
    '20220506174128',
    1,
    NOW()
);`

func TestInsertTableMigration(t *testing.T) {
	result := insertTableMigration("immudb_migrations", "20220506174128", 1)

	if result != testInsertTableMigration {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(result, testInsertTableMigration, false)

		t.Errorf("invalid statement:\n%s", dmp.DiffPrettyText(diffs))
	}
}

const testSelectAppliedTableMigrations = `
SELECT id, name, group_id, migrated_at FROM immudb_migrations;`

func TestSelectAppliedTableMigrations(t *testing.T) {
	result := selectAppliedTableMigrations("immudb_migrations")

	if result != testSelectAppliedTableMigrations {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(result, testSelectAppliedTableMigrations, false)

		t.Errorf("invalid statement:\n%s", dmp.DiffPrettyText(diffs))
	}
}
