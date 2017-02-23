package rossete

import (
	"testing"
)

func TestEntities(t *testing.T) {
	text := "Bill Murray will appear in new Ghostbusters film: Dr. Peter Venkman was spotted filming a cameo in Boston $12"
	_, err := Entities(text)
	if err != nil {
		t.Errorf("Error retrieveing entities: %v", err)
	}
}
