package slack

import (
	"crypto/md5"
	"encoding/hex"
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Fake struct
type Slack struct {
	filestorer.FileStore
	Endpoint proto.Endpoint
}

// Validate slack
func (s *Slack) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &s.Endpoint, err
	}
	s.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	s.Endpoint.Id = getMD5Hash(s.Endpoint.Url + s.Endpoint.UserId)

	return &s.Endpoint, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
