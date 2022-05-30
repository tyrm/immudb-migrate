package migrate

import "sort"

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
