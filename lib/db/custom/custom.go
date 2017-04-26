package custom

import (
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"golang.org/x/net/context"
)

type DBCustom interface {
	Init() error
	Custom
}

type Custom interface {
	ScrollUnprocessedFiles(ctx context.Context, req *proto_custom.ScrollUnprocessedFilesRequest) (*proto_custom.ScrollUnprocessedFilesResponse, error)
}

var (
	custom DBCustom
)

func Register(storage DBCustom) {
	custom = storage
}

func Init() error {
	return custom.Init()
}

func ScrollUnprocessedFiles(ctx context.Context, req *proto_custom.ScrollUnprocessedFilesRequest) (*proto_custom.ScrollUnprocessedFilesResponse, error) {
	return custom.ScrollUnprocessedFiles(ctx, req)
}
