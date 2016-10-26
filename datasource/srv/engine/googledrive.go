package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

// Googledrive struct
type Googledrive struct {
	Endpoint datasource_proto.Endpoint
}

// Validate google drive data source
func (g *Googledrive) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	s, err := globals.NewUUID()
	if err != nil {
		return &g.Endpoint, err
	}
	g.Endpoint.Index = "index" + strings.Replace(s, "-", "", 1)
	g.Endpoint.Id = globals.GetMD5Hash(g.Endpoint.Url + g.Endpoint.UserId)

	return &g.Endpoint, nil
}

// Save google drive data source
func (g *Googledrive) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete google drive data source
func (g *Googledrive) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &g.Endpoint)
}

// Scan google drive data source
func (g *Googledrive) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &g.Endpoint)
}

// CreateIndeWithAlias creates a index for google drive datasource
func (g *Googledrive) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &g.Endpoint)
}
