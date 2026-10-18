package notifications

import (
	"context"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type recordingDB struct {
	dal.DB
	args  []any
	query string
}

func (db *recordingDB) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	db.query, db.args = query, args
	return pgconn.NewCommandTag("INSERT 0 2"), nil
}
func (db *recordingDB) QueryRow(_ context.Context, query string, args ...any) pgx.Row {
	db.query, db.args = query, args
	return notificationRow{}
}

type notificationRow struct{}

func (notificationRow) Scan(values ...any) error { *(values[0].(*int64)) = 1; return nil }

func TestBuilderPreservesRecipientAndValidatesTitle(t *testing.T) {
	id := uuid.NewV7()
	n, err := NewBuilder("Title", id).Content(strings.Repeat("x", 101)).Data("test", "value").Build()
	require.NoError(t, err)
	require.Equal(t, id, n.UserID())
	require.Equal(t, "value", n.Data["test"])
	_, err = NewBuilder(strings.Repeat("界", 101), id).Build()
	require.ErrorIs(t, err, ErrTitleTooLong)
}

func TestSinkPersistsRecipientAndMetadata(t *testing.T) {
	db := &recordingDB{}
	n, err := NewBuilder("Title", uuid.NewV7()).Content("Body").Data("test", "value").Build()
	require.NoError(t, err)
	require.NoError(t, NewSimpleSink(NewSQLStore(db)).Push(context.Background(), n))
	require.Equal(t, store.PgtypeUUID(n.UserID()), db.args[1])
	require.JSONEq(t, `{"test":"value"}`, string(db.args[6].([]byte)))
	restored, err := NewSQLStore(db).(*sqlStore).newNotificationFromRow(store.Notification{
		Key: n.Key(), Type: n.Type(), Title: n.Title, Content: n.Content, UserID: store.PgtypeUUID(n.UserID()), Metadata: db.args[6].([]byte), CreatedAt: pgtype.Timestamptz{Time: time.Unix(1234, 0), Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, time.Unix(1234, 0), restored.CreatedAt)
}

func TestSendRecipientRules(t *testing.T) {
	id := uuid.NewV7().String()
	for _, tc := range []struct{ user, role, title, content string }{
		{"", "", "Title", ""}, {id, "admin", "Title", ""}, {"invalid", "", "Title", ""}, {"", "unknown", "Title", ""}, {id, "", "", ""}, {id, "", strings.Repeat("x", 101), ""}, {id, "", "Title", strings.Repeat("x", 1001)},
	} {
		_, err := NewStreamReader(nil).Send(context.Background(), tc.user, tc.role, tc.title, tc.content)
		require.Error(t, err)
	}
	db := &recordingDB{}
	count, err := NewStreamReader(db).Send(context.Background(), "", "moderator", "Title", "Body")
	require.NoError(t, err)
	require.EqualValues(t, 2, count)
	require.False(t, db.args[3].(pgtype.UUID).Valid)
	require.Equal(t, "moderator", db.args[4].(pgtype.Text).String)
	// PostgreSQL cannot compare its user_role enum directly with a text parameter.
	require.Contains(t, db.query, "users.role::text = $5::text")
	count, err = NewStreamReader(db).Send(context.Background(), id, "", "Title", "Body")
	require.NoError(t, err)
	require.EqualValues(t, 2, count)
	require.True(t, db.args[3].(pgtype.UUID).Valid)
	require.False(t, db.args[4].(pgtype.Text).Valid)
}

type emptyNotificationRows struct{ pgx.Rows }

func (emptyNotificationRows) Next() bool { return false }
func (emptyNotificationRows) Close()     {}
func (emptyNotificationRows) Err() error { return nil }
func (db *recordingDB) Query(_ context.Context, query string, args ...any) (pgx.Rows, error) {
	db.query, db.args = query, args
	return emptyNotificationRows{}, nil
}
func TestReadersScopeQueriesToRecipientAndCursor(t *testing.T) {
	db := &recordingDB{}
	id := uuid.NewV7()
	reader := NewStreamReader(db)
	events, err := reader.Read(context.Background(), id, 123)
	require.NoError(t, err)
	require.Empty(t, events)
	require.Contains(t, db.query, "where user_id = $1 and id > $2")
	require.Equal(t, store.PgtypeUUID(id), db.args[0])
	require.EqualValues(t, 123, db.args[1])
	_, err = reader.History(context.Background(), id, 321)
	require.NoError(t, err)
	require.Contains(t, db.query, "where user_id = $1")
	require.Equal(t, store.PgtypeUUID(id), db.args[0])
	require.EqualValues(t, 321, db.args[1])
}
