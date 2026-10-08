package public

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/bookfont"
	"github.com/ggicci/httpin"
	"github.com/stretchr/testify/require"
)

func TestAPIBookManagerGetFontPolicy(t *testing.T) {
	controller := &apiControllerBM{fontPolicy: bookfont.Policy{
		MaxPerChapter: 7,
		Whitelist:     []string{"Poppins", "Literata"},
	}}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/_api/books-manager/font-policy", nil)

	controller.getFontPolicy(w, r)

	require.Equal(t, 200, w.Code)
	require.JSONEq(t, `{"maxPerChapter":7,"whitelist":["Poppins","Literata"]}`, w.Body.String())
}

// With no service configured, reaching persistence would panic. Rejected legacy
// requests must return before any book update or authentication lookup.
func TestBookDirectUpdateRejectsLegacyAdultFieldBeforeSaving(t *testing.T) {
	for _, rating := range []string{"G", "R", "NC-17"} {
		for _, adult := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%v", rating, adult), func(t *testing.T) {
				var input apiPayloadBookDirectUpdate
				body := fmt.Sprintf(`{"ageRating":%q,"isAdult":%v,"name":"changed"}`, rating, adult)
				require.NoError(t, json.Unmarshal([]byte(body), &input.Body))
				r := httptest.NewRequest("POST", "/", nil)
				r = r.WithContext(context.WithValue(r.Context(), httpin.Input, &input))
				w := httptest.NewRecorder()
				(&apiControllerBM{}).bookDirectUpdate(w, r)
				require.Equal(t, 400, w.Code)
				require.Contains(t, w.Body.String(), "omit isAdult")
			})
		}
	}
}
