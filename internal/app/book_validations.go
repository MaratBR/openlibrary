package app

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
)

const MaxChapterNameLength = 200

var (
	ErrEmptyBookName   = errors.New("empty book name")
	ErrBookNameTooLong = errors.New("invalid book name")
	BookSummaryTooLong = errors.New("book summary is too long")
)

func validateBookName(name string) error {
	if name == "" {
		return ErrEmptyBookName
	}
	if len(name) > 500 {
		return ErrBookNameTooLong
	}
	return nil
}

func validateBookSummary(summary string) error {
	if len(summary) > 100_000 {
		return BookSummaryTooLong
	}
	return nil
}

func validateChapterName(name string) error {
	if strings.TrimSpace(name) == "" {
		return apperror.ValidationError.New("chapter name is required")
	}
	if utf8.RuneCountInString(name) > MaxChapterNameLength {
		return apperror.ValidationError.New("chapter name cannot exceed %d characters", MaxChapterNameLength)
	}
	return nil
}
