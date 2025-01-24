package app

import "errors"

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
