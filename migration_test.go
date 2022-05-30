package migrate

import (
	"fmt"
	"testing"
)

func TestMigration_IsApplied(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  Migration
		output bool
	}{
		{
			Migration{},
			false,
		},
		{
			Migration{
				ID: 1,
			},
			true,
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running Migration.IsApplied", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			isApplied := table.input.IsApplied()

			if isApplied != table.output {
				t.Errorf("[%d] invalid isApplied, got: '%t', want: '%t'", i, isApplied, table.output)
			}
		})
	}
}

func TestMigration_String(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  Migration
		output string
	}{
		{
			Migration{
				Name: "test migration 1",
			},
			"test migration 1",
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running Migration.String", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			str := table.input.String()

			if str != table.output {
				t.Errorf("[%d] invalid string, got: '%s', want: '%s'", i, str, table.output)
			}
		})
	}
}
