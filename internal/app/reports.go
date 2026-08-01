package app

import (
	"context"
	"strconv"
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
)

type Report struct {
	ID             string
	Time           time.Time
	ReporterUserID uuid.UUID
	TargetType     ReportTargetType
	TargetID       string
	Reason         string
	Description    string
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
	Create(context.Context, Report) error
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
		ID: strconv.FormatInt(GenID(), 10), Time: s.now(), ReporterUserID: cmd.ReporterUserID,
		TargetType: cmd.TargetType, TargetID: cmd.TargetID, Reason: cmd.Reason, Description: cmd.Description,
	}
	if err := s.repo.Create(ctx, report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func NewReportService(repo ReportRepository) ReportService {
	return &reportService{repo: repo, now: time.Now}
}
