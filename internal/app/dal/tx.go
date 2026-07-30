package dal

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func RollbackTx(ctx context.Context, tx pgx.Tx) {
	err := tx.Rollback(ctx)
	if err != nil {
		if err == pgx.ErrTxClosed {
			return
		}
		slog.Error("failed to rollback transaction", "err", err)
	}
}
