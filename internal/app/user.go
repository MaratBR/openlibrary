package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type UserDetailsDto struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar struct {
		LG string `json:"lg"`
		MD string `json:"md"`
	} `json:"avatar"`
	JoinedAt       time.Time       `json:"joinedAt"`
	IsBanned       bool            `json:"isBlocked"`
	IsAdmin        bool            `json:"isAdmin"`
	HasCustomTheme bool            `json:"hasCustomTheme"`
	About          UserAboutDto    `json:"about"`
	Books          []AuthorBookDto `json:"books"`
}

type UserAboutDto struct {
	Status string `json:"status"`
	Bio    string `json:"bio"`
	Gender string `json:"gender"`
}

type GetUserQuery struct {
	ID     uuid.UUID
	UserID uuid.NullUUID
}

type UserService interface {
	GetUser(ctx context.Context, query GetUserQuery) (*UserDetailsDto, error)
}
