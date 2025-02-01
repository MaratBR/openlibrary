package main

import (
	"context"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/gofrs/uuid"
)

type settingsController struct {
	userService    app.UserService
	sessionService app.SessionService
}

func (c *settingsController) GetCustomizationSettings(w http.ResponseWriter, r *http.Request) {
	getUserSettings(w, r, c.userService.GetUserCustomizationSettings)
}

func (c *settingsController) GetPrivacySettings(w http.ResponseWriter, r *http.Request) {
	getUserSettings(w, r, c.userService.GetUserPrivacySettings)
}

func (c *settingsController) GetModerationSettings(w http.ResponseWriter, r *http.Request) {
	getUserSettings(w, r, c.userService.GetUserModerationSettings)
}

func (c *settingsController) GetAboutSettings(w http.ResponseWriter, r *http.Request) {
	getUserSettings(w, r, c.userService.GetUserAboutSettings)
}

func (c *settingsController) UpdateCustomizationSettings(w http.ResponseWriter, r *http.Request) {
	updateUserSettings(w, r, c.userService.UpdateUserCustomizationSettings)
}

func (c *settingsController) UpdatePrivacySettings(w http.ResponseWriter, r *http.Request) {
	updateUserSettings(w, r, c.userService.UpdateUserPrivacySettings)
}

func (c *settingsController) UpdateModerationSettings(w http.ResponseWriter, r *http.Request) {
	updateUserSettings(w, r, c.userService.UpdateUserModerationSettings)
}

func (c *settingsController) UpdateAboutSettings(w http.ResponseWriter, r *http.Request) {
	updateUserSettings(w, r, c.userService.UpdateUserAboutSettings)
}

func updateUserSettings[T any](w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, userID uuid.UUID, settings T) error) {
	session := auth.RequireSession(r.Context())
	settings, err := getJSON[T](r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	err = fn(r.Context(), session.UserID, settings)
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, settings)
}

func getUserSettings[T any](w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, userID uuid.UUID) (*T, error)) {
	session := auth.RequireSession(r.Context())
	settings, err := fn(r.Context(), session.UserID)
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, settings)
}

type sessionsResponse struct {
	Sessions        []app.SessionPublicInfo `json:"sessions"`
	ActiveSessionID int64                   `json:"activeSessionId,string"`
}

func (c *settingsController) GetSessions(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	sessions, err := c.sessionService.GetByUserID(r.Context(), session.UserID)
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, sessionsResponse{
		Sessions:        app.SessionPublicInfoArray(sessions),
		ActiveSessionID: session.ID,
	})
}

func newSettingsController(service app.UserService, sessionService app.SessionService) settingsController {
	return settingsController{
		userService:    service,
		sessionService: sessionService,
	}
}
