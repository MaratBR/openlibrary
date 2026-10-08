package app

import (
	"context"
	"testing"

	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/zap"
)

func TestAdultContentClassification(t *testing.T) {
	service := NewAdultContentService()
	for _, rating := range append(append([]AgeRating{}, AllRatings...), AgeRating("invalid")) {
		t.Run(string(rating), func(t *testing.T) {
			if got := service.IsAdult(rating, nil); got != (rating == AgeRatingNC17) {
				t.Fatalf("rating %q: adult=%v", rating, got)
			}
			warnings := []DefinedTagDto{{IsAdult: false}, {IsAdult: false}}
			if got := service.IsAdult(rating, warnings); got != (rating == AgeRatingNC17) {
				t.Fatalf("non-adult warnings changed classification: %v", got)
			}
			warnings = append(warnings, DefinedTagDto{IsAdult: true})
			if !service.IsAdult(rating, warnings) {
				t.Fatal("adult tag did not raise classification")
			}
		})
	}
}

type customAdultContentService struct{ adult bool }

func (s customAdultContentService) IsAdult(AgeRating, []DefinedTagDto) bool { return s.adult }

func TestBookDetailsUseCustomAdultContentService(t *testing.T) {
	db := &overrideReadDB{book: store.Book_GetRow{IsPubliclyVisible: true}}
	service := &bookService{
		queries: store.New(db), tagsService: overrideTags{},
		adultContentService: customAdultContentService{adult: true}, log: zap.NewNop().Sugar(),
	}
	book, err := service.GetBookDetails(context.Background(), GetBookQuery{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !book.IsAdult || !book.GetAdultWarning().ShouldShowWarning() {
		t.Fatal("book details or warning ignored the custom policy")
	}
}

func TestManagerDetailsUseCustomAdultContentService(t *testing.T) {
	db := &overrideReadDB{book: store.Book_GetRow{IsPubliclyVisible: true}}
	service := &bookManagerService{
		queries: store.New(db),
		deps: BookManagerServiceDeps{
			TagsService: overrideTags{}, AdultContentService: customAdultContentService{adult: true},
		},
	}
	result, err := service.GetBook(context.Background(), ManagerGetBookQuery{BookID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Book.IsAdult {
		t.Fatal("manager details ignored the custom policy")
	}
}

func TestAdultWarningDoesNotOverrideServiceClassification(t *testing.T) {
	book := BookDetailsDto{IsAdult: false, Tags: []DefinedTagDto{{IsAdult: true}}}
	if book.GetAdultWarning().ShouldShowWarning() {
		t.Fatal("warning reapplied hardcoded tag rules to the service result")
	}
}
