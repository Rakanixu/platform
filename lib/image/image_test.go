package image

import (
	"io"
	"os"
	"testing"
)

func TestThumbnail_JPG(t *testing.T) {
	f, err := os.Open("img_test.jpg")
	if err != nil {
		t.Fatal("Error reading test file")
	}

	r, err := Thumbnail(f, 50)
	if err != nil {
		t.Fatal("Error reading test file")
	}

	of, err := os.Create("img_test_thumb_jpg.png")
	if err != nil {
		t.Fatal("Error creating thumb test file")
	}

	_, err = io.Copy(of, r)
	if err != nil {
		t.Fatal("Error copying thumb test file")
	}
}

func TestThumbnail_PNG(t *testing.T) {
	f, err := os.Open("img_test.png")
	if err != nil {
		t.Fatal("Error reading test file")
	}

	r, err := Thumbnail(f, 50)
	if err != nil {
		t.Fatal("Error reading test file")
	}

	of, err := os.Create("img_test_thumb_png.png")
	if err != nil {
		t.Fatal("Error creating thumb test file")
	}

	_, err = io.Copy(of, r)
	if err != nil {
		t.Fatal("Error copying thumb test file")
	}
}
