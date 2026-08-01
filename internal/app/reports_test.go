package app

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
)

type reportRepoStub struct {
	exists  bool
	created Report
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
