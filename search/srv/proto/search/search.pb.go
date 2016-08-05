// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/search/srv/proto/search/search.proto
// DO NOT EDIT!

/*
Package search is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/search/srv/proto/search/search.proto

It has these top-level messages:
	SearchRequest
	SearchResponse
*/
package search

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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type SearchRequest struct {
	Term     string `protobuf:"bytes,1,opt,name=term" json:"term,omitempty"`
	From     int64  `protobuf:"varint,2,opt,name=from" json:"from,omitempty"`
	Size     int64  `protobuf:"varint,3,opt,name=size" json:"size,omitempty"`
	Category string `protobuf:"bytes,4,opt,name=category" json:"category,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*SearchRequest)(nil), "SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "SearchResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Search service

type SearchClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
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
		serviceName = "search"
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

// Server API for Search service

type SearchHandler interface {
	Search(context.Context, *SearchRequest, *SearchResponse) error
}

func RegisterSearchHandler(s server.Server, hdlr SearchHandler) {
	s.Handle(s.NewHandler(&Search{hdlr}))
}

type Search struct {
	SearchHandler
}

func (h *Search) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.SearchHandler.Search(ctx, in, out)
}

var fileDescriptor0 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x8f, 0xcd, 0x4a, 0xc5, 0x30,
	0x10, 0x85, 0xad, 0x2d, 0x45, 0x03, 0x56, 0xc8, 0x42, 0x42, 0x57, 0xd2, 0x55, 0x41, 0x48, 0x40,
	0x71, 0xed, 0x3b, 0xd4, 0x27, 0x48, 0xcb, 0xf4, 0x07, 0x1b, 0x27, 0x4e, 0x12, 0xc1, 0x3e, 0xbd,
	0xb7, 0x69, 0x6f, 0xa1, 0xab, 0x7c, 0xe7, 0x23, 0xc9, 0x99, 0x61, 0x1f, 0xc3, 0xe4, 0xc7, 0xd0,
	0xca, 0x0e, 0x8d, 0xfa, 0xd2, 0x0b, 0x06, 0xab, 0xec, 0xac, 0x7d, 0x8f, 0x64, 0x94, 0x03, 0x4d,
	0xdd, 0xa8, 0x1c, 0xfd, 0x2a, 0x4b, 0xe8, 0xf1, 0x10, 0xf1, 0x90, 0xd1, 0x55, 0x03, 0x7b, 0xf8,
	0x8c, 0xb9, 0x81, 0x9f, 0x00, 0xce, 0x73, 0xce, 0x32, 0x0f, 0x64, 0x44, 0xf2, 0x9c, 0xd4, 0xf7,
	0x4d, 0xe4, 0xd5, 0xf5, 0x84, 0x46, 0xdc, 0x5e, 0x5c, 0xda, 0x44, 0x5e, 0x9d, 0x9b, 0x16, 0x10,
	0xe9, 0xe6, 0x56, 0xe6, 0x25, 0xbb, 0xeb, 0xb4, 0x87, 0x01, 0xe9, 0x4f, 0x64, 0xf1, 0xfd, 0x91,
	0xab, 0x9a, 0x15, 0xd7, 0x22, 0x67, 0xf1, 0xdb, 0x01, 0x7f, 0x62, 0x39, 0x81, 0x0b, 0xb3, 0xdf,
	0xbb, 0xf6, 0xf4, 0xfa, 0xce, 0xf2, 0xed, 0x26, 0x7f, 0x39, 0xa8, 0x90, 0xa7, 0x29, 0xcb, 0x47,
	0x79, 0xfe, 0xac, 0xba, 0x69, 0xf3, 0xb8, 0xd0, 0xdb, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x89,
	0x9b, 0xaa, 0x40, 0x13, 0x01, 0x00, 0x00,
}
