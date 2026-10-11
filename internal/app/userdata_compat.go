package app

// Compatibility with existing callers; new feature code belongs in userdata.
import (
	"github.com/MaratBR/openlibrary/internal/app/userdata"
)

const UserDataMaxSize = userdata.UserDataMaxSize

var UserDataErrors = userdata.UserDataErrors
var ErrTypeUserDataKeyNotAllowed = userdata.ErrTypeUserDataKeyNotAllowed
var ErrUserDataKeyNotAllowed = userdata.ErrUserDataKeyNotAllowed
var ErrUserDataTooLarge = userdata.ErrUserDataTooLarge

type GetUserDataQuery = userdata.GetUserDataQuery
type SetUserDataQuery = userdata.SetUserDataQuery
type UserDataService = userdata.UserDataService

func NewUserDataService(db DB) UserDataService { return userdata.NewUserDataService(db) }
