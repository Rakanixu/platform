package slack

import (
	"crypto/md5"
	"encoding/hex"
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"strconv"
	"time"
)

// Fake struct
type Slack struct {
	filestorer.FileStore
	Endpoint proto.Endpoint
}

// Validate slack
func (s *Slack) Validate(datasources string) (*proto.Endpoint, error) {
	s.Endpoint.Index = "index" + strconv.Itoa(int(time.Now().UnixNano()))
	s.Endpoint.Id = getMD5Hash(s.Endpoint.Url)

	return &s.Endpoint, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
