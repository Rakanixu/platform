package mock

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"github.com/kazoup/platform/lib/file"
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

// ScrollUnprocessedAudio mock
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
