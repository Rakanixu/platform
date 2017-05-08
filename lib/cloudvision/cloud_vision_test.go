package cloudvision

import (
	"os"
	"testing"
)

func TestTag(t *testing.T) {
	file, err := os.Open("img_test.jpg")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	_, err = Tag(file)
	if err != nil {
		t.Fatalf("Failed to tag image: %v", err)
	}
}
