package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Book struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	CoverURL    string `json:"CoverURL"`
	Description string `json:"Description"`
	Author      struct {
		Name string `json:"Name"`
		ID   int    `json:"ID"`
	} `json:"Author"`
	Chapters []Chapter `json:"Chapters"`
	Tags     []string  `json:"Tags"`
}

type Chapter struct {
	ID         int    `json:"ID"`
	Title      string `json:"Title"`
	ReleasedAt string `json:"ReleasedAt"`
	Content    string `json:"Content,omitempty"`
}

type BookData struct {
	Book     Book      `json:"Book"`
	Chapters []Chapter `json:"Chapters"`
}

type ScrapedData struct {
	AuthorImage string `json:"author_image,omitempty"`
	AuthorName  string `json:"author_name,omitempty"`
	AuthorURL   string `json:"author_url,omitempty"`
}

func main() {
	// Create a rate limiter channel
	rateLimiter := time.NewTicker(250 * time.Millisecond) // 4 requests per second
	defer rateLimiter.Stop()

	// Create HTTP client
	client := &http.Client{}

	// Read all JSON files from rr-books directory
	files, err := filepath.Glob("rr-books/*.json")
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		// Read the JSON file
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}

		var bookData BookData
		if err := json.Unmarshal(data, &bookData); err != nil {
			fmt.Printf("Error parsing JSON in %s: %v\n", file, err)
			continue
		}

		// Skip if we already have author information
		if bookData.Book.Author.Name != "" && bookData.Book.Author.ID != 0 {
			fmt.Printf("Skipping %s - already has author information\n", file)
			continue
		}

		// Wait for rate limiter
		<-rateLimiter.C

		// Fetch the webpage
		url := fmt.Sprintf("https://www.royalroad.com/fiction/%d", bookData.Book.ID)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
			continue
		}

		// Parse the HTML
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error parsing HTML for %s: %v\n", url, err)
			continue
		}

		// Extract author information
		var scrapedData ScrapedData

		// Find author image
		doc.Find(".fic-header img.thumbnail.inline-block").Each(func(i int, s *goquery.Selection) {
			if src, exists := s.Attr("src"); exists {
				scrapedData.AuthorImage = src
			}
		})

		// Find author name and URL
		doc.Find(".fic-header a").Each(func(i int, s *goquery.Selection) {
			if href, exists := s.Attr("href"); exists && strings.HasPrefix(href, "/profile/") {
				scrapedData.AuthorName = s.Text()
				scrapedData.AuthorURL = "https://www.royalroad.com" + href
			}
		})

		// Update the book data
		if scrapedData.AuthorName != "" {
			// Extract author ID from URL
			var authorID int
			if strings.Contains(scrapedData.AuthorURL, "/profile/") {
				fmt.Sscanf(scrapedData.AuthorURL, "https://www.royalroad.com/profile/%d", &authorID)
			}

			bookData.Book.Author.Name = scrapedData.AuthorName
			bookData.Book.Author.ID = authorID
		}

		// Write back to file
		updatedData, err := json.MarshalIndent(bookData, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling updated data for %s: %v\n", file, err)
			continue
		}

		if err := os.WriteFile(file, updatedData, 0644); err != nil {
			fmt.Printf("Error writing updated data to %s: %v\n", file, err)
			continue
		}

		fmt.Printf("Successfully updated %s\n", file)
	}
}
