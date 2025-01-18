package admin

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/knadh/koanf/v2"
)

type tagsController struct {
	db  app.DB
	cfg *koanf.Koanf
}

func newTagsController(db app.DB, cfg *koanf.Koanf) *tagsController {
	return &tagsController{db: db, cfg: cfg}
}

func (*tagsController) Home(w http.ResponseWriter, r *http.Request) {
	templates.TagsHome(r.Context()).Render(r.Context(), w)
}
