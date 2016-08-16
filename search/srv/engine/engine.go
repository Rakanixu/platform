package engine

import (
	"github.com/kazoup/platform/search/srv/proto/search"
)

type Engine interface {
	Init() error
	Search(req *search.SearchRequest) (res *search.SearchResponse, err error)
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

func Search(req *search.SearchRequest) (res *search.SearchResponse, err error) {
	return engine.Search(req)
}
