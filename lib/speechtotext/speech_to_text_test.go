package speechtotext

import (
	"testing"
)

func TestAsyncContent(t *testing.T) {
	// This file may not be in GCS
	_, err := AsyncContent("gs://kazoup-audio-bucket/2e60baae5504c6ada1f1f7a561758cd0")
	if err != nil {
		if err.Error() != "rpc error: code = NotFound desc = Requested entity was not found." {
			t.Fatalf("Failed with error: %v", err)
		}
	}
}

func TestContent(t *testing.T) {
	/*	f, err := os.Open("test01.flac")
		if err != nil {
			t.Fatal("Error reading test file")
		}

		_, err = Content(f)
		if err != nil {
			t.Fatalf("Failed with error: %v", err)
		}*/

}
