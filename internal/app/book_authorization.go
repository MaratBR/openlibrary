package app

import "github.com/gofrs/uuid"

type userBookPermissionState struct {
	CanView bool
	IsOwner bool
}

type getUserBookPermissionsStateRequest struct {
	UserID            uuid.NullUUID
	BookAuthorID      uuid.UUID
	IsPubliclyVisible bool
}

func getUserBookPermissionsState(
	req getUserBookPermissionsStateRequest,
) userBookPermissionState {
	if req.UserID.Valid {
		return userBookPermissionState{
			CanView: req.IsPubliclyVisible || req.UserID.UUID == req.BookAuthorID,
			IsOwner: req.UserID.UUID == req.BookAuthorID,
		}
	} else {
		return userBookPermissionState{
			CanView: req.IsPubliclyVisible,
			IsOwner: false,
		}
	}
}
