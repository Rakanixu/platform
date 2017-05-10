package categories

import (
	"reflect"
	"testing"
)

func TestSetMap(t *testing.T) {
	result := SetMap()

	if nil != result {
		t.Errorf("Expected: %v, got %v", nil, result)
	}

	if len(categoryMap) == 0 {
		t.Error("categoryMap empty")
	}
}

func TestGetMap(t *testing.T) {
	result := GetMap()

	if !reflect.DeepEqual(categoryMap, result) {
		t.Errorf("Expected: %v, got: %v", categoryMap, result)
	}
}

func TestGetDocType(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{
			in:  "",
			out: "None",
		},
		{
			in:  ".pdf",
			out: "Documents",
		},
	}

	for _, tt := range testData {
		result := GetDocType(tt.in)

		if tt.out != result {
			t.Errorf("Expected: %v, got: %v", tt.out, result)
		}
	}
}
