package mock

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
)

const (
	MOCK_USER_ID = "test_user"
	MOCK_INDEX   = "test_index"
	MOCK_ID      = "test_index"
)

type mock struct{}

func init() {
	custom.Register(new(mock))
}

func (e *mock) Init() error {
	return nil
}

// ScrollUnprocessedFiles mock
func (e *mock) ScrollUnprocessedFiles(ctx context.Context, req *proto_custom.ScrollUnprocessedFilesRequest) (*proto_custom.ScrollUnprocessedFilesResponse, error) {
	f := &file.KazoupFile{
		Index:  MOCK_INDEX,
		UserId: MOCK_USER_ID,
		ID:     MOCK_ID,
	}

	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return &proto_custom.ScrollUnprocessedFilesResponse{
		Result: `[` + string(b) + `]`,
	}, nil
}

func (e *mock) ScrollDatasources(ctx context.Context, req *proto_custom.ScrollDatasourcesRequest) (*proto_custom.ScrollDatasourcesResponse, error) {
	f := &proto_datasource.Endpoint{
		Index:  MOCK_INDEX,
		UserId: MOCK_USER_ID,
		Url:    globals.Mock,
	}

	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return &proto_custom.ScrollDatasourcesResponse{
		Result: `[` + string(b) + `]`,
	}, nil
}
