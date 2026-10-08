package olhttp

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func TestPaginationPages(t *testing.T) {
	for _, tc := range []struct {
		page, total uint32
		want        []uint32
	}{
		{1, 500, []uint32{1, 2, 3, 4, 5, 0, 500}},
		{5, 500, []uint32{1, 0, 4, 5, 6, 0, 500}},
		{499, 500, []uint32{1, 0, 496, 497, 498, 499, 500}},
		{500, 500, []uint32{1, 0, 496, 497, 498, 499, 500}},
		{2, 3, []uint32{1, 2, 3}},
		{0, 3, []uint32{1, 2, 3}},
		{99, 3, []uint32{1, 2, 3}},
	} {
		if got := paginationPages(tc.page, tc.total, 9); !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("page %d/%d: got %v want %v", tc.page, tc.total, got, tc.want)
		}
	}
}
func TestPaginationPreservesFiltersAndMarksCurrentPage(t *testing.T) {
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../translations/en.toml"}}, zap.NewNop().Sugar())
	var out bytes.Buffer
	handler := ReqCtxMiddleware(provider.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := PaginationImpl(PaginationImplProps{Page: 5, TotalPages: 500, Size: 9, UseAlpineJsAjax: true}).Render(r.Context(), &out); err != nil {
			t.Fatal(err)
		}
	})))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/search?q=book&admin.link=1&__fragment=1&p=5", nil))
	html := out.String()
	for _, value := range []string{`aria-current="page"`, `aria-label="Previous page"`, `aria-label="Next page"`, `p=4`, `p=6`, `admin.link=1`, `q=book`, `x-target.push="PaginationAjaxSlot"`, `Pagination-ellipsis`} {
		if !strings.Contains(html, value) {
			t.Fatalf("missing %s in %s", value, html)
		}
	}
	if strings.Contains(html, "__fragment") || strings.Contains(html, `role="listbox"`) || strings.Contains(html, `href="#"`) {
		t.Fatal("pagination kept invalid links or semantics")
	}
}
