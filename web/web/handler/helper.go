package handler

import (
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}
