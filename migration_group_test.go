package migrate

import (
	"fmt"
	"testing"
)

func TestMigrationGroup_IsZero(t *testing.T) {
	t.Parallel()

	tables := []struct {
		mg     MigrationGroup
		output bool
	}{
		{
			MigrationGroup{},
			true,
		},
		{
			MigrationGroup{
				ID: 1,
				Migrations: MigrationSlice{
					testMigration1,
				},
			},
			false,
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationGroup.IsZero", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			isZero := table.mg.IsZero()

			if isZero != table.output {
				t.Errorf("[%d] invalid isZero, got: '%t', want: '%t'", i, isZero, table.output)
			}
		})
	}
}

func TestMigrationGroup_String(t *testing.T) {
	t.Parallel()

	tables := []struct {
		mg     MigrationGroup
		output string
	}{
		{
			MigrationGroup{},
			"nil",
		},
		{
			MigrationGroup{
				ID: 1,
				Migrations: MigrationSlice{
					testMigration1,
				},
			},
			"group #1 (20220504174129)",
		},
		{
			MigrationGroup{
				ID: 2,
				Migrations: MigrationSlice{
					testMigration2,
					testMigration3,
				},
			},
			"group #2 (20220506174129, 20220508174129)",
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationGroup.String", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			output := table.mg.String()

			if output != table.output {
				t.Errorf("[%d] invalid string, got: '%s', want: '%s'", i, output, table.output)
			}
		})
	}
}
