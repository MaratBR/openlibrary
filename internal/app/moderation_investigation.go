package app

import (
	"context"
	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"strings"
	"time"
	"uuid"
)

// go2tsdef:generate
type ModerationContentEntry struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   Nullable[time.Time] `json:"updatedAt"`
	Visible     bool                `json:"visible"`
	BookID      string              `json:"bookId"`
	BookName    string              `json:"bookName"`
	AuthorID    string              `json:"authorId"`
	AuthorName  string              `json:"authorName"`
	ChapterID   string              `json:"chapterId"`
	ChapterName string              `json:"chapterName"`
	ParentID    string              `json:"parentId"`
	Content     string              `json:"content"`
	Replies     []string            `json:"replies"`
}

// go2tsdef:generate
type ModerationContentPage struct {
	Entries    []ModerationContentEntry `json:"entries"`
	Page       uint32                   `json:"page"`
	PageSize   uint32                   `json:"pageSize"`
	TotalPages uint32                   `json:"totalPages"`
	Total      int64                    `json:"total"`
}

// go2tsdef:generate
type ModerationOverview struct {
	Unreviewed             int64               `json:"unreviewed"`
	Escalated              int64               `json:"escalated"`
	Resolved               int64               `json:"resolved"`
	OldestPending          Nullable[time.Time] `json:"oldestPending"`
	TriageAvailable        bool                `json:"triageAvailable"`
	NotificationsAvailable bool                `json:"notificationsAvailable"`
}

// ModerationTriage and ModerationNotifications are integration ports only.
// There is deliberately no default assignment policy or delivery transport.
type ModerationTriage interface {
	AssignReport(context.Context, int64) error
}
type ModerationNotifications interface {
	ReportChanged(context.Context, int64) error
}

type ModerationContentQuery struct {
	ActorUserID                                      uuid.UUID
	Kind, Search, BookID, AuthorID, ChapterID, State string
	Page, PageSize                                   uint32
}
type ModerationInvestigationService interface {
	Search(context.Context, ModerationContentQuery) (ModerationContentPage, error)
	Get(context.Context, ModerationContentQuery, int64) (ModerationContentEntry, error)
	Overview(context.Context, uuid.UUID) (ModerationOverview, error)
}
type moderationInvestigationService struct {
	db   DB
	auth ModerationAuthorizer
}

func NewModerationInvestigationService(db DB, auth ModerationAuthorizer) ModerationInvestigationService {
	return &moderationInvestigationService{db, auth}
}

func (s *moderationInvestigationService) Overview(ctx context.Context, actor uuid.UUID) (ModerationOverview, error) {
	if err := s.auth.AuthorizeModerator(ctx, actor); err != nil {
		return ModerationOverview{}, err
	}
	row, err := store.New(s.db).Report_Overview(ctx)
	if err != nil {
		return ModerationOverview{}, apperror.WrapUnexpectedDBError(err)
	}
	return ModerationOverview{Unreviewed: row.Unreviewed, Escalated: row.Escalated, Resolved: row.Resolved, OldestPending: timeNullableDbToDomain(row.OldestPending)}, nil
}
func (s *moderationInvestigationService) Search(ctx context.Context, f ModerationContentQuery) (ModerationContentPage, error) {
	if err := s.auth.AuthorizeModerator(ctx, f.ActorUserID); err != nil {
		return ModerationContentPage{}, err
	}
	if !oneOf(f.Kind, "chapter", "comment") || !oneOf(f.State, "", "visible", "hidden") {
		return ModerationContentPage{}, ErrModerationFilter
	}
	page, size, limit, offset := normalizeModerationPage(f.Page, f.PageSize)
	f.Search = strings.TrimSpace(f.Search)
	q := store.New(s.db)
	entries := []ModerationContentEntry{}
	var total int64
	var err error
	if f.Kind == "chapter" {
		var rows []store.Moderation_SearchChaptersRow
		rows, err = q.Moderation_SearchChapters(ctx, store.Moderation_SearchChaptersParams{Search: f.Search, BookID: f.BookID, AuthorID: f.AuthorID, State: f.State, PageLimit: limit, PageOffset: offset})
		if err == nil {
			total, err = q.Moderation_CountChapters(ctx, store.Moderation_CountChaptersParams{Search: f.Search, BookID: f.BookID, AuthorID: f.AuthorID, State: f.State})
		}
		for _, r := range rows {
			entries = append(entries, contentEntry(r.ID, r.Title, r.CreatedAt, r.UpdatedAt, r.Visible, r.BookID, r.BookName, r.AuthorID, r.AuthorName))
		}
	} else {
		var rows []store.Moderation_SearchCommentsRow
		rows, err = q.Moderation_SearchComments(ctx, store.Moderation_SearchCommentsParams{Search: f.Search, BookID: f.BookID, AuthorID: f.AuthorID, ChapterID: f.ChapterID, State: f.State, PageLimit: limit, PageOffset: offset})
		if err == nil {
			total, err = q.Moderation_CountComments(ctx, store.Moderation_CountCommentsParams{Search: f.Search, BookID: f.BookID, AuthorID: f.AuthorID, ChapterID: f.ChapterID, State: f.State})
		}
		for _, r := range rows {
			e := contentEntry(r.ID, r.Title, r.CreatedAt, r.UpdatedAt, r.Visible, r.BookID, r.BookName, r.AuthorID, r.AuthorName)
			e.ChapterID = strconv.FormatInt(r.ChapterID, 10)
			e.ChapterName = r.ChapterName
			entries = append(entries, e)
		}
	}
	if err != nil {
		return ModerationContentPage{}, apperror.WrapUnexpectedDBError(err)
	}
	result := moderationPage(entries, page, size, total)
	return ModerationContentPage{entries, result.Page, result.PageSize, result.TotalPages, total}, nil
}
func contentEntry(id int64, title string, created, updated pgtype.Timestamptz, visible bool, book int64, bookName string, author pgtype.UUID, authorName string) ModerationContentEntry {
	return ModerationContentEntry{ID: strconv.FormatInt(id, 10), Title: title, CreatedAt: timeDbToDomain(created), UpdatedAt: timeNullableDbToDomain(updated), Visible: visible, BookID: strconv.FormatInt(book, 10), BookName: bookName, AuthorID: uuidDbToDomain(author).String(), AuthorName: authorName, Replies: []string{}}
}
func (s *moderationInvestigationService) Get(ctx context.Context, f ModerationContentQuery, id int64) (ModerationContentEntry, error) {
	if err := s.auth.AuthorizeModerator(ctx, f.ActorUserID); err != nil {
		return ModerationContentEntry{}, err
	}
	q := store.New(s.db)
	var e ModerationContentEntry
	var err error
	if f.Kind == "chapter" {
		var r store.Moderation_ChapterDetailRow
		r, err = q.Moderation_ChapterDetail(ctx, id)
		e = contentEntry(r.ID, r.Title, r.CreatedAt, r.UpdatedAt, r.Visible, r.BookID, r.BookName, r.AuthorID, r.AuthorName)
		e.Content = r.Content
	} else if f.Kind == "comment" {
		var r store.Moderation_CommentDetailRow
		r, err = q.Moderation_CommentDetail(ctx, id)
		e = contentEntry(r.ID, r.Title, r.CreatedAt, r.UpdatedAt, r.Visible, r.BookID, r.BookName, r.AuthorID, r.AuthorName)
		e.Content = r.Content
		e.ChapterID = strconv.FormatInt(r.ChapterID, 10)
		e.ChapterName = r.ChapterName
		if r.ParentID.Valid {
			e.ParentID = strconv.FormatInt(r.ParentID.Int64, 10)
		}
		if err == nil {
			var replies []int64
			replies, err = q.Moderation_CommentReplies(ctx, pgtype.Int8{Int64: id, Valid: true})
			for _, reply := range replies {
				e.Replies = append(e.Replies, strconv.FormatInt(reply, 10))
			}
		}
	} else {
		return e, ErrModerationFilter
	}
	if err == store.ErrNoRows {
		return e, ErrReportTargetMissing
	}
	if err != nil {
		return e, apperror.WrapUnexpectedDBError(err)
	}
	return e, nil
}
