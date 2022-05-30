package migrate

import (
	"sort"
	"time"
)

const (
	microsecondsToSeconds     = 1000000
	microsecondsToNanoseconds = 1000
)

func migrationMap(ms MigrationSlice) map[string]*Migration {
	mp := make(map[string]*Migration)
	for i := range ms {
		m := &ms[i]
		mp[m.Name] = m
	}

	return mp
}

func sortAsc(ms MigrationSlice) {
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Name < ms[j].Name
	})
}

func tsToTime(ts int64) time.Time {
	sec := ts / microsecondsToSeconds
	ns := (ts % microsecondsToSeconds) * microsecondsToNanoseconds

	return time.Unix(sec, ns).UTC()
}
