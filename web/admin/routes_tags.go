package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"

	"github.com/ggicci/httpin"
)

type tagsController struct {
	service app.TagsService
}

func newTagsController(service app.TagsService) *tagsController {
	return &tagsController{service: service}
}

type tagEditBody struct {
	Adult       string `in:"form=adult"`
	Spoiler     string `in:"form=spoiler"`
	Name        string `in:"form=name,required"`
	Type        string `in:"form=type,required"`
	Description string `in:"form=description"`
	SynonymOf   string `in:"form=synonymOf"`
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
		var synonymOf app.Nullable[int64]
		if integer, err := strconv.ParseInt(body.SynonymOf, 10, 64); err == nil {
			synonymOf = app.Value(integer)
		}
		err := c.service.UpdateTag(r.Context(), app.UpdateTagCommand{
			ID:             id,
			Name:           body.Name,
			Description:    body.Description,
			IsAdult:        body.Adult == "on",
			IsSpoiler:      body.Spoiler == "on",
			SynonymOfTagID: synonymOf,
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

		http.Redirect(w, r, fmt.Sprintf("/admin#/tags/%d", id), http.StatusSeeOther)
		return
	}

}
