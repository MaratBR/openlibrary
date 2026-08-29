package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestValidateUserDataKey(t *testing.T) {
	for _, key := range []string{"bm:fonts", "bm:tools"} {
		if err := validateUserDataKey(key, false); err != nil {
			t.Fatalf("expected %q to be allowed: %v", key, err)
		}
	}
	if err := validateUserDataKey("other", false); !errors.Is(err, ErrUserDataKeyNotAllowed) {
		t.Fatalf("expected non-whitelisted key to be rejected, got %v", err)
	}
	if err := validateUserDataKey("other", true); err != nil {
		t.Fatalf("expected whitelist bypass to allow key: %v", err)
	}
}

func TestUserDataSizeLimitCannotBeBypassed(t *testing.T) {
	service := &userDataService{}
	err := service.Set(context.Background(), SetUserDataQuery{
		Key:             "not-whitelisted",
		Data:            json.RawMessage(`"` + string(make([]byte, UserDataMaxSize)) + `"`),
		IgnoreWhitelist: true,
	})
	if !errors.Is(err, ErrUserDataTooLarge) {
		t.Fatalf("expected size limit error, got %v", err)
	}
}
