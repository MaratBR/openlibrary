package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type GetModerationReportQuery struct {
	ActorUserID uuid.UUID
	ReportID    int64
}

type SearchModerationReportsQuery struct {
	ActorUserID uuid.UUID
	Search      string
	TargetType  string
	Page        uint32
	PageSize    uint32
}

type ModerationReportListEntry struct {
	Report
	ReporterUserName string
}

type ModerationReportActivity struct {
	Time        time.Time
	Actor       string
	Description string
	Kind        string
}

type ModerationReportDetail struct {
	Report
	ReporterUserName string
	Status           string
	Priority         string
	AssignedTo       string
	AssignedTeam     string
	Channel          string
	SLADeadline      time.Time
	Tags             []string
	Activities       []ModerationReportActivity
}

type ModerationReportService interface {
	GetReport(context.Context, GetModerationReportQuery) (ModerationReportDetail, error)
	SearchReports(context.Context, SearchModerationReportsQuery) (ModerationPage[ModerationReportListEntry], error)
}

func (s *moderationReportService) SearchReports(ctx context.Context, query SearchModerationReportsQuery) (ModerationPage[ModerationReportListEntry], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationReportListEntry]{}, err
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	entries, total, err := s.repo.Search(ctx, query.Search, query.TargetType, limit, offset)
	if err != nil {
		return ModerationPage[ModerationReportListEntry]{}, err
	}
	return moderationPage(entries, page, size, total), nil
}

type moderationReportService struct {
	auth ModerationAuthorizer
	repo ReportRepository
}

func (s *moderationReportService) GetReport(ctx context.Context, query GetModerationReportQuery) (ModerationReportDetail, error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationReportDetail{}, err
	}
	report, reporterName, err := s.repo.GetByID(ctx, query.ReportID)
	if err != nil {
		return ModerationReportDetail{}, err
	}

	// TODO: Persist report workflow fields and activities instead of returning placeholder support-ticket data.
	return ModerationReportDetail{
		Report: report, ReporterUserName: reporterName,
		Status: "open", Priority: "high", AssignedTo: "Maya Chen", AssignedTeam: "Trust & Safety",
		Channel: "In-product report", SLADeadline: report.Time.Add(24 * time.Hour),
		Tags: []string{"needs-review", string(report.TargetType), "community-safety"},
		Activities: []ModerationReportActivity{
			{Time: report.Time, Actor: reporterName, Description: "submitted this report", Kind: "created"},
			{Time: report.Time.Add(12 * time.Minute), Actor: "Triage automation", Description: "set priority to High", Kind: "priority"},
			{Time: report.Time.Add(18 * time.Minute), Actor: "Maya Chen", Description: "was assigned to this report", Kind: "assignment"},
		},
	}, nil
}

func NewModerationReportService(auth ModerationAuthorizer, repo ReportRepository) ModerationReportService {
	return &moderationReportService{auth: auth, repo: repo}
}
