package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/gravatar"
	"github.com/MaratBR/openlibrary/internal/store"
)

type userService struct {
	queries *store.Queries
	db      DB
}

// GetUser implements UserService.
func (u *userService) GetUser(ctx context.Context, query GetUserQuery) (*UserDetailsDto, error) {
	user, err := u.queries.GetUser(ctx, uuidDomainToDb(query.ID))
	if err != nil {
		return nil, err
	}

	details := &UserDetailsDto{
		ID:             uuidDbToDomain(user.ID),
		Name:           user.Name,
		JoinedAt:       timeDbToDomain(user.JoinedAt),
		IsBanned:       false,
		IsAdmin:        false,
		HasCustomTheme: true,
		About: UserAboutDto{
			Status: "Feelin good",
			Bio:    "lorem ipsum",
			Gender: "male",
		},
	}

	details.Avatar.MD = getUserAvatar(user.Name, 84)
	details.Avatar.LG = getUserAvatar(user.Name, 256)

	return details, nil
}

func NewUserService(db DB) UserService {
	return &userService{
		queries: store.New(db),
		db:      db,
	}
}

func getUserAvatar(name string, size int) string {
	g := gravatar.NewGravatarFromEmail(name)
	g.Size = size
	return g.GetURL()
}
