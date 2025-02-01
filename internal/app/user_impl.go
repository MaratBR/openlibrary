package app

import (
	"context"
	"math"

	"github.com/MaratBR/openlibrary/internal/app/gravatar"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type userService struct {
	queries       *store.Queries
	db            DB
	uploadService *UploadService
}

// ListUsers implements UserService.
func (u *userService) ListUsers(ctx context.Context, req UsersQuery) (UserListResponse, error) {
	var (
		limit  uint = uint(max(1, req.PageSize))
		offset uint = uint((max(1, req.Page) - 1) * req.PageSize)
	)

	dbQuery := store.UsersQuery{
		Query:  req.Query,
		Banned: req.Banned,
		Limit:  limit,
		Offset: offset,
		Roles: commonutil.MapSlice(req.Role, func(role UserRole) store.UserRole {
			return store.UserRole(role)
		}),
	}

	count, err := store.CountUsers(ctx, u.db, &dbQuery)
	if err != nil {
		return UserListResponse{}, wrapUnexpectedDBError(err)
	}

	users, err := store.ListUsers(ctx, u.db, dbQuery)
	if err != nil {
		return UserListResponse{}, wrapUnexpectedDBError(err)
	}

	return UserListResponse{
		Users: commonutil.MapSlice(users, func(user store.UserRow) UserSearchItem {
			return UserSearchItem{
				ID:       uuidDbToDomain(user.ID),
				Name:     user.Name,
				Role:     UserRole(user.Role),
				IsBanned: user.IsBanned,
				JoinedAt: user.JoinedAt,
			}
		}),
		Total:      int32(count),
		Page:       req.Page,
		TotalPages: int32(math.Ceil(float64(count) / float64(req.PageSize))),
	}, nil
}

// FollowUser implements UserService.
func (u *userService) FollowUser(ctx context.Context, cmd FollowUserCommand) error {
	if cmd.Follower == cmd.UserID {
		return ErrFollowYourself
	}

	isFollowing, err := u.queries.IsFollowing(ctx, store.IsFollowingParams{
		FollowerID: uuidDomainToDb(cmd.Follower),
		FollowedID: uuidDomainToDb(cmd.UserID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	if isFollowing {
		return nil
	}

	err = u.queries.InsertUserFollow(ctx, store.InsertUserFollowParams{
		FollowerID: uuidDomainToDb(cmd.Follower),
		FollowedID: uuidDomainToDb(cmd.UserID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// UnfollowUser implements UserService.
func (u *userService) UnfollowUser(ctx context.Context, cmd UnfollowUserCommand) error {
	err := u.queries.DeleteUserFollow(ctx, store.DeleteUserFollowParams{
		FollowerID: uuidDomainToDb(cmd.Follower),
		FollowedID: uuidDomainToDb(cmd.UserID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// GetUserModerationSettings implements UserService.
func (u *userService) GetUserModerationSettings(ctx context.Context, userID uuid.UUID) (*UserModerationSettings, error) {
	user, err := u.queries.GetUserModerationSettings(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}
	return &UserModerationSettings{
		CensoredTags:     user.CensoredTags,
		CensoredTagsMode: CensorMode(user.CensoredTagsMode),
		ShowAdultContent: user.ShowAdultContent,
	}, nil
}

// GetUserPrivacySettings implements UserService.
func (u *userService) GetUserPrivacySettings(ctx context.Context, userID uuid.UUID) (*UserPrivacySettings, error) {
	user, err := u.queries.GetUserPrivacySettings(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}
	return &UserPrivacySettings{
		HideStats:      user.PrivacyHideStats,
		HideComments:   user.PrivacyHideComments,
		HideEmail:      user.PrivacyHideEmail,
		AllowSearching: user.PrivacyAllowSearching,
	}, nil
}

// GetUserAboutSettings implements UserService.
func (u *userService) GetUserAboutSettings(ctx context.Context, userID uuid.UUID) (*UserAboutSettings, error) {
	user, err := u.queries.GetUserAboutSettings(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}
	return &UserAboutSettings{
		About:  user.About,
		Gender: user.Gender,
	}, nil
}

// GetUserCustomizationSettings implements UserService.
func (u *userService) GetUserCustomizationSettings(ctx context.Context, userID uuid.UUID) (*UserCustomizationSetting, error) {
	user, err := u.queries.GetUserCustomizationSettings(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}
	return &UserCustomizationSetting{
		ProfileCSS:       user.ProfileCss,
		EnableProfileCSS: user.EnableProfileCss,
		DefaultTheme:     user.DefaultTheme,
	}, nil
}

// UpdateUserAboutSettings implements UserService.
func (u *userService) UpdateUserAboutSettings(ctx context.Context, userID uuid.UUID, settings UserAboutSettings) error {
	err := u.queries.UpdateUserAboutSettings(ctx, store.UpdateUserAboutSettingsParams{
		About:  settings.About,
		Gender: settings.Gender,
		ID:     uuidDomainToDb(userID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// UpdateUserCustomizationSettings implements UserService.
func (u *userService) UpdateUserCustomizationSettings(ctx context.Context, userID uuid.UUID, settings UserCustomizationSetting) error {
	err := u.queries.UpdateUserCustomizationSettings(ctx, store.UpdateUserCustomizationSettingsParams{
		ProfileCss:       settings.ProfileCSS,
		EnableProfileCss: settings.EnableProfileCSS,
		DefaultTheme:     settings.DefaultTheme,
		ID:               uuidDomainToDb(userID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// UpdateUserModerationSettings implements UserService.
func (u *userService) UpdateUserModerationSettings(ctx context.Context, userID uuid.UUID, settings UserModerationSettings) error {
	err := u.queries.UpdateUserModerationSettings(ctx, store.UpdateUserModerationSettingsParams{
		CensoredTags:     settings.CensoredTags,
		CensoredTagsMode: store.CensorMode(settings.CensoredTagsMode),
		ShowAdultContent: settings.ShowAdultContent,
		ID:               uuidDomainToDb(userID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// UpdateUserPrivacySettings implements UserService.
func (u *userService) UpdateUserPrivacySettings(ctx context.Context, userID uuid.UUID, settings UserPrivacySettings) error {
	err := u.queries.UpdateUserPrivacySettings(ctx, store.UpdateUserPrivacySettingsParams{
		PrivacyHideStats:      settings.HideStats,
		PrivacyHideComments:   settings.HideComments,
		PrivacyHideEmail:      settings.HideEmail,
		PrivacyAllowSearching: settings.AllowSearching,
		ID:                    uuidDomainToDb(userID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(ctx context.Context, cmd UpdateUserCommand) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := u.queries.WithTx(tx)

	if cmd.Password != "" {
		hash, err := hashPassword(cmd.Password)
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedAppError(err)
		}

		err = queries.UpdateUserPassword(ctx, store.UpdateUserPasswordParams{
			ID:           uuidDomainToDb(cmd.UserID),
			PasswordHash: hash,
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedDBError(err)
		}
	}

	err = queries.UpdateUserAboutSettings(ctx, store.UpdateUserAboutSettingsParams{
		ID:     uuidDomainToDb(cmd.UserID),
		About:  cmd.About,
		Gender: cmd.Gender,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	if cmd.Role.Valid {
		println(cmd.Role.Value)
		err = u.updateUserRole(ctx, queries, cmd.UserID, cmd.Role.Value)
		if err != nil {
			rollbackTx(ctx, tx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func (u *userService) updateUserRole(ctx context.Context, queries *store.Queries, userID uuid.UUID, role UserRole) error {
	err := queries.UpdateUserRole(ctx, store.UpdateUserRoleParams{
		ID:   uuidDomainToDb(userID),
		Role: store.UserRole(role),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// GetUserSelfData implements UserService.
func (u *userService) GetUserSelfData(ctx context.Context, userID uuid.UUID) (*SelfUserDto, error) {
	user, err := u.queries.GetUser(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, err
	}

	details := &SelfUserDto{
		ID:                uuidDbToDomain(user.ID),
		Name:              user.Name,
		JoinedAt:          timeDbToDomain(user.JoinedAt),
		IsBanned:          false,
		PreferredTheme:    "dark",
		Role:              UserRole(user.Role),
		ShowAdultContent:  user.ShowAdultContent,
		BookCensoredTags:  user.CensoredTags,
		BookCensoringMode: CensorMode(user.CensoredTagsMode),
	}

	details.Avatar.MD = getUserAvatar(user.Name, 84)
	details.Avatar.LG = getUserAvatar(user.Name, 256)

	return details, nil
}

// GetUserDetails implements UserService.
func (u *userService) GetUserDetails(ctx context.Context, query GetUserQuery) (*UserDetailsDto, error) {
	user, err := u.queries.GetUserWithDetails(ctx, store.GetUserWithDetailsParams{
		ID:          uuidDomainToDb(query.ID),
		ActorUserID: uuidNullableDomainToDb(query.UserID),
	})
	if err != nil {
		if err == store.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	details := &UserDetailsDto{
		ID:             uuidDbToDomain(user.ID),
		Name:           user.Name,
		Email:          user.Name,
		JoinedAt:       timeDbToDomain(user.JoinedAt),
		IsBanned:       user.IsBanned,
		Role:           UserRole(user.Role),
		HasCustomTheme: true,
		About: UserAboutDto{
			Bio:    user.About,
			Gender: user.Gender,
		},
		Following:   int32(user.Following),
		Followers:   int32(user.Followers),
		BooksTotal:  int32(user.BooksTotal),
		HideEmail:   user.PrivacyHideEmail,
		HideStats:   user.PrivacyHideStats,
		IsFollowing: user.IsFollowing,
	}

	if !query.UserID.Valid || details.ID != query.UserID.UUID {
		if user.PrivacyHideStats {
			details.Followers = -1
			details.Following = -1
		}

		if user.PrivacyHideEmail {
			details.Email = ""
		}
	}

	details.Avatar.MD = getUserAvatar(user.Name, 84)
	details.Avatar.LG = getUserAvatar(user.Name, 256)

	return details, nil
}

func NewUserService(db DB) UserService {
	return &userService{
		queries: store.New(db),
		db:      db,
	}
}

func getUserAvatar(name string, size int) string {
	g := gravatar.NewGravatarFromEmail(name)
	g.Default = "robohash"
	g.Size = size
	return g.GetURL()
}
