package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
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
	var totalSize int64
	var activeWorkers int32

	pausedNotifCh := make(chan struct{})
	pauseChannels := make([]chan struct{}, 0, workers)
	resumeChannels := make([]chan struct{}, 0, workers)

	for range workers {
		activeWorkers++
		pauseChannel := make(chan struct{})
		resumeChannel := make(chan struct{})

		pauseChannels = append(pauseChannels, pauseChannel)
		resumeChannels = append(resumeChannels, resumeChannel)

		wg.Add(1)
		go func() {
			defer atomic.AddInt32(&activeWorkers, -1)
			for {
				select {
				case <-pauseChannel:
					pausedNotifCh <- struct{}{}
					// pause requested - waiting until all
					<-resumeChannel
					continue
				default:
				}

				bookEntry, ok := <-ch
				if !ok {
					break
				}
				fileName := fmt.Sprintf("rr-books/%d.json", bookEntry.ID)

				_, err := os.Stat(fileName)
				if err == nil {
					slog.Debug("book already downloaded", "bookID", bookEntry.ID)
					continue
				} else if !os.IsNotExist(err) {
					slog.Error("failed to check if file exists", "err", err)
					continue
				}

				slog.Debug("downloading book", "bookID", bookEntry.ID)
				book, err := c.GetBookWithChapters(bookEntry.ID)
				if err != nil {
					slog.Error("failed to get book with chapters", "err", err)
					continue
				}

				f, err := os.Create(fileName)
				if err != nil {
					slog.Error("failed to create file", "err", err, "bookID", book.Book.ID)
					continue
				}

				enc := json.NewEncoder(f)
				enc.SetIndent("", "  ")
				err = enc.Encode(book)
				f.Close()
				if err != nil {
					slog.Error("failed to encode book", "err", err, "bookID", book.Book.ID)
					continue
				}

				fileInfo, err := os.Stat(fileName)
				if err != nil {
					slog.Error("failed to stat file", "file", fileName, "err", err)
				} else {
					atomic.AddInt64(&totalSize, fileInfo.Size())
				}

				time.Sleep(time.Millisecond)
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				time.Sleep(time.Second * 10)
				activeWorkersCurrent := atomic.LoadInt32(&activeWorkers)
				if activeWorkersCurrent == 0 {
					break
				}

				totalSizeCurrent := atomic.LoadInt64(&totalSize)
				if totalSizeCurrent < 2_000_000_000 {
					continue
				}

				// request all works pause
				slog.Info("reached file threshold, requesting all workers pause")
				for _, ch := range pauseChannels {
					ch <- struct{}{}
				}

				// wait until all workers are paused
				for range pauseChannels {
					<-pausedNotifCh
				}

				slog.Info("all workers are paused gzipping JSON files...")
				time.Sleep(time.Second)
				tarGzAllFiles()
				atomic.StoreInt64(&totalSize, 0)

				slog.Info("gzip done, resuming workers")
				for _, ch := range resumeChannels {
					ch <- struct{}{}
				}

				time.Sleep(time.Second)
			}
		}()
	}

	return wg
}

func tarGzAllFiles() {
	err := runCommand("mkdir", "-p", "tar")
	if err != nil {
		slog.Error("Error running mkdir", "err", err)
	}

	files, err := filepath.Glob("./rr-books/*")
	if err != nil {
		slog.Error("Error finding files", "err", err)
	}

	if len(files) == 0 {
		slog.Error("No files found to archive")
		return
	}

	tarFile := fmt.Sprintf("tar/%s.tar.gz", time.Now().Format(time.RFC3339))
	args := append([]string{"-czf", tarFile}, files...)
	err = runCommand("tar", args...)
	if err != nil {
		slog.Error("Error running tar", "err", err)
	}

	for _, file := range files {
		os.Remove(file)
	}

}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
