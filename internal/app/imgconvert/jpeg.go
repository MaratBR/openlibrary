package imgconvert

import (
	"github.com/h2non/bimg"
)

func ConvertToJPEG(buf []byte) ([]byte, error) {
	img := bimg.NewImage(buf)

	if img.Type() == "jpeg" {
		return buf, nil
	}

	return img.Convert(bimg.JPEG)
}

func Resize(
	buf []byte,
	maxWidth, maxHeight int,
) ([]byte, error) {
	img := bimg.NewImage(buf)
	size, err := img.Size()
	if err != nil {
		return nil, err
	}
	if size.Height < maxHeight && size.Width < maxWidth {
		return buf, nil
	}
	ratio := float64(maxWidth) / float64(size.Width)
	if ratio*float64(size.Height) > float64(maxHeight) {
		ratio = float64(maxHeight) / float64(size.Height)
	}
	return img.Resize(int(float64(size.Width)*ratio), int(float64(size.Height)*ratio))
}
