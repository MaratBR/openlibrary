package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type UserRole string

var (
	ErrUserNotFound   = AppErrors.NewType("user_not_found", ErrTraitEntityNotFound).New("user not found")
	ErrFollowYourself = AppErrors.NewType("follow_yourself").New("you can't follow yourself")
)

type UserDetailsDto struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar struct {
		LG string `json:"lg"`
		MD string `json:"md"`
	} `json:"avatar"`
	JoinedAt       time.Time       `json:"joinedAt"`
	IsBanned       bool            `json:"isBlocked"`
	Role           UserRole        `json:"role"`
	HasCustomTheme bool            `json:"hasCustomTheme"`
	About          UserAboutDto    `json:"about"`
	Books          []AuthorBookDto `json:"books"`
	Followers      int32           `json:"followers"`
	Following      int32           `json:"following"`
	Favorites      int32           `json:"favorites"`
	BooksTotal     int32           `json:"booksTotal"`
	HideEmail      bool            `json:"hideEmail"`
	HideStats      bool            `json:"hideStats"`
	HideFavorites  bool            `json:"hideFavorites"`
	IsFollowing    bool            `json:"isFollowing"`
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

type UserService interface {
	GetUserPrivacySettings(ctx context.Context, userID uuid.UUID) (*UserPrivacySettings, error)
	GetUserModerationSettings(ctx context.Context, userID uuid.UUID) (*UserModerationSettings, error)
	GetUserCustomizationSettings(ctx context.Context, userID uuid.UUID) (*UserCustomizationSetting, error)
	GetUserAboutSettings(ctx context.Context, userID uuid.UUID) (*UserAboutSettings, error)

	UpdateUserPrivacySettings(ctx context.Context, userID uuid.UUID, settings UserPrivacySettings) error
	UpdateUserModerationSettings(ctx context.Context, userID uuid.UUID, settings UserModerationSettings) error
	UpdateUserCustomizationSettings(ctx context.Context, userID uuid.UUID, settings UserCustomizationSetting) error
	UpdateUserAboutSettings(ctx context.Context, userID uuid.UUID, settings UserAboutSettings) error

	GetUserDetails(ctx context.Context, query GetUserQuery) (*UserDetailsDto, error)
	GetUserSelfData(ctx context.Context, userID uuid.UUID) (*SelfUserDto, error)

	FollowUser(ctx context.Context, cmd FollowUserCommand) error
	UnfollowUser(ctx context.Context, cmd UnfollowUserCommand) error
}
