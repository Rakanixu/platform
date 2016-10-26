package engine

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

// Dropbox struct
type Dropbox struct {
	Endpoint proto.Endpoint
}

// Validate dropbox datasource
func (s *Dropbox) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &s.Endpoint, err
	}
	s.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	s.Endpoint.Id = globals.GetMD5Hash(s.Endpoint.Url + s.Endpoint.UserId)

	return &s.Endpoint, nil
}

// Save dropbox data source
func (s *Dropbox) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete dropbox data source
func (s *Dropbox) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &s.Endpoint)
}

// Scan dropbox data source
func (s *Dropbox) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &s.Endpoint)
}
