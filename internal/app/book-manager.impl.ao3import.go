package app

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/MaratBR/openlibrary/internal/ao3import"
	"github.com/gofrs/uuid"
)

type ManagerCreateBookFromAo3Command struct {
	Ao3ID  string
	UserID uuid.UUID
}

func (s *bookManagerService) ImportFromBookAo3(ctx context.Context, command ManagerCreateBookFromAo3Command) (int64, error) {
	client := ao3import.NewClient()

	if strings.Trim(command.Ao3ID, " \n\t") == "" {
		return 0, errors.New("ao3 is is empty")
	}

	book, err := client.GetBook(command.Ao3ID)
	if err != nil {
		return 0, err
	}

	fandoms := book.Tags["Fandom"]
	characters, ok := book.Tags["Characters"]
	if !ok {
		characters = []string{}
	}
	other, ok := book.Tags["Additional Tags"]
	if !ok {
		other = []string{}
	}
	relationships, ok := book.Tags["Relationships"]
	if !ok {
		characters = []string{}
	}
	relationshipTypes, ok := book.Tags["Category"]
	if !ok {
		relationshipTypes = []string{}
	}
	warnings, ok := book.Tags["Archive Warning"]
	if !ok {
		warnings = []string{}
	}

	newTags := []TagDescriptor{}

	for _, fandom := range fandoms {
		newTags = append(newTags, TagDescriptor{
			Name:     fandom,
			Category: TagsCategoryFandom,
		})
	}

	for _, character := range characters {
		newTags = append(newTags, TagDescriptor{
			Name:     character,
			Category: TagsCategoryOther,
		})
	}

	for _, other := range other {
		newTags = append(newTags, TagDescriptor{
			Name:     other,
			Category: TagsCategoryOther,
		})
	}

	for _, relationship := range relationships {
		newTags = append(newTags, TagDescriptor{
			Name:     relationship,
			Category: TagsCategoryRelationship,
		})
	}

	for _, relationshipType := range relationshipTypes {
		newTags = append(newTags, TagDescriptor{
			Name:     relationshipType,
			Category: TagsCategoryRelationshipType,
		})
	}

	for _, warning := range warnings {
		newTags = append(newTags, TagDescriptor{
			Name:      warning,
			Category:  TagsCategoryWarning,
			IsWarning: true,
		})
	}

	tags, err := s.tagsService.CreateTags(ctx, CreateTagsCommand{
		Tags: newTags,
	})
	if err != nil {
		return 0, err
	}

	bookID, err := s.CreateBook(ctx, CreateBookCommand{
		UserID:            command.UserID,
		Name:              book.Name,
		Tags:              mapSlice(tags, func(tag DefinedTagDto) int64 { return tag.ID }),
		AgeRating:         getAgeRatingFromAo3Book(book.Rating),
		Summary:           book.Summary,
		IsPubliclyVisible: true,
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
