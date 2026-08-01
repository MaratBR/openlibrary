package app

import (
	"context"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type reportRepository struct{ db DB }

func (r *reportRepository) TargetExists(ctx context.Context, targetType ReportTargetType, targetID string) (bool, error) {
	q := store.New(r.db)
	var (
		exists bool
		err    error
	)
	switch targetType {
	case ReportTargetUser:
		id, parseErr := uuid.FromString(targetID)
		if parseErr != nil {
			return false, nil
		}
		exists, err = q.Report_UserExists(ctx, uuidDomainToDb(id))
	case ReportTargetBook, ReportTargetComment:
		id, parseErr := strconv.ParseInt(targetID, 10, 64)
		if parseErr != nil {
			return false, nil
		}
		if targetType == ReportTargetBook {
			exists, err = q.Report_BookExists(ctx, id)
		} else {
			exists, err = q.Report_CommentExists(ctx, id)
		}
	default:
		return false, nil
	}
	if err != nil {
		return false, apperror.WrapUnexpectedDBError(err)
	}
	return exists, nil
}

func (r *reportRepository) Create(ctx context.Context, report Report) (Report, error) {
	created, err := store.New(r.db).Report_Create(ctx, store.Report_CreateParams{
		Time: timeToTimestamptz(report.Time), ReporterUserID: uuidDomainToDb(report.ReporterUserID),
		TargetType: string(report.TargetType), TargetID: report.TargetID, Reason: report.Reason, Description: report.Description,
	})
	if err != nil {
		return Report{}, apperror.WrapUnexpectedDBError(err)
	}
	report.ID = created.ID
	report.Number = created.Number
	return report, nil
}

func NewReportRepository(db DB) ReportRepository { return &reportRepository{db: db} }
