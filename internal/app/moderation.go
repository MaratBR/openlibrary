package app

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

var (
	ErrModerationForbidden = apperror.AppErrors.NewType("moderation_forbidden", apperror.ErrTraitAuthorizationIssue).New("moderator privileges are required")
	ErrModerationReason    = apperror.AppErrors.NewType("moderation_reason_required").New("a moderation reason is required")
)

// ModerationAuthorizer centralizes the role check shared by every moderation
// service. Keeping it in the application layer prevents controllers from being
// able to accidentally bypass authorization.
type ModerationAuthorizer interface {
	AuthorizeModerator(ctx context.Context, actorUserID uuid.UUID) error
}

type moderationAuthorizer struct {
	users UserService
}

func (a *moderationAuthorizer) AuthorizeModerator(ctx context.Context, actorUserID uuid.UUID) error {
	user, err := a.users.GetUserSelfData(ctx, actorUserID)
	if err != nil {
		return err
	}
	if !user.Role.IsModOrHigher() {
		return ErrModerationForbidden
	}
	return nil
}

func NewModerationAuthorizer(users UserService) ModerationAuthorizer {
	return &moderationAuthorizer{users: users}
}

func validateModerationReason(reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrModerationReason
	}
	return nil
}

type ModerationAction string

const (
	ModerationActionHideChapter    ModerationAction = "hide_chapter"
	ModerationActionRestoreChapter ModerationAction = "restore_chapter"
	ModerationActionRemoveComment  ModerationAction = "remove_comment"
	ModerationActionRestoreComment ModerationAction = "restore_comment"
	ModerationActionBanUser        ModerationAction = "ban_user"
	ModerationActionUnbanUser      ModerationAction = "unban_user"
	ModerationActionRenameUser     ModerationAction = "rename_user"
	ModerationActionChangeAbout    ModerationAction = "change_user_about"
)

type ModerationAuditEntry struct {
	ID          int64
	Time        time.Time
	Action      ModerationAction
	ActorUserID uuid.UUID
	TargetType  string
	TargetID    string
	Reason      string
	Payload     []byte
}

// ContentModerationRepository is the persistence port for application-level
// moderation. Every mutation must update its target and append its audit entry
// atomically in the generalized moderation log.
type ContentModerationRepository interface {
	SetChapterVisibilityAndLog(ctx context.Context, chapterID int64, visible bool, log ModerationAuditEntry) error
	SetCommentRemovedAndLog(ctx context.Context, commentID int64, removed bool, log ModerationAuditEntry) error
	SetUserBanAndLog(ctx context.Context, userID uuid.UUID, expiresAt time.Time, banned bool, log ModerationAuditEntry) error
	RenameUserAndLog(ctx context.Context, userID uuid.UUID, name string, log ModerationAuditEntry) error
	ChangeUserAboutAndLog(ctx context.Context, userID uuid.UUID, about string, log ModerationAuditEntry) error
}

type ModerateChapterCommand struct {
	ActorUserID uuid.UUID
	ChapterID   int64
	Reason      string
}
type ModerateCommentCommand struct {
	ActorUserID uuid.UUID
	CommentID   int64
	Reason      string
}
type BanUserCommand struct {
	ActorUserID, UserID uuid.UUID
	Until               time.Time
	Reason              string
}
type ModerateUserProfileCommand struct {
	ActorUserID, UserID uuid.UUID
	Value, Reason       string
}

type ContentModerationService interface {
	HideChapter(context.Context, ModerateChapterCommand) error
	RestoreChapter(context.Context, ModerateChapterCommand) error
	RemoveComment(context.Context, ModerateCommentCommand) error
	RestoreComment(context.Context, ModerateCommentCommand) error
	BanUser(context.Context, BanUserCommand) error
	PermanentlyBanUser(context.Context, BanUserCommand) error
	UnbanUser(context.Context, BanUserCommand) error
	RenameUser(context.Context, ModerateUserProfileCommand) error
	ChangeUserAbout(context.Context, ModerateUserProfileCommand) error
}

type contentModerationService struct {
	auth ModerationAuthorizer
	repo ContentModerationRepository
	now  func() time.Time
}

func (s *contentModerationService) entry(actor uuid.UUID, action ModerationAction, targetType, targetID, reason string) ModerationAuditEntry {
	return ModerationAuditEntry{ID: GenID(), Time: s.now(), ActorUserID: actor, Action: action, TargetType: targetType, TargetID: targetID, Reason: strings.TrimSpace(reason)}
}

func (s *contentModerationService) authorize(ctx context.Context, actor uuid.UUID, reason string) error {
	if err := validateModerationReason(reason); err != nil {
		return err
	}
	return s.auth.AuthorizeModerator(ctx, actor)
}

func (s *contentModerationService) HideChapter(ctx context.Context, c ModerateChapterCommand) error {
	return s.chapter(ctx, c, false, ModerationActionHideChapter)
}
func (s *contentModerationService) RestoreChapter(ctx context.Context, c ModerateChapterCommand) error {
	return s.chapter(ctx, c, true, ModerationActionRestoreChapter)
}
func (s *contentModerationService) chapter(ctx context.Context, c ModerateChapterCommand, visible bool, action ModerationAction) error {
	if err := s.authorize(ctx, c.ActorUserID, c.Reason); err != nil {
		return err
	}
	return s.repo.SetChapterVisibilityAndLog(ctx, c.ChapterID, visible, s.entry(c.ActorUserID, action, "chapter", strconv.FormatInt(c.ChapterID, 10), c.Reason))
}
func (s *contentModerationService) RemoveComment(ctx context.Context, c ModerateCommentCommand) error {
	return s.comment(ctx, c, true, ModerationActionRemoveComment)
}
func (s *contentModerationService) RestoreComment(ctx context.Context, c ModerateCommentCommand) error {
	return s.comment(ctx, c, false, ModerationActionRestoreComment)
}
func (s *contentModerationService) comment(ctx context.Context, c ModerateCommentCommand, removed bool, action ModerationAction) error {
	if err := s.authorize(ctx, c.ActorUserID, c.Reason); err != nil {
		return err
	}
	return s.repo.SetCommentRemovedAndLog(ctx, c.CommentID, removed, s.entry(c.ActorUserID, action, "comment", strconv.FormatInt(c.CommentID, 10), c.Reason))
}
func (s *contentModerationService) BanUser(ctx context.Context, c BanUserCommand) error {
	return s.ban(ctx, c, c.Until, true, ModerationActionBanUser)
}
func (s *contentModerationService) PermanentlyBanUser(ctx context.Context, c BanUserCommand) error {
	// PostgreSQL supports this date and it is deliberately finite so existing
	// expiry-based ban handling needs no special case.
	return s.ban(ctx, c, time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC), true, ModerationActionBanUser)
}
func (s *contentModerationService) UnbanUser(ctx context.Context, c BanUserCommand) error {
	return s.ban(ctx, c, s.now(), false, ModerationActionUnbanUser)
}
func (s *contentModerationService) ban(ctx context.Context, c BanUserCommand, until time.Time, banned bool, action ModerationAction) error {
	if err := s.authorize(ctx, c.ActorUserID, c.Reason); err != nil {
		return err
	}
	if banned && !until.After(s.now()) {
		return apperror.AppErrors.NewType("invalid_ban_expiry").New("ban expiry must be in the future")
	}
	return s.repo.SetUserBanAndLog(ctx, c.UserID, until, banned, s.entry(c.ActorUserID, action, "user", c.UserID.String(), c.Reason))
}
func (s *contentModerationService) RenameUser(ctx context.Context, c ModerateUserProfileCommand) error {
	if err := s.authorize(ctx, c.ActorUserID, c.Reason); err != nil {
		return err
	}
	if err := ValidateUserName(c.Value); err != nil {
		return err
	}
	return s.repo.RenameUserAndLog(ctx, c.UserID, c.Value, s.entry(c.ActorUserID, ModerationActionRenameUser, "user", c.UserID.String(), c.Reason))
}
func (s *contentModerationService) ChangeUserAbout(ctx context.Context, c ModerateUserProfileCommand) error {
	if err := s.authorize(ctx, c.ActorUserID, c.Reason); err != nil {
		return err
	}
	return s.repo.ChangeUserAboutAndLog(ctx, c.UserID, c.Value, s.entry(c.ActorUserID, ModerationActionChangeAbout, "user", c.UserID.String(), c.Reason))
}

func NewContentModerationService(auth ModerationAuthorizer, repo ContentModerationRepository) ContentModerationService {
	return &contentModerationService{auth: auth, repo: repo, now: time.Now}
}

type LoginHistoryEntry struct {
	IPAddress, UserAgent string
	LoggedInAt           time.Time
}
type GetLoginHistoryQuery struct{ ActorUserID, UserID uuid.UUID }
type LoginHistoryService interface {
	GetUserLoginHistory(context.Context, GetLoginHistoryQuery) ([]LoginHistoryEntry, error)
}
type loginHistoryService struct {
	auth     ModerationAuthorizer
	sessions SessionService
}

func (s *loginHistoryService) GetUserLoginHistory(ctx context.Context, q GetLoginHistoryQuery) ([]LoginHistoryEntry, error) {
	if err := s.auth.AuthorizeModerator(ctx, q.ActorUserID); err != nil {
		return nil, err
	}
	rows, err := s.sessions.GetByUserID(ctx, q.UserID)
	if err != nil {
		return nil, err
	}
	result := make([]LoginHistoryEntry, len(rows))
	for i, row := range rows {
		result[i] = LoginHistoryEntry{IPAddress: row.IpAddress, UserAgent: row.UserAgent, LoggedInAt: row.CreatedAt}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].LoggedInAt.After(result[j].LoggedInAt) })
	return result, nil
}
func NewLoginHistoryService(auth ModerationAuthorizer, sessions SessionService) LoginHistoryService {
	return &loginHistoryService{auth: auth, sessions: sessions}
}
