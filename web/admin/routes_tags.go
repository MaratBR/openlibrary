package admin

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/MaratBR/openlibrary/web/olresponse"
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
		olresponse.Write500(w, r, err)
		return
	}

	tag, err := c.service.GetTag(r.Context(), id)
	if err != nil {
		olresponse.Write500(w, r, err)
		return
	}

	templates.Tag(tag).Render(r.Context(), w)
}
