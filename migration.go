package migrate

import (
	"context"
	"time"

	"github.com/codenotary/immudb/pkg/client"
)

type MigrationHandler func(ctx context.Context, db client.Tx) error

type Migration struct {
	ID         int64
	Name       string
	GroupID    int64
	MigratedAt time.Time

	Up MigrationHandler
}

func (m *Migration) IsApplied() bool {
	return m.ID > 0
}

func (m *Migration) String() string {
	return m.Name
}
