package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type UserRole string

var (
	RoleUser      = UserRole("user")
	RoleAdmin     = UserRole("admin")
	RoleSystem    = UserRole("system")
	RoleModerator = UserRole("moderator")

	AllRoles = []UserRole{
		RoleUser,
		RoleAdmin,
		RoleSystem,
		RoleModerator,
	}
)

func (r UserRole) IsAdmin() bool {
	return r == RoleAdmin || r == RoleSystem
}

func (r UserRole) IsModOrHigher() bool {
	return r.IsAdmin() || r == RoleModerator
}

func ParseUserRole(role string) (UserRole, error) {
	switch role {
	case "user":
		return RoleUser, nil
	case "admin":
		return RoleAdmin, nil
	case "system":
		return RoleSystem, nil
	case "moderator":
		return RoleModerator, nil
	default:
		return "", ErrInvalidUserRole
	}
}

var (
	ErrUserNotFound    = AppErrors.NewType("user_not_found", ErrTraitEntityNotFound).New("user not found")
	ErrFollowYourself  = AppErrors.NewType("follow_yourself").New("you can't follow yourself")
	ErrInvalidUserRole = AppErrors.NewType("invalid_user_role").New("invalid user role")
)

type UserDetailsDto struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar struct {
		LG string `json:"lg"`
		MD string `json:"md"`
	} `json:"avatar"`
	JoinedAt       time.Time    `json:"joinedAt"`
	IsBanned       bool         `json:"isBlocked"`
	Role           UserRole     `json:"role"`
	HasCustomTheme bool         `json:"hasCustomTheme"`
	About          UserAboutDto `json:"about"`
	Followers      int32        `json:"followers"`
	Following      int32        `json:"following"`
	BooksTotal     int32        `json:"booksTotal"`
	HideEmail      bool         `json:"hideEmail"`
	HideStats      bool         `json:"hideStats"`
	HideFavorites  bool         `json:"hideFavorites"`
	IsFollowing    bool         `json:"isFollowing"`
}

type UserAboutDto struct {
	Bio    string `json:"bio"`
	Gender string `json:"gender"`
}

type SelfUserDto struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Role   UserRole  `json:"role"`
	Avatar struct {
		LG string `json:"lg"`
		MD string `json:"md"`
	} `json:"avatar"`
	JoinedAt          time.Time  `json:"joinedAt"`
	IsBanned          bool       `json:"isBlocked"`
	PreferredTheme    string     `json:"preferredTheme"`
	ShowAdultContent  bool       `json:"showAdultContent"`
	BookCensoredTags  []string   `json:"bookCensoredTags"`
	BookCensoringMode CensorMode `json:"bookCensoringMode"`
}

type GetUserQuery struct {
	ID     uuid.UUID
	UserID uuid.NullUUID
}

type UserPrivacySettings struct {
	HideStats      bool `json:"hideStats"`
	HideFavorites  bool `json:"hideFavorites"`
	HideComments   bool `json:"hideComments"`
	HideEmail      bool `json:"hideEmail"`
	AllowSearching bool `json:"allowSearching"`
}

type CensorMode store.CensorMode

type UserModerationSettings struct {
	ShowAdultContent bool       `json:"showAdultContent"`
	CensoredTags     []string   `json:"censoredTags"`
	CensoredTagsMode CensorMode `json:"censoredTagsMode"`
}

type UserAboutSettings struct {
	About  string `json:"about"`
	Status string `json:"status"`
	Gender string `json:"gender"`
}

type UserCustomizationSetting struct {
	ProfileCSS       string `json:"profileCss"`
	DefaultTheme     string `json:"defaultTheme"`
	EnableProfileCSS bool   `json:"enableProfileCss"`
}

type FollowUserCommand struct {
	UserID   uuid.UUID
	Follower uuid.UUID
}

type UnfollowUserCommand struct {
	UserID   uuid.UUID
	Follower uuid.UUID
}

type UsersQuery struct {
	Banned bool
	Role   []UserRole
	Query  string

	Page     uint32
	PageSize uint32
}

type UserSearchItem struct {
	ID       uuid.UUID
	Name     string
	Role     UserRole
	IsBanned bool
	JoinedAt time.Time
	Avatar   string
}

type UserListResponse struct {
	Users      []UserSearchItem `json:"users"`
	Total      int32            `json:"total"`
	Page       uint32           `json:"page"`
	TotalPages int32            `json:"totalPages"`
}

type UpdateUserCommand struct {
	UserID      uuid.UUID
	ActorUserID uuid.NullUUID
	Password    string
	Role        Nullable[UserRole]
	About       string
	Gender      string
}

type UserService interface {
	GetUserPrivacySettings(ctx context.Context, userID uuid.UUID) (*UserPrivacySettings, error)
	GetUserModerationSettings(ctx context.Context, userID uuid.UUID) (*UserModerationSettings, error)
	GetUserCustomizationSettings(ctx context.Context, userID uuid.UUID) (*UserCustomizationSetting, error)
	GetUserAboutSettings(ctx context.Context, userID uuid.UUID) (*UserAboutSettings, error)

	UpdateUserPrivacySettings(ctx context.Context, userID uuid.UUID, settings UserPrivacySettings) error
	UpdateUserModerationSettings(ctx context.Context, userID uuid.UUID, settings UserModerationSettings) error
	UpdateUserCustomizationSettings(ctx context.Context, userID uuid.UUID, settings UserCustomizationSetting) error
	UpdateUserAboutSettings(ctx context.Context, userID uuid.UUID, settings UserAboutSettings) error

	UpdateUser(ctx context.Context, cmd UpdateUserCommand) error

	GetUserDetails(ctx context.Context, query GetUserQuery) (*UserDetailsDto, error)
	GetUserSelfData(ctx context.Context, userID uuid.UUID) (*SelfUserDto, error)

	FollowUser(ctx context.Context, cmd FollowUserCommand) error
	UnfollowUser(ctx context.Context, cmd UnfollowUserCommand) error

	ListUsers(ctx context.Context, req UsersQuery) (UserListResponse, error)
}
