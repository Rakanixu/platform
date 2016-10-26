package engine

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

// Slack struct
type Slack struct {
	Endpoint proto.Endpoint
}

// Validate slack data source
func (s *Slack) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &s.Endpoint, err
	}
	s.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	s.Endpoint.Id = globals.GetMD5Hash(s.Endpoint.Url + s.Endpoint.UserId)

	return &s.Endpoint, nil
}

// Save slack datasource
func (s *Slack) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete slack data source
func (s *Slack) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &s.Endpoint)
}
