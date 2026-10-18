package app

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type moderationTransactionKey struct{}

// All moderation repositories participate in the case transaction. Nested Begin
// calls create savepoints; only the outer case transaction can commit enforcement.
func moderationDB(ctx context.Context, db DB) DB {
	if tx, ok := ctx.Value(moderationTransactionKey{}).(pgx.Tx); ok {
		return tx
	}
	return db
}
