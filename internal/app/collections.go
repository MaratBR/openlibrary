package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type GetRecentCollections struct {
	UserID uuid.UUID
	Limit  int32
}

type CollectionDto struct {
	ID            int64
	Name          string
	BooksCount    int
	LastUpdatedAt Nullable[time.Time]
}

type CreateCollectionCommand struct {
	Name   string
	UserID uuid.UUID
}

type AddToCollectionsCommand struct {
	ActorUserID  uuid.UUID
	CollectionID []int64
	BookID       int64
}

type GetBookCollectionsQuery struct {
	ActorUserID uuid.UUID
	BookID      int64
}

type CollectionService interface {
	GetRecentUserCollections(ctx context.Context, query GetRecentCollections) ([]CollectionDto, error)
	CreateCollection(ctx context.Context, cmd CreateCollectionCommand) (int64, error)
	AddToCollections(ctx context.Context, cmd AddToCollectionsCommand) error
	GetBookCollections(ctx context.Context, query GetBookCollectionsQuery) ([]CollectionDto, error)
}
