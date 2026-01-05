package mockeddata

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	royalroadapi "github.com/MaratBR/openlibrary/internal/royalroadapi"
	"github.com/gofrs/uuid"
)

func massImport(ctx context.Context, dir string, userIds []uuid.UUID, service app.BookManagerService, tagsService app.TagsService, workers int) error {
	var ch chan os.DirEntry
	{
		entries, err := os.ReadDir(dir)

		if err != nil {
			return err
		}

		if workers < 1 {
			workers = 1
		}

		ch = make(chan os.DirEntry, len(entries))
		for _, entry := range entries {
			ch <- entry
		}
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				entry, ok := <-ch
				if !ok {
					break
				}

				if entry.IsDir() {
					continue
				}

				if !strings.HasSuffix(entry.Name(), ".json") {
					continue
				}

				fileName := fmt.Sprintf("%s/%s", dir, entry.Name())
				file, err := os.Open(fileName)
				if err != nil {
					slog.Error("failed to open file", "file", fileName, "err", err)
					continue
				}

				var book royalroadapi.BookWithChapters
				err = json.NewDecoder(file).Decode(&book)
				if err != nil {
					slog.Error("failed to decode file", "file", fileName, "err", err)
					continue
				}

				randInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userIds))))
				randomUserId := userIds[int(randInt.Int64())]
				err = importBook(ctx, service, tagsService, book, randomUserId)

				if err != nil {
					slog.Error("failed to import book", "file", fileName, "err", err)
					continue
				}
			}
		}()
	}
	close(ch)
	wg.Wait()

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
			BookID:            bookID,
			Name:              chapter.Title,
			Content:           book.Chapters[i].Content,
			IsAdultOverride:   false,
			Summary:           fmt.Sprintf("Original chapter: https://www.royalroad.com/fiction/%d/chapter/%d", book.Book.ID, chapter.ID),
			IsPubliclyVisible: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
