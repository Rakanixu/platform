package file

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	"reflect"
	"testing"
)

func TestNewFileFromString(t *testing.T) {

	slackFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.Slack,
	})
	if err != nil {
		t.Fatal(err)
	}

	type out struct {
		file File
		err  error
	}

	testData := []struct {
		in  string
		out out
	}{
		{
			string(slackFileBytes),
			out{
				&KazoupSlackFile{
					KazoupFile{
						FileType: globals.Slack,
					},
					nil,
				},
				nil,
			},
		},
	}

	for _, tt := range testData {
		result, err := NewFileFromString(tt.in)

		if !reflect.DeepEqual(tt.out.file, result) {
			t.Errorf("Expected: %v, got: %v", tt.out.file, result)
		}

		if tt.out.err != err {
			t.Errorf("Expected error: %v, got: %v", tt.out.err, err)
		}
	}
}
