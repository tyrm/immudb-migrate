package migrate

import (
	"fmt"
	"strings"
)

type MigrationSlice []Migration

func (ms MigrationSlice) String() string {
	if len(ms) == 0 {
		return "empty"
	}

	if len(ms) > 5 {
		return fmt.Sprintf("%d migrations (%s ... %s)", len(ms), ms[0].Name, ms[len(ms)-1].Name)
	}

	var sb strings.Builder

	for i := range ms {
		if i > 0 {
			sb.WriteString(", ") //nolint
		}
		sb.WriteString(ms[i].Name) //nolint
	}

	return sb.String()
}

// LastGroupID returns the last applied migration group id.
// The id is 0 when there are no migration groups.
func (ms MigrationSlice) LastGroupID() int64 {
	var lastGroupID int64
	for i := range ms {
		groupID := ms[i].GroupID
		if groupID > lastGroupID {
			lastGroupID = groupID
		}
	}

	return lastGroupID
}

// Unapplied returns unapplied migrations in ascending order
// (the order is important and is used in Migrate).
func (ms MigrationSlice) Unapplied() MigrationSlice {
	var unapplied MigrationSlice
	for i := range ms {
		if !ms[i].IsApplied() {
			unapplied = append(unapplied, ms[i])
		}
	}
	sortAsc(unapplied)

	return unapplied
}
