package csrf_test

import (
	"testing"

	"github.com/MaratBR/openlibrary/internal/csrf"
)

func TestHMACCsrfToken(t *testing.T) {
	secret := "mysupersecretthing"
	sid := "uhuirh2wuierfhweiurgl;oeither"

	csrfToken := csrf.GenerateHMACCsrfToken(sid, secret)
	valid := csrf.VerifyHMACCsrfToken(csrfToken, sid, secret)
	if !valid {
		t.Fatal("invalid token")
	}
}

func TestHMACCsrfTokenInvalidSecret(t *testing.T) {
	secret := "mysupersecretthing"
	sid := "uhuirh2wuierfhweiurgl;oeither"

	csrfToken := csrf.GenerateHMACCsrfToken(sid, secret)
	valid := csrf.VerifyHMACCsrfToken(csrfToken, sid, "someothersecret")
	if valid {
		t.Fatal("expected token to be invalid")
	}
}

func TestHMACCsrfTokenInvalidSid(t *testing.T) {
	secret := "mysupersecretthing"
	sid := "uhuirh2wuierfhweiurgl;oeither"

	csrfToken := csrf.GenerateHMACCsrfToken(sid, secret)
	valid := csrf.VerifyHMACCsrfToken(csrfToken, "someothersidvalue", secret)
	if valid {
		t.Fatal("expected token to be invalid")
	}
}

func TestHMACCsrfTokenInvalidToken(t *testing.T) {
	secret := "mysupersecretthing"
	sid := "uhuirh2wuierfhweiurgl;oeither"

	csrfToken := csrf.GenerateHMACCsrfToken(sid, secret)
	valid := csrf.VerifyHMACCsrfToken(csrfToken+"asd", sid, secret)
	if valid {
		t.Fatal("expected token to be invalid")
	}
}
