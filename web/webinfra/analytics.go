package webinfra

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
)

func GetAnalyticsViewMetadata(r *http.Request) analytics.ViewMetadata {
	userID := auth.GetNullableUserID(r.Context())
	ip := olhttp.GetIP(r)

	return analytics.ViewMetadata{
		IP:     ip,
		UserID: userID,
	}
}
