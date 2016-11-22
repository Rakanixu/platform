package engine

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Local struct
type Local struct {
	Endpoint   proto.Endpoint
	DataOrigin string
}

// Validate local datasource (directory exists) and check for intersections between local datasources
func (l *Local) Validate(ctx context.Context, c client.Client, datasources string) (*proto_datasource.Endpoint, error) {
	i := strings.LastIndex(l.Endpoint.Url, "//")

	l.DataOrigin = l.Endpoint.Url[i+1 : len(l.Endpoint.Url)] // Local filesystem path
	if _, err := os.Stat(l.DataOrigin); os.IsNotExist(err) {
		return nil, err
	}

	var endpoints []*proto.Endpoint

	if err := json.Unmarshal([]byte(datasources), &endpoints); err != nil {
		return nil, err
	}

	for _, v := range endpoints {
		if len(v.Url) >= len(l.Endpoint.Url) {
			if strings.Contains(v.Url, l.Endpoint.Url) {
				return nil, errors.New("Datasource trying to create is parent of existing ones. Delete them to create a parent datasource.")
			}
		} else {
			if strings.Contains(l.Endpoint.Url, v.Url) {
				// Datasource tying to create is a child of an existing one
				return nil, errors.New("Datasource trying to create is being covered by an existing one. Kick off scan if data not present.")
			}
		}

	}

	var err error
	l.Endpoint, err = GenerateEndpoint(ctx, c, l.Endpoint)
	if err != nil {
		return nil, err
	}

	return &l.Endpoint, nil
}

// Save local datasource
func (l *Local) Save(ctx context.Context, c client.Client, data interface{}, id string) error {
	return SaveDataSource(ctx, c, data, id)
}

// Delete local data source
func (l *Local) Delete(ctx context.Context, c client.Client) error {
	if err := DeleteDataSource(ctx, c, &l.Endpoint); err != nil {
		return err
	}

	// Specific clean up for local datasources ()
	if strings.Contains(l.Endpoint.Url, localEndpoint) {
		// Remove records from helper index that only belongs to the datasource
		if err := cleanFilesHelperIndex(ctx, c, &l.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

// Scan local data source
func (l *Local) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &l.Endpoint)
}

// ScheduleScan register a chron task
func (l *Local) ScheduleScan(ctx context.Context, c client.Client, sc *scheduler_proto.CreateScheduledTaskRequest) error {
	return ScheduleScanDataSource(ctx, c, sc)
}

// CreateIndeWithAlias creates a index for local datasource
func (l *Local) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &l.Endpoint)
}