package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type BookDetailsDto struct {
	ID              int64                 `json:"id,string"`
	Name            string                `json:"name"`
	AgeRating       AgeRating             `json:"ageRating"`
	IsAdult         bool                  `json:"adult"`
	Tags            []DefinedTagDto       `json:"tags"`
	Words           int                   `json:"words"`
	WordsPerChapter int                   `json:"wordsPerChapter"`
	CreatedAt       time.Time             `json:"createdAt"`
	Collections     []BookCollectionDto   `json:"collections"`
	Author          BookDetailsAuthorDto  `json:"author"`
	Permissions     BookUserPermissions   `json:"permissions"`
	Summary         string                `json:"summary"`
	Favorites       int32                 `json:"favorites"`
	IsFavorite      bool                  `json:"isFavorite"`
	Notifications   []GenericNotification `json:"notifications,omitempty"`
	Cover           string                `json:"cover"`
	Rating          Nullable[float64]     `json:"rating"`
	Votes           int32                 `json:"votes"`
	Reviews         int32                 `json:"reviews"`
}

type GetBookQuery struct {
	ID          int64
	ActorUserID uuid.NullUUID
}

type GetBookChaptersQuery struct {
	ID int64
}

type BookChapterDto struct {
	ID        int64     `json:"id,string"`
	Order     int       `json:"order"`
	Name      string    `json:"name"`
	Words     int       `json:"words"`
	CreatedAt time.Time `json:"createdAt"`
	Summary   string    `json:"summary"`
}

type BookDetailsAuthorDto struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BookUserPermissions struct {
	CanEdit bool `json:"canEdit"`
}

type BookService interface {
	GetBook(ctx context.Context, query GetBookQuery) (BookDetailsDto, error)
	GetBookChapters(ctx context.Context, query GetBookChaptersQuery) ([]BookChapterDto, error)
	GetBookChapter(ctx context.Context, query GetBookChapterQuery) (GetBookChapterResult, error)
}
