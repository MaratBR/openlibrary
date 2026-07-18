package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

const defaultSettingsType = "about"

type SettingsController struct {
	userService app.UserService
}

func NewSettingsController(userService app.UserService) *SettingsController {
	return &SettingsController{userService: userService}
}

func (c *SettingsController) Register(r chi.Router) {
	r.Get("/settings", c.redirectToDefault)
	r.Get("/settings/", c.redirectToDefault)
	r.Get("/settings/{settingsType}", c.page)
	r.Patch("/settings/{settingsType}", c.update)
}

func (c *SettingsController) redirectToDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/account/settings/"+defaultSettingsType, http.StatusFound)
}

func (c *SettingsController) page(w http.ResponseWriter, r *http.Request) {
	settingsType := chi.URLParam(r, "settingsType")
	userID := auth.RequireSession(r.Context()).UserID

	var (
		settings any
		err      error
	)

	switch settingsType {
	case "about":
		settings, err = c.userService.GetUserAboutSettings(r.Context(), userID)
	case "privacy":
		settings, err = c.userService.GetUserPrivacySettings(r.Context(), userID)
	case "moderation":
		settings, err = c.userService.GetUserModerationSettings(r.Context(), userID)
	case "customization":
		settings, err = c.userService.GetUserCustomizationSettings(r.Context(), userID)
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		olhttp.Write500(w, r, err)
		return
	}

	if err := templates.AccountSettings(settingsType, settings).Render(r.Context(), w); err != nil {
		olhttp.Write500(w, r, err)
	}
}

func (c *SettingsController) update(w http.ResponseWriter, r *http.Request) {
	settingsType := chi.URLParam(r, "settingsType")
	userID := auth.RequireSession(r.Context()).UserID
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var err error
	switch settingsType {
	case "about":
		var settings app.UserAboutSettings
		if err := decoder.Decode(&settings); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		err = c.userService.UpdateUserAboutSettings(r.Context(), userID, settings)
	case "privacy":
		var settings app.UserPrivacySettings
		if err := decoder.Decode(&settings); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		err = c.userService.UpdateUserPrivacySettings(r.Context(), userID, settings)
	case "moderation":
		var settings app.UserModerationSettings
		if err := decoder.Decode(&settings); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if !isValidCensorMode(settings.CensoredTagsMode) {
			writeError(w, http.StatusBadRequest, errors.New("invalid censored tags mode"))
			return
		}
		err = c.userService.UpdateUserModerationSettings(r.Context(), userID, settings)
	case "customization":
		var settings app.UserCustomizationSetting
		if err := decoder.Decode(&settings); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		err = c.userService.UpdateUserCustomizationSettings(r.Context(), userID, settings)
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	olhttp.NewAPIResponseOK().Write(w)
}

func isValidCensorMode(mode app.CensorMode) bool {
	switch mode {
	case "none", "hide", "censor":
		return true
	default:
		return false
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	olhttp.NewAPIError(err).Write(w)
}
