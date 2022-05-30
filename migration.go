package migrate

import (
	"context"
	"github.com/codenotary/immudb/pkg/client"
	"time"
)

type MigrationHandler func(ctx context.Context, db client.ImmuClient) error

type Migration struct {
	ID         int64
	Name       string
	GroupID    int64
	MigratedAt time.Time

	Up MigrationHandler
}

func (m *Migration) String() string {
	return m.Name
}

func (m *Migration) IsApplied() bool {
	return m.ID > 0
}