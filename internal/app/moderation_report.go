package app

import (
	"context"
	"fmt"
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
	BookContext      *ModerationReportBookContext
}

type ModerationReportBookContext struct {
	Scope            string
	Title            string
	Author           string
	CoverURL         string
	Chapter          string
	Excerpt          string
	Rating           string
	Warnings         []string
	PublicationState string
	LastUpdated      time.Time
	RelatedReports   int
	EditedAfter      time.Duration
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
	detail := ModerationReportDetail{
		Report: report, ReporterUserName: reporterName,
		Status: "unreviewed", Priority: "medium", AssignedTo: "", AssignedTeam: "Trust & Safety",
		Channel: "In-product report", SLADeadline: report.Time.Add(24 * time.Hour),
		Tags: []string{"needs-review", string(report.TargetType), "community-safety"},
		Activities: []ModerationReportActivity{
			{Time: report.Time, Actor: reporterName, Description: "submitted this report", Kind: "created"},
			{Time: report.Time.Add(12 * time.Minute), Actor: "Triage automation", Description: "set severity to Medium", Kind: "priority"},
			{Time: report.Time.Add(18 * time.Minute), Actor: "Triage automation", Description: "added this report to the review queue", Kind: "assignment"},
		},
	}
	if report.TargetType == ReportTargetBook {
		// TODO: Replace this book snapshot with persisted report-time content metadata.
		detail.BookContext = &ModerationReportBookContext{
			Scope: "book",
			Title: "Ashes of the Northern Crown", Author: "Mira Vale",
			CoverURL: fmt.Sprintf("/_/embed-assets/cover/%d.h300.webp", (report.ID%5)+1),
			Chapter:  "Chapter 18 · The Siege",
			Excerpt:  "He brought the axe down in a single, sickening arc. The guard's head split open, blood and bone spraying across the stones. Mira stumbled back, her hands over her mouth. Another soldier fell, crushed beneath the battering ram.",
			Rating:   "Teen", Warnings: []string{}, PublicationState: "Published",
			LastUpdated: report.Time.Add(-2 * time.Hour), RelatedReports: 3,
			EditedAfter: 24 * time.Minute,
		}
		switch report.ID % 3 {
		case 1:
			detail.BookContext.Scope = "chapter"
			detail.BookContext.Excerpt = ""
		case 2:
			detail.BookContext.Scope = "text"
		default:
			detail.BookContext.Chapter = ""
			detail.BookContext.Excerpt = ""
		}
	}
	return detail, nil
}

func NewModerationReportService(auth ModerationAuthorizer, repo ReportRepository) ModerationReportService {
	return &moderationReportService{auth: auth, repo: repo}
}
