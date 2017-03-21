// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/docenrich/srv/proto/docenrich/docenrich.proto
// DO NOT EDIT!

/*
Package go_micro_srv_docenrich is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/docenrich/srv/proto/docenrich/docenrich.proto

It has these top-level messages:
	CreateRequest
	CreateResponse
	HealthRequest
	HealthResponse
*/
package go_micro_srv_docenrich

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
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CreateRequest struct {
	Type  string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Index string `protobuf:"bytes,2,opt,name=index" json:"index,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CreateResponse struct {
	Info string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *CreateResponse) Reset()                    { *m = CreateResponse{} }
func (m *CreateResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()               {}
func (*CreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type HealthResponse struct {
	Status int64  `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *HealthResponse) Reset()                    { *m = HealthResponse{} }
func (m *HealthResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()               {}
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*CreateRequest)(nil), "go.micro.srv.docenrich.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "go.micro.srv.docenrich.CreateResponse")
	proto.RegisterType((*HealthRequest)(nil), "go.micro.srv.docenrich.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "go.micro.srv.docenrich.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for DocEnrich service

type DocEnrichClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error)
}

type docEnrichClient struct {
	c           client.Client
	serviceName string
}

func NewDocEnrichClient(serviceName string, c client.Client) DocEnrichClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.docenrich"
	}
	return &docEnrichClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *docEnrichClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DocEnrich.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *docEnrichClient) Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DocEnrich.Health", in)
	out := new(HealthResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DocEnrich service

type DocEnrichHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterDocEnrichHandler(s server.Server, hdlr DocEnrichHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&DocEnrich{hdlr}, opts...))
}

type DocEnrich struct {
	DocEnrichHandler
}

func (h *DocEnrich) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.DocEnrichHandler.Create(ctx, in, out)
}

func (h *DocEnrich) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.DocEnrichHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/docenrich/srv/proto/docenrich/docenrich.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x91, 0x3f, 0x4f, 0xf3, 0x30,
	0x10, 0xc6, 0x9b, 0xf4, 0x7d, 0x23, 0xf5, 0xa4, 0x66, 0xb0, 0x50, 0x15, 0x31, 0x21, 0xf3, 0x47,
	0x4c, 0xb6, 0x04, 0x2b, 0x1b, 0x20, 0x95, 0x35, 0x1b, 0xa3, 0xeb, 0xb8, 0x8d, 0x45, 0x93, 0x33,
	0xb6, 0x83, 0x80, 0x6f, 0xc7, 0x37, 0xa3, 0x71, 0x9a, 0x40, 0x24, 0x0a, 0xdb, 0xe3, 0xbb, 0x47,
	0xbf, 0xbb, 0xe7, 0x0c, 0xcb, 0x8d, 0xf6, 0x65, 0xb3, 0x62, 0x12, 0x2b, 0xfe, 0x24, 0xde, 0xb1,
	0x31, 0xdc, 0x6c, 0x85, 0x5f, 0xa3, 0xad, 0x78, 0x81, 0x52, 0xd5, 0x56, 0xcb, 0x92, 0x3b, 0xfb,
	0xc2, 0x8d, 0x45, 0x8f, 0xdf, 0x6a, 0x83, 0x62, 0xa1, 0x43, 0x16, 0x1b, 0x64, 0x95, 0x96, 0x16,
	0xd9, 0xce, 0xcd, 0x86, 0x2e, 0x7d, 0x80, 0xf9, 0xad, 0x55, 0xc2, 0xab, 0x5c, 0x3d, 0x37, 0xca,
	0x79, 0x42, 0xe0, 0x9f, 0x7f, 0x33, 0x2a, 0x8b, 0x4e, 0xa2, 0xcb, 0x59, 0x1e, 0x34, 0x39, 0x82,
	0xff, 0xba, 0x2e, 0xd4, 0x6b, 0x16, 0x87, 0x62, 0xf7, 0x20, 0x29, 0xc4, 0xba, 0xc8, 0xa6, 0xa1,
	0xb4, 0x53, 0xf4, 0x0c, 0xd2, 0x1e, 0xe5, 0x0c, 0xd6, 0x4e, 0xb5, 0x2c, 0x5d, 0xaf, 0xb1, 0x67,
	0xb5, 0x9a, 0x9e, 0xc2, 0x7c, 0xa9, 0xc4, 0xd6, 0x97, 0xbf, 0x0c, 0xa4, 0x37, 0x90, 0xf6, 0xa6,
	0x3d, 0x6a, 0x01, 0x89, 0xf3, 0xc2, 0x37, 0x2e, 0xf8, 0xa6, 0xf9, 0xfe, 0x35, 0x8c, 0x88, 0xbf,
	0x46, 0x5c, 0x7d, 0x44, 0x30, 0xbb, 0x43, 0x79, 0x1f, 0x12, 0x92, 0x47, 0x48, 0xba, 0xb5, 0xc8,
	0x39, 0xfb, 0xf9, 0x08, 0x6c, 0x74, 0x81, 0xe3, 0x8b, 0xbf, 0x6c, 0xdd, 0x4a, 0x74, 0xd2, 0xa2,
	0xbb, 0x35, 0x0f, 0xa3, 0x47, 0x59, 0x0f, 0xa3, 0xc7, 0x69, 0xe9, 0x64, 0x95, 0x84, 0x6f, 0xbb,
	0xfe, 0x0c, 0x00, 0x00, 0xff, 0xff, 0xef, 0x4c, 0xb5, 0x00, 0x02, 0x02, 0x00, 0x00,
}