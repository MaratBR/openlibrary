package app

import (
	"errors"
	"testing"
)

func TestValidateUserRoleChange(t *testing.T) {
	for _, current := range AllRoles {
		for _, next := range AllRoles {
			t.Run(string(current)+"_to_"+string(next), func(t *testing.T) {
				err := ValidateUserRoleChange(current, next)
				forbidden := current != next && (current == RoleSystem || next == RoleSystem)
				if forbidden {
					if !errors.Is(err, ErrInvalidUserRole) {
						t.Fatalf("expected invalid role error, got %v", err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
	if !errors.Is(ValidateUserRoleChange(RoleUser, UserRole("invalid")), ErrInvalidUserRole) {
		t.Fatal("unknown roles must be rejected")
	}
}
