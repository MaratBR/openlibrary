package siteconfig

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
)

var (
	PasswordError = apperror.AppErrors.NewType("password")
)

var (
	regexDigits    = regexp.MustCompile(`\d`)
	regexUppercase = regexp.MustCompile("[A-Z]")
	regexLowercase = regexp.MustCompile("[a-z]")
)

func ValidatePassword(pwd string, r PasswordRequirements) error {
	if r.Digits && !regexDigits.Match([]byte(pwd)) {
		return PasswordError.New("password must have digits")
	}
	if r.MinLength > 0 && len(pwd) < r.MinLength {
		return PasswordError.New("%s", fmt.Sprintf("password must be at least %d characters long", r.MinLength))
	}
	if r.DifferentCases && !(regexLowercase.Match([]byte(pwd)) && regexUppercase.Match([]byte(pwd))) {
		return PasswordError.New("password must contain at least one uppercase and one lowercase letter")
	}
	if r.SymbolsEnabled && len(r.Symbols) > 0 {
		var containsSymbol bool
		for _, rune := range r.Symbols {
			if strings.ContainsRune(pwd, rune) {
				containsSymbol = true
				break
			}
		}

		if !containsSymbol {
			return PasswordError.New("password")
		}
	}
	return nil
}
