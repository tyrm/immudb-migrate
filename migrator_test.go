package migrate

import "testing"

func TestWithTableName(t *testing.T) {
	t.Parallel()

	migOption := WithTableName("new_table_name")
	migrator := &Migrator{}

	migOption(migrator)

	if migrator.table != "new_table_name" {
		t.Errorf("WithTableName returned invalid table, got: '%s', want: '%s'", migrator.table, "new_table_name")
	}
}

func TestWithLocksKeyName(t *testing.T) {
	t.Parallel()

	migOption := WithLocksKeyName("new_table_lock_name")
	migrator := &Migrator{}

	migOption(migrator)

	if migrator.locksKey != "new_table_lock_name" {
		t.Errorf("WithTableName returned invalid table, got: '%s', want: '%s'", migrator.table, "new_table_lock_name")
	}
}
