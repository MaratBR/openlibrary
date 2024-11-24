package server

import (
	"fmt"
	"net/http"
	"time"

	uiserver "github.com/MaratBR/openlibrary/cmd/server/ui-server"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/knadh/koanf/v2"
)

type serverSessionData struct {
	ExpiresAt time.Time `json:"expiresAt"`
}

type devStaticHandler struct {
	devFrontEndServerProxy http.Handler
	userService            app.UserService
}

func NewDevStaticHandler(
	config *koanf.Koanf,
	enableDevProxy bool,
	userService app.UserService,
) http.Handler {
	var c devStaticHandler

	c.userService = userService

	if enableDevProxy {
		address := fmt.Sprintf("%s://%s:%d", config.String("frontend-proxy.target-protocol"), config.String("frontend-proxy.target-host"), config.Int("frontend-proxy.target-port"))
		c.devFrontEndServerProxy = uiserver.NewDevServerProxy(address, uiserver.DevServerOptions{
			GetInjectedHTMLSegment: c.getInjectedHTMLSegment,
		})
	}

	return c
}

func (c *devStaticHandler) getInjectedHTMLSegment(r *http.Request) []byte {

	data := &serverData{
		ServerMetadata: map[string]any{
			"session":       nil,
			"user":          nil,
			"_preload":      map[string]any{},
			"clientPreload": true,
			"serverPreload": true,
		},
	}

	session, ok := getSession(r)

	if ok {
		user, err := c.userService.GetUserSelfData(r.Context(), session.UserID)
		if err == nil {
			data.ServerMetadata["session"] = serverSessionData{
				ExpiresAt: session.ExpiresAt,
			}
			data.ServerMetadata["user"] = user
		}
	}

	if r.URL.Path == "/user/__profile" {
		data.ServerMetadata["iframeAllowed"] = true
	}

	return getInjectedHTMLSegment(r, data)
}

func (c devStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.devFrontEndServerProxy == nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("this endpoint is disabled in the current environment"))
		return
	}

	c.devFrontEndServerProxy.ServeHTTP(w, r)
}
