package migrate

import (
	"context"

	"github.com/codenotary/immudb/pkg/client"
)

func Rollback(ctx context.Context, tx client.Tx, topErr error) error {
	return NewRollbackError(tx.Rollback(ctx), topErr)
}
