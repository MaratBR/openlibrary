package public

import (
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

type apiBookManagerController struct {
	service       app.BookManagerService
	fileValidator upload.FileValidator
}

func newAPIBookManagerController(service app.BookManagerService, fileValidator upload.FileValidator) *apiBookManagerController {
	return &apiBookManagerController{service: service, fileValidator: fileValidator}
}

func (c *apiBookManagerController) Register(r chi.Router) {
	r.Route("/books-manager", func(r chi.Router) {
		r.With(httpin.NewInput(&uploadCoverInput{})).Post("/book/{bookID}/cover", c.uploadCover)
		r.With(httpin.NewInput(&updateBookChaptersOrderRequest{})).Post("/book/{bookID}/chapters-order", c.updateBookChaptersOrder)
		r.Post("/book/{bookID}/{chapterID}/{draftID}", c.updateDraftContent)
		r.Post("/book/{bookID}/{chapterID}/{draftID}/publish", c.updateDraftContentAndPublish)
		r.Post("/book/{bookID}/{chapterID}/{draftID}/chapterName", c.updateDraftChapterName)
		r.With(httpin.NewInput(&createChapterRequest{})).Post("/book/{bookID}/createChapter", c.createChapter)
	})
}

type uploadCoverInput struct {
	File *httpin.File `in:"form=file"`
}

func (c *apiBookManagerController) uploadCover(w http.ResponseWriter, r *http.Request) {

	session := auth.RequireSession(r.Context())
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	input := r.Context().Value(httpin.Input).(*uploadCoverInput)

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

type updateBookChaptersOrderRequest struct {
	Chapters []app.Int64String `in:"body=json"`
}

func (c *apiBookManagerController) updateBookChaptersOrder(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*updateBookChaptersOrderRequest)
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	err = c.service.UpdateBookChaptersOrder(r.Context(), app.UpdateBookChaptersOrders{
		BookID:     bookID,
		ChapterIDs: app.ArrInt64StringToInt64(input.Chapters),
	})

	if err != nil {
		apiWriteApplicationError(w, err)
	} else {
		apiWriteOK(w)
	}
}

func (c *apiBookManagerController) updateDraftContent(w http.ResponseWriter, r *http.Request) {
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

	apiWriteOK(w)
}

func (c *apiBookManagerController) updateDraftChapterName(w http.ResponseWriter, r *http.Request) {
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

func (c *apiBookManagerController) updateDraftContentAndPublish(w http.ResponseWriter, r *http.Request) {
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

	apiWriteOK(w)
}

type createChapterRequest struct {
	Body struct {
		Name            string `json:"name"`
		Summary         string `json:"summary"`
		IsAdultOverride bool   `json:"isAdultOverride"`
		Content         string `json:"content"`
	} `in:"body=json"`
}

func (c *apiBookManagerController) createChapter(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	input := r.Context().Value(httpin.Input).(*createChapterRequest)

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
