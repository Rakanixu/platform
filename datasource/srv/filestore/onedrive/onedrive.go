package onedrive

import (
	"crypto/md5"
	"encoding/hex"
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"strconv"
	"time"
)

// Fake struct
type Onedrive struct {
	Endpoint datasource_proto.Endpoint
	filestorer.FileStore
}

// Validate fake, always fine
func (o *Onedrive) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	o.Endpoint.Index = "index" + strconv.Itoa(int(time.Now().UnixNano()))
	o.Endpoint.Id = getMD5Hash(o.Endpoint.Url)

	return &o.Endpoint, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
