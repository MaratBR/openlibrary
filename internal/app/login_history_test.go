package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

type loginHistoryRepoStub struct {
	entries       []LoginHistoryEntry
	total         int64
	locations     []LoginLocation
	query         GetLoginHistoryQuery
	limit, offset int32
}

func (s *loginHistoryRepoStub) Search(_ context.Context, query GetLoginHistoryQuery, limit, offset int32) ([]LoginHistoryEntry, int64, error) {
	s.query, s.limit, s.offset = query, limit, offset
	return s.entries, s.total, nil
}
func (s *loginHistoryRepoStub) RecentLocations(context.Context, uuid.UUID) ([]LoginLocation, error) {
	return s.locations, nil
}

func TestLoginHistorySearchNormalizesAndPaginates(t *testing.T) {
	repo := &loginHistoryRepoStub{entries: []LoginHistoryEntry{{ID: 1}}, total: 25}
	result, err := NewLoginHistoryService(moderationAuthStub{}, repo).GetUserLoginHistory(context.Background(), GetLoginHistoryQuery{Search: "  firefox ", Status: "active", Page: 2, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if repo.query.Search != "firefox" || repo.limit != 10 || repo.offset != 10 || result.Page != 2 || result.TotalPages != 3 || result.Total != 25 {
		t.Fatalf("unexpected query/page: query=%#v result=%#v", repo.query, result)
	}
}

func TestLoginHistoryRejectsInvalidFilters(t *testing.T) {
	svc := NewLoginHistoryService(moderationAuthStub{}, &loginHistoryRepoStub{})
	for _, query := range []GetLoginHistoryQuery{
		{Status: "unknown"},
		{DateFrom: Value(time.Now()), DateTo: Value(time.Now().Add(-time.Hour))},
		{UserIDs: make([]uuid.UUID, 21)},
	} {
		_, err := svc.GetUserLoginHistory(context.Background(), query)
		if !errors.Is(err, ErrInvalidLoginHistoryFilter) {
			t.Fatalf("expected invalid filter error, got %v", err)
		}
	}
}

func TestLoginHistoryPromotesScopedUserToUserFilter(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	repo := &loginHistoryRepoStub{}
	_, err := NewLoginHistoryService(moderationAuthStub{}, repo).GetUserLoginHistory(context.Background(), GetLoginHistoryQuery{UserID: userID})
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.query.UserIDs) != 1 || repo.query.UserIDs[0] != userID {
		t.Fatalf("expected scoped user filter, got %#v", repo.query.UserIDs)
	}
}

func TestRecentLoginLocationsReturnsThreeNewestDistinctLocations(t *testing.T) {
	now := time.Now()
	berlin := IPLocation{Country: "Germany", Region: "Berlin", City: "Berlin"}
	tokyo := IPLocation{Country: "Japan", Region: "Tokyo", City: "Tokyo"}
	sydney := IPLocation{Country: "Australia", Region: "New South Wales", City: "Sydney"}
	toronto := IPLocation{Country: "Canada", Region: "Ontario", City: "Toronto"}
	repo := &loginHistoryRepoStub{locations: []LoginLocation{
		{IPLocation: berlin, LastSeenAt: now},
		{IPLocation: berlin, LastSeenAt: now.Add(-time.Hour)},
		{IPLocation: tokyo, LastSeenAt: now.Add(-2 * time.Hour)},
		{IPLocation: IPLocation{}, LastSeenAt: now.Add(-3 * time.Hour)},
		{IPLocation: sydney, LastSeenAt: now.Add(-4 * time.Hour)},
		{IPLocation: toronto, LastSeenAt: now.Add(-5 * time.Hour)},
	}}
	locations, err := NewLoginHistoryService(moderationAuthStub{}, repo).GetRecentLoginLocations(context.Background(), GetLoginHistoryQuery{})
	if err != nil {
		t.Fatal(err)
	}
	if len(locations) != 3 || locations[0].IPLocation != berlin || locations[1].IPLocation != tokyo || locations[2].IPLocation != sydney {
		t.Fatalf("unexpected recent locations: %#v", locations)
	}
}
