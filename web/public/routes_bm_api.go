package public

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/upload"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type apiControllerBM struct {
	service       app.BookManagerService
	fileValidator upload.FileValidator
}

func newAPIBookManagerController(service app.BookManagerService, fileValidator upload.FileValidator) *apiControllerBM {
	return &apiControllerBM{service: service, fileValidator: fileValidator}
}

func (c *apiControllerBM) Register(r chi.Router) {
	r.Route("/books-manager", func(r chi.Router) {
		r.With(httpin.NewInput(&apiPayloadUploadCover{})).Post("/book/{bookID}/cover", c.uploadCover)
		r.With(httpin.NewInput(&apiPayloadChangeChaptersOrder{})).Post("/book/{bookID}/chapters-order", c.changeChaptersOrder)
		r.With(httpin.NewInput(&apiPayloadCreateChapter{})).Post("/book/{bookID}/create-chapter", c.createChapter)
		r.Get("/book/{bookID}/chapters", c.getChapters)

		r.Post("/book/{bookID}/{chapterID}/{draftID}", c.updateDraftContent)
		r.Post("/book/{bookID}/{chapterID}/{draftID}/publish", c.updateDraftContentAndPublish)
		r.Post("/book/{bookID}/{chapterID}/{draftID}/chapterName", c.updateDraftChapterName)
		r.With(httpin.NewInput(&apiPayloadGetBooks{})).Get("/books", c.getBooks)
		r.Get("/books/{bookID}", c.getBook)
		r.With(httpin.NewInput(&apiPayloadTrashBook{})).Post("/books/trash", c.trashBook)

	})
}

type apiPayloadUploadCover struct {
	File *httpin.File `in:"form=file"`
}

func (c *apiControllerBM) uploadCover(w http.ResponseWriter, r *http.Request) {

	session := auth.RequireSession(r.Context())
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	input := r.Context().Value(httpin.Input).(*apiPayloadUploadCover)

	err = c.fileValidator.Validate(input.File.Size())
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	file, err := input.File.Open()
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	defer file.Close()

	result, err := c.service.UploadBookCover(r.Context(), app.UploadBookCoverCommand{
		UserID: session.UserID,
		BookID: bookID,
		File:   file,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	response := olhttp.NewAPIResponse(result.URL)
	response.Write(w)
}

type responseModifiedChapterOrderPositions map[int64]int

func (r responseModifiedChapterOrderPositions) MarshalJSON() ([]byte, error) {
	m := make(map[app.Int64String]int, len(r))
	for chapterID, pos := range r {
		m[app.Int64String(chapterID)] = pos
	}
	return json.Marshal(m)
}

type apiPayloadChangeChaptersOrder struct {
	Data struct {
		Modifications []struct {
			ChapterID int64 `json:"chapterId,string"`
			NewIndex  int   `json:"newIndex"`
		} `json:"modifications"`
	} `in:"body=json"`
}

func (c *apiControllerBM) changeChaptersOrder(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*apiPayloadChangeChaptersOrder)

	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	modifications := make([]app.ChapterOrderModification, 0, len(input.Data.Modifications))

	for _, modification := range input.Data.Modifications {
		modifications = append(modifications, app.ChapterOrderModification{
			ChapterID:        modification.ChapterID,
			NewPositionIndex: modification.NewIndex,
		})
	}

	result, err := c.service.UpdateBookChaptersOrder(r.Context(), app.UpdateBookChapterOrdersCommand{
		BookID:        bookID,
		Modifications: modifications,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponse(responseModifiedChapterOrderPositions(result.ModifiedPositions))
}

func (c *apiControllerBM) updateDraftContent(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	chapterID, err := olhttp.URLParamInt64(r, "chapterID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	draftID, err := olhttp.URLParamInt64(r, "draftID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		apiWriteBadRequest(w, errors.New("Content-Type must be text/plain"))
		return
	}

	contentBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	session := auth.RequireSession(r.Context())

	err = c.service.UpdateDraftContent(r.Context(), app.UpdateDraftContentCommand{
		BookID:    bookID,
		ChapterID: chapterID,
		DraftID:   draftID,
		UserID:    session.UserID,
		Content:   string(contentBytes),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	draft, err := c.service.GetDraft(r.Context(), app.GetDraftQuery{
		DraftID:   draftID,
		ChapterID: chapterID,
		BookID:    bookID,
		UserID:    session.UserID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponse(draft).Write(w)
}

func (c *apiControllerBM) updateDraftChapterName(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	chapterID, err := olhttp.URLParamInt64(r, "chapterID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	draftID, err := olhttp.URLParamInt64(r, "draftID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		apiWriteBadRequest(w, errors.New("Content-Type must be text/plain"))
		return
	}

	contentBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	session := auth.RequireSession(r.Context())

	err = c.service.UpdateDraftChapterName(r.Context(), app.UpdateDraftChapterNameCommand{
		BookID:      bookID,
		ChapterID:   chapterID,
		DraftID:     draftID,
		UserID:      session.UserID,
		ChapterName: string(contentBytes),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	apiWriteOK(w)
}

func (c *apiControllerBM) updateDraftContentAndPublish(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	chapterID, err := olhttp.URLParamInt64(r, "chapterID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	draftID, err := olhttp.URLParamInt64(r, "draftID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	makePublic := olhttp.GetBool(r.URL.Query(), "makePublic").Or(false)

	if r.Header.Get("Content-Type") != "text/plain" {
		apiWriteBadRequest(w, errors.New("Content-Type must be text/plain"))
		return
	}

	contentBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	session := auth.RequireSession(r.Context())

	err = c.service.UpdateDraftContent(r.Context(), app.UpdateDraftContentCommand{
		BookID:    bookID,
		ChapterID: chapterID,
		DraftID:   draftID,
		UserID:    session.UserID,
		Content:   string(contentBytes),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	err = c.service.PublishDraft(r.Context(), app.PublishDraftCommand{
		DraftID:    draftID,
		UserID:     session.UserID,
		MakePublic: makePublic,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
	}

	draft, err := c.service.GetDraft(r.Context(), app.GetDraftQuery{
		DraftID:   draftID,
		ChapterID: chapterID,
		BookID:    bookID,
		UserID:    session.UserID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponse(draft).Write(w)
}

type apiPayloadCreateChapter struct {
	Body struct {
		Name            string `json:"name"`
		Summary         string `json:"summary"`
		IsAdultOverride bool   `json:"isAdultOverride"`
		Content         string `json:"content"`
	} `in:"body=json"`
}

func (c *apiControllerBM) createChapter(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	input := r.Context().Value(httpin.Input).(*apiPayloadCreateChapter)

	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	command := app.CreateBookChapterCommand{
		BookID:            bookID,
		Name:              input.Body.Name,
		Summary:           input.Body.Summary,
		IsAdultOverride:   input.Body.IsAdultOverride,
		Content:           input.Body.Content,
		UserID:            session.UserID,
		IsPubliclyVisible: false,
	}

	result, err := c.service.CreateBookChapter(r.Context(), command)

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	olhttp.NewAPIResponse(app.Int64String(result.ID)).Write(w)
}

func (c *apiControllerBM) getChapters(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	s := auth.RequireSession(r.Context())

	result, err := c.service.GetBookChapters(r.Context(), app.ManagerGetBookChaptersQuery{
		BookID: bookID,
		UserID: s.UserID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(result.Chapters).Write(w)
}
