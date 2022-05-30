package migrate

import (
	"fmt"
	"testing"
	"time"
)

//revive:disable:add-constant

func TestMigrationMap(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  MigrationSlice
		output map[string]*Migration
	}{
		{
			input: MigrationSlice{
				Migration{
					Name: "20220506174128",
				},
				Migration{
					Name: "20220508174128",
				},
				Migration{
					Name: "20220504174128",
				},
				Migration{
					Name: "20220510174128",
				},
			},
			output: map[string]*Migration{
				"20220506174128": {
					Name: "20220506174128",
				},
				"20220508174128": {
					Name: "20220508174128",
				},
				"20220504174128": {
					Name: "20220504174128",
				},
				"20220510174128": {
					Name: "20220510174128",
				},
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running migrationMap", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			output := migrationMap(table.input)

			j := 0
			for k := range output {
				if k != table.input[j].Name {
					t.Errorf("[%d] invalid key at %d, got: '%s', want: '%s'", i, j, k, table.input[j].Name)
				}

				if output[k].Name != table.input[j].Name {
					t.Errorf("[%d] invalid migration at %d, got: '%s', want: '%s'", i, j, output[k].Name, table.input[j].Name)
				}

				j++
			}
		})
	}
}

func TestSortAsc(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  MigrationSlice
		output MigrationSlice
	}{
		{
			input: MigrationSlice{
				Migration{
					Name: "20220506174128",
				},
				Migration{
					Name: "20220508174128",
				},
				Migration{
					Name: "20220504174128",
				},
				Migration{
					Name: "20220510174128",
				},
			},
			output: MigrationSlice{
				Migration{
					Name: "20220504174128",
				},
				Migration{
					Name: "20220506174128",
				},
				Migration{
					Name: "20220508174128",
				},
				Migration{
					Name: "20220510174128",
				},
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running sortAsc", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sortAsc(table.input)

			for j := range table.input {
				if table.input[j].Name != table.output[j].Name {
					t.Errorf("[%d] invalid order at %d, got: '%s', want: '%s'", i, j, table.input[i].Name, table.output[i].Name)
				}
			}
		})
	}
}

func TestTsToTime(t *testing.T) {
	t.Parallel()

	tables := []struct {
		ts   int64
		time time.Time
	}{
		{1653885379429868, time.Date(2022, 05, 30, 4, 36, 19, 429868000, time.UTC)},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running tsToTime for %d", i, table.ts)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newTime := tsToTime(table.ts)
			if !newTime.Equal(table.time) {
				t.Errorf("[%d] invalid time, got: '%s', want: '%s'", i, newTime, table.time)
			}
		})
	}
}

//revive:enable:add-constant
