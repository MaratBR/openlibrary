package analytics

import "context"

type analyticsViewsDummyService struct{}

// CommitPendingViewsToDB implements [ViewsService].
func (a *analyticsViewsDummyService) CommitPendingViewsToDB(ctx context.Context) {
}

// GetBookViews implements [ViewsService].
func (a *analyticsViewsDummyService) GetBookViews(ctx context.Context, bookID int64) (Views, error) {
	return Views{}, nil
}

// GetBooksViews implements [ViewsService].
func (a *analyticsViewsDummyService) GetBooksViews(ctx context.Context, bookIDs []int64) (map[int64]Views, error) {
	m := make(map[int64]Views)
	for _, id := range bookIDs {
		m[id] = Views{}
	}
	return m, nil
}

// GetMostViewedBooks implements [ViewsService].
func (a *analyticsViewsDummyService) GetMostViewedBooks(ctx context.Context, period AnalyticsPeriod) ([]BookViewEntry, error) {
	return []BookViewEntry{}, nil
}

// IncrBookView implements [ViewsService].
func (a *analyticsViewsDummyService) IncrBookView(ctx context.Context, bookID int64, meta ViewMetadata) error {
	return nil
}

// IncrChapterView implements [ViewsService].
func (a *analyticsViewsDummyService) IncrChapterView(ctx context.Context, bookID int64, chapterID int64, meta ViewMetadata) error {
	return nil
}

func NewAnalyticsDummyViewsService() ViewsService {
	return &analyticsViewsDummyService{}
}
