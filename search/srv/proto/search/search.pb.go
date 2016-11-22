// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/search/srv/proto/search/search.proto
// DO NOT EDIT!

package go_micro_srv_search_search

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Search service

type SearchClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	SearchProxy(ctx context.Context, in *SearchProxyRequest, opts ...client.CallOption) (*SearchProxyResponse, error)
	Aggregate(ctx context.Context, in *AggregateRequest, opts ...client.CallOption) (*AggregateResponse, error)
}

type searchClient struct {
	c           client.Client
	serviceName string
}

func NewSearchClient(serviceName string, c client.Client) SearchClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.search.search"
	}
	return &searchClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *searchClient) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Search.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) SearchProxy(ctx context.Context, in *SearchProxyRequest, opts ...client.CallOption) (*SearchProxyResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Search.SearchProxy", in)
	out := new(SearchProxyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) Aggregate(ctx context.Context, in *AggregateRequest, opts ...client.CallOption) (*AggregateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Search.Aggregate", in)
	out := new(AggregateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Search service

type SearchHandler interface {
	Search(context.Context, *SearchRequest, *SearchResponse) error
	SearchProxy(context.Context, *SearchProxyRequest, *SearchProxyResponse) error
	Aggregate(context.Context, *AggregateRequest, *AggregateResponse) error
}

func RegisterSearchHandler(s server.Server, hdlr SearchHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Search{hdlr}, opts...))
}

type Search struct {
	SearchHandler
}

func (h *Search) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.SearchHandler.Search(ctx, in, out)
}

func (h *Search) SearchProxy(ctx context.Context, in *SearchProxyRequest, out *SearchProxyResponse) error {
	return h.SearchHandler.SearchProxy(ctx, in, out)
}

func (h *Search) Aggregate(ctx context.Context, in *AggregateRequest, out *AggregateResponse) error {
	return h.SearchHandler.Aggregate(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/search/srv/proto/search/search.proto", fileDescriptor1)
}

var fileDescriptor1 = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x90, 0x41, 0x6e, 0xc2, 0x30,
	0x10, 0x45, 0xab, 0x2e, 0x22, 0xd5, 0xdd, 0x79, 0x99, 0x65, 0x77, 0xad, 0xda, 0xb1, 0x54, 0x0e,
	0x80, 0x72, 0x03, 0x04, 0x27, 0x70, 0xa2, 0xc1, 0x09, 0x60, 0xc6, 0x78, 0x1c, 0x04, 0xdc, 0x90,
	0x5b, 0x81, 0x6c, 0xe2, 0x1d, 0x10, 0xb1, 0x1a, 0xeb, 0xcf, 0x9b, 0xff, 0x24, 0x8b, 0xa9, 0xe9,
	0x42, 0xdb, 0xd7, 0xd0, 0x90, 0x55, 0x6b, 0x7d, 0xa2, 0xde, 0x29, 0xb7, 0xd1, 0x61, 0x49, 0xde,
	0x2a, 0x46, 0xed, 0x9b, 0x56, 0xb1, 0xdf, 0x2b, 0xe7, 0x29, 0x50, 0x0e, 0xe2, 0x80, 0x98, 0xc9,
	0xd2, 0x10, 0xd8, 0xae, 0xf1, 0x04, 0x57, 0x0e, 0x6e, 0xab, 0x34, 0xca, 0xea, 0x85, 0x72, 0x8b,
	0xcc, 0xda, 0x20, 0xa7, 0xfa, 0xff, 0xf3, 0xbb, 0x28, 0x16, 0x71, 0x23, 0x75, 0x7e, 0x7d, 0xc3,
	0x7d, 0x29, 0x24, 0x66, 0x8e, 0xbb, 0x1e, 0x39, 0x94, 0x3f, 0x63, 0x50, 0x76, 0xb4, 0x65, 0xfc,
	0x7a, 0x93, 0x4e, 0x7c, 0xa6, 0x6c, 0xe6, 0xe9, 0x70, 0x94, 0xf0, 0xfc, 0x38, 0x82, 0x83, 0x4c,
	0x8d, 0xe6, 0xb3, 0x71, 0x25, 0x3e, 0x2a, 0x63, 0x3c, 0x1a, 0x1d, 0x50, 0xfe, 0x3e, 0xba, 0xcf,
	0xd8, 0x60, 0xfb, 0x1b, 0x49, 0x0f, 0xae, 0xba, 0x88, 0x5f, 0x3a, 0xb9, 0x04, 0x00, 0x00, 0xff,
	0xff, 0x37, 0xc9, 0x6a, 0x87, 0xf4, 0x01, 0x00, 0x00,
}