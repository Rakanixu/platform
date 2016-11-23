package handler

import (
	proto "github.com/kazoup/platform/crawler/srv/proto/crawler"

	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Crawler struct
type Crawler struct {
	Client client.Client
}

// Start handler, retrieve kazoup appliance status
func (c *Crawler) Start(ctx context.Context, req *proto.StartRequest, rsp *proto.StartResponse) error {

	return nil
}

// Stop handler, retrieve kazoup appliance status
func (c *Crawler) Stop(ctx context.Context, req *proto.StopRequest, rsp *proto.StopResponse) error {

	return nil
}

// Status handler, retrieve kazoup appliance status
func (c *Crawler) Status(ctx context.Context, req *proto.StatusRequest, rsp proto.Crawl_StatusStream) error {

	return nil
}

// Search handler, retrieve kazoup appliance status
func (c *Crawler) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {

	return nil
}
