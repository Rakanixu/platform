package tika

import (
	"log"
	"os"
	"testing"
)

func TestExtractContent(t *testing.T) {
	f, err := os.Open("test_file.docx")
	if err != nil {
		t.Fatal("Error reading test file")
	}

	tika, err := ExtractContent(f)
	if err != nil {
		t.Fatalf("Failed with error: %v", err)
	}

	if tika == nil {
		t.Fatalf("Interface is nil: %v", err)
	}

	log.Println(tika.Content())
}
