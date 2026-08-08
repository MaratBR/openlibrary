package app

import (
	"context"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

type ReportTargetType string

const (
	ReportTargetUser    ReportTargetType = "user"
	ReportTargetBook    ReportTargetType = "book"
	ReportTargetComment ReportTargetType = "comment"
)

var (
	ErrInvalidReportTarget = apperror.AppErrors.NewType("invalid_report_target").New("report target must be a user, book, or comment")
	ErrReportReason        = apperror.AppErrors.NewType("report_reason_required").New("a report reason is required")
	ErrReportTargetMissing = apperror.AppErrors.NewType("report_target_not_found", apperror.ErrTraitEntityNotFound).New("report target not found")
	ErrReportNotFound      = apperror.AppErrors.NewType("report_not_found", apperror.ErrTraitEntityNotFound).New("report not found")
)

type Report struct {
	ID             int64
	Number         string
	Time           time.Time
	ReporterUserID uuid.UUID
	TargetType     ReportTargetType
	TargetID       string
	Reason         string
	Description    string
	Status         string
	Priority       string
}

type ModerationReportBookContextData struct {
	Title, Author, CoverID, Excerpt, Rating   string
	BookID                                    int64
	ChapterID                                 Nullable[int64]
	Chapter                                   Nullable[string]
	Warnings                                  []string
	IsPermanentlyRemoved, IsBanned, IsTrashed bool
	IsPubliclyVisible                         bool
	BookCreatedAt                             time.Time
	ChapterCreatedAt, ChapterUpdatedAt        Nullable[time.Time]
	ChapterContentUpdatedAt                   Nullable[time.Time]
	RelatedReports                            int
}

type CreateReportCommand struct {
	ReporterUserID uuid.UUID
	TargetType     ReportTargetType
	TargetID       string
	Reason         string
	Description    string
}

type ReportRepository interface {
	TargetExists(context.Context, ReportTargetType, string) (bool, error)
	Create(context.Context, Report) (Report, error)
	GetByID(context.Context, int64) (Report, string, error)
	GetBookContext(context.Context, int64) (*ModerationReportBookContextData, error)
	Search(context.Context, string, string, int32, int32) ([]ModerationReportListEntry, int64, error)
}

type ReportService interface {
	Create(context.Context, CreateReportCommand) (Report, error)
}

type reportService struct {
	repo ReportRepository
	now  func() time.Time
}

func (s *reportService) Create(ctx context.Context, cmd CreateReportCommand) (Report, error) {
	cmd.TargetID = strings.TrimSpace(cmd.TargetID)
	cmd.Reason = strings.TrimSpace(cmd.Reason)
	cmd.Description = strings.TrimSpace(cmd.Description)
	if cmd.TargetType != ReportTargetUser && cmd.TargetType != ReportTargetBook && cmd.TargetType != ReportTargetComment {
		return Report{}, ErrInvalidReportTarget
	}
	if cmd.TargetID == "" {
		return Report{}, ErrInvalidReportTarget
	}
	if cmd.Reason == "" {
		return Report{}, ErrReportReason
	}
	exists, err := s.repo.TargetExists(ctx, cmd.TargetType, cmd.TargetID)
	if err != nil {
		return Report{}, err
	}
	if !exists {
		return Report{}, ErrReportTargetMissing
	}
	report := Report{
		Time: s.now(), ReporterUserID: cmd.ReporterUserID,
		TargetType: cmd.TargetType, TargetID: cmd.TargetID, Reason: cmd.Reason, Description: cmd.Description,
	}
	report, err = s.repo.Create(ctx, report)
	if err != nil {
		return Report{}, err
	}
	return report, nil
}

func NewReportService(repo ReportRepository) ReportService {
	return &reportService{repo: repo, now: time.Now}
}
