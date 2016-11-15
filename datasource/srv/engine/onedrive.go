package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"time"
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
func (o *Onedrive) ScheduleScan(ctx context.Context, c client.Client) error {
	return ScheduleScanDataSource(ctx, c, &scheduler_proto.CreateScheduledTaskRequest{
		Task: &scheduler_proto.Task{
			Id:     o.Endpoint.Id,
			Action: globals.StartScanTask,
		},
		Schedule: &scheduler_proto.Schedule{
			IntervalSeconds: int64(time.Hour.Seconds()),
		},
	})
}

// CreateIndeWithAlias creates a index for local datasource
func (o *Onedrive) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &o.Endpoint)
}
