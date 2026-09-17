package htmlbag

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/MaratBR/openlibrary/internal/olhttp/webcomponents"
)

type ImageType byte

const (
	ImageTypeURL ImageType = iota
	ImageTypeFontAwesome
)

type Image struct {
	Type    ImageType `json:"type"`
	Payload string    `json:"payload"`
}

func NewURLImage(url string) Image {
	return Image{Type: ImageTypeURL, Payload: url}
}

func NewFontAwesomeImage(fontAwesomeIcon string) Image {
	return Image{Type: ImageTypeFontAwesome, Payload: fontAwesomeIcon}
}

func (img Image) Render(ctx context.Context, w io.Writer) error {
	switch img.Type {
	case ImageTypeURL:
		return webcomponents.Image(img.Payload).Render(ctx, w)
	case ImageTypeFontAwesome:
		return webcomponents.ImageFontAwesome(FontAwesomeIconType(img.Payload).ToClass()).Render(ctx, w)
	default:
		return nil
	}
}

func (img Image) MarshalText() string {
	b, err := json.Marshal(img)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FontAwesomeIconType string

func (faic FontAwesomeIconType) ToClass() string {
	terms := strings.Split(string(faic), ",")

	for i := range terms {
		term := terms[i]
		if !strings.HasPrefix(term, "fa-") {
			term = "fa-" + term
		}

	}
	return strings.Join(terms, " ")
}
