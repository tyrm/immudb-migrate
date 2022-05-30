package migrate

import "fmt"

type MigrationGroup struct {
	ID         int64
	Migrations MigrationSlice
}

func (g *MigrationGroup) IsZero() bool {
	return g.ID == 0 && len(g.Migrations) == 0
}

func (g *MigrationGroup) String() string {
	if g.IsZero() {
		return "nil"
	}

	return fmt.Sprintf("group #%d (%s)", g.ID, g.Migrations)
}
