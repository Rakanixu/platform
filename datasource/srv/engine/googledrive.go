package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Googledrive struct
type Googledrive struct {
	Endpoint datasource_proto.Endpoint
}

// Validate google drive data source
func (g *Googledrive) Validate(ctx context.Context, c client.Client, datasources string) (*datasource_proto.Endpoint, error) {
	var err error

	g.Endpoint, err = GenerateEndpoint(ctx, c, g.Endpoint)
	if err != nil {
		return nil, err
	}

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

// ScheduleScan register a chron task
func (g *Googledrive) ScheduleScan(ctx context.Context, c client.Client, sc *scheduler_proto.CreateScheduledTaskRequest) error {
	return ScheduleScanDataSource(ctx, c, sc)
}

// CreateIndeWithAlias creates a index for google drive datasource
func (g *Googledrive) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &g.Endpoint)
}
