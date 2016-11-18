package image

// Required the _ imports to be able to decode different formats
import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
)

func Thumbnail(file []byte, width int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	ni := resize.Resize(uint(width), 0, img, resize.MitchellNetravali)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, ni, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
