package handler

import (
	"go_appengine/goroot/src/appengine_internal/search"

	"golang.org/x/net/context"
)

type SearchBleve struct {
}

func (bs *BleveSearch) Search(context.Context, *search.SearchRequest, *search.SearchResponse) error {
	panic("not implemented")
}
