package migrate

import (
	"fmt"
	"testing"
)

var (
	testMigration1 = Migration{
		ID:      1,
		Name:    "20220504174129",
		GroupID: 1,
	}
	testMigration2 = Migration{
		ID:      2,
		Name:    "20220506174129",
		GroupID: 2,
	}
	testMigration3 = Migration{
		ID:      3,
		Name:    "20220508174129",
		GroupID: 2,
	}
	testMigration4 = Migration{
		ID:      4,
		Name:    "20220510174129",
		GroupID: 3,
	}
	testMigration5 = Migration{
		ID:      5,
		Name:    "20220512174129",
		GroupID: 3,
	}
	testMigration6 = Migration{
		ID:      6,
		Name:    "20220514174129",
		GroupID: 3,
	}
	testMigration7 = Migration{
		Name: "20220516174129",
	}
	testMigration8 = Migration{
		Name: "20220518174129",
	}
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
