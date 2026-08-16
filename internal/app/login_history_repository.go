package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type loginHistoryRepository struct{ db DB }

func loginHistoryTime(value Nullable[time.Time]) pgtype.Timestamptz {
	if !value.Valid {
		return pgtype.Timestamptz{}
	}
	return timeToTimestamptz(value.Value)
}

func (r *loginHistoryRepository) Search(ctx context.Context, query GetLoginHistoryQuery, limit, offset int32) ([]LoginHistoryEntry, int64, error) {
	q := store.New(r.db)
	params := store.Moderation_SearchUserLoginHistoryParams{UserIds: arrUuidDomainToDb(query.UserIDs), Search: query.Search, DateFrom: loginHistoryTime(query.DateFrom), DateTo: loginHistoryTime(query.DateTo), SessionStatus: query.Status, PageLimit: limit, PageOffset: offset}
	rows, err := q.Moderation_SearchUserLoginHistory(ctx, params)
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	total, err := q.Moderation_CountUserLoginHistory(ctx, store.Moderation_CountUserLoginHistoryParams{UserIds: params.UserIds, Search: params.Search, DateFrom: params.DateFrom, DateTo: params.DateTo, SessionStatus: params.SessionStatus})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	entries := MapSlice(rows, func(row store.Moderation_SearchUserLoginHistoryRow) LoginHistoryEntry {
		return LoginHistoryEntry{ID: row.ID, UserID: uuidDbToDomain(row.UserID), UserName: row.UserName, IPAddress: row.IpAddress, UserAgent: row.UserAgent, Location: IPLocation{Country: row.LocationCountry, Region: row.LocationRegion, City: row.LocationCity}, LoggedInAt: timeDbToDomain(row.CreatedAt), ExpiresAt: timeDbToDomain(row.ExpiresAt), IsTerminated: row.IsTerminated}
	})
	return entries, total, nil
}

func (r *loginHistoryRepository) RecentLocations(ctx context.Context, userID uuid.UUID) ([]LoginLocation, error) {
	rows, err := store.New(r.db).Moderation_GetRecentLoginSessions(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return MapSlice(rows, func(row store.Moderation_GetRecentLoginSessionsRow) LoginLocation {
		return LoginLocation{IPLocation: IPLocation{Country: row.LocationCountry, Region: row.LocationRegion, City: row.LocationCity}, LastSeenAt: timeDbToDomain(row.CreatedAt)}
	}), nil
}

func NewLoginHistoryRepository(db DB) LoginHistoryRepository { return &loginHistoryRepository{db: db} }
