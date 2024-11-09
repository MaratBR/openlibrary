package app

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app/ao3import"
	"github.com/gofrs/uuid"
)

type ManagerCreateBookFromAo3Command struct {
	Ao3ID  string
	UserID uuid.UUID
}

func (s *BookManagerService) ImportFromBookAo3(ctx context.Context, command ManagerCreateBookFromAo3Command) (int64, error) {
	client := ao3import.NewClient()

	if strings.Trim(command.Ao3ID, " \n\t") == "" {
		return 0, errors.New("ao3 is is empty")
	}

	book, err := client.GetBook(command.Ao3ID)
	if err != nil {
		return 0, err
	}

	bookID, err := s.CreateBook(ctx, CreateBookCommand{
		UserID:            command.UserID,
		Name:              book.Name,
		Tags:              []string{},
		AgeRating:         getAgeRatingFromAo3Book(book.Rating),
		Summary:           book.Summary,
		IsPubliclyVisible: false,
	})
	if err != nil {
		return 0, err
	}

	for _, chapter := range book.Chapters {
		_, err = s.CreateBookChapter(ctx, CreateBookChapterCommand{
			BookID:          bookID,
			Name:            chapter.Title,
			IsAdultOverride: false,
			Content:         chapter.Content,
			Summary:         chapter.Summary,
		})
		if err != nil {
			slog.Error("failed to import chapter", "err", err)
		}
	}

	return bookID, nil
}

func getAgeRatingFromAo3Book(ao3Rating string) AgeRating {
	switch ao3Rating {
	case "Explicit":
		return AgeRatingNC17
	case "Mature":
		return AgeRatingR
	case "Teen And Up Audiences":
		return AgeRatingPG13
	case "General Audiences":
		return AgeRatingG
	default:
		return AgeRatingUnknown
	}
}
