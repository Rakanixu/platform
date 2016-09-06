package googledrive

import (
	"crypto/md5"
	"encoding/hex"
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"strconv"
	"time"
)

// Fake struct
type Googledrive struct {
	Endpoint datasource_proto.Endpoint
	filestorer.FileStore
}

// Validate fake, always fine
func (g *Googledrive) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	g.Endpoint.Index = "index" + strconv.Itoa(int(time.Now().UnixNano()))
	g.Endpoint.Id = getMD5Hash(g.Endpoint.Url)

	return &g.Endpoint, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
