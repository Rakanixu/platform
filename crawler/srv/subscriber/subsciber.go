package subscriber

import (
	"github.com/kazoup/platform/crawler/srv/handler"
	"golang.org/x/net/context"

	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

func Scans(ctx context.Context, endpoint *datasource.Endpoint) error {
	l := int64(len(handler.Crawls)) + 1

	s, err := handler.MapScanner(l, endpoint)
	if err != nil {
		return err
	}

	handler.Crawls[l] = s
	s.Start()

	return nil
}
