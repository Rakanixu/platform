package local

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/datasource/srv/engine"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"os"
	"strings"
)

// Local struct
type Local struct {
	engine.DataSource
	Endpoint   proto.Endpoint
	DataOrigin string
}

// Validate local datasource (directory exists) and check for intersections between local datasources
func (l *Local) Validate(datasources string) (*proto_datasource.Endpoint, error) {
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
	s, err := globals.NewUUID()
	if err != nil {
		return &l.Endpoint, err
	}
	l.Endpoint.Index = "index" + strings.Replace(s, "|", "", 1)
	l.Endpoint.Id = globals.GetMD5Hash(l.Endpoint.Url + l.Endpoint.UserId)

	return &l.Endpoint, nil
}
