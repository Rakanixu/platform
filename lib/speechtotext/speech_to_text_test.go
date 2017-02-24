package speechtotext

import (
	"os"
	"testing"
)

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

func TestWav(t *testing.T) {
	f, err := os.Open("test01.wav")
	if err != nil {
		t.Fatal("Error reading test file")
	}

	Wav(f)
}
