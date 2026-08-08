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
	auth    ModerationAuthorizer
	repo    ReportRepository
	uploads *UploadService
}

func (s *moderationReportService) GetReport(ctx context.Context, query GetModerationReportQuery) (ModerationReportDetail, error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationReportDetail{}, err
	}
	report, reporterName, err := s.repo.GetByID(ctx, query.ReportID)
	if err != nil {
		return ModerationReportDetail{}, err
	}

	detail := ModerationReportDetail{
		Report: report, ReporterUserName: reporterName,
		Status: report.Status, Priority: report.Priority,
		Activities: []ModerationReportActivity{
			{Time: report.Time, Actor: reporterName, Description: "submitted this report", Kind: "created"},
		},
	}
	if report.TargetType == ReportTargetBook {
		err = s.loadBookReportContext(ctx, report, &detail)
		if err != nil {
			return ModerationReportDetail{}, err
		}
	}
	return detail, nil
}

func (s *moderationReportService) loadBookReportContext(ctx context.Context, report Report, detail *ModerationReportDetail) error {
	book, err := s.repo.GetBookContext(ctx, report.ID)
	if err != nil {
		return err
	}
	if book != nil {
		scope := "book"
		if book.ChapterID.Valid {
			scope = "chapter"
			if book.Excerpt != "" {
				scope = "text"
			}
		}
		publicationState := "Unpublished"
		switch {
		case book.IsPermanentlyRemoved:
			publicationState = "Permanently removed"
		case book.IsBanned:
			publicationState = "Banned"
		case book.IsTrashed:
			publicationState = "Trashed"
		case book.IsPubliclyVisible:
			publicationState = "Published"
		}
		lastUpdated := book.BookCreatedAt
		for _, candidate := range []Nullable[time.Time]{book.ChapterCreatedAt, book.ChapterContentUpdatedAt, book.ChapterUpdatedAt} {
			if candidate.Valid && candidate.Value.After(lastUpdated) {
				lastUpdated = candidate.Value
			}
		}
		editedAfter := lastUpdated.Sub(report.Time)
		if editedAfter < 0 {
			editedAfter = 0
		}
		detail.BookContext = &ModerationReportBookContext{
			Scope:            scope,
			Title:            book.Title,
			Author:           book.Author,
			CoverURL:         getBookCover(s.uploads, book.CoverID, book.BookID).URL,
			Chapter:          book.Chapter.Value,
			Excerpt:          book.Excerpt,
			Rating:           book.Rating,
			Warnings:         book.Warnings,
			PublicationState: publicationState,
			LastUpdated:      lastUpdated,
			RelatedReports:   book.RelatedReports,
			EditedAfter:      editedAfter,
		}
	}
	return nil
}

func NewModerationReportService(auth ModerationAuthorizer, repo ReportRepository, uploads *UploadService) ModerationReportService {
	return &moderationReportService{auth: auth, repo: repo, uploads: uploads}
}
