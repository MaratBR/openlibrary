package bookfont

import (
	"strings"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/knadh/koanf/v2"
)

const defaultMaxPerChapter = 10

var ErrLimitExceeded = apperror.AppErrors.NewType("chapter_font_limit_exceeded", apperror.ErrTraitValidationError)

type Policy struct {
	MaxPerChapter int
	Whitelist     []string
}

func NewPolicy(cfg *koanf.Koanf) Policy {
	policy := Policy{
		MaxPerChapter: cfg.Int("chapter-fonts.max-per-chapter"),
		Whitelist:     cfg.Strings("chapter-fonts.whitelist"),
	}
	if policy.MaxPerChapter == 0 && !cfg.Exists("chapter-fonts.max-per-chapter") {
		policy.MaxPerChapter = defaultMaxPerChapter
	}
	return policy
}

func (p Policy) Validate(fonts []string) error {
	limit := p.MaxPerChapter
	if limit == 0 && p.Whitelist == nil {
		limit = defaultMaxPerChapter
	}

	whitelisted := normalizedSet(p.Whitelist)
	used := normalizedSet(fonts)
	for font := range whitelisted {
		delete(used, font)
	}
	if len(used) > limit {
		return ErrLimitExceeded.New(
			"chapter uses %d limited fonts, maximum is %d", len(used), limit,
		)
	}
	return nil
}

func normalizedSet(fonts []string) map[string]struct{} {
	result := make(map[string]struct{}, len(fonts))
	for _, font := range fonts {
		font = strings.ToLower(strings.TrimSpace(font))
		if font != "" {
			result[font] = struct{}{}
		}
	}
	return result
}
