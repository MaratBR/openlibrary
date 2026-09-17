package app

import "uuid"

type userBookPermissionState struct {
	CanView bool
	IsOwner bool
}

type getUserBookPermissionsStateRequest struct {
	UserID            Nullable[uuid.UUID]
	BookAuthorID      uuid.UUID
	IsPubliclyVisible bool
}

func getUserBookPermissionsState(
	req getUserBookPermissionsStateRequest,
) userBookPermissionState {
	if req.UserID.Valid {
		return userBookPermissionState{
			CanView: req.IsPubliclyVisible || req.UserID.Value == req.BookAuthorID,
			IsOwner: req.UserID.Value == req.BookAuthorID,
		}
	} else {
		return userBookPermissionState{
			CanView: req.IsPubliclyVisible,
			IsOwner: false,
		}
	}
}
