package onedrive

import (
	"crypto/md5"
	"encoding/hex"
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Fake struct
type Onedrive struct {
	Endpoint datasource_proto.Endpoint
	filestorer.FileStore
}

// Validate
func (o *Onedrive) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	if len(o.Endpoint.Index) == 0 {
		s, err := globals.NewUUID()
		if err != nil {
			return &o.Endpoint, err
		}
		o.Endpoint.Index = "index" + strings.Replace(s, "-", "", 1)
	}
	o.Endpoint.Id = getMD5Hash(o.Endpoint.Url + o.Endpoint.UserId)

	return &o.Endpoint, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
