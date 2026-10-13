package public

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// renderBookPreviewWidgets fills validated widget containers before sending the
// page. The returned server fragments also seed saved widgets in the editor.
// The cache belongs to this request so access decisions never cross viewers.
func renderBookPreviewWidgets(r *http.Request, markup string, books app.BookService) (string, map[string]string, error) {
	contextNode := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	nodes, err := html.ParseFragment(strings.NewReader(markup), contextNode)
	if err != nil {
		return "", nil, err
	}
	previews := make(map[string]string)
	var walk func(*html.Node) error
	walk = func(node *html.Node) error {
		var widget, reference string
		for _, attr := range node.Attr {
			if attr.Namespace != "" {
				continue
			}
			switch attr.Key {
			case "data-widget":
				widget = attr.Val
			case "data-book-id":
				reference = attr.Val
			}
		}
		if node.Type == html.ElementNode && widget == "book-preview" {
			bookID, err := strconv.ParseInt(reference, 10, 64)
			if err != nil || bookID <= 0 {
				return nil
			}
			preview, cached := previews[reference]
			if !cached {
				book, lookupErr := books.GetBookDetails(r.Context(), app.GetBookQuery{
					ID: bookID, ActorUserID: auth.GetNullableUserID(r.Context()), AdminOverride: adminBookOverride(r),
				})
				var rendered strings.Builder
				if lookupErr != nil {
					err = templates.BookPreviewUnavailable().Render(r.Context(), &rendered)
				} else {
					err = templates.BookPreviewPartial(book).Render(r.Context(), &rendered)
				}
				if err != nil {
					return err
				}
				preview = rendered.String()
				previews[reference] = preview
			}
			var rendered strings.Builder
			if err := templates.BookPreviewWidgetContent(preview, bookID).Render(r.Context(), &rendered); err != nil {
				return err
			}
			children, err := html.ParseFragment(strings.NewReader(rendered.String()), contextNode)
			if err != nil {
				return err
			}
			for node.FirstChild != nil {
				node.RemoveChild(node.FirstChild)
			}
			for _, child := range children {
				node.AppendChild(child)
			}
			// Server fragments are final; do not recursively expand references
			// appearing inside a book's summary (which may reference itself).
			return nil
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if err := walk(child); err != nil {
				return err
			}
		}
		return nil
	}
	var rendered strings.Builder
	for _, node := range nodes {
		if err := walk(node); err != nil {
			return "", nil, err
		}
		if err := html.Render(&rendered, node); err != nil {
			return "", nil, err
		}
	}
	return rendered.String(), previews, nil
}
