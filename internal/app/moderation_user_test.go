package app

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
)

type moderationUserRepoStub struct {
	result       ModerationUserInfo
	err          error
	called       bool
	reportsTotal int64
	seeded       int
	searchID     *uuid.UUID
}

func (r *moderationUserRepoStub) SearchUsers(_ context.Context, _ string, searchID *uuid.UUID, _, _ string, _, _ int32) ([]ModerationUserListEntry, int64, error) {
	r.searchID = searchID
	return nil, 0, nil
}

func (r *moderationUserRepoStub) GetUserInfo(context.Context, uuid.UUID) (ModerationUserInfo, error) {
	r.called = true
	return r.result, r.err
}
func (r *moderationUserRepoStub) GetBooks(context.Context, uuid.UUID, int32, int32) ([]ModerationUserBook, error) {
	return nil, nil
}
func (r *moderationUserRepoStub) GetComments(context.Context, uuid.UUID, int32, int32) ([]ModerationUserComment, error) {
	return nil, nil
}
func (r *moderationUserRepoStub) GetHistory(context.Context, uuid.UUID, int32, int32) ([]ModerationUserHistoryEntry, int64, error) {
	return nil, 0, nil
}
func (r *moderationUserRepoStub) GetReports(context.Context, uuid.UUID, int32, int32) ([]ModerationUserReport, int64, error) {
	return nil, r.reportsTotal + int64(r.seeded), nil
}
func (r *moderationUserRepoStub) CreateTemporaryRandomReports(_ context.Context, _ uuid.UUID, count int) error {
	r.seeded += count
	return nil
}

func TestModerationUserReportsTemporarilySeedsFourBelowThreshold(t *testing.T) {
	repo := &moderationUserRepoStub{reportsTotal: 49}
	svc := NewModerationUserService(moderationAuthStub{}, repo)
	result, err := svc.GetReports(context.Background(), ModerationUserPageQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if repo.seeded != 4 || result.Total != 53 {
		t.Fatalf("expected four generated reports and total 53, got seeded=%d total=%d", repo.seeded, result.Total)
	}
}

func TestModerationUserServiceRequiresModerator(t *testing.T) {
	repo := &moderationUserRepoStub{}
	svc := NewModerationUserService(moderationAuthStub{err: ErrModerationForbidden}, repo)
	_, err := svc.GetUserInfo(context.Background(), GetModerationUserQuery{})
	if !errors.Is(err, ErrModerationForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
	if repo.called {
		t.Fatal("repository was called before authorization")
	}
}

func TestModerationUserServiceReturnsRepositoryData(t *testing.T) {
	want := ModerationUserInfo{Name: "reader", BooksTotal: 4, CommentsTotal: 9, FollowersTotal: 2}
	repo := &moderationUserRepoStub{result: want}
	svc := NewModerationUserService(moderationAuthStub{}, repo)
	got, err := svc.GetUserInfo(context.Background(), GetModerationUserQuery{})
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

func TestModerationUserSearchRecognizesExactUUID(t *testing.T) {
	repo := &moderationUserRepoStub{}
	svc := NewModerationUserService(moderationAuthStub{}, repo)
	want := uuid.Must(uuid.NewV4())
	if _, err := svc.SearchUsers(context.Background(), ModerationUsersQuery{Search: want.String(), Page: 1, PageSize: 20}); err != nil {
		t.Fatal(err)
	}
	if repo.searchID == nil || *repo.searchID != want {
		t.Fatalf("expected exact UUID search %s, got %v", want, repo.searchID)
	}
}
