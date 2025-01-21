package mockeddata

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	royalroadapi "github.com/MaratBR/openlibrary/internal/royal-road-api"
	"github.com/gofrs/uuid"
)

func massImport(ctx context.Context, dir string, userIds []uuid.UUID, service app.BookManagerService, tagsService app.TagsService) error {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		file, err := os.Open(fmt.Sprintf("%s/%s", dir, entry.Name()))
		if err != nil {
			return err
		}

		var book royalroadapi.BookWithChapters
		err = json.NewDecoder(file).Decode(&book)
		if err != nil {
			return err
		}

		randInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userIds))))
		randomUserId := userIds[int(randInt.Int64())]
		err = importBook(ctx, service, tagsService, book, randomUserId)

		if err != nil {
			return err
		}
	}

	return nil
}

func importBook(
	ctx context.Context,
	service app.BookManagerService,
	tagsService app.TagsService,
	book royalroadapi.BookWithChapters,
	userID uuid.UUID,
) error {
	tags, err := tagsService.CreateTags(ctx, app.CreateTagsCommand{
		Tags: commonutil.MapSlice(book.Book.Tags, func(tag string) app.TagDescriptor {
			return app.TagDescriptor{
				Name:     tag,
				Category: app.TagsCategoryOther,
			}
		}),
	})

	if err != nil {
		return err
	}

	tagIds := commonutil.MapSlice(tags, func(tag app.DefinedTagDto) int64 { return tag.ID })

	bookID, err := service.CreateBook(ctx, app.CreateBookCommand{
		Name:              book.Book.Name,
		UserID:            userID,
		Summary:           fmt.Sprintf("<p>Original book: https://www.royalroad.com/fiction/%d</p>%s", book.Book.ID, book.Book.Description),
		IsPubliclyVisible: true,
		AgeRating:         app.AgeRatingPG,
		Tags:              tagIds,
	})

	if err != nil {
		return err
	}

	for i := 0; i < len(book.Chapters); i++ {
		chapter := book.Book.Chapters[i]
		_, err = service.CreateBookChapter(ctx, app.CreateBookChapterCommand{
			BookID:          bookID,
			Name:            chapter.Title,
			Content:         book.Chapters[i].Content,
			IsAdultOverride: false,
			Summary:         fmt.Sprintf("Original chapter: https://www.royalroad.com/fiction/%d/chapter/%d", book.Book.ID, chapter.ID),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
