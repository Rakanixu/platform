package fs_test

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/fs"
	"log"
	"testing"
)

func TestLocalFsList(t *testing.T) {
	f := fs.NewLocalFsFromEndpoint(&datasource_proto.Endpoint{
		Url: "local:///tmp",
	})

	c, r, err := f.List()
	if err != nil {
		t.Error(err)
	}
	for {
		select {
		case <-r:

			close(c)
			close(r)
			return

		case file := <-c:
			log.Println(file.PreviewURL())
		}
	}
}
