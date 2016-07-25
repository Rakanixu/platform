package handler

import (
	crawler "github.com/kazoup/platform/crawler/srv/proto/crawler"
	scanner "github.com/kazoup/platform/crawler/srv/scan"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"strconv"
)

// Crawl ...
type Crawl struct{}

var crawls = make(map[int64]scanner.Scanner)

// Start ...
func (c *Crawl) Start(ctx context.Context, req *crawler.StartRequest, rsp *crawler.StartResponse) error {
	l := int64(len(crawls)) + 1

	// mapScanner ensures data validation
	s, err := mapScanner(l, &datasource.Endpoint{
		Url: "local://",
	})
	if err != nil {
		return errors.InternalServerError("go.micro.srv.crawler.Crawl.Start", err.Error())
	}

	crawls[l] = s
	s.Start()

	return nil
}

// Stop ...
func (c *Crawl) Stop(ctx context.Context, req *crawler.StopRequest, rsp *crawler.StopResponse) error {
	scan, ok := crawls[int64(req.Id)]

	if !ok {
		return errors.BadRequest("go.micro.srv.crawler.Crawl.Stop", "Crawler not found")
	}

	scan.Stop()
	delete(crawls, int64(req.Id))

	return nil
}

// Search ...
func (c *Crawl) Search(ctx context.Context, req *crawler.SearchRequest, rsp *crawler.SearchResponse) error {
	r := make(map[string]*crawler.Status)

	for k, v := range crawls {
		inf, err := v.Info()
		if err != nil {
			return errors.InternalServerError("go.micro.srv.crawler.Crawl.Search", err.Error())
		}

		r[strconv.FormatInt(k, 10)] = &crawler.Status{
			Id:          inf.Id,
			Type:        inf.Type,
			Description: inf.Description,
			Config:      inf.Config,
		}
	}

	rsp.Crawls = r

	return nil
}

func Subscriber(ctx context.Context, endpoint *datasource.Endpoint) error {
	l := int64(len(crawls)) + 1

	s, err := mapScanner(l, endpoint)
	if err != nil {
		return err
	}

	crawls[l] = s
	s.Start()

	return nil
}
