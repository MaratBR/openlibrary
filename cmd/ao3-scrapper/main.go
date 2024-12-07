package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app/ao3import"
)

var (
	maxPages int
	url      string
)

func main() {
	flag.IntVar(&maxPages, "max-pages", 1, "max pages to scrape")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	url, _ := reader.ReadString('\n')
	url = strings.TrimSuffix(url, "\n")

	var wg sync.WaitGroup

	client := ao3import.NewClient()
	out := make(chan string, 10000)

	go func() {
		wg.Add(1)
		for {
			id, ok := <-out
			if !ok {
				break
			}
			fmt.Printf("%s\n", id)
		}
		wg.Done()
	}()

	wg.Add(1)
	err := client.ScrapeBookIDs(url, maxPages, out)
	if err != nil {
		panic(err)
	}
	wg.Done()
	close(out)

	wg.Wait()

}
