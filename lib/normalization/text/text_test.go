package text

import (
	"testing"
)

func TestNormalizeBytes(t *testing.T) {
	expected := `a b cn
	迪
	d e

	f g`
	str := `a b c\n
	迪
	d e

	f g`

	result, _ := NormalizeBytes([]byte(str))

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestNormalize(t *testing.T) {
	expected := `a b cn
	迪
	d e

	f g`
	str := `a b c\n
	迪
	d e

	f g`

	result, _ := Normalize(str)

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReplaceTabs(t *testing.T) {
	expected := `ab  cb`
	str := `ab		cb`

	result, _ := ReplaceTabs(str)

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReplaceNewLines(t *testing.T) {
	expected := `ab cb`
	str := `ab
cb`

	result, _ := ReplaceNewLines(str)

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReplaceDoubleQuotes(t *testing.T) {
	expected := `'ab' cb`
	str := `"ab" cb`

	result, _ := ReplaceDoubleQuotes(str)

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSkipper_Transform(t *testing.T) {
	// TODO
}

func TestSkipper_Reset(t *testing.T) {
	// TODO
}
