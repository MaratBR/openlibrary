package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type ModerationAuditLogEntry struct {
	ID                                                  int64
	Time                                                time.Time
	Action, TargetType, TargetID, Reason, ActorUserName string
	Payload                                             json.RawMessage
	ActorUserID                                         uuid.UUID
}

type ModerationAuditLogQuery struct {
	ActorUserID    uuid.UUID
	TargetType     string
	Page, PageSize uint32
}
type ModerationAuditLogService interface {
	GetAuditLog(context.Context, ModerationAuditLogQuery) (ModerationPage[ModerationAuditLogEntry], error)
}

type moderationAuditLogService struct {
	auth    ModerationAuthorizer
	queries *store.Queries
}

func (s *moderationAuditLogService) GetAuditLog(ctx context.Context, query ModerationAuditLogQuery) (ModerationPage[ModerationAuditLogEntry], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationAuditLogEntry]{}, err
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	total, err := s.queries.Moderation_CountAuditLog(ctx, query.TargetType)
	if err != nil {
		return ModerationPage[ModerationAuditLogEntry]{}, apperror.WrapUnexpectedDBError(err)
	}
	rows, err := s.queries.Moderation_GetAuditLog(ctx, store.Moderation_GetAuditLogParams{TargetType: query.TargetType, PageLimit: limit, PageOffset: offset})
	if err != nil {
		return ModerationPage[ModerationAuditLogEntry]{}, apperror.WrapUnexpectedDBError(err)
	}
	entries := MapSlice(rows, func(row store.Moderation_GetAuditLogRow) ModerationAuditLogEntry {
		return ModerationAuditLogEntry{ID: row.ID, Time: timeDbToDomain(row.Time), Action: row.Type, TargetType: row.TargetType, TargetID: row.TargetID, Reason: row.Reason, Payload: row.Payload, ActorUserID: uuidDbToDomain(row.ActorUserID), ActorUserName: row.ActorUserName.String}
	})
	return moderationPage(entries, page, size, total), nil
}
func NewModerationAuditLogService(db DB, auth ModerationAuthorizer) ModerationAuditLogService {
	return &moderationAuditLogService{auth: auth, queries: store.New(db)}
}
