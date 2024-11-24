package app

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

type Int32 struct {
	Valid bool
	Int32 int32
}

func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

func (i *Int32) UnmarshalJSON(b []byte) error {
	v := new(int32)
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*i = Int32{}
	} else {
		*i = Int32{Valid: true, Int32: *v}
	}
	return nil
}

type Int32Range struct {
	Min Int32 `json:"min"`
	Max Int32 `json:"max"`
}

type BookSearchQuery struct {
	UserID uuid.NullUUID

	Words           Int32Range
	Chapters        Int32Range
	WordsPerChapter Int32Range
	Favorites       Int32Range

	IncludeTags  []string
	ExcludeTags  []string
	IncludeUsers []uuid.UUID
	ExcludeUsers []uuid.UUID

	IncludeBanned bool
	IncludeHidden bool
	IncludeEmpty  bool

	Limit  uint
	Offset uint
}

type BookSearchItem struct {
	ID              int64                `json:"id,string"`
	Name            string               `json:"name"`
	CreatedAt       time.Time            `json:"createdAt"`
	AgeRating       AgeRating            `json:"ageRating"`
	Words           int                  `json:"words"`
	WordsPerChapter int                  `json:"wordsPerChapter"`
	Chapters        int                  `json:"chapters"`
	Summary         string               `json:"summary"`
	Favorites       int32                `json:"favorites"`
	Author          BookDetailsAuthorDto `json:"author"`
	Cover           string               `json:"cover"`
	Tags            []DefinedTagDto      `json:"tags"`
	Collections     []BookCollectionDto  `json:"collections"`
}

type BookSearchResult struct {
	TookUS int64 `json:"took"`
	Cache  struct {
		Key string `json:"key"`
		Hit bool   `json:"hit"`
	} `json:"cache"`
	Books []BookSearchItem `json:"books"`
}

type BookExtremes struct {
	Words           Int32Range `json:"words"`
	Chapters        Int32Range `json:"chapters"`
	WordsPerChapter Int32Range `json:"wordsPerChapter"`
	Favorites       Int32Range `json:"favorites"`
}

type SearchService interface {
	SearchBooks(ctx context.Context, req BookSearchQuery) (*BookSearchResult, error)
	GetBookExtremes(ctx context.Context) (*BookExtremes, error)
}
