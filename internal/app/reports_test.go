package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

type reportRepoStub struct {
	exists   bool
	created  Report
	result   Report
	book     *ModerationReportBookContextData
	user     *ModerationReportUserContextData
	comment  *ModerationReportCommentContextData
	events   []ReportDecision
	decision *ReportDecision
	status   string
	getErr   error
}

type reportActionExecutorStub struct{ err error }

func (s reportActionExecutorStub) Execute(context.Context, Report, DecideModerationReportCommand) error {
	return s.err
}
func (s reportActionExecutorStub) AvailableActions(report Report) []string {
	return NewReportActionExecutor(nil, nil).AvailableActions(report)
}

func (r *reportRepoStub) GetEvents(context.Context, int64) ([]ReportDecision, error) {
	return r.events, nil
}
func (r *reportRepoStub) AddDecision(_ context.Context, _ int64, status string, decision ReportDecision) error {
	r.status, r.decision = status, &decision
	return nil
}

func (r *reportRepoStub) GetByID(context.Context, int64) (Report, string, error) {
	return r.result, "reporter", r.getErr
}

func (r *reportRepoStub) GetBookContext(context.Context, int64) (*ModerationReportBookContextData, error) {
	return r.book, nil
}
func (r *reportRepoStub) GetUserContext(context.Context, int64) (*ModerationReportUserContextData, error) {
	return r.user, nil
}
func (r *reportRepoStub) GetCommentContext(context.Context, int64) (*ModerationReportCommentContextData, error) {
	return r.comment, nil
}

func (r *reportRepoStub) Search(context.Context, string, string, int32, int32) ([]ModerationReportListEntry, int64, error) {
	return nil, 0, nil
}

func TestModerationReportUsesPersistedWorkflow(t *testing.T) {
	createdAt := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	repo := &reportRepoStub{result: Report{ID: 42, Number: "R-2026-0807-42", Time: createdAt, TargetType: ReportTargetComment, TargetID: "99", Reason: "Harassment", Status: "reviewing", Priority: "high"}}
	result, err := NewModerationReportService(moderationAuthStub{}, repo, nil, reportActionExecutorStub{}).GetReport(context.Background(), GetModerationReportQuery{ReportID: 42})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != 42 || result.Status != "reviewing" || result.Priority != "high" || len(result.Activities) != 1 {
		t.Fatalf("expected report core and persisted workflow, got %#v", result)
	}
	if result.BookContext != nil {
		t.Fatalf("comment report should not include book context, got %#v", result.BookContext)
	}
}

func TestModerationReportRequiresModerator(t *testing.T) {
	repo := &reportRepoStub{getErr: errors.New("repository should not be called")}
	_, err := NewModerationReportService(moderationAuthStub{err: ErrModerationForbidden}, repo, nil, reportActionExecutorStub{}).GetReport(context.Background(), GetModerationReportQuery{})
	if !errors.Is(err, ErrModerationForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestModerationReportDecisionValidationAndPersistence(t *testing.T) {
	repo := &reportRepoStub{result: Report{ID: 42, Status: "unreviewed"}}
	svc := NewModerationReportService(moderationAuthStub{}, repo, nil, reportActionExecutorStub{})
	actor := uuid.Must(uuid.NewV4())
	err := svc.DecideReport(context.Background(), DecideModerationReportCommand{ActorUserID: actor, ReportID: 42, Disposition: ReportDispositionNoViolation, PolicyReason: " No policy violation ", InternalNote: " reviewed "})
	if err != nil {
		t.Fatal(err)
	}
	if repo.status != "resolved" || repo.decision == nil || repo.decision.PolicyReason != "No policy violation" || repo.decision.InternalNote != "reviewed" {
		t.Fatalf("decision was not normalized and persisted: status=%q decision=%#v", repo.status, repo.decision)
	}
}

func TestModerationReportDecisionRejectsInvalidTransitions(t *testing.T) {
	tests := []DecideModerationReportCommand{
		{Disposition: "unknown", PolicyReason: "reason"},
		{Disposition: ReportDispositionNoViolation, Action: "comment.remove", PolicyReason: "reason"},
		{Disposition: ReportDispositionActionTaken, PolicyReason: "reason"},
		{Disposition: ReportDispositionEscalated, Action: "book.ban", PolicyReason: "reason"},
		{Disposition: ReportDispositionNoViolation},
	}
	for _, cmd := range tests {
		repo := &reportRepoStub{result: Report{ID: 42, Status: "unreviewed"}}
		cmd.ReportID = 42
		if err := NewModerationReportService(moderationAuthStub{}, repo, nil, reportActionExecutorStub{}).DecideReport(context.Background(), cmd); err == nil {
			t.Fatalf("expected invalid command to fail: %#v", cmd)
		}
		if repo.decision != nil {
			t.Fatalf("invalid command was persisted: %#v", cmd)
		}
	}
	repo := &reportRepoStub{result: Report{ID: 42, Status: "resolved"}}
	err := NewModerationReportService(moderationAuthStub{}, repo, nil, reportActionExecutorStub{}).DecideReport(context.Background(), DecideModerationReportCommand{ReportID: 42, Disposition: ReportDispositionNoViolation, PolicyReason: "reason"})
	if !errors.Is(err, ErrReportAlreadyResolved) {
		t.Fatalf("expected resolved error, got %v", err)
	}
}

func TestModerationReportFailedEnforcementDoesNotResolve(t *testing.T) {
	expected := errors.New("enforcement failed")
	repo := &reportRepoStub{result: Report{ID: 42, Status: "unreviewed", TargetType: ReportTargetComment, TargetID: "9"}}
	err := NewModerationReportService(moderationAuthStub{}, repo, nil, reportActionExecutorStub{err: expected}).DecideReport(context.Background(), DecideModerationReportCommand{ReportID: 42, Disposition: ReportDispositionActionTaken, Action: "comment.remove", PolicyReason: "Harassment"})
	if !errors.Is(err, expected) {
		t.Fatalf("expected enforcement error, got %v", err)
	}
	if repo.decision != nil || repo.status != "" {
		t.Fatalf("failed enforcement changed report: status=%q decision=%#v", repo.status, repo.decision)
	}
}

func TestReportAvailableActionsMatchTargetAndScope(t *testing.T) {
	executor := NewReportActionExecutor(nil, nil)
	chapterActions := executor.AvailableActions(Report{TargetType: ReportTargetBook, BookChapterID: Value[int64](7)})
	wholeBookActions := executor.AvailableActions(Report{TargetType: ReportTargetBook})
	if len(chapterActions) != len(wholeBookActions)+2 || chapterActions[len(chapterActions)-1] != "chapter.restore" {
		t.Fatalf("chapter actions do not include scoped actions: %#v", chapterActions)
	}
	if actions := executor.AvailableActions(Report{TargetType: ReportTargetComment}); len(actions) != 2 || actions[0] != "comment.remove" {
		t.Fatalf("unexpected comment actions: %#v", actions)
	}
}

func TestModerationWholeBookReportHandlesMissingChapterData(t *testing.T) {
	createdAt := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	repo := &reportRepoStub{
		result: Report{ID: 42, Time: createdAt, TargetType: ReportTargetBook, TargetID: "99"},
		book:   &ModerationReportBookContextData{BookID: 99, Title: "Book", BookCreatedAt: createdAt.Add(-time.Hour)},
	}
	result, err := NewModerationReportService(moderationAuthStub{}, repo, nil, reportActionExecutorStub{}).GetReport(context.Background(), GetModerationReportQuery{ReportID: 42})
	if err != nil {
		t.Fatal(err)
	}
	if result.BookContext == nil || result.BookContext.Scope != "book" || !result.BookContext.LastUpdated.Equal(createdAt.Add(-time.Hour)) {
		t.Fatalf("unexpected whole-book context: %#v", result.BookContext)
	}
}

func TestModerationReportLoadsUserAndCommentSnapshots(t *testing.T) {
	user := &ModerationReportUserContextData{Name: "reported user", IsBanned: true}
	userRepo := &reportRepoStub{result: Report{ID: 1, TargetType: ReportTargetUser}, user: user}
	userResult, err := NewModerationReportService(moderationAuthStub{}, userRepo, nil, reportActionExecutorStub{}).GetReport(context.Background(), GetModerationReportQuery{ReportID: 1})
	if err != nil || userResult.UserContext != user || userResult.CommentContext != nil {
		t.Fatalf("unexpected user report snapshot: result=%#v err=%v", userResult, err)
	}
	comment := &ModerationReportCommentContextData{ID: 9, Content: "reported comment"}
	commentRepo := &reportRepoStub{result: Report{ID: 2, TargetType: ReportTargetComment}, comment: comment}
	commentResult, err := NewModerationReportService(moderationAuthStub{}, commentRepo, nil, reportActionExecutorStub{}).GetReport(context.Background(), GetModerationReportQuery{ReportID: 2})
	if err != nil || commentResult.CommentContext != comment || commentResult.UserContext != nil {
		t.Fatalf("unexpected comment report snapshot: result=%#v err=%v", commentResult, err)
	}
}

func (r *reportRepoStub) TargetExists(context.Context, ReportTargetType, string) (bool, error) {
	return r.exists, nil
}

func (r *reportRepoStub) BookChapterExists(context.Context, int64, int64) (bool, error) {
	return r.exists, nil
}

func (r *reportRepoStub) Create(_ context.Context, report Report) (Report, error) {
	report.ID = 123
	report.Number = "R-2026-0712-1"
	r.created = report
	return report, nil
}

func TestReportServiceCreatesServerIDAndNormalizesText(t *testing.T) {
	repo := &reportRepoStub{exists: true}
	svc := NewReportService(repo)
	report, err := svc.Create(context.Background(), CreateReportCommand{
		ReporterUserID: uuid.Must(uuid.NewV4()), TargetType: ReportTargetBook, TargetID: " 42 ",
		Reason: " Spam ", Description: " repeated links ", BookChapterID: Value[int64](7), BookExcerpt: " selected text ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.ID == 0 || report.ID != repo.created.ID || report.Number != "R-2026-0712-1" {
		t.Fatalf("expected database-generated identity and number, got %#v", report)
	}
	if report.TargetID != "42" || report.Reason != "Spam" || report.Description != "repeated links" {
		t.Fatalf("report text was not normalized: %#v", report)
	}
	if !report.BookChapterID.Valid || report.BookChapterID.Value != 7 || report.BookExcerpt != "selected text" {
		t.Fatalf("book report scope was not preserved: %#v", report)
	}
}

func TestReportServiceRejectsUnsupportedOrMissingTargets(t *testing.T) {
	for _, test := range []struct {
		name string
		cmd  CreateReportCommand
		err  error
	}{
		{"unsupported type", CreateReportCommand{TargetType: "chapter", TargetID: "1", Reason: "Spam"}, ErrInvalidReportTarget},
		{"blank reason", CreateReportCommand{TargetType: ReportTargetBook, TargetID: "1"}, ErrReportReason},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewReportService(&reportRepoStub{exists: true}).Create(context.Background(), test.cmd)
			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}
		})
	}
	_, err := NewReportService(&reportRepoStub{}).Create(context.Background(), CreateReportCommand{TargetType: ReportTargetComment, TargetID: "1", Reason: "Spam"})
	if err != ErrReportTargetMissing {
		t.Fatalf("expected missing target error, got %v", err)
	}
}
