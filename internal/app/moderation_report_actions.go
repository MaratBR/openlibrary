package app

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

var ErrInvalidReportAction = apperror.AppErrors.NewType("invalid_report_action").New("action is not valid for this report target")

type ReportActionExecutor interface {
	Execute(context.Context, Report, DecideModerationReportCommand) error
	AvailableActions(Report) []string
}

type reportActionPayload struct {
	Value string    `json:"value"`
	Until time.Time `json:"until"`
}

type reportActionExecutor struct {
	books   ModerationBookService
	content ContentModerationService
}

func (e *reportActionExecutor) AvailableActions(report Report) []string {
	switch report.TargetType {
	case ReportTargetUser:
		return []string{"user.temporary_ban", "user.permanent_ban", "user.unban", "user.rename", "user.change_about"}
	case ReportTargetComment:
		return []string{"comment.remove", "comment.restore"}
	case ReportTargetBook:
		actions := []string{"book.restrict", "book.restore", "book.shadow_restrict", "book.restore_shadow", "book.permanent_remove", "book.change_age_rating", "book.change_summary"}
		if report.BookChapterID.Valid {
			actions = append(actions, "chapter.hide", "chapter.restore")
		}
		return actions
	default:
		return []string{}
	}
}

func (e *reportActionExecutor) Execute(ctx context.Context, report Report, cmd DecideModerationReportCommand) error {
	var payload reportActionPayload
	if len(cmd.Payload) > 0 {
		if err := json.Unmarshal(cmd.Payload, &payload); err != nil {
			return ErrInvalidReportAction
		}
	}
	reason := cmd.PolicyReason
	switch report.TargetType {
	case ReportTargetUser:
		id, err := uuid.FromString(report.TargetID)
		if err != nil {
			return ErrInvalidReportAction
		}
		ban := BanUserCommand{ActorUserID: cmd.ActorUserID, UserID: id, Until: payload.Until, Reason: reason}
		profile := ModerateUserProfileCommand{ActorUserID: cmd.ActorUserID, UserID: id, Value: payload.Value, Reason: reason}
		switch cmd.Action {
		case "user.temporary_ban":
			return e.content.BanUser(ctx, ban)
		case "user.permanent_ban":
			return e.content.PermanentlyBanUser(ctx, ban)
		case "user.unban":
			return e.content.UnbanUser(ctx, ban)
		case "user.rename":
			return e.content.RenameUser(ctx, profile)
		case "user.change_about":
			return e.content.ChangeUserAbout(ctx, profile)
		}
	case ReportTargetBook:
		id, err := strconv.ParseInt(report.TargetID, 10, 64)
		if err != nil {
			return ErrInvalidReportAction
		}
		book := ModerationPerformBookActionCommand{ActorUserID: cmd.ActorUserID, BookID: id, Reason: reason, Value: payload.Value}
		switch cmd.Action {
		case "book.restrict":
			return e.books.BanBook(ctx, book)
		case "book.restore":
			return e.books.UnBanBook(ctx, book)
		case "book.shadow_restrict":
			return e.books.ShadowBanBook(ctx, book)
		case "book.restore_shadow":
			return e.books.UnShadowBanBook(ctx, book)
		case "book.permanent_remove":
			return e.books.PermanentlyRemoveBook(ctx, book)
		case "book.change_age_rating":
			return e.books.ChangeAgeRating(ctx, book)
		case "book.change_summary":
			return e.books.ChangeSummary(ctx, book)
		case "chapter.hide", "chapter.restore":
			if !report.BookChapterID.Valid {
				return ErrInvalidReportAction
			}
			chapter := ModerateChapterCommand{ActorUserID: cmd.ActorUserID, ChapterID: report.BookChapterID.Value, Reason: reason}
			if cmd.Action == "chapter.hide" {
				return e.content.HideChapter(ctx, chapter)
			}
			return e.content.RestoreChapter(ctx, chapter)
		}
	case ReportTargetComment:
		id, err := strconv.ParseInt(report.TargetID, 10, 64)
		if err != nil {
			return ErrInvalidReportAction
		}
		comment := ModerateCommentCommand{ActorUserID: cmd.ActorUserID, CommentID: id, Reason: reason}
		if cmd.Action == "comment.remove" {
			return e.content.RemoveComment(ctx, comment)
		}
		if cmd.Action == "comment.restore" {
			return e.content.RestoreComment(ctx, comment)
		}
	}
	return ErrInvalidReportAction
}

func NewReportActionExecutor(books ModerationBookService, content ContentModerationService) ReportActionExecutor {
	return &reportActionExecutor{books: books, content: content}
}
