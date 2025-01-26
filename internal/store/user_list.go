package store

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UsersQuery struct {
	Query  string
	Roles  []UserRole
	Banned bool

	Limit  uint
	Offset uint
}

type UserRow struct {
	ID       pgtype.UUID
	Name     string
	Role     UserRole
	IsBanned bool
	JoinedAt time.Time
}

func applyUsersQuery(
	query *goqu.SelectDataset,
	req *UsersQuery,
) *goqu.SelectDataset {
	if req.Query != "" {
		query = query.Where(goqu.I("users.name").ILike("%" + req.Query + "%"))
	}

	if len(req.Roles) > 0 {
		query = query.Where(goqu.I("users.role").In(req.Roles))
	}

	if req.Banned {
		query = query.Where(goqu.I("users.is_banned").IsTrue())
	}

	return query
}

func CountUsers(ctx context.Context, db DBTX, req *UsersQuery) (int64, error) {
	query := postgresQuery.
		Select(
			goqu.COUNT("*").As("count"),
		).
		From("users")

	query = applyUsersQuery(query, req)

	sql, params, err := query.ToSQL()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		var count int64
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}

	return 0, nil
}

func ListUsers(ctx context.Context, db DBTX, req UsersQuery) ([]UserRow, error) {
	query := postgresQuery.
		Select(
			goqu.I("users.id"),
			goqu.I("users.name"),
			goqu.I("users.role"),
			goqu.I("users.is_banned"),
			goqu.I("users.joined_at"),
		).
		From("users")

	query = applyUsersQuery(query, &req)

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}

	sql, params, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var users []UserRow
	for rows.Next() {
		row, err := pgx.RowToStructByName[UserRow](rows)
		if err != nil {
			return nil, err
		}
		users = append(users, row)
	}

	return users, nil
}
