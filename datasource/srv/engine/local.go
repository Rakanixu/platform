package engine

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
	"os"
	"strings"
)

// Local struct
type Local struct {
	Endpoint   proto_datasource.Endpoint
	DataOrigin string
}

// Validate local datasource (directory exists) and check for intersections between local datasources
func (l *Local) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	i := strings.LastIndex(l.Endpoint.Url, "//")

	l.DataOrigin = l.Endpoint.Url[i+1 : len(l.Endpoint.Url)] // Local filesystem path
	if _, err := os.Stat(l.DataOrigin); os.IsNotExist(err) {
		return nil, err
	}

	var endpoints []*proto_datasource.Endpoint

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
	l.Endpoint, err = GenerateEndpoint(ctx, l.Endpoint)
	if err != nil {
		return nil, err
	}

	return &l.Endpoint, nil
}

// Save local datasource
func (l *Local) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete local data source
func (l *Local) Delete(ctx context.Context) error {
	if err := DeleteDataSource(ctx, &l.Endpoint); err != nil {
		return err
	}

	// Specific clean up for local datasources ()
	if strings.Contains(l.Endpoint.Url, localEndpoint) {
		// Remove records from helper index that only belongs to the datasource
		if err := cleanFilesHelperIndex(ctx, &l.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

// CreateIndeWithAlias creates a index for local datasource
func (l *Local) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &l.Endpoint)
}
