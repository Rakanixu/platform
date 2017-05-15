package validate

import (
	"github.com/kazoup/platform/lib/errors"
	"testing"
)

func TestExists(t *testing.T) {
	testData := []struct {
		in  string
		out error
	}{
		{
			"",
			errors.ErrMissingParams,
		},
		{
			"asdf",
			nil,
		},
	}

	for _, tt := range testData {
		result := Exists(tt.in)

		if tt.out != result {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}
