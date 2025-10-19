package admin

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"

	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type tagsController struct {
	db      app.DB
	cfg     *koanf.Koanf
	service app.TagsService
}

func newTagsController(db app.DB, cfg *koanf.Koanf, service app.TagsService) *tagsController {
	return &tagsController{db: db, cfg: cfg, service: service}
}

func (c *tagsController) Setup(r chi.Router) {
	r.Get("/", c.Home)
	r.Get("/tag-details/{id}", c.Tag)
	r.Get("/tag-details/{id}/edit", c.TagEdit)
	r.With(httpin.NewInput(&tagEditBody{})).Post("/tag-details/{id}/edit", c.TagEdit)

}

func (c *tagsController) Home(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := olhttp.GetPage(query, "p")
	searchQuery := query.Get("q")
	onlyAdultTags := olhttp.GetBoolDefault(query, "onlyAdultTags", false)
	onlyParentTags := olhttp.GetBoolDefault(query, "onlyParentTags", false)

	tags, err := c.service.List(r.Context(), app.ListTagsQuery{
		Page:           page,
		PageSize:       50,
		SearchQuery:    searchQuery,
		OnlyParentTags: onlyParentTags,
		OnlyAdultTags:  onlyAdultTags,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	templates.TagsList(tags, templates.TagsSearchRequest{
		SearchQuery:    searchQuery,
		OnlyParentTags: onlyParentTags,
		OnlyAdultTags:  onlyAdultTags,
	}).Render(r.Context(), w)
}

func (c *tagsController) Tag(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamInt64(r, "id")
	if err != nil {
		olhttp.Write500(w, r, err)
		return
	}

	tag, err := c.service.GetTag(r.Context(), id)
	if err != nil {
		olhttp.Write500(w, r, err)
		return
	}

	templates.Tag(tag).Render(r.Context(), w)
}

type tagEditBody struct {
	Adult       string `in:"form=adult,required"`
	Spoiler     string `in:"form=spoiler,required"`
	Name        string `in:"form=name,required"`
	Type        string `in:"form=type,required"`
	Description string `in:"form=description,required"`
	SynonymOf   *int64 `in:"form=synonymOf"`
}

func (c *tagsController) TagEdit(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamInt64(r, "id")
	if err != nil {
		olhttp.Write500(w, r, err)
		return
	}

	session := auth.RequireSession(r.Context())

	if r.Method == http.MethodPost {
		body := r.Context().Value(httpin.Input).(*tagEditBody)

		err := c.service.UpdateTag(r.Context(), app.UpdateTagCommand{
			ID:             id,
			Name:           body.Name,
			Description:    body.Description,
			IsAdult:        body.Adult == "on",
			IsSpoiler:      body.Spoiler == "on",
			SynonymOfTagID: app.NullableFromPtr(body.SynonymOf),
			UserID:         session.UserID,
			Type:           app.TagsCategoryFromName(body.Type),
		})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}

		l := i18n.GetLocalizer(r.Context())
		flash.Add(r, flash.Text(
			l.T("admin.tags.updatedSuccessfully"),
		))

		http.Redirect(w, r, fmt.Sprintf("/admin/tags/tag-details/%d", id), http.StatusFound)
		return
	}

	tag, err := c.service.GetTag(r.Context(), id)
	if err != nil {
		olhttp.Write500(w, r, err)
		return
	}

	templates.TagEdit(tag).Render(r.Context(), w)
}
