package app

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"uuid"
)

type GetModerationReportQuery struct {
	ActorUserID uuid.UUID
	ReportID    int64
}

type SearchModerationReportsQuery struct {
	ActorUserID                                 uuid.UUID
	Search                                      string
	TargetType                                  string
	TargetID, ChapterID, Status, Priority, Sort string
	DateFrom, DateTo                            Nullable[time.Time]
	Page                                        uint32
	PageSize                                    uint32
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
	Decision    *ReportDecision
}

type ModerationReportDetail struct {
	Report
	ReporterUserName string
	Status           string
	Priority         string
	Activities       []ModerationReportActivity
	BookContext      *ModerationReportBookContext
	UserContext      *ModerationReportUserContextData
	CommentContext   *ModerationReportCommentContextData
	AvailableActions []string
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
	DecideReport(context.Context, DecideModerationReportCommand) error
	AddNote(context.Context, ReportNoteCommand) error
}

const (
	ReportDispositionNoViolation    = "no_violation"
	ReportDispositionRequestChanges = "request_changes"
	ReportDispositionActionTaken    = "action_taken"
	ReportDispositionEscalated      = "escalated"
)

var (
	ErrInvalidReportDisposition = apperror.AppErrors.NewType("invalid_report_disposition").New("invalid report disposition")
	ErrReportAlreadyResolved    = apperror.AppErrors.NewType("report_already_resolved").New("report is already resolved")
	ErrReportPolicyReason       = apperror.AppErrors.NewType("report_policy_reason_required").New("a policy reason is required")
	ErrReportInternalNote       = apperror.AppErrors.NewType("report_internal_note_too_long").New("internal note must not exceed 1000 characters")
)

type DecideModerationReportCommand struct {
	OperationID  string
	ActorUserID  uuid.UUID
	ReportID     int64
	Disposition  string
	Action       string
	PolicyReason string
	InternalNote string
	NotifyTarget bool
	Payload      json.RawMessage
}

func (s *moderationReportService) SearchReports(ctx context.Context, query SearchModerationReportsQuery) (ModerationPage[ModerationReportListEntry], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationReportListEntry]{}, err
	}
	if err := validateReportSearch(&query); err != nil {
		return ModerationPage[ModerationReportListEntry]{}, err
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	entries, total, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return ModerationPage[ModerationReportListEntry]{}, err
	}
	return moderationPage(entries, page, size, total), nil
}

type moderationReportService struct {
	auth    ModerationAuthorizer
	repo    ReportRepository
	uploads *UploadService
	actions ReportActionExecutor
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
		AvailableActions: s.actions.AvailableActions(report),
	}
	events, err := s.repo.GetEvents(ctx, report.ID)
	if err != nil {
		return ModerationReportDetail{}, err
	}
	for i := range events {
		event := events[i]
		detail.Activities = append(detail.Activities, ModerationReportActivity{Time: event.Time, Actor: event.ActorUserName, Description: "recorded a report decision", Kind: event.Kind, Decision: &event})
	}
	if report.TargetType == ReportTargetBook {
		err = s.loadBookReportContext(ctx, report, &detail)
		if err != nil {
			return ModerationReportDetail{}, err
		}
	}
	if report.TargetType == ReportTargetUser {
		detail.UserContext, err = s.repo.GetUserContext(ctx, report.ID)
	}
	if report.TargetType == ReportTargetComment {
		detail.CommentContext, err = s.repo.GetCommentContext(ctx, report.ID)
	}
	if err != nil {
		return ModerationReportDetail{}, err
	}
	if report.Status == "resolved" {
		detail.AvailableActions = []string{}
	} else {
		filtered := []string{}
		for _, action := range detail.AvailableActions {
			keep := true
			switch report.TargetType {
			case ReportTargetUser:
				if detail.UserContext == nil {
					keep = false
				} else if action == "user.unban" {
					keep = detail.UserContext.IsBanned
				}
			case ReportTargetComment:
				if detail.CommentContext == nil {
					keep = false
				} else if action == "comment.remove" {
					keep = !detail.CommentContext.DeletedAt.Valid
				} else if action == "comment.restore" {
					keep = detail.CommentContext.DeletedAt.Valid
				}
			case ReportTargetBook:
				if detail.BookContext == nil || detail.BookContext.PublicationState == "Permanently removed" {
					keep = false
				}
				if detail.BookContext != nil && action == "book.restrict" {
					keep = detail.BookContext.PublicationState != "Banned" && detail.BookContext.PublicationState != "Permanently removed"
				}
				if detail.BookContext != nil && action == "book.restore" {
					keep = detail.BookContext.PublicationState == "Banned"
				}
			}
			if keep {
				filtered = append(filtered, action)
			}
		}
		detail.AvailableActions = filtered
	}
	return detail, nil
}

func (s *moderationReportService) DecideReport(ctx context.Context, cmd DecideModerationReportCommand) error {
	if err := s.auth.AuthorizeModerator(ctx, cmd.ActorUserID); err != nil {
		return err
	}
	return s.repo.WithLockedReport(ctx, cmd.ReportID, func(ctx context.Context) error { return s.decideLocked(ctx, cmd) })
}

func (s *moderationReportService) decideLocked(ctx context.Context, cmd DecideModerationReportCommand) error {
	if len(cmd.OperationID) > 128 {
		return ErrInvalidReportDisposition
	}
	encoded, err := json.Marshal(cmd)
	if err != nil {
		return ErrInvalidReportDisposition
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(encoded))
	if cmd.OperationID != "" {
		actor, previousHash, found, err := s.repo.FindOperation(ctx, cmd.ReportID, cmd.OperationID)
		if err != nil {
			return err
		}
		if found {
			if actor == cmd.ActorUserID && previousHash == hash {
				return nil
			}
			return ErrModerationStateChanged
		}
	}

	report, _, err := s.repo.GetByID(ctx, cmd.ReportID)
	if err != nil {
		return err
	}
	if report.Status == "resolved" {
		return ErrReportAlreadyResolved
	}
	cmd.PolicyReason = strings.TrimSpace(cmd.PolicyReason)
	cmd.InternalNote = strings.TrimSpace(cmd.InternalNote)
	if cmd.PolicyReason == "" {
		return ErrReportPolicyReason
	}
	if len([]rune(cmd.InternalNote)) > 1000 {
		return ErrReportInternalNote
	}
	status := "resolved"
	switch cmd.Disposition {
	case ReportDispositionNoViolation, ReportDispositionRequestChanges:
		if cmd.Action != "" {
			return ErrInvalidReportDisposition
		}
	case ReportDispositionActionTaken:
		if cmd.Action == "" {
			return ErrInvalidReportDisposition
		}
		if err = s.actions.Execute(ctx, report, cmd); err != nil {
			return err
		}
	case ReportDispositionEscalated:
		if cmd.Action != "" {
			return ErrInvalidReportDisposition
		}
		status = "escalated"
	default:
		return ErrInvalidReportDisposition
	}
	return s.repo.AddDecision(ctx, cmd.ReportID, status, ReportDecision{OperationID: cmd.OperationID, RequestHash: hash, Time: time.Now(), ActorUserID: cmd.ActorUserID, Disposition: cmd.Disposition, Action: cmd.Action, PolicyReason: cmd.PolicyReason, InternalNote: cmd.InternalNote, NotifyTarget: cmd.NotifyTarget, Payload: cmd.Payload})
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

func NewModerationReportService(auth ModerationAuthorizer, repo ReportRepository, uploads *UploadService, actions ReportActionExecutor) ModerationReportService {
	return &moderationReportService{auth: auth, repo: repo, uploads: uploads, actions: actions}
}

var ErrModerationFilter = apperror.AppErrors.NewType("invalid_moderation_filter").New("invalid moderation filter or date range")

func validateReportSearch(q *SearchModerationReportsQuery) error {
	q.Search = strings.TrimSpace(q.Search)
	if q.Status == "" {
		q.Status = "pending"
	}
	if q.Sort == "" {
		q.Sort = "priority"
	}
	if !oneOf(q.Status, "pending", "all", "unreviewed", "escalated", "resolved") ||
		!oneOf(q.Sort, "priority", "oldest", "newest") ||
		!oneOf(q.Priority, "", "low", "medium", "high", "critical") ||
		!oneOf(q.TargetType, "", "user", "book", "comment") ||
		(q.DateFrom.Valid && q.DateTo.Valid && !q.DateFrom.Value.Before(q.DateTo.Value)) {
		return ErrModerationFilter
	}
	return nil
}
func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

type ReportNoteCommand struct {
	ActorUserID            uuid.UUID
	ReportID               int64
	Note, Priority, Reason string
}

func (s *moderationReportService) AddNote(ctx context.Context, cmd ReportNoteCommand) error {
	if err := s.auth.AuthorizeModerator(ctx, cmd.ActorUserID); err != nil {
		return err
	}
	cmd.Note, cmd.Reason = strings.TrimSpace(cmd.Note), strings.TrimSpace(cmd.Reason)
	if len([]rune(cmd.Note)) > 1000 {
		return ErrReportInternalNote
	}
	if cmd.Priority == "" && cmd.Note == "" {
		return ErrModerationReason
	}
	if !oneOf(cmd.Priority, "", "low", "medium", "high", "critical") {
		return ErrModerationFilter
	}
	if cmd.Priority != "" && cmd.Reason == "" {
		return ErrReportPolicyReason
	}
	return s.repo.WithLockedReport(ctx, cmd.ReportID, func(ctx context.Context) error {
		report, _, err := s.repo.GetByID(ctx, cmd.ReportID)
		if err != nil {
			return err
		}
		if report.Status == "resolved" {
			return ErrReportAlreadyResolved
		}
		return s.repo.AddNote(ctx, cmd.ReportID, ReportDecision{Time: time.Now(), ActorUserID: cmd.ActorUserID, InternalNote: cmd.Note, PolicyReason: cmd.Reason}, cmd.Priority)
	})
}
