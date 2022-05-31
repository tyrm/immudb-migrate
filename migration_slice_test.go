package migrate

import (
	"fmt"
	"testing"
)

func TestMigrationSlice_LastGroupID(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  MigrationSlice
		output int64
	}{
		{
			MigrationSlice{},
			0,
		},
		{
			MigrationSlice{
				testMigration1,
			},
			1,
		},
		{
			MigrationSlice{
				testMigration2,
				testMigration3,
				testMigration1,
			},
			2,
		},
		{
			MigrationSlice{
				testMigration1,
				testMigration2,
				testMigration3,
				testMigration4,
				testMigration5,
				testMigration6,
			},
			3,
		},
		{
			MigrationSlice{
				testMigration6,
				testMigration5,
				testMigration4,
				testMigration3,
				testMigration2,
				testMigration1,
			},
			3,
		},
		{
			MigrationSlice{
				testMigration3,
				testMigration4,
				testMigration1,
				testMigration5,
				testMigration2,
				testMigration6,
			},
			3,
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationSlice.LastGroupID", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			groupID := table.input.LastGroupID()

			if groupID != table.output {
				t.Errorf("[%d] invalid group id, got: '%d', want: '%d'", i, groupID, table.output)
			}
		})
	}
}

func TestMigrationSlice_String(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  MigrationSlice
		output string
	}{
		{
			MigrationSlice{},
			"empty",
		},
		{
			MigrationSlice{
				testMigration1,
			},
			"20220504174129",
		},
		{
			MigrationSlice{
				testMigration1,
				testMigration2,
			},
			"20220504174129, 20220506174129",
		},
		{
			MigrationSlice{
				testMigration1,
				testMigration2,
				testMigration3,
				testMigration4,
				testMigration5,
				testMigration6,
			},
			"6 migrations (20220504174129 ... 20220514174129)",
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationSlice.String", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			str := table.input.String()

			if str != table.output {
				t.Errorf("[%d] invalid string, got: '%s', want: '%s'", i, str, table.output)
			}
		})
	}
}

func TestMigrationSlice_Unapplied(t *testing.T) {
	t.Parallel()

	tables := []struct {
		input  MigrationSlice
		output MigrationSlice
	}{
		{
			MigrationSlice{},
			MigrationSlice{},
		},
		{
			MigrationSlice{
				testMigration1,
			},
			MigrationSlice{},
		},
		{
			MigrationSlice{
				testMigration2,
				testMigration3,
				testMigration1,
			},
			MigrationSlice{},
		},
		{
			MigrationSlice{
				testMigration1,
				testMigration2,
				testMigration3,
				testMigration4,
				testMigration5,
				testMigration6,
			},
			MigrationSlice{},
		},
		{
			MigrationSlice{
				testMigration6,
				testMigration5,
				testMigration4,
				testMigration3,
				testMigration2,
				testMigration1,
			},
			MigrationSlice{},
		},
		{
			MigrationSlice{
				testMigration3,
				testMigration4,
				testMigration1,
				testMigration5,
				testMigration2,
				testMigration6,
			},
			MigrationSlice{},
		},
		{
			MigrationSlice{
				testMigration3,
				testMigration4,
				testMigration7,
				testMigration1,
				testMigration5,
				testMigration2,
				testMigration6,
			},
			MigrationSlice{
				testMigration7,
			},
		},
		{
			MigrationSlice{
				testMigration7,
				testMigration8,
			},
			MigrationSlice{
				testMigration7,
				testMigration8,
			},
		},
	}

	for i, table := range tables {
		i := i
		table := table

		name := fmt.Sprintf("[%d] Running MigrationSlice.LastGroupID", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			unapplied := table.input.Unapplied()

			if len(unapplied) != len(table.output) {
				t.Errorf("[%d] invalid length, got: '%d', want: '%d'", i, len(unapplied), len(table.output))

				return
			}

			testCompareMigrationSlice(t, i, unapplied, table.output)
		})
	}
}

func testCompareMigrationSlice(t *testing.T, run int, got, want MigrationSlice) {
	t.Helper()

	runStr := ""
	if run > 0 {
		runStr = fmt.Sprintf("[%d] ", run)
	}

	for i, eMig := range want {
		oMig := got[i]

		if oMig.ID != eMig.ID {
			t.Errorf("%sinvalid migration id at %d, got: '%d', want: '%d'", runStr, i, oMig.ID, eMig.ID)
		}
		if oMig.Name != eMig.Name {
			t.Errorf("%sinvalid migration name at %d, got: '%+v', want: '%+v'", runStr, i, oMig.Name, eMig.Name)
		}
		if oMig.GroupID != eMig.GroupID {
			t.Errorf("%sinvalid migration group id at %d, got: '%d', want: '%d'", runStr, i, oMig.GroupID, eMig.GroupID)
		}
	}
}
