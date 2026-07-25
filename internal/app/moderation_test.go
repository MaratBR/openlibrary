package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

type moderationAuthStub struct{ err error }

func (s moderationAuthStub) AuthorizeModerator(context.Context, uuid.UUID) error { return s.err }

type moderationRepoStub struct {
	called                   string
	visible, removed, banned bool
	until                    time.Time
	log                      ModerationAuditEntry
}

func (r *moderationRepoStub) SetChapterVisibilityAndLog(_ context.Context, _ int64, visible bool, log ModerationAuditEntry) error {
	r.called, r.visible, r.log = "chapter", visible, log
	return nil
}
func (r *moderationRepoStub) SetCommentRemovedAndLog(_ context.Context, _ int64, removed bool, log ModerationAuditEntry) error {
	r.called, r.removed, r.log = "comment", removed, log
	return nil
}
func (r *moderationRepoStub) SetUserBanAndLog(_ context.Context, _ uuid.UUID, until time.Time, banned bool, log ModerationAuditEntry) error {
	r.called, r.until, r.banned, r.log = "ban", until, banned, log
	return nil
}
func (r *moderationRepoStub) RenameUserAndLog(_ context.Context, _ uuid.UUID, _ string, log ModerationAuditEntry) error {
	r.called, r.log = "rename", log
	return nil
}
func (r *moderationRepoStub) ChangeUserAboutAndLog(_ context.Context, _ uuid.UUID, _ string, log ModerationAuditEntry) error {
	r.called, r.log = "about", log
	return nil
}

func TestModerationRejectsUnauthorizedActorBeforeMutation(t *testing.T) {
	repo := &moderationRepoStub{}
	svc := NewContentModerationService(moderationAuthStub{err: ErrModerationForbidden}, repo)
	err := svc.RemoveComment(context.Background(), ModerateCommentCommand{Reason: "spam"})
	if !errors.Is(err, ErrModerationForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	if repo.called != "" {
		t.Fatalf("repository was called: %s", repo.called)
	}
}

func TestPermanentBanUsesYear3000AndLogs(t *testing.T) {
	repo := &moderationRepoStub{}
	svc := NewContentModerationService(moderationAuthStub{}, repo)
	err := svc.PermanentlyBanUser(context.Background(), BanUserCommand{Reason: "abuse"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.called != "ban" || !repo.banned {
		t.Fatalf("ban mutation not made: %#v", repo)
	}
	if repo.until.Year() != 3000 {
		t.Fatalf("expiry year = %d", repo.until.Year())
	}
	if repo.log.Action != ModerationActionBanUser || repo.log.Reason != "abuse" {
		t.Fatalf("invalid log: %#v", repo.log)
	}
}

func TestModerationRequiresReason(t *testing.T) {
	repo := &moderationRepoStub{}
	svc := NewContentModerationService(moderationAuthStub{}, repo)
	err := svc.HideChapter(context.Background(), ModerateChapterCommand{Reason: " \n "})
	if !errors.Is(err, ErrModerationReason) {
		t.Fatalf("expected reason error, got %v", err)
	}
	if repo.called != "" {
		t.Fatalf("repository was called: %s", repo.called)
	}
}
