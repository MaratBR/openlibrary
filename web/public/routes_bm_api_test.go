package public

import (
	"net/http/httptest"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/bookfont"
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
