package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type Int4Range struct {
	Min pgtype.Int4
	Max pgtype.Int4
}

type NullInt4Range struct {
	Range Int4Range
	Valid bool
}

func applyRange(builder squirrel.SelectBuilder, column string, int4Range Int4Range) squirrel.SelectBuilder {
	if int4Range.Min.Valid {
		builder = builder.Where(sq.GtOrEq{column: int4Range.Min.Int32})
	}
	if int4Range.Min.Valid {
		builder = builder.Where(sq.LtOrEq{column: int4Range.Max.Int32})
	}

	return builder
}

type BookSearchRequest struct {
	Words           NullInt4Range
	WordsPerChapter NullInt4Range

	IncludeAuthors []pgtype.UUID
	ExcludeAuthors []pgtype.UUID

	IncludeParentTags []int64
	ExcludeParentTags []int64

	Limit  uint64
	Offset uint64

	IncludeBanned bool
	IncludeHidden bool
	IncludeEmpty  bool
}

func createBookSearchSelect(req *BookSearchRequest) sq.SelectBuilder {
	if req == nil {
		panic("BookSearchRequest is nil")
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.
		Select("books.*", "author.name as author_name").
		From("books").
		Join("users author on author.id = books.author_user_id")

	if req.Words.Valid {
		query = applyRange(query, "words", req.Words.Range)
	}

	if req.WordsPerChapter.Valid {
		query = applyRange(query, "cast(words as real) / chapters", req.WordsPerChapter.Range)
	}

	if len(req.IncludeAuthors) > 0 {
		query = query.Where(sq.Eq{"author_user_id": req.IncludeAuthors})
	}

	if len(req.ExcludeAuthors) > 0 {
		query = query.Where(sq.NotEq{"author_user_id": req.ExcludeAuthors})
	}

	if len(req.IncludeParentTags) > 0 {
		query = query.Where("cached_parent_tag_ids @> ?", req.IncludeParentTags)
	}

	if !req.IncludeBanned {
		query = query.Where("is_banned = false")
	}

	if !req.IncludeHidden {
		query = query.Where("is_publicly_visible = true")
	}

	if !req.IncludeEmpty {
		query = query.Where("chapters > 0")
	}

	query = query.Limit(req.Limit).Offset(req.Offset)

	return query
}

type BookSearchRow struct {
	ID                 int64              `db:"id"`
	Name               string             `db:"name"`
	Summary            string             `db:"summary"`
	CreatedAt          pgtype.Timestamptz `db:"created_at"`
	AgeRating          AgeRating          `db:"age_rating"`
	Words              int                `db:"words"`
	Chapters           int                `db:"chapters"`
	TagIds             []int64            `db:"tag_ids"`
	CachedParentTagIds []int64            `db:"cached_parent_tag_ids"`
	AuthorUserID       pgtype.UUID        `db:"author_user_id"`
	AuthorName         string             `db:"author_name"`
}

func SearchBooks(ctx context.Context, db DBTX, req BookSearchRequest) ([]BookSearchRow, error) {
	query := createBookSearchSelect(&req).RunWith(wrapDBTX(db, ctx))
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	var books []BookSearchRow
	err = sqlx.StructScan(rows, &books)
	if err != nil {
		return nil, err
	}
	return books, nil
}
