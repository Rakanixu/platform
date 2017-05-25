package tika

import (
	"log"
	"os"
	"testing"
)

var extract = []struct {
	in string
}{
	{
		in: "test_file.docx",
	},
	{
		in: "test_file.xls",
	},
}

func TestExtractContent(t *testing.T) {
	for _, v := range extract {
		f, err := os.Open(v.in)
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
	}
}

func TestExtractPlainContent(t *testing.T) {
	for _, v := range extract {
		f, err := os.Open(v.in)
		if err != nil {
			t.Fatal("Error reading test file")
		}

		tika, err := ExtractPlainContent(f)
		if err != nil {
			t.Fatalf("Failed with error: %v", err)
		}

		if tika == nil {
			t.Fatalf("Interface is nil: %v", err)
		}
	}
}
