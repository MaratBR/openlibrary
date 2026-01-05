package upload

import (
	"math"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/joomcode/errorx"
	"github.com/knadh/koanf/v2"
)

var (
	errNamespace = errorx.NewNamespace("file_validator")

	ErrFileTooLarge = errNamespace.NewType("file_too_large", apperror.ErrTraitValidationError).New("file is too large to be uploaded")
)

type FileValidator interface {
	Validate(size int64) error
}

type sizeValidator struct {
	maxSize int64
}

// Validate implements FileValidator.
func (s *sizeValidator) Validate(size int64) error {
	if size > s.maxSize {
		return ErrFileTooLarge
	}
	return nil
}

func NewSizeValidator(maxSize int64) FileValidator {
	return &sizeValidator{maxSize: maxSize}
}

func NewFileValidator(cfg *koanf.Koanf) FileValidator {
	v := cfg.Int64("upload.max-file-size")
	if v <= 0 {
		return NewSizeValidator(math.MaxInt64)
	}

	return NewSizeValidator(v)
}
