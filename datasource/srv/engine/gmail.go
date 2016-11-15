package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"time"
)

// Gmail struct
type Gmail struct {
	Endpoint datasource_proto.Endpoint
}

// Validate gmail data  source
func (g *Gmail) Validate(ctx context.Context, c client.Client, datasources string) (*datasource_proto.Endpoint, error) {
	var err error

	g.Endpoint, err = GenerateEndpoint(ctx, c, g.Endpoint)
	if err != nil {
		return nil, err
	}

	return &g.Endpoint, nil
}

// Save gmail data source
func (g *Gmail) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete gmail data source
func (g *Gmail) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &g.Endpoint)
}

// Scan gmail data source
func (g *Gmail) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &g.Endpoint)
}

// ScheduleScan register a chron task
func (g *Gmail) ScheduleScan(ctx context.Context, c client.Client) error {
	return ScheduleScanDataSource(ctx, c, &scheduler_proto.CreateScheduledTaskRequest{
		Task: &scheduler_proto.Task{
			Id:     g.Endpoint.Id,
			Action: globals.StartScanTask,
		},
		Schedule: &scheduler_proto.Schedule{
			IntervalSeconds: int64(time.Hour.Seconds()),
		},
	})
}

// CreateIndeWithAlias creates a index for gmail datasource
func (g *Gmail) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &g.Endpoint)
}
