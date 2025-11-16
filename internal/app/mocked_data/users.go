package mockeddata

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/go-faker/faker/v4"
	"github.com/gofrs/uuid"
)

type fakeUser struct {
	ID       uuid.UUID
	IPV4     string `faker:"ipv4"`
	Name     string `faker:"username"`
	Password string `faker:"password"`
}

func CreateUsers(
	service app.AuthService,
	signUpService app.SignUpService,
	count int,
) ([]uuid.UUID, error) {
	var user fakeUser
	userIds := make([]uuid.UUID, count)

	for i := 0; i < count; i++ {
		err := faker.FakeData(&user)
		if err != nil {
			panic(err)
		}
		user.ID, err = uuid.NewV4()
		if err != nil {
			panic(err)
		}

		result, err := signUpService.SignUp(context.Background(), app.SignUpCommand{
			Username:                  user.Name,
			Password:                  user.Password,
			UserAgent:                 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
			IpAddress:                 user.IPV4,
			BypassEmailRequirement:    true,
			BypassPasswordRequirement: true,
		})
		if err != nil {
			return nil, err
		}

		userIds[i] = result.CreatedUserID
	}

	return userIds, nil
}
