package operations

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/file"
	"golang.org/x/net/context"
)

type DBOperations interface {
	Init() error
	Operations
}

type Operations interface {
	Create(ctx context.Context, req *proto_operations.CreateRequest) (*proto_operations.CreateResponse, error)
	Read(ctx context.Context, req *proto_operations.ReadRequest) (*proto_operations.ReadResponse, error)
	Update(ctx context.Context, req *proto_operations.UpdateRequest) (*proto_operations.UpdateResponse, error)
	Delete(ctx context.Context, req *proto_operations.DeleteRequest) (*proto_operations.DeleteResponse, error)
	DeleteByQuery(ctx context.Context, req *proto_operations.DeleteByQueryRequest) (*proto_operations.DeleteByQueryResponse, error)
	Search(ctx context.Context, req *proto_operations.SearchRequest) (*proto_operations.SearchResponse, error)
	SearchById(ctx context.Context, req *proto_operations.SearchByIdRequest) (*proto_operations.SearchByIdResponse, error)
}

var (
	operations DBOperations
)

const (
	File       = "file"
	Datasource = "datasource"
)

func Register(storage DBOperations) {
	operations = storage
}

func Init() error {
	return operations.Init()
}

func Create(ctx context.Context, req *proto_operations.CreateRequest) (*proto_operations.CreateResponse, error) {
	return operations.Create(ctx, req)
}

func Read(ctx context.Context, req *proto_operations.ReadRequest) (*proto_operations.ReadResponse, error) {
	return operations.Read(ctx, req)
}

func Update(ctx context.Context, req *proto_operations.UpdateRequest) (*proto_operations.UpdateResponse, error) {
	return operations.Update(ctx, req)
}

func Delete(ctx context.Context, req *proto_operations.DeleteRequest) (*proto_operations.DeleteResponse, error) {
	return operations.Delete(ctx, req)
}

func DeleteByQuery(ctx context.Context, req *proto_operations.DeleteByQueryRequest) (*proto_operations.DeleteByQueryResponse, error) {
	return operations.DeleteByQuery(ctx, req)
}

func Search(ctx context.Context, req *proto_operations.SearchRequest) (*proto_operations.SearchResponse, error) {
	return operations.Search(ctx, req)
}

func SearchById(ctx context.Context, req *proto_operations.SearchByIdRequest) (*proto_operations.SearchByIdResponse, error) {
	return operations.SearchById(ctx, req)
}

func TypeFactory(typ string, data string) (interface{}, error) {
	switch typ {
	case File:
		file, err := file.NewFileFromString(data)
		if err != nil {
			return nil, err
		}
		return file, nil
	case Datasource:
		var e *datasource_proto.Endpoint
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}

		return e, nil
	}

	return nil, nil
}
