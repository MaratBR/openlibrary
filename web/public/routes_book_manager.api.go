package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	i18nProvider "github.com/MaratBR/openlibrary/internal/i18n-provider"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/upload"
	"github.com/MaratBR/openlibrary/web/olresponse"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type apiBookManagerController struct {
	service       app.BookManagerService
	fileValidator upload.FileValidator
}

func newApiBookManagerController(service app.BookManagerService, fileValidator upload.FileValidator) *apiBookManagerController {
	return &apiBookManagerController{service: service, fileValidator: fileValidator}
}

func (c *apiBookManagerController) Register(r chi.Router) {
	r.Route("/_api/books-manager", func(r chi.Router) {
		r.With(httpin.NewInput(&updateBookRequest{})).Post("/book/{bookID}", c.updateBook)
		r.With(httpin.NewInput(&uploadCoverInput{})).Post("/book/{bookID}/cover", c.uploadCover)

	})
}

type updateBookRequest struct {
	Payload struct {
		Name              string            `json:"name"`
		Summary           string            `json:"summary"`
		Tags              []app.Int64String `json:"tags"`
		Rating            string            `json:"rating"`
		IsPubliclyVisible bool              `json:"isPubliclyVisible"`
	} `in:"body=json"`
}

func (c *apiBookManagerController) updateBook(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	input := r.Context().Value(httpin.Input).(*updateBookRequest)

	rating := app.AsRating(input.Payload.Rating)
	tags := app.ArrInt64StringToInt64(input.Payload.Tags)
	name := input.Payload.Name
	summary := input.Payload.Summary

	err = c.service.UpdateBook(r.Context(), app.UpdateBookCommand{
		BookID:            bookID,
		UserID:            session.UserID,
		Tags:              tags,
		Name:              name,
		Summary:           summary,
		AgeRating:         rating,
		IsPubliclyVisible: input.Payload.IsPubliclyVisible,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	l := i18nProvider.GetLocalizer(r.Context())

	bookResult, err := c.service.GetBook(r.Context(), app.ManagerGetBookQuery{
		ActorUserID: session.UserID,
		BookID:      bookID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	response := olresponse.NewJSONResponse(bookResult.Book)
	response.AddNotification(olresponse.NewNotification(
		l.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "bookManager.edit.editedSuccessfully",
			TemplateData: map[string]string{
				"Name": name,
			},
		}),
		olresponse.NotificationInfo,
	))
	response.Write(w)
}

type uploadCoverInput struct {
	ClientCropped bool         `in:"form=clientCropped"`
	File          *httpin.File `in:"form=file"`
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

	response := olresponse.NewJSONResponse(result.URL)
	response.Write(w)
}
