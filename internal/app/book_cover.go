package app

import "fmt"

type BookCover struct {
	URL string `json:"url"`
}

func getBookCover(
	uploadService *UploadService,
	coverID string,
	bookID int64,
) BookCover {
	var url string

	if coverID != "" {
		url = uploadService.GetPublicURL(fmt.Sprintf("%s/%s.jpeg", BOOK_COVER_DIRECTORY, coverID))
	} else {
		url = fmt.Sprintf("/_/embed-assets/cover/%d.h300.webp", (bookID%5)+1)
	}

	return BookCover{
		URL: url,
	}
}
