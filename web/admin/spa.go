package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/go-chi/chi/v5"
)

// All SPA endpoints are mounted inside the same authorization group as the UI.
func setupSPA(r chi.Router, users *usersController, tags *tagsController, debug *debugController) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		olhttp.WriteTemplate(w, r.Context(), templates.AdminSPA())
	})
	for _, path := range []string{"/users", "/books", "/tags", "/debug", "/users/{id}", "/tags/tag-details/{id}", "/tags/tag-details/{id}/edit"} {
		r.Get(path, redirectAdminPage)
		if !strings.Contains(path, "{id}") {
			r.Get(path+"/", redirectAdminPage)
		}
	}
	r.Route("/api", func(r chi.Router) {
		r.Get("/users", users.listJSON)
		r.Get("/users/{id}", users.userJSON)
		r.Post("/users/{id}", users.updateJSON)
		r.Get("/tags", tags.listJSON)
		r.Get("/tags/{id}", tags.tagJSON)
		r.Post("/tags/{id}", tags.updateJSON)
		r.Get("/debug", debug.listJSON)
		r.Post("/debug", debug.runJSON)
	})
}

func redirectAdminPage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/admin"), "/")
	path = strings.Replace(path, "/tags/tag-details/", "/tags/", 1)
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, "/admin#"+path, http.StatusFound)
}

func requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUser(r.Context())
		isAPI := strings.HasPrefix(r.URL.Path, "/admin/api/")
		if !ok {
			if isAPI {
				adminJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			} else {
				http.Redirect(w, r, "/admin/login", http.StatusFound)
			}
			return
		}
		// Use the freshly loaded user role rather than a cached session role.
		if !user.Role.IsAdmin() {
			if isAPI {
				adminJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
			} else {
				w.WriteHeader(http.StatusForbidden)
				templates.Forbidden().Render(r.Context(), w)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func adminJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func adminAPIError(w http.ResponseWriter, status int) {
	adminJSON(w, status, map[string]string{"error": http.StatusText(status)})
}

func adminReadError(w http.ResponseWriter, err error) {
	if errors.Is(err, app.ErrUserNotFound) || errors.Is(err, app.ErrTagNotFound) {
		adminAPIError(w, http.StatusNotFound)
	} else {
		adminAPIError(w, http.StatusInternalServerError)
	}
}

type adminUser struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Role     app.UserRole `json:"role"`
	IsBanned bool         `json:"isBanned"`
	Avatar   string       `json:"avatar"`
	JoinedAt string       `json:"joinedAt"`
	Bio      string       `json:"bio"`
	Gender   string       `json:"gender"`
}

func (c *usersController) listJSON(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result, err := c.service.ListUsers(r.Context(), app.UsersQuery{Page: olhttp.GetPage(q, "p"), PageSize: 20, Query: q.Get("q"), Role: getUserRoles(q, "usersFilter.role")})
	if err != nil {
		adminReadError(w, err)
		return
	}
	items := make([]adminUser, 0, len(result.Users))
	for _, u := range result.Users {
		items = append(items, adminUser{ID: u.ID.String(), Name: u.Name, Role: u.Role, IsBanned: u.IsBanned, Avatar: u.Avatar, JoinedAt: u.JoinedAt.Format("2006-01-02")})
	}
	adminJSON(w, http.StatusOK, struct {
		Users      []adminUser `json:"users"`
		Page       uint32      `json:"page"`
		TotalPages int32       `json:"totalPages"`
		Total      int32       `json:"total"`
	}{items, result.Page, result.TotalPages, result.Total})
}

func (c *usersController) userJSON(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	u, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{ID: id, UserID: auth.GetNullableUserID(r.Context())})
	if err != nil {
		adminReadError(w, err)
		return
	}
	adminJSON(w, http.StatusOK, adminUser{ID: u.ID.String(), Name: u.Name, Role: u.Role, IsBanned: u.IsBanned, Avatar: u.Avatar.LG, JoinedAt: u.JoinedAt.Format("2006-01-02"), Bio: u.About.Bio, Gender: u.About.Gender})
}

func readAdminForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	return r.ParseForm()
}

func (c *usersController) updateJSON(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	if err = readAdminForm(w, r); err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	role, err := app.ParseUserRole(r.PostForm.Get("role"))
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	if _, ok := r.PostForm["about"]; !ok {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	if _, ok := r.PostForm["gender"]; !ok {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	// An empty password means no reset; the optional field is never generated on load.
	if _, err = c.service.GetUserDetails(r.Context(), app.GetUserQuery{ID: id, UserID: auth.GetNullableUserID(r.Context())}); err != nil {
		adminReadError(w, err)
		return
	}
	err = c.service.UpdateUser(r.Context(), app.UpdateUserCommand{UserID: id, ActorUserID: auth.GetNullableUserID(r.Context()), Role: app.Value(role), About: r.PostForm.Get("about"), Gender: r.PostForm.Get("gender"), Password: r.PostForm.Get("password")})
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	adminJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type adminTag struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Adult       bool            `json:"adult"`
	Spoiler     bool            `json:"spoiler"`
	CreatedAt   string          `json:"createdAt"`
	IsDefault   bool            `json:"isDefault"`
	SynonymOf   *adminTagParent `json:"synonymOf"`
}
type adminTagParent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func tagToAdmin(tag app.TagDetailsItemDto) adminTag {
	result := adminTag{ID: strconv.FormatInt(tag.ID, 10), Name: tag.Name, Description: tag.Description, Category: tag.Category.String(), Adult: tag.IsAdult, Spoiler: tag.IsSpoiler, CreatedAt: tag.CreatedAt.Format("2006-01-02"), IsDefault: tag.IsDefault}
	if tag.SynonymOf.Valid {
		result.SynonymOf = &adminTagParent{ID: strconv.FormatInt(tag.SynonymOf.Value.ID, 10), Name: tag.SynonymOf.Value.Name}
	}
	return result
}
func (c *tagsController) listJSON(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result, err := c.service.List(r.Context(), app.ListTagsQuery{Page: olhttp.GetPage(q, "p"), PageSize: 50, SearchQuery: q.Get("q"), OnlyAdultTags: olhttp.GetBoolDefault(q, "onlyAdultTags", false), OnlyParentTags: olhttp.GetBoolDefault(q, "onlyParentTags", false)})
	if err != nil {
		adminReadError(w, err)
		return
	}
	items := make([]adminTag, 0, len(result.Tags))
	for _, t := range result.Tags {
		items = append(items, tagToAdmin(t))
	}
	adminJSON(w, http.StatusOK, struct {
		Tags       []adminTag `json:"tags"`
		Page       uint32     `json:"page"`
		TotalPages uint32     `json:"totalPages"`
		Total      uint32     `json:"total"`
	}{items, result.Page, result.TotalPages, result.TotalCount})
}
func positiveTagID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid tag ID")
	}
	return id, nil
}
func (c *tagsController) tagJSON(w http.ResponseWriter, r *http.Request) {
	id, err := positiveTagID(chi.URLParam(r, "id"))
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	tag, err := c.service.GetTag(r.Context(), id)
	if err != nil {
		adminReadError(w, err)
		return
	}
	adminJSON(w, http.StatusOK, tagToAdmin(tag))
}

func parseTagUpdate(r *http.Request, id int64) (app.UpdateTagCommand, error) {
	f := r.PostForm
	cat := f.Get("type")
	validCat := false
	for _, c := range app.TagsCategoryList {
		if c.String() == cat {
			validCat = true
			break
		}
	}
	if !validCat || strings.TrimSpace(f.Get("name")) == "" {
		return app.UpdateTagCommand{}, errors.New("invalid tag fields")
	}
	if _, ok := f["description"]; !ok {
		return app.UpdateTagCommand{}, errors.New("missing description")
	}
	for _, key := range []string{"adult", "spoiler"} {
		if v := f.Get(key); v != "" && v != "on" {
			return app.UpdateTagCommand{}, errors.New("invalid flag")
		}
	}
	parent := app.Null[int64]()
	if f.Get("synonymOf") != "" {
		value, err := positiveTagID(f.Get("synonymOf"))
		if err != nil || value == id {
			return app.UpdateTagCommand{}, errors.New("invalid parent")
		}
		parent = app.Value(value)
	}
	return app.UpdateTagCommand{ID: id, Name: strings.TrimSpace(f.Get("name")), Description: f.Get("description"), Type: app.TagsCategoryFromName(cat), IsAdult: f.Get("adult") == "on", IsSpoiler: f.Get("spoiler") == "on", SynonymOfTagID: parent}, nil
}
func (c *tagsController) updateJSON(w http.ResponseWriter, r *http.Request) {
	id, err := positiveTagID(chi.URLParam(r, "id"))
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	if err = readAdminForm(w, r); err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	command, err := parseTagUpdate(r, id)
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	// Resolve both IDs before updating so missing targets and parents cannot report success.
	if _, err = c.service.GetTag(r.Context(), id); err != nil {
		adminReadError(w, err)
		return
	}
	if command.SynonymOfTagID.Valid {
		if _, err = c.service.GetTag(r.Context(), command.SynonymOfTagID.Value); err != nil {
			adminReadError(w, err)
			return
		}
	}
	command.UserID = auth.RequireUser(r.Context()).ID
	if err = c.service.UpdateTag(r.Context(), command); err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	adminJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (c *debugController) listJSON(w http.ResponseWriter, r *http.Request) {
	actions := make([]string, 0, len(c.actions))
	for id := range c.actions {
		actions = append(actions, id)
	}
	sort.Strings(actions)
	adminJSON(w, http.StatusOK, actions)
}
func (c *debugController) runJSON(w http.ResponseWriter, r *http.Request) {
	if err := readAdminForm(w, r); err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	act, ok := c.actions[r.PostForm.Get("act")]
	if !ok {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	if err := act.Func(); err != "" {
		adminAPIError(w, http.StatusInternalServerError)
		return
	}
	adminJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
