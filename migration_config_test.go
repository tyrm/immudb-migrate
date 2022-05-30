package migrate

import (
	"fmt"
	"testing"
)

func TestNewMigrationConfig(t *testing.T) {
	t.Parallel()

	tables := []struct {
		opts   []MigrationOption
		output migrationConfig
	}{
		{
			[]MigrationOption{},
			migrationConfig{
				nop: false,
			},
		},
		{
			[]MigrationOption{
				WithNopMigration(),
			},
			migrationConfig{
				nop: true,
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running newMigrationConfig", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			migConfig := newMigrationConfig(table.opts)

			if *migConfig != table.output {
				t.Errorf("[%d] invalid config, got: '%+v', want: '%+v'", i, *migConfig, table.output)
			}
		})
	}
}

func TestWithNopMigration(t *testing.T) {
	t.Parallel()

	migOption := WithNopMigration()
	migConfig := &migrationConfig{}

	migOption(migConfig)

	if !migConfig.nop {
		t.Errorf("WithNopMigration returned invalid isApplied, got: '%t', want: '%t'", migConfig.nop, true)
	}
}
