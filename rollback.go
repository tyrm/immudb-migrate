package migrate

import (
	"context"
	"fmt"
	"github.com/codenotary/immudb/pkg/client"
)

func Rollback(ctx context.Context, tx client.Tx, topErr error) error {
	err := tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("%s; %s", err.Error(), topErr.Error())
	}

	return topErr
}
