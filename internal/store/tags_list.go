package store

import (
	"context"
	"log/slog"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ListTagsQuery struct {
	Query          string
	OnlyParentTags bool
	OnlyAdultTags  bool
	Limit          uint
	Offset         uint
}

type TagRow struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	IsAdult     bool      `db:"is_adult"`
	IsSpoiler   bool      `db:"is_spoiler"`
	TagType     TagType   `db:"tag_type"`
	Description string    `db:"description"`
	IsDefault   bool      `db:"is_default"`
	CreatedAt   time.Time `db:"created_at"`

	SynonymOf     pgtype.Int8 `db:"synonym_of"`
	SynonymOfName pgtype.Text `db:"synonym_name"`
}

func applyTagsQuery(query *goqu.SelectDataset, req *ListTagsQuery) *goqu.SelectDataset {
	if req.OnlyParentTags {
		query = query.Where(goqu.I("defined_tags.synonym_of").IsNull())
	}

	if req.OnlyParentTags {
		query = query.Where(goqu.I("defined_tags.synonym_of").IsNull())
	}

	if req.OnlyAdultTags {
		query = query.Where(goqu.I("defined_tags.is_adult").IsTrue())
	}

	if req.Query != "" {
		query = query.Where(goqu.I("defined_tags.name").ILike("%" + req.Query + "%"))
	}

	return query
}

func ListTags(
	ctx context.Context,
	db DBTX,
	req ListTagsQuery,
) ([]TagRow, error) {

	query := postgresQuery.
		Select(
			goqu.I("defined_tags.id"),
			goqu.I("defined_tags.name"),
			goqu.I("defined_tags.description"),
			goqu.I("defined_tags.is_spoiler"),
			goqu.I("defined_tags.is_adult"),
			goqu.I("defined_tags.created_at"),
			goqu.I("defined_tags.tag_type"),
			goqu.I("defined_tags.synonym_of"),
			goqu.I("defined_tags.is_default"),
			goqu.I("synonym.name").As("synonym_name"),
		).
		From("defined_tags").
		LeftJoin(goqu.T("defined_tags").As("synonym"), goqu.On(
			goqu.I("synonym.id").Eq(goqu.I("defined_tags.synonym_of")),
		)).
		Order(goqu.I("defined_tags.name").Asc())

	query = applyTagsQuery(query, &req)

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
		slog.Error("failed to execute search query", "err", err, "sql", sql)
		return nil, err
	}

	var tags []TagRow
	for rows.Next() {
		row, err := pgx.RowToStructByName[TagRow](rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, row)
	}

	return tags, nil
}

func CountTags(
	ctx context.Context,
	db DBTX,
	req ListTagsQuery,
) (int64, error) {

	query := postgresQuery.
		Select(
			goqu.COUNT("*").As("count"),
		).
		From("defined_tags")

	query = applyTagsQuery(query, &req)

	sql, params, err := query.ToSQL()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		slog.Error("failed to execute search query", "err", err, "sql", sql)
		return 0, err
	}

	var count int64
	if rows.Next() {
		rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}
