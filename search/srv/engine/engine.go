package engine

import (
	search "github.com/kazoup/platform/search/srv/proto/search"
)

type Engine interface {
	Init() error
	Search(req *search.SearchRequest) (*search.SearchResponse, error)
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

func Search(req *search.SearchRequest) (*search.SearchResponse, error) {
	return engine.Search(req)
}
