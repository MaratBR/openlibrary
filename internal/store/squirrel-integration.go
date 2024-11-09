package store

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

type wrappedCommandTag pgconn.CommandTag

var (
	errLastInsertIdNotSupported = errors.New("LastInsertId not supported")
)

// LastInsertId implements sql.Result.
func (w wrappedCommandTag) LastInsertId() (int64, error) {
	return 0, errLastInsertIdNotSupported
}

// RowsAffected implements sql.Result.
func (w wrappedCommandTag) RowsAffected() (int64, error) {
	return pgconn.CommandTag(w).RowsAffected(), nil
}

func wrapCommandTag(tag pgconn.CommandTag) sql.Result {
	return wrappedCommandTag(tag)
}

type wrappedDBTX struct {
	db  DBTX
	ctx context.Context
}

// Exec implements squirrel.BaseRunner.
func (w *wrappedDBTX) Exec(query string, args ...interface{}) (sql.Result, error) {
	tag, err := w.db.Exec(w.ctx, query, args...)
	return wrapCommandTag(tag), err
}

// Query implements squirrel.BaseRunner.
func (w *wrappedDBTX) Query(query string, args ...interface{}) (*sql.Rows, error) {
	panic("unimplemented")
}

func wrapDBTX(db DBTX, ctx context.Context) sq.BaseRunner {
	return &wrappedDBTX{db: db, ctx: ctx}
}
