package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/ggicci/httpin"
)

type apiPayloadGetBooks struct {
	Page   uint32 `in:"query=page"`
	Size   uint32 `in:"query=size"`
	Search string `in:"query=search"`
}

type apiResponseGetBooks struct {
	Books      []app.ManagerBookDto `json:"books"`
	TotalPages uint32               `json:"totalPages"`
	Page       uint32               `json:"page"`
}

func (c *apiControllerBM) getBooks(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*apiPayloadGetBooks)

	s := auth.RequireSession(r.Context())
	booksResult, err := c.service.GetUserBooks(r.Context(), app.ManagerGetUserBooksQuery{
		UserID:      s.UserID,
		ActorUserID: app.Value(s.UserID),
		Page:        input.Page,
		PageSize:    input.Size,
		SearchQuery: input.Search,
	})
	if err != nil {
		olhttp.NewAPIError(err).Write(w)
		return
	}

	resp := apiResponseGetBooks{
		Books:      booksResult.Books,
		TotalPages: booksResult.TotalPages,
		Page:       booksResult.Page,
	}

	olhttp.NewAPIResponse(resp).Write(w)
}

type apiPayloadTrashBook struct {
	ID    int64 `in:"query=id"`
	Trash bool  `in:"query=trash"`
}

func (c *apiControllerBM) trashBook(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*apiPayloadTrashBook)

	s := auth.RequireSession(r.Context())

	var err error

	if input.Trash {
		err = c.service.TrashBook(r.Context(), app.TrashBookCommand{
			ActorUserID: s.UserID,
			BookID:      input.ID,
		})
	} else {
		err = c.service.UntrashBook(r.Context(), app.UntrashBookCommand{
			ActorUserID: s.UserID,
			BookID:      input.ID,
		})
	}

	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	apiWriteOK(w)
}
