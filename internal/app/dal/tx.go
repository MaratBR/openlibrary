package dal

import (
	"context"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
)

func RollbackTx(ctx context.Context, tx pgx.Tx) {
	err := tx.Rollback(ctx)
	if err != nil {
		if err == pgx.ErrTxClosed {
			return
		}
		zap.S().Errorw("failed to rollback transaction", "err", err)
	}
}
