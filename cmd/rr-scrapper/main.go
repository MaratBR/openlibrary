package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	royalroadapi "github.com/MaratBR/openlibrary/internal/royal-road-api"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	booksParallelism := 1

	c := royalroadapi.NewClient()
	c.Run()
	defer c.Close()

	err := os.MkdirAll("rr-books", os.ModePerm)
	if err != nil {
		panic(err)
	}

	bestRatedBooks := make(chan royalroadapi.BestRatedBook, booksParallelism)
	go loadAllBestRatedBooks(bestRatedBooks, c)

	// open json file for all best rated books
	f, err := os.Create("rr-best-rated-books.json")
	if err != nil {
		slog.Error("failed to open file", "err", err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	downloadBooks := make(chan royalroadapi.BestRatedBook)
	dwnldWg := downloadAllBooks(downloadBooks, c, booksParallelism)

	// put all books into a file and then put them into a queue to be downloaded
	for book := range bestRatedBooks {
		err = enc.Encode(book)
		if err != nil {
			slog.Error("failed to encode book", "err", err)
			return
		}
		downloadBooks <- book
	}
	close(downloadBooks)
	dwnldWg.Wait()

	// book, err := c.GetBookWithChapters(21220)
	// if err != nil {
	// 	slog.Error("failed to get book page", "err", err)
	// 	return
	// }

	// slog.Info("got book page", "title", book.Book.Name, "description", book.Book.Description)

	// for _, chapter := range book.Chapters {
	// 	slog.Info("got chapter", "content", chapter.Content)
	// }

}

func loadAllBestRatedBooks(ch chan royalroadapi.BestRatedBook, c *royalroadapi.Client) {
	defer close(ch)

	var (
		page    = 1
		maxPage = 1
	)

	for page <= maxPage {
		resp, err := c.GetBestRatedBooks(page)
		if err != nil {
			slog.Error("failed to get best rated books", "err", err)
			time.Sleep(time.Millisecond)
			continue
		}
		page++
		maxPage = resp.LastPage

		for _, book := range resp.Books {
			ch <- book
		}
	}
}

func downloadAllBooks(
	ch chan royalroadapi.BestRatedBook,
	c *royalroadapi.Client,
	workers int,
) *sync.WaitGroup {
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for book := range ch {
				fileName := fmt.Sprintf("rr-books/%d.json", book.ID)

				_, err := os.Stat(fileName)
				if err == nil {
					slog.Debug("book already downloaded", "bookID", book.ID)
					continue
				} else if !os.IsNotExist(err) {
					slog.Error("failed to check if file exists", "err", err)
					continue
				}

				slog.Debug("downloading book", "bookID", book.ID)
				book, err := c.GetBookWithChapters(book.ID)
				if err != nil {
					slog.Error("failed to get book with chapters", "err", err)
					continue
				}

				f, err := os.Create(fileName)
				if err != nil {
					slog.Error("failed to create file", "err", err, "bookID", book.Book.ID)
					continue
				}
				defer f.Close()

				enc := json.NewEncoder(f)
				enc.SetIndent("", "  ")
				err = enc.Encode(book)
				if err != nil {
					slog.Error("failed to encode book", "err", err, "bookID", book.Book.ID)
					continue
				}

				time.Sleep(time.Millisecond)
			}
			wg.Done()
		}()
	}
	return wg
}
