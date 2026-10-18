package notifications

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type statusDB struct {
	store.DBTX
	user               pgtype.UUID
	through            pgtype.Timestamptz
	execErr, statusErr error
}

func (db *statusDB) Exec(_ context.Context, _ string, args ...any) (pgconn.CommandTag, error) {
	db.user, db.through = args[0].(pgtype.UUID), args[1].(pgtype.Timestamptz)
	return pgconn.NewCommandTag("INSERT 0 1"), db.execErr
}
func (db *statusDB) QueryRow(_ context.Context, _ string, args ...any) pgx.Row {
	db.user = args[0].(pgtype.UUID)
	return statusRow{err: db.statusErr}
}

type statusRow struct{ err error }

func (row statusRow) Scan(values ...any) error {
	if row.err != nil {
		return row.err
	}
	*(values[0].(*int64)) = 4
	*(values[1].(*pgtype.Timestamptz)) = pgtype.Timestamptz{}
	*(values[2].(*pgtype.Timestamptz)) = store.PgtypeTimestamptz(time.Unix(1700000000, 0))
	return nil
}

func TestReadStateServiceScopeAndErrors(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewV7()
	db := &statusDB{}
	service := NewStreamReader(db)
	status, err := service.Status(ctx, id)
	require.NoError(t, err)
	require.EqualValues(t, 4, status.UnreadCount)
	require.Nil(t, status.LastReadAt)
	require.Equal(t, store.PgtypeUUID(id), db.user)
	_, err = service.MarkRead(ctx, id, status.ObservedAt)
	require.NoError(t, err)
	require.Equal(t, status.ObservedAt, db.through.Time)
	require.Equal(t, store.PgtypeUUID(id), db.user)
	for _, through := range []time.Time{{}, time.Unix(-1, 0)} {
		_, err := service.MarkRead(ctx, id, through)
		require.ErrorIs(t, err, ErrInvalidReadTime)
	}
	db.execErr = errors.New("write failed")
	_, err = service.MarkRead(ctx, id, status.ObservedAt)
	require.ErrorIs(t, err, db.execErr)
	db.execErr, db.statusErr = nil, errors.New("read failed")
	_, err = service.Status(ctx, id)
	require.ErrorIs(t, err, db.statusErr)
}

// Set NOTIFICATION_TEST_DATABASE_URL to run against PostgreSQL. All fixtures live
// in transaction-local temporary tables; no application data is read or changed.
func TestReadStatePostgres(t *testing.T) {
	connectionURL := os.Getenv("NOTIFICATION_TEST_DATABASE_URL")
	if connectionURL == "" {
		t.Skip("NOTIFICATION_TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connectionURL)
	require.NoError(t, err)
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `
 create temporary table users (id uuid primary key, role text) on commit drop;
 create temporary table notifications (
  id bigint generated always as identity primary key, key text not null,
  user_id uuid not null references users(id), created_at timestamptz not null,
  type text not null, title text not null, content text not null, metadata jsonb,
  unique(user_id, key)
 ) on commit drop;
 create temporary table user_data (
  user_id uuid references users(id), key text, data jsonb, primary key(key, user_id)
 ) on commit drop;`)
	require.NoError(t, err)
	first, second := uuid.NewV7(), uuid.NewV7()
	_, err = tx.Exec(ctx, "insert into users values ($1, 'admin'), ($2, 'user')", first.String(), second.String())
	require.NoError(t, err)
	service := NewStreamReader(tx)
	_, err = service.Send(ctx, first.String(), "", "First", "<b>Body</b>", Options{Category: CategoryModeration, Urgent: true})
	require.NoError(t, err)
	_, err = service.Send(ctx, second.String(), "", "Other account", "Body")
	require.NoError(t, err)
	snapshot, err := service.Status(ctx, first)
	require.NoError(t, err)
	require.EqualValues(t, 1, snapshot.UnreadCount)
	require.Nil(t, snapshot.LastReadAt)
	// Explicit future creation simulates a notification newer than the snapshot.
	_, err = store.New(tx).Notification_Upsert(ctx, store.Notification_UpsertParams{
		UserID: store.PgtypeUUID(first), Key: "newer", Title: "Newer", Content: "body", Type: "reply",
		Metadata: []byte(`{}`), CreatedAt: store.PgtypeTimestamptz(snapshot.ObservedAt.Add(time.Microsecond)),
	})
	require.NoError(t, err)
	marked, err := service.MarkRead(ctx, first, snapshot.ObservedAt)
	require.NoError(t, err)
	require.EqualValues(t, 1, marked.UnreadCount)
	require.WithinDuration(t, snapshot.ObservedAt, *marked.LastReadAt, time.Microsecond)
	// An older acknowledgement cannot regress state.
	marked, err = service.MarkRead(ctx, first, snapshot.ObservedAt.Add(-time.Hour))
	require.NoError(t, err)
	require.EqualValues(t, 1, marked.UnreadCount)
	require.WithinDuration(t, snapshot.ObservedAt, *marked.LastReadAt, time.Microsecond)
	// A future acknowledgement is clamped to database time.
	marked, err = service.MarkRead(ctx, first, snapshot.ObservedAt.Add(24*time.Hour))
	require.NoError(t, err)
	require.False(t, marked.LastReadAt.After(marked.ObservedAt))
	require.EqualValues(t, 0, marked.UnreadCount)
	other, err := service.Status(ctx, second)
	require.NoError(t, err)
	require.Nil(t, other.LastReadAt)
	require.EqualValues(t, 1, other.UnreadCount)
	history, err := service.History(ctx, first, 0)
	require.NoError(t, err)
	require.Len(t, history, 2)
	require.Equal(t, CategoryModeration, history[1].Category)
	require.True(t, history[1].Urgent)
	require.Equal(t, "<b>Body</b>", history[1].Content)
}
