package app

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
)

type reportBookActionsStub struct{ called string }

func (s *reportBookActionsStub) SearchBooks(context.Context, SearchModerationBooksQuery) (ModerationPage[ModerationBookListEntry], error) {
	return ModerationPage[ModerationBookListEntry]{}, nil
}
func (s *reportBookActionsStub) GetBookInfo(context.Context, GetBookInfoQuery) (BookModerationInfo, error) {
	return BookModerationInfo{}, nil
}
func (s *reportBookActionsStub) GetBookLog(context.Context, GetBookLogQuery) (BookLogResult, error) {
	return BookLogResult{}, nil
}
func (s *reportBookActionsStub) GetBookChapters(context.Context, GetBookInfoQuery) ([]BookModerationChapter, error) {
	return nil, nil
}
func (s *reportBookActionsStub) mark(value string) error { s.called = value; return nil }
func (s *reportBookActionsStub) ChangeAgeRating(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.change_age_rating")
}
func (s *reportBookActionsStub) ChangeSummary(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.change_summary")
}
func (s *reportBookActionsStub) BanBook(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.restrict")
}
func (s *reportBookActionsStub) ShadowBanBook(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.shadow_restrict")
}
func (s *reportBookActionsStub) PermanentlyRemoveBook(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.permanent_remove")
}
func (s *reportBookActionsStub) UnBanBook(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.restore")
}
func (s *reportBookActionsStub) UnShadowBanBook(context.Context, ModerationPerformBookActionCommand) error {
	return s.mark("book.restore_shadow")
}

type reportContentActionsStub struct{ called string }

func (s *reportContentActionsStub) mark(value string) error { s.called = value; return nil }
func (s *reportContentActionsStub) HideChapter(context.Context, ModerateChapterCommand) error {
	return s.mark("chapter.hide")
}
func (s *reportContentActionsStub) RestoreChapter(context.Context, ModerateChapterCommand) error {
	return s.mark("chapter.restore")
}
func (s *reportContentActionsStub) RemoveComment(context.Context, ModerateCommentCommand) error {
	return s.mark("comment.remove")
}
func (s *reportContentActionsStub) RestoreComment(context.Context, ModerateCommentCommand) error {
	return s.mark("comment.restore")
}
func (s *reportContentActionsStub) BanUser(context.Context, BanUserCommand) error {
	return s.mark("user.temporary_ban")
}
func (s *reportContentActionsStub) PermanentlyBanUser(context.Context, BanUserCommand) error {
	return s.mark("user.permanent_ban")
}
func (s *reportContentActionsStub) UnbanUser(context.Context, BanUserCommand) error {
	return s.mark("user.unban")
}
func (s *reportContentActionsStub) RenameUser(context.Context, ModerateUserProfileCommand) error {
	return s.mark("user.rename")
}
func (s *reportContentActionsStub) ChangeUserAbout(context.Context, ModerateUserProfileCommand) error {
	return s.mark("user.change_about")
}

func TestReportActionExecutorDispatchesEveryAvailableAction(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	tests := []struct {
		action string
		report Report
	}{
		{"user.temporary_ban", Report{TargetType: ReportTargetUser, TargetID: userID.String()}},
		{"user.permanent_ban", Report{TargetType: ReportTargetUser, TargetID: userID.String()}},
		{"user.unban", Report{TargetType: ReportTargetUser, TargetID: userID.String()}},
		{"user.rename", Report{TargetType: ReportTargetUser, TargetID: userID.String()}},
		{"user.change_about", Report{TargetType: ReportTargetUser, TargetID: userID.String()}},
		{"book.restrict", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"book.restore", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"book.shadow_restrict", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"book.restore_shadow", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"book.permanent_remove", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"book.change_age_rating", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"book.change_summary", Report{TargetType: ReportTargetBook, TargetID: "42"}},
		{"chapter.hide", Report{TargetType: ReportTargetBook, TargetID: "42", BookChapterID: Value[int64](7)}},
		{"chapter.restore", Report{TargetType: ReportTargetBook, TargetID: "42", BookChapterID: Value[int64](7)}},
		{"comment.remove", Report{TargetType: ReportTargetComment, TargetID: "9"}},
		{"comment.restore", Report{TargetType: ReportTargetComment, TargetID: "9"}},
	}
	for _, test := range tests {
		t.Run(test.action, func(t *testing.T) {
			books, content := &reportBookActionsStub{}, &reportContentActionsStub{}
			executor := NewReportActionExecutor(books, content)
			if err := executor.Execute(context.Background(), test.report, DecideModerationReportCommand{Action: test.action, PolicyReason: "reason", Payload: []byte(`{"value":"G","until":"2999-01-01T00:00:00Z"}`)}); err != nil {
				t.Fatal(err)
			}
			called := books.called
			if called == "" {
				called = content.called
			}
			if called != test.action {
				t.Fatalf("expected %q, dispatched %q", test.action, called)
			}
		})
	}
}
