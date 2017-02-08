package engine

import (
	search "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type Engine interface {
	Init() error
	Read(ctx context.Context, req *search.ReadRequest, client client.Client) (*search.ReadResponse, error)
	Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error)
	SearchById(ctx context.Context, req *search.SearchByIdRequest, client client.Client) (*search.SearchByIdResponse, error)
}

var (
	engine Engine
)

func Register(backend Engine) {
	engine = backend
}

func Init() error {
	return engine.Init()
}

func Read(ctx context.Context, req *search.ReadRequest, client client.Client) (*search.ReadResponse, error) {
	return engine.Read(ctx, req, client)
}

func Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error) {
	return engine.Search(ctx, req, client)
}

func SearchById(ctx context.Context, req *search.SearchByIdRequest, client client.Client) (*search.SearchByIdResponse, error) {
	return engine.SearchById(ctx, req, client)
}
