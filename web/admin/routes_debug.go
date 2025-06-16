package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/web/admin/templates"
)

type debugActionDescriptor struct {
	Func func() string
	Name string
}

type debugController struct {
	fullReindexService *app.BookFullReindexService
	actions            map[string]debugActionDescriptor
}

func newDebugController(fullReindexService *app.BookFullReindexService) *debugController {
	c := &debugController{
		fullReindexService: fullReindexService,
		actions: map[string]debugActionDescriptor{
			"books:elastic:reindex": {
				Name: "Reindex all books",
				Func: func() string {
					err := fullReindexService.ScheduleReindexAll()
					if err == nil {
						return ""
					}

					return err.Error()
				},
			},
		},
	}
	return c
}

func (c *debugController) Actions(w http.ResponseWriter, r *http.Request) {
	actions := make([]templates.DebugAction, 0, len(c.actions))
	for actID, act := range c.actions {
		actions = append(actions, templates.DebugAction{
			Name: act.Name,
			ID:   actID,
		})
	}
	if r.Method == http.MethodGet {
		templates.DebugActions(actions).Render(r.Context(), w)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			writeBadRequest(w, r, err)
			return
		}

		act := r.FormValue("act")
		if act == "" {
			writeBadRequest(w, r, errors.New("act form field is empty"))
			return
		}

		if actData, ok := c.actions[act]; ok {
			errString := actData.Func()
			if errString == "" {
				flash.Add(r, flash.Text(fmt.Sprintf("Action %s was executed, it might take some time for some actions to finish", actData.Name)))
			} else {
				flash.Add(r, flash.Text(fmt.Sprintf("Action %s failed with an error: %s", actData.Name, errString)))
			}
			templates.DebugActions(actions).Render(r.Context(), w)
		} else {
			writeApplicationError(w, r, errors.New("unknown action"))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
