package app

import "fmt"

func getBookCoverURL(
	uploadService *UploadService,
	bookID int64,
	hasCover bool,
) string {
	if !hasCover {
		return ""
	}

	return uploadService.GetPublicURL(fmt.Sprintf("book-covers/%d.jpg", bookID))
}
