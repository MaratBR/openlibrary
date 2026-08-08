package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

type reportRepoStub struct {
	exists  bool
	created Report
	result  Report
	getErr  error
}

func (r *reportRepoStub) GetByID(context.Context, int64) (Report, string, error) {
	return r.result, "reporter", r.getErr
}

func (r *reportRepoStub) Search(context.Context, string, string, int32, int32) ([]ModerationReportListEntry, int64, error) {
	return nil, 0, nil
}

func TestModerationReportUsesRealReportAndPlaceholderWorkflow(t *testing.T) {
	createdAt := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	repo := &reportRepoStub{result: Report{ID: 42, Number: "R-2026-0807-42", Time: createdAt, TargetType: ReportTargetComment, TargetID: "99", Reason: "Harassment"}}
	result, err := NewModerationReportService(moderationAuthStub{}, repo).GetReport(context.Background(), GetModerationReportQuery{ReportID: 42})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != 42 || result.Status == "" || result.AssignedTeam == "" || len(result.Activities) == 0 {
		t.Fatalf("expected report core and placeholder workflow, got %#v", result)
	}
	if result.BookContext != nil {
		t.Fatalf("comment report should not include book context, got %#v", result.BookContext)
	}
	if !result.SLADeadline.Equal(createdAt.Add(24 * time.Hour)) {
		t.Fatalf("unexpected SLA deadline: %s", result.SLADeadline)
	}
}

func TestModerationBookReportIncludesPlaceholderBookScopes(t *testing.T) {
	createdAt := time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		id      int64
		scope   string
		chapter bool
		excerpt bool
	}{{42, "book", false, false}, {43, "chapter", true, false}, {44, "text", true, true}} {
		repo := &reportRepoStub{result: Report{ID: test.id, Time: createdAt, TargetType: ReportTargetBook, TargetID: "99"}}
		result, err := NewModerationReportService(moderationAuthStub{}, repo).GetReport(context.Background(), GetModerationReportQuery{ReportID: test.id})
		if err != nil {
			t.Fatal(err)
		}
		book := result.BookContext
		if book == nil || book.Title == "" || book.Scope != test.scope || (book.Chapter != "") != test.chapter || (book.Excerpt != "") != test.excerpt {
			t.Fatalf("expected %s placeholder scope, got %#v", test.scope, book)
		}
	}
}

func TestModerationReportRequiresModerator(t *testing.T) {
	repo := &reportRepoStub{getErr: errors.New("repository should not be called")}
	_, err := NewModerationReportService(moderationAuthStub{err: ErrModerationForbidden}, repo).GetReport(context.Background(), GetModerationReportQuery{})
	if !errors.Is(err, ErrModerationForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func (r *reportRepoStub) TargetExists(context.Context, ReportTargetType, string) (bool, error) {
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
		Reason: " spam ", Description: " repeated links ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.ID == 0 || report.ID != repo.created.ID || report.Number != "R-2026-0712-1" {
		t.Fatalf("expected database-generated identity and number, got %#v", report)
	}
	if report.TargetID != "42" || report.Reason != "spam" || report.Description != "repeated links" {
		t.Fatalf("report text was not normalized: %#v", report)
	}
}

func TestReportServiceRejectsUnsupportedOrMissingTargets(t *testing.T) {
	for _, test := range []struct {
		name string
		cmd  CreateReportCommand
		err  error
	}{
		{"unsupported type", CreateReportCommand{TargetType: "chapter", TargetID: "1", Reason: "reason"}, ErrInvalidReportTarget},
		{"blank reason", CreateReportCommand{TargetType: ReportTargetBook, TargetID: "1"}, ErrReportReason},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewReportService(&reportRepoStub{exists: true}).Create(context.Background(), test.cmd)
			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}
		})
	}
	_, err := NewReportService(&reportRepoStub{}).Create(context.Background(), CreateReportCommand{TargetType: ReportTargetComment, TargetID: "1", Reason: "reason"})
	if err != ErrReportTargetMissing {
		t.Fatalf("expected missing target error, got %v", err)
	}
}
