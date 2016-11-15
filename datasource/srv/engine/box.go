package engine

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"time"
)

// Box struct
type Box struct {
	Endpoint proto.Endpoint
}

// Validate box datasource
func (b *Box) Validate(ctx context.Context, c client.Client, datasources string) (*proto.Endpoint, error) {
	var err error

	b.Endpoint, err = GenerateEndpoint(ctx, c, b.Endpoint)
	if err != nil {
		return nil, err
	}

	return &b.Endpoint, nil
}

// Save box data source
func (b *Box) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete box data source
func (b *Box) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &b.Endpoint)
}

// Scan box data source
func (b *Box) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &b.Endpoint)
}

// ScheduleScan register a chron task
func (b *Box) ScheduleScan(ctx context.Context, c client.Client) error {
	return ScheduleScanDataSource(ctx, c, &scheduler_proto.CreateScheduledTaskRequest{
		Task: &scheduler_proto.Task{
			Id:     b.Endpoint.Id,
			Action: globals.StartScanTask,
		},
		Schedule: &scheduler_proto.Schedule{
			IntervalSeconds: int64(time.Hour.Seconds()),
		},
	})
}

// CreateIndeWithAlias creates a index for box datasource
func (b *Box) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &b.Endpoint)
}
