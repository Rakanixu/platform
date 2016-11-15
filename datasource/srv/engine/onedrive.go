package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Onedrive struct
type Onedrive struct {
	Endpoint datasource_proto.Endpoint
}

// Validate
func (o *Onedrive) Validate(ctx context.Context, c client.Client, datasources string) (*datasource_proto.Endpoint, error) {
	var err error

	o.Endpoint, err = GenerateEndpoint(ctx, c, o.Endpoint)
	if err != nil {
		return nil, err
	}

	return &o.Endpoint, nil
}

// Save one drive datasource
func (o *Onedrive) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete one drive data source
func (o *Onedrive) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &o.Endpoint)
}

// Scan one drive data source
func (o *Onedrive) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &o.Endpoint)
}

// ScheduleScan register a chron task
func (o *Onedrive) ScheduleScan(ctx context.Context, c client.Client, sc *scheduler_proto.CreateScheduledTaskRequest) error {
	return ScheduleScanDataSource(ctx, c, sc)
}

// CreateIndeWithAlias creates a index for local datasource
func (o *Onedrive) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &o.Endpoint)
}
