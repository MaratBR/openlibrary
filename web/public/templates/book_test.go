package templates

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/web/frontend"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"uuid"
)

func TestBookExternalLinksOnlyUsesProvidedSources(t *testing.T) {
	for _, links := range [][]app.BookExternalLinkDto{nil, {{Label: "Author website", URL: "https://example.org/book"}}} {
		var out bytes.Buffer
		if err := bookExternalLinks(links, "Sources").Render(context.Background(), &out); err != nil {
			t.Fatal(err)
		}
		html := out.String()
		if len(links) == 0 && strings.TrimSpace(html) != "" {
			t.Fatalf("absent sources rendered a section: %s", html)
		}
		if strings.Contains(html, "Patreon") {
			t.Fatal("illustrative source rendered")
		}
		if len(links) > 0 && (!strings.Contains(html, links[0].URL) || !strings.Contains(html, links[0].Label)) {
			t.Fatal("provided source missing")
		}
	}
}

func TestBookPagePanelsAndOptionalContent(t *testing.T) {
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../../translations/en.toml"}}, zap.NewNop().Sugar())
	book := app.BookDetailsDto{ID: 1, Name: "A book", AgeRating: "PG", CreatedAt: time.Now(), IsPubliclyAvailable: true, Author: app.BookDetailsAuthorDto{Name: "An author"}}
	for _, scenario := range []string{"empty", "long", "populated", "authenticated", "warning"} {
		t.Run(scenario, func(t *testing.T) {
			current := book
			var topReviews []app.ReviewDto
			var ratingAndReview app.RatingAndReview
			readingList := app.Nullable[app.BookReadingListDto]{}
			authenticated := scenario == "authenticated" || scenario == "warning"
			if authenticated {
				current.ID = 6707299548257058898
				readingList = app.Value(app.BookReadingListDto{Status: app.ReadingListStatusReading, ChapterID: app.Value(app.Int64String(6707299651336274143)), ChapterName: "Saved chapter"})
				current.IsPubliclyAvailable = false
				current.FirstChapterID = app.Value[int64](1)
				ratingAndReview.Review = app.Value(app.ReviewDto{User: app.ReviewUserDto{Name: "An author", Avatar: "/_/embed-assets/logo.svg"}, Rating: 8, Content: "<p>My review</p>"})
			}
			if scenario == "populated" {
				current.Cover.URL = "/_/embed-assets/logo.svg"
				current.Chapters = 1
				current.Words = 100
				current.Reviews = 1
				topReviews = []app.ReviewDto{{User: app.ReviewUserDto{Name: strings.Repeat("Reviewer", 30), Avatar: "/_/embed-assets/logo.svg"}, Rating: 8, Content: "<p>" + strings.Repeat("Review content ", 100) + "</p>"}}
			}
			if scenario == "long" {
				current.Name = strings.Repeat("LongTitle", 20)
				current.Author.Name = strings.Repeat("LongAuthor", 20)
				current.Summary = "<p>" + strings.Repeat("A longer summary. ", 100) + "</p>"
				current.ExternalLinks = []app.BookExternalLinkDto{{Label: "Author website", URL: "https://example.org/book"}}
				for i := 0; i < 20; i++ {
					current.Tags = append(current.Tags, app.DefinedTagDto{ID: int64(i + 1), Name: strings.Repeat("LongTag", 8)})
				}
			}
			router := chi.NewRouter()
			router.Use(olhttp.ReqCtxMiddleware, provider.Middleware)
			if authenticated {
				router.Use(session.Middleware(bookFixtureStore{}, zap.NewNop().Sugar()), auth.NewAuthorizationMiddleware(bookFixtureSessions{}, bookFixtureUsers{}, auth.MiddlewareOptions{}, zap.NewNop().Sugar()))
			}
			frontend.AttachAssetsInliningHandler(os.DirFS("../../frontend/embed-assets"), "embed-assets", router)
			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				if err := BookPage(current, analytics.MetricValue{}, ratingAndReview, readingList, topReviews, false, scenario == "warning").Render(r.Context(), w); err != nil {
					t.Error(err)
				}
			})
			out := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)
			if authenticated {
				request.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
			}
			router.ServeHTTP(out, request)
			html := out.Body.String()
			for _, expected := range []string{`role="tab"`, `aria-controls="slot-book-toc"`, `aria-controls="book-panel-reviews"`, `aria-labelledby="book-tab-toc"`, `aria-labelledby="book-tab-reviews"`} {
				if !strings.Contains(html, expected) {
					t.Errorf("missing %s", expected)
				}
			}
			if scenario != "populated" && !strings.Contains(html, `BookPage-coverFallback`) {
				t.Fatal("missing cover fallback")
			}
			if scenario == "empty" && (strings.Contains(html, `class="BookPage-summarySection"`) || strings.Contains(html, `class="BookPage-externalLinks"`)) {
				t.Fatal("absent content produced an empty section")
			}
			if dir := os.Getenv("BOOK_PAGE_FIXTURE_DIR"); dir != "" {
				if err := os.WriteFile(filepath.Join(dir, scenario+".html"), out.Body.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	var populated bytes.Buffer
	if err := BookTOC(context.Background(), 1, []app.ChapterListDto{{ID: 1, Name: strings.Repeat("LongChapter", 30), Words: 100, CreatedAt: time.Now()}}, 1).Render(localizedContext(provider), &populated); err != nil {
		t.Fatal(err)
	}
	if dir := os.Getenv("BOOK_PAGE_FIXTURE_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "toc-populated.html"), populated.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
	}
	var out bytes.Buffer
	if err := BookTOC(context.Background(), 1, nil, 0).Render(localizedContext(provider), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `aria-labelledby="book-tab-toc"`) || !strings.Contains(out.String(), `role="tabpanel"`) {
		t.Fatal("chapter fragment lost panel semantics")
	}
	if dir := os.Getenv("BOOK_PAGE_FIXTURE_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "toc.html"), out.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func localizedContext(provider *i18n.LocaleProvider) context.Context {
	var ctx context.Context
	provider.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { ctx = r.Context() })).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	return ctx
}

// Exercise the real authorization middleware without persistence or a live account.
type bookFixtureStore struct{}

func (bookFixtureStore) Get(context.Context, string) (session.Session, error) {
	return bookFixtureSession{}, nil
}

type bookFixtureSession struct{}

func (bookFixtureSession) Get(string) (string, bool)  { return "", false }
func (bookFixtureSession) Put(string, string)         {}
func (bookFixtureSession) Save(context.Context) error { return nil }
func (bookFixtureSession) ID() string                 { return "fixture" }

type bookFixtureSessions struct{ app.SessionService }

func (bookFixtureSessions) GetBySID(context.Context, string) (*app.SessionInfo, error) {
	return &app.SessionInfo{}, nil
}

type bookFixtureUsers struct{ app.UserService }

func (bookFixtureUsers) GetUserSelfData(context.Context, uuid.UUID) (*app.SelfUserDto, error) {
	user := &app.SelfUserDto{Name: "An author"}
	user.Avatar.MD = "/_/embed-assets/logo.svg"
	return user, nil
}
