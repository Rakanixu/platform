package engine

import (
	search "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type Engine interface {
	Init() error
	Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error)
	SearchProxy(ctx context.Context, req *search.SearchProxyRequest, client client.Client) (*search.SearchProxyResponse, error)
	Aggregate(ctx context.Context, req *search.AggregateRequest, client client.Client) (*search.AggregateResponse, error)
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

func Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error) {
	return engine.Search(ctx, req, client)
}

func SearchProxy(ctx context.Context, req *search.SearchProxyRequest, client client.Client) (*search.SearchProxyResponse, error) {
	return engine.SearchProxy(ctx, req, client)
}

func Aggregate(ctx context.Context, req *search.AggregateRequest, client client.Client) (*search.AggregateResponse, error) {
	return engine.Aggregate(ctx, req, client)
}