package bookfont

import (
	"testing"

	"github.com/joomcode/errorx"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/require"
)

func TestNewPolicyDefaults(t *testing.T) {
	policy := NewPolicy(koanf.New("."))
	require.Equal(t, defaultMaxPerChapter, policy.MaxPerChapter)
}

func TestNewPolicyFromKoanf(t *testing.T) {
	cfg := koanf.New(".")
	require.NoError(t, cfg.Set("chapter-fonts.max-per-chapter", 3))
	require.NoError(t, cfg.Set("chapter-fonts.whitelist", []string{"Poppins"}))

	policy := NewPolicy(cfg)
	require.Equal(t, 3, policy.MaxPerChapter)
	require.Equal(t, []string{"Poppins"}, policy.Whitelist)
}

func TestPolicyValidate(t *testing.T) {
	policy := Policy{MaxPerChapter: 2, Whitelist: []string{"Whitelisted"}}

	require.NoError(t, policy.Validate([]string{"New", "new", "WHITELISTED"}))
	err := policy.Validate([]string{"First", "Second", "Third", "Whitelisted"})
	require.True(t, errorx.IsOfType(err, ErrLimitExceeded))
}
