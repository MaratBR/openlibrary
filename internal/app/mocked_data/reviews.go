package mockeddata

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

//go:embed reviews/short.json
var shortJSON []byte

//go:embed reviews/long.json
var longJSON []byte

type review struct {
	Rating  app.RatingValue `json:"rating"`
	Content string          `json:"content"`
}

func getAllReviews() []review {
	arr := []review{}

	tempArr := []review{}
	if err := json.Unmarshal(shortJSON, &tempArr); err != nil {
		panic(err)
	}
	arr = append(arr, tempArr...)

	tempArr = []review{}
	if err := json.Unmarshal(longJSON, &arr); err != nil {
		panic(err)
	}
	arr = append(arr, tempArr...)

	return arr
}

func CreateReviews(
	service app.ReviewsService,
	userIDs []uuid.UUID,
	bookID int64,
) error {
	reviews := getAllReviews()
	if len(reviews) == 0 {
		return errors.New("no reviews embedded in the binary")
	}
	i := 0
	for _, userID := range userIDs {
		review := reviews[i]
		i++
		if i == len(reviews) {
			i = 0
		}
		if _, err := service.UpdateReview(context.Background(), app.UpdateReviewCommand{
			UserID:  userID,
			BookID:  bookID,
			Rating:  review.Rating,
			Content: review.Content,
		}); err != nil {
			return err
		}
	}

	return nil
}
