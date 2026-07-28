package webinfra

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
)

func GetAnalyticsViewMetadata(r *http.Request) analytics.EventMetadata {
	userID := auth.GetNullableUserID(r.Context())
	ip := olhttp.GetIP(r)

	return analytics.EventMetadata{
		IP:     ip,
		UserID: userID,
	}
}
