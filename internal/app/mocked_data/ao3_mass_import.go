package mockeddata

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/MaratBR/openlibrary/internal/ao3import"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/gofrs/uuid"
)

// Ao3 mass import of fanfics
// because I need actual real user-generated content
// for testing, not stealing anything pinky promise

func MassImportAo3(
	bookManagerService app.BookManagerService,
	tagsService app.TagsService,
	location string,
	userIds []uuid.UUID,
) ([]int64, error) {
	// step 1: get list of all json files in the directory using golang stdlib
	files, err := os.ReadDir(location)
	if err != nil {
		return nil, err
	}

	inputChannel := make(chan string, len(files))
	filesCount := 0

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			inputChannel <- file.Name()
			filesCount++
		}
	}

	close(inputChannel)

	// step 2: read each book

	ids := make([]int64, 0, filesCount)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	userIdIndex := 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for htmlFile := range inputChannel {
				file, err := os.Open(fmt.Sprintf("%s/%s", location, htmlFile))
				if err != nil {
					slog.Error("failed to open file", "file", htmlFile, "err", err.Error())
					continue
				}
				defer file.Close()

				book, err := ao3import.ParseBook(file)
				if err != nil {
					slog.Error("failed to parse book", "file", htmlFile, "err", err.Error())
					continue
				}

				bookId, err := createBook(bookManagerService, tagsService, book, userIds[userIdIndex])
				if err != nil {
					slog.Error("failed to create book", "file", htmlFile, "err", err.Error())
					continue
				}

				userIdIndex++
				if userIdIndex >= len(userIds) {
					userIdIndex = 0
				}

				slog.Info("imported book", "file", htmlFile, "id", bookId)

				mutex.Lock()
				ids = append(ids, bookId)
				mutex.Unlock()

			}
		}()
	}

	wg.Wait()

	return ids, nil
}

func createBook(
	bookManagerService app.BookManagerService,
	tagsService app.TagsService,
	book *ao3import.Ao3Book,
	userId uuid.UUID,
) (int64, error) {
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

	newTags := []app.TagDescriptor{}

	for _, fandom := range fandoms {
		newTags = append(newTags, app.TagDescriptor{
			Name:     fandom,
			Category: app.TagsCategoryFandom,
		})
	}

	for _, character := range characters {
		newTags = append(newTags, app.TagDescriptor{
			Name:     character,
			Category: app.TagsCategoryOther,
		})
	}

	for _, other := range other {
		newTags = append(newTags, app.TagDescriptor{
			Name:     other,
			Category: app.TagsCategoryOther,
		})
	}

	for _, relationship := range relationships {
		newTags = append(newTags, app.TagDescriptor{
			Name:     relationship,
			Category: app.TagsCategoryRelationship,
		})
	}

	for _, relationshipType := range relationshipTypes {
		newTags = append(newTags, app.TagDescriptor{
			Name:     relationshipType,
			Category: app.TagsCategoryRelationshipType,
		})
	}

	for _, warning := range warnings {
		newTags = append(newTags, app.TagDescriptor{
			Name:      warning,
			Category:  app.TagsCategoryWarning,
			IsWarning: true,
		})
	}

	tags, err := tagsService.CreateTags(context.Background(), app.CreateTagsCommand{
		Tags: newTags,
	})
	if err != nil {
		return 0, err
	}

	bookID, err := bookManagerService.CreateBook(context.Background(), app.CreateBookCommand{
		UserID:            userId,
		Name:              book.Name,
		Tags:              commonutil.MapSlice(tags, func(tag app.DefinedTagDto) int64 { return tag.ID }),
		AgeRating:         getAgeRatingFromAo3Book(book.Rating),
		Summary:           book.Summary,
		IsPubliclyVisible: true,
	})
	if err != nil {
		return 0, err
	}

	for _, chapter := range book.Chapters {
		_, err = bookManagerService.CreateBookChapter(context.Background(), app.CreateBookChapterCommand{
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

func getAgeRatingFromAo3Book(ao3Rating string) app.AgeRating {
	switch ao3Rating {
	case "Explicit":
		return app.AgeRatingNC17
	case "Mature":
		return app.AgeRatingR
	case "Teen And Up Audiences":
		return app.AgeRatingPG13
	case "General Audiences":
		return app.AgeRatingG
	default:
		return app.AgeRatingUnknown
	}
}
