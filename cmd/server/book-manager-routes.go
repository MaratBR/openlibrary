package main

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
)

type bookManagerController struct {
	service app.BookManagerService
}

func newBookManagerController(service app.BookManagerService) *bookManagerController {
	return &bookManagerController{service: service}
}

type createBookRequest struct {
	Name              string            `json:"name"`
	AgeRating         app.AgeRating     `json:"ageRating"`
	Tags              []app.Int64String `json:"tags"`
	Summary           string            `json:"summary"`
	IsPubliclyVisible bool              `json:"isPubliclyVisible"`
}

type createBookResponse struct {
	ID int64 `json:"id,string"`
}

func (c *bookManagerController) CreateBook(w http.ResponseWriter, r *http.Request) {
	session, ok := auth.GetSession(r.Context())
	if !ok {
		writeUnauthorizedError(w)
		return
	}

	req, err := getJSON[createBookRequest](r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	bookID, err := c.service.CreateBook(r.Context(), app.CreateBookCommand{
		UserID:            session.UserID,
		Name:              req.Name,
		Tags:              unwrapInt64StringArr(req.Tags),
		AgeRating:         req.AgeRating,
		IsPubliclyVisible: req.IsPubliclyVisible,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, createBookResponse{ID: bookID})
}

type updateBookRequest struct {
	Name              string            `json:"name"`
	AgeRating         app.AgeRating     `json:"ageRating"`
	Tags              []app.Int64String `json:"tags"`
	Summary           string            `json:"summary"`
	IsPubliclyVisible bool              `json:"isPubliclyVisible"`
}

type updateBookResponse struct {
	Book app.ManagerBookDetailsDto `json:"book"`
}

func (c *bookManagerController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	session, ok := auth.GetSession(r.Context())
	if !ok {
		writeUnauthorizedError(w)
		return
	}

	req, err := getJSON[updateBookRequest](r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	err = c.service.UpdateBook(r.Context(), app.UpdateBookCommand{
		BookID:            bookID,
		UserID:            session.UserID,
		Name:              req.Name,
		Tags:              unwrapInt64StringArr(req.Tags),
		AgeRating:         req.AgeRating,
		IsPubliclyVisible: req.IsPubliclyVisible,
		Summary:           req.Summary,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	book, err := c.service.GetBook(r.Context(), app.ManagerGetBookQuery{BookID: bookID, ActorUserID: session.UserID})
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, updateBookResponse{Book: book.Book})
}

type getBookResponse struct {
	app.ManagerBookDetailsDto
}

func (c *bookManagerController) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	session, ok := auth.GetSession(r.Context())
	if !ok {
		writeUnauthorizedError(w)
		return
	}

	book, err := c.service.GetBook(r.Context(), app.ManagerGetBookQuery{BookID: id, ActorUserID: session.UserID})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	response := getBookResponse{
		book.Book,
	}
	writeJSON(w, response)
}

type createChapterRequest struct {
	Content         string `json:"content"`
	Name            string `json:"name"`
	IsAdultOverride bool   `json:"isAdultOverride"`
	Summary         string `json:"summary"`
}

type createChapterResponse struct {
	ID int64 `json:"id,string"`
}

func (c *bookManagerController) CreateChapter(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	body, err := getJSON[createChapterRequest](r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	if len(body.Content) == 0 || len(body.Content) > 50000 {
		writeUnprocessableEntity(w, "chapter content must be between 1 and 50000 characters")
		return
	}
	if len(body.Name) > 500 {
		writeUnprocessableEntity(w, "chapter name must not be over 50 characters")
		return
	}

	chapter, err := c.service.CreateBookChapter(r.Context(), app.CreateBookChapterCommand{
		BookID:          bookID,
		Name:            body.Name,
		Content:         body.Content,
		IsAdultOverride: body.IsAdultOverride,
		Summary:         body.Summary,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, createChapterResponse{ID: chapter.ID})
}

type updateChapterRequest struct {
	Content         string `json:"content"`
	Name            string `json:"name"`
	IsAdultOverride bool   `json:"isAdultOverride"`
	Summary         string `json:"summary"`
}

func (c *bookManagerController) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	_, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	chapterID, err := urlParamInt64(r, "chapterID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	body, err := getJSON[updateChapterRequest](r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	if len(body.Content) == 0 || len(body.Content) > 50000 {
		writeUnprocessableEntity(w, "chapter content must be between 1 and 50000 characters")
		return
	}
	if len(body.Name) > 500 {
		writeUnprocessableEntity(w, "chapter name must not be over 50 characters")
		return
	}

	err = c.service.UpdateBookChapter(r.Context(), app.UpdateBookChapterCommand{
		ID:              chapterID,
		Name:            body.Name,
		Content:         body.Content,
		IsAdultOverride: body.IsAdultOverride,
		Summary:         body.Summary,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeOK(w)
}

type reorderChapterRequest struct {
	Sequence []app.Int64String `json:"sequence"`
}

func (c *bookManagerController) UpdateChaptersOrder(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	sessionInfo := auth.RequireSession(r.Context())

	body, err := getJSON[reorderChapterRequest](r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	err = c.service.ReorderChapters(r.Context(), app.ReorderChaptersCommand{
		UserID:     sessionInfo.UserID,
		BookID:     bookID,
		ChapterIDs: unwrapInt64StringArr(body.Sequence),
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeOK(w)
}

func (c *bookManagerController) GetChapters(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	sessionInfo := auth.RequireSession(r.Context())

	chapters, err := c.service.GetBookChapters(r.Context(), app.ManagerGetBookChaptersQuery{
		UserID: sessionInfo.UserID,
		BookID: bookID})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, chapters.Chapters)
}

func (c *bookManagerController) GetChapter(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	chapterID, err := urlParamInt64(r, "chapterID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	sessionInfo := auth.RequireSession(r.Context())

	result, err := c.service.GetChapter(r.Context(), app.ManagerGetChapterQuery{
		UserID:    sessionInfo.UserID,
		BookID:    bookID,
		ChapterID: chapterID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, result.Chapter)
}

func (c *bookManagerController) GetMyBooks(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	books, err := c.service.GetUserBooks(r.Context(), app.GetUserBooksQuery{
		UserID: session.UserID,
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, myBooksResponse{
		Books: books.Books,
	})
}

type uploadBookCoverResponse struct {
	URL string `json:"url"`
}

func (c *bookManagerController) UploadBookCover(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	session := auth.RequireSession(r.Context())

	file, _, err := r.FormFile("file")
	if err != nil {
		writeBadRequest(err, w)
		return
	}
	defer file.Close()

	result, err := c.service.UploadBookCover(r.Context(), app.UploadBookCoverCommand{
		UserID: session.UserID,
		BookID: bookID,
		File:   file,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, uploadBookCoverResponse{
		URL: result.URL,
	})
}
