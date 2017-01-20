package image

// Required the _ imports to be able to decode different formats
import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
)

func Thumbnail(rd io.ReadCloser, width int) (io.Reader, error) {
	defer rd.Close()

	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, resize.Resize(uint(width), 0, img, resize.MitchellNetravali))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
