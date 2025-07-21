package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type libraryController struct {
	service app.ReadingListService
}

func newLibraryController(service app.ReadingListService) *libraryController {
	return &libraryController{service: service}
}

func (c *libraryController) Register(r chi.Router) {
	r.Get("/library/", c.index)
	r.Get("/library/archive", c.archive)
}

func (c *libraryController) index(w http.ResponseWriter, r *http.Request) {
	session, isAuthorized := auth.GetSession(r.Context())

	if !isAuthorized {
		templates.LibraryAnon().Render(r.Context(), w)
		return
	}

	wantToRead, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusWantToRead,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	reading, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusReading,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	paused, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusPaused,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.Library(wantToRead, reading, paused).Render(r.Context(), w)
}

func (c *libraryController) archive(w http.ResponseWriter, r *http.Request) {
	session, isAuthorized := auth.GetSession(r.Context())

	if !isAuthorized {
		templates.LibraryAnon().Render(r.Context(), w)
		return
	}

	read, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusRead,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	dnf, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusDnf,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.LibraryArchive(read, dnf).Render(r.Context(), w)
}
