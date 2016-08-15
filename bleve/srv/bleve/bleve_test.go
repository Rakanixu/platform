package bleve_test

import (
	"os"
	"testing"

	"github.com/kazoup/platform/bleve/srv/bleve"
)

func TestBleveInit(t *testing.T) {
	// Init bleve index
	idx := bleve.Init()

	if idx == nil {
		t.Errorf("No index")
	}
	if _, err := os.Stat(bleve.IndexPath); err != nil {
		t.Errorf("Can't create index : %s", err)
	}
	//Clean up
	if err := os.RemoveAll(bleve.IndexPath); err != nil {
		t.Error(err)
	}
}
