package handler

import (
	"log"
	"strconv"
	"time"

	crawler "github.com/kazoup/platform/crawler/srv/proto/crawler"
	scanner "github.com/kazoup/platform/crawler/srv/scan"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Crawl ...
type Crawl struct{}

var Crawls = make(map[int64]scanner.Scanner)

// Start ... Just for debugging sync, all start / stop to a crawler must be done async with messages
func (c *Crawl) Start(ctx context.Context, req *crawler.StartRequest, rsp *crawler.StartResponse) error {
	l := int64(len(Crawls)) + 1

	// mapScanner ensures data validation
	s, err := MapScanner(l, &datasource.Endpoint{
		Url: req.Url,
	})
	if err != nil {
		log.Print("got error")
		return errors.InternalServerError("go.micro.srv.crawler.Crawl.Start", err.Error())
	}

	Crawls[l] = s
	s.Start(Crawls, l)

	return nil
}

// Stop ... Just for debugging sync, all start / stop to a crawler must be done async with messages
func (c *Crawl) Stop(ctx context.Context, req *crawler.StopRequest, rsp *crawler.StopResponse) error {
	scan, ok := Crawls[int64(req.Id)]

	if !ok {
		return errors.BadRequest("go.micro.srv.crawler.Crawl.Stop", "Crawler not found")
	}

	scan.Stop()
	delete(Crawls, int64(req.Id))

	return nil
}

// Search ...
func (c *Crawl) Search(ctx context.Context, req *crawler.SearchRequest, rsp *crawler.SearchResponse) error {
	r := make(map[string]*crawler.Status)

	for k, v := range Crawls {
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

// Status srv handler
func (c *Crawl) Status(ctx context.Context, req *crawler.StatusRequest, stream crawler.Crawl_StatusStream) error {
	log.Printf("Recived Crawl.Status req: %", req)
	for {
		if err := stream.Send(&crawler.StatusResponse{
			Counter: int64(len(Crawls)),
		}); err != nil {
			log.Printf("Error Crawl.Status %s", err)
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}
