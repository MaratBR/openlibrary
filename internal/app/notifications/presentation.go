package notifications

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
)

const (
	CategorySystem          = "system"
	CategoryModeration      = "moderation"
	CategoryAdmin           = "admin"
	IconImage          byte = 0
	IconFontAwesome    byte = 1
	IconEmoji          byte = 2
)

// Icon keeps the existing image/Font Awesome metadata representation compatible.
type Icon struct {
	Type    byte   `json:"type"`
	Payload string `json:"payload"`
}

type Options struct {
	Category string `json:"category"`
	Urgent   bool   `json:"urgent"`
	Type     string `json:"type"`
	Icon     *Icon  `json:"icon,omitempty"`
}

var inlinePolicy = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("b", "strong", "i", "em", "u", "s", "del", "sub", "sup", "br", "code")
	p.AllowAttrs("href", "title").OnElements("a")
	p.AllowRelativeURLs(true)
	p.AllowURLSchemes("http", "https", "mailto")
	p.RequireParseableURLs(true)
	p.RequireNoFollowOnLinks(true)
	return p
}()

// SanitizeContent is applied to new writes and historical rows before delivery.
func SanitizeContent(content string) string { return inlinePolicy.Sanitize(content) }

var fontAwesomeIcon = regexp.MustCompile(`^fa-(solid|regular|brands) fa-[a-z0-9-]+$`)

func (icon Icon) validate() error {
	if icon.Payload == "" || utf8.RuneCountInString(icon.Payload) > 256 {
		return errors.New("invalid notification icon")
	}
	switch icon.Type {
	case IconEmoji:
		if utf8.RuneCountInString(icon.Payload) > 32 {
			return errors.New("emoji icon is too long")
		}
	case IconFontAwesome:
		if !fontAwesomeIcon.MatchString(icon.Payload) {
			return errors.New("invalid Font Awesome icon")
		}
	case IconImage:
		u, err := url.Parse(icon.Payload)
		if err != nil || u.User != nil || strings.ContainsAny(icon.Payload, "\\\r\n\t") {
			return errors.New("invalid notification image URL")
		}
		absolute := (u.Scheme == "https" || u.Scheme == "http") && u.Host != ""
		relative := u.Scheme == "" && u.Host == "" && strings.HasPrefix(icon.Payload, "/") && !strings.HasPrefix(icon.Payload, "//")
		if !absolute && !relative {
			return errors.New("invalid notification image URL")
		}
	default:
		return errors.New("invalid notification icon type")
	}
	return nil
}

func (o Options) metadata() (map[string]string, error) {
	category := o.Category
	if category == "" {
		category = CategorySystem
	}
	switch category {
	case CategorySystem, CategoryModeration, CategoryAdmin:
	default:
		return nil, errors.New("invalid notification category")
	}
	data := map[string]string{"category": category, "urgent": strconv.FormatBool(o.Urgent)}
	if o.Icon != nil {
		if err := o.Icon.validate(); err != nil {
			return nil, err
		}
		raw, _ := json.Marshal(o.Icon)
		data[dataKeyIcon] = string(raw)
	}
	return data, nil
}

func presentation(data map[string]string) (category string, urgent bool, icon *Icon) {
	category = data["category"]
	if category != CategoryAdmin && category != CategoryModeration {
		category = CategorySystem
	}
	urgent, _ = strconv.ParseBool(data["urgent"])
	var custom Icon
	if json.Unmarshal([]byte(data[dataKeyIcon]), &custom) == nil && custom.validate() == nil {
		return category, urgent, &custom
	}
	switch category {
	case CategoryModeration:
		icon = &Icon{Type: IconFontAwesome, Payload: "fa-solid fa-shield"}
	case CategoryAdmin:
		icon = &Icon{Type: IconFontAwesome, Payload: "fa-solid fa-gauge-high"}
	}
	return
}
