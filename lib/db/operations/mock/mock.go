package mock

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/file"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
)

type mock struct{}

func init() {
	operations.Register(new(mock))
}

func (e *mock) Init() error {
	return nil
}

// Create record
func (e *mock) Create(ctx context.Context, req *proto_operations.CreateRequest) (*proto_operations.CreateResponse, error) {
	return &proto_operations.CreateResponse{}, nil
}

// Read record
func (e *mock) Read(ctx context.Context, req *proto_operations.ReadRequest) (*proto_operations.ReadResponse, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, errors.New("Invalid context")
	}

	if len(md["Wanted-Type"]) == 0 {
		return nil, errors.New("Wanted-Type not set in context")
	}

	if md["Wanted-Type"] == "file" {
		f := &file.KazoupFile{}

		b, err := json.Marshal(f)
		if err != nil {
			return nil, err
		}

		return &proto_operations.ReadResponse{
			Result: string(b),
		}, nil
	}

	if md["Wanted-Type"] == "datasource" {
		e := &proto_datasource.Endpoint{}

		b, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}

		return &proto_operations.ReadResponse{
			Result: string(b),
		}, nil
	}

	return &proto_operations.ReadResponse{}, nil
}

func (e *mock) Update(ctx context.Context, req *proto_operations.UpdateRequest) (*proto_operations.UpdateResponse, error) {
	return &proto_operations.UpdateResponse{}, nil
}

func (e *mock) Delete(ctx context.Context, req *proto_operations.DeleteRequest) (*proto_operations.DeleteResponse, error) {
	return &proto_operations.DeleteResponse{}, nil
}

func (e *mock) DeleteByQuery(ctx context.Context, req *proto_operations.DeleteByQueryRequest) (*proto_operations.DeleteByQueryResponse, error) {
	return &proto_operations.DeleteByQueryResponse{}, nil
}

func (e *mock) Search(ctx context.Context, req *proto_operations.SearchRequest) (*proto_operations.SearchResponse, error) {
	return &proto_operations.SearchResponse{
	//Result: rstr,
	//Info:   info.String(),
	}, nil
}
