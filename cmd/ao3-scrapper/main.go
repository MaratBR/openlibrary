package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/ao3import"
)

var (
	maxPages     int
	url          string
	outputFolder string
	logLevel     string
	urls         []string
)

func main() {
	urls = []string{
		"https://archiveofourown.org/works?commit=Sort+and+Filter&work_search%5Bsort_column%5D=revised_at&work_search%5Bother_tag_names%5D=&exclude_work_search%5Brating_ids%5D%5B%5D=13&exclude_work_search%5Barchive_warning_ids%5D%5B%5D=20&exclude_work_search%5Barchive_warning_ids%5D%5B%5D=19&work_search%5Bexcluded_tag_names%5D=&work_search%5Bcrossover%5D=&work_search%5Bcomplete%5D=&work_search%5Bwords_from%5D=&work_search%5Bwords_to%5D=&work_search%5Bdate_from%5D=&work_search%5Bdate_to%5D=&work_search%5Bquery%5D=&work_search%5Blanguage_id%5D=&tag_id=Choose+Not+To+Use+Archive+Warnings",
	}

	flag.IntVar(&maxPages, "max-pages", 5000, "max pages to scrape")
	flag.StringVar(&outputFolder, "output-folder", "ao3-books", "output folder")
	flag.StringVar(&logLevel, "log", "info", "log level, possible values: debug, info, warn, error")
	flag.Parse()

	initLogger()

	// create output folder
	err := os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	client := ao3import.NewClient()
	client.Run()
	defer client.Close()

	chIds := make(chan string, 100)
	for i := 0; i < 10; i++ {
		go downloadBooks(chIds, client)
	}

	fmt.Fprintf(os.Stderr, "starting ao3 scraper\n")

	for _, url := range urls {
		if url == "" {
			break
		}

		slog.Info("scraping url", "url", url, "max-pages", maxPages)

		for {
			err = client.ScrapeBookIDs(url, maxPages, chIds)
			if err != nil {
				slog.Error("failed to scrape url", "url", url, "err", err.Error())
				time.Sleep(time.Second * 2)
				continue
			}
			break
		}

	}

	close(chIds)
	wg.Wait()

}

func downloadBooks(ch <-chan string, c *ao3import.Client) {
	for id := range ch {
		fileName := fmt.Sprintf("%s/%s.html", outputFolder, id)

		_, err := os.Stat(fileName)
		if err == nil {
			slog.Warn("skipping book - file already exists", "id", id)
			continue
		}

		resp, err := c.DownloadBook(id)
		if err != nil {
			slog.Error("failed to get book", "id", id, "err", err.Error())
			continue
		}
		defer resp.Body.Close()

		f, err := os.Create(fileName)
		if err != nil {
			slog.Error("failed to create file", "err", err, "id", id)
			continue
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			slog.Error("failed to copy to file", "err", err, "id", id)
			continue
		}
	}
}

func initLogger() {

	var (
		level slog.Level
	)

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		fmt.Fprintf(os.Stderr, "unknown log level: %s\n", logLevel)
		level = slog.LevelDebug
		break
	}

	fmt.Fprintf(os.Stderr, "log level: %s (%d)\n", logLevel, level)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

}
