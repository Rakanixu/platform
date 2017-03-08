// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/imgenrich/srv/proto/imgenrich/imgenrich.proto
// DO NOT EDIT!

/*
Package go_micro_srv_imgenrich is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/imgenrich/srv/proto/imgenrich/imgenrich.proto

It has these top-level messages:
	CreateRequest
	CreateResponse
	HealthRequest
	HealthResponse
*/
package go_micro_srv_imgenrich

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
	proto.RegisterType((*CreateRequest)(nil), "go.micro.srv.imgenrich.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "go.micro.srv.imgenrich.CreateResponse")
	proto.RegisterType((*HealthRequest)(nil), "go.micro.srv.imgenrich.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "go.micro.srv.imgenrich.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ImgEnrich service

type ImgEnrichClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error)
}

type imgEnrichClient struct {
	c           client.Client
	serviceName string
}

func NewImgEnrichClient(serviceName string, c client.Client) ImgEnrichClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.imgenrich"
	}
	return &imgEnrichClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *imgEnrichClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ImgEnrich.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imgEnrichClient) Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ImgEnrich.Health", in)
	out := new(HealthResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ImgEnrich service

type ImgEnrichHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterImgEnrichHandler(s server.Server, hdlr ImgEnrichHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&ImgEnrich{hdlr}, opts...))
}

type ImgEnrich struct {
	ImgEnrichHandler
}

func (h *ImgEnrich) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.ImgEnrichHandler.Create(ctx, in, out)
}

func (h *ImgEnrich) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.ImgEnrichHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/imgenrich/srv/proto/imgenrich/imgenrich.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x91, 0x4f, 0x4b, 0x03, 0x31,
	0x10, 0xc5, 0xfb, 0x47, 0x17, 0x3a, 0xd0, 0x45, 0x82, 0x94, 0xc5, 0x93, 0x44, 0x14, 0x4f, 0x09,
	0xe8, 0xd5, 0x9b, 0x08, 0xed, 0x75, 0x6f, 0x1e, 0xd3, 0x6d, 0xba, 0x1b, 0x6c, 0x76, 0x62, 0x32,
	0x2b, 0xea, 0xb7, 0xf3, 0x9b, 0xd9, 0x66, 0xb7, 0xc5, 0x05, 0xab, 0xb7, 0x97, 0x99, 0x37, 0x8f,
	0xf7, 0x23, 0x30, 0x2f, 0x0d, 0x55, 0xcd, 0x52, 0x14, 0x68, 0xe5, 0x8b, 0xfa, 0xc4, 0xc6, 0x49,
	0xb7, 0x51, 0xb4, 0x46, 0x6f, 0xa5, 0xb1, 0xa5, 0xae, 0xbd, 0x29, 0x2a, 0x19, 0xfc, 0x9b, 0x74,
	0x1e, 0x09, 0x7f, 0xcc, 0x0e, 0x4a, 0xc4, 0x0d, 0x9b, 0x95, 0x28, 0xac, 0x29, 0x3c, 0x8a, 0xad,
	0x5b, 0x1c, 0xb6, 0x7c, 0x01, 0xd3, 0x47, 0xaf, 0x15, 0xe9, 0x5c, 0xbf, 0x36, 0x3a, 0x10, 0x63,
	0x70, 0x42, 0x1f, 0x4e, 0x67, 0xc3, 0xcb, 0xe1, 0xed, 0x24, 0x8f, 0x9a, 0x9d, 0xc3, 0xa9, 0xa9,
	0x57, 0xfa, 0x3d, 0x1b, 0xc5, 0x61, 0xfb, 0x60, 0x29, 0x8c, 0xcc, 0x2a, 0x1b, 0xc7, 0xd1, 0x56,
	0xf1, 0x33, 0x48, 0xf7, 0x51, 0xc1, 0x61, 0x1d, 0x34, 0xbf, 0x82, 0xe9, 0x5c, 0xab, 0x0d, 0x55,
	0x7f, 0x84, 0xf3, 0x07, 0x48, 0xf7, 0xa6, 0xf6, 0x8c, 0xcd, 0x20, 0x09, 0xa4, 0xa8, 0x09, 0xd1,
	0x37, 0xce, 0xbb, 0xd7, 0xee, 0xda, 0xd4, 0x6b, 0xec, 0x5a, 0x44, 0x7d, 0xf7, 0x35, 0x84, 0xc9,
	0xc2, 0x96, 0x4f, 0x91, 0x86, 0x3d, 0x43, 0xd2, 0x56, 0x60, 0xd7, 0xe2, 0x77, 0x60, 0xd1, 0xa3,
	0xbd, 0xb8, 0xf9, 0xcf, 0xd6, 0x91, 0x0c, 0x76, 0xd1, 0x6d, 0xcd, 0xe3, 0xd1, 0x3d, 0xd6, 0xe3,
	0xd1, 0x7d, 0x5a, 0x3e, 0x58, 0x26, 0xf1, 0x8b, 0xee, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x91,
	0x4f, 0xc1, 0x14, 0xee, 0x01, 0x00, 0x00,
}
