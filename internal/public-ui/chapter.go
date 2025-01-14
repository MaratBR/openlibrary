package publicui

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/public-ui/templates"
)

type chaptersController struct {
	service app.BookService
}

func newChaptersController(service app.BookService) *chaptersController {
	return &chaptersController{service: service}
}

func (c *chaptersController) GetChapter(w http.ResponseWriter, r *http.Request) {

	bookID, err := commonutil.URLParamInt64(r, "bookID")
	if err != nil {
		writeRequestError(w, r, err)
		return
	}

	chapterID, err := commonutil.URLParamInt64(r, "chapterID")
	if err != nil {
		writeRequestError(w, r, err)
		return
	}

	userID := auth.GetNullableUserID(r.Context())
	book, err := c.service.GetBook(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	chapter, err := c.service.GetBookChapter(r.Context(), app.GetBookChapterQuery{
		BookID:    bookID,
		ChapterID: chapterID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	uiSettings := getUIBookSettings(r)
	templates.Chapter(r.Context(), chapter.Chapter, book, templates.ReaderSettings{
		FontSize: uiSettings.FontSize,
	}).Render(r.Context(), w)
}
