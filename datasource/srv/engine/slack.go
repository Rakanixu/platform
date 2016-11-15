package engine

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"time"
)

// Slack struct
type Slack struct {
	Endpoint proto.Endpoint
}

// Validate slack data source
func (s *Slack) Validate(ctx context.Context, c client.Client, datasources string) (*proto.Endpoint, error) {
	var err error

	s.Endpoint, err = GenerateEndpoint(ctx, c, s.Endpoint)
	if err != nil {
		return nil, err
	}

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

// Scan slack data source
func (s *Slack) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &s.Endpoint)
}

// ScheduleScan register a chron task
func (s *Slack) ScheduleScan(ctx context.Context, c client.Client) error {
	return ScheduleScanDataSource(ctx, c, &scheduler_proto.CreateScheduledTaskRequest{
		Task: &scheduler_proto.Task{
			Id:     s.Endpoint.Id,
			Action: globals.StartScanTask,
		},
		Schedule: &scheduler_proto.Schedule{
			IntervalSeconds: int64(time.Hour.Seconds()),
		},
	})
}

// CreateIndeWithAlias creates a index for local datasource
func (s *Slack) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &s.Endpoint)
}
