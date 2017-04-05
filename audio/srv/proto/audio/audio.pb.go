// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/audio/srv/proto/audio/audio.proto
// DO NOT EDIT!

/*
Package proto_audio is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/audio/srv/proto/audio/audio.proto

It has these top-level messages:
	EnrichFileRequest
	EnrichFileResponse
	EnrichDatasourceRequest
	EnrichDatasourceResponse
	HealthRequest
	HealthResponse
*/
package proto_audio

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

type EnrichFileRequest struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Index string `protobuf:"bytes,2,opt,name=index" json:"index,omitempty"`
}

func (m *EnrichFileRequest) Reset()                    { *m = EnrichFileRequest{} }
func (m *EnrichFileRequest) String() string            { return proto.CompactTextString(m) }
func (*EnrichFileRequest) ProtoMessage()               {}
func (*EnrichFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type EnrichFileResponse struct {
	Info string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *EnrichFileResponse) Reset()                    { *m = EnrichFileResponse{} }
func (m *EnrichFileResponse) String() string            { return proto.CompactTextString(m) }
func (*EnrichFileResponse) ProtoMessage()               {}
func (*EnrichFileResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type EnrichDatasourceRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *EnrichDatasourceRequest) Reset()                    { *m = EnrichDatasourceRequest{} }
func (m *EnrichDatasourceRequest) String() string            { return proto.CompactTextString(m) }
func (*EnrichDatasourceRequest) ProtoMessage()               {}
func (*EnrichDatasourceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type EnrichDatasourceResponse struct {
	Info string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *EnrichDatasourceResponse) Reset()                    { *m = EnrichDatasourceResponse{} }
func (m *EnrichDatasourceResponse) String() string            { return proto.CompactTextString(m) }
func (*EnrichDatasourceResponse) ProtoMessage()               {}
func (*EnrichDatasourceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type HealthResponse struct {
	Status int64  `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *HealthResponse) Reset()                    { *m = HealthResponse{} }
func (m *HealthResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()               {}
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*EnrichFileRequest)(nil), "proto.audio.EnrichFileRequest")
	proto.RegisterType((*EnrichFileResponse)(nil), "proto.audio.EnrichFileResponse")
	proto.RegisterType((*EnrichDatasourceRequest)(nil), "proto.audio.EnrichDatasourceRequest")
	proto.RegisterType((*EnrichDatasourceResponse)(nil), "proto.audio.EnrichDatasourceResponse")
	proto.RegisterType((*HealthRequest)(nil), "proto.audio.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "proto.audio.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Service service

type ServiceClient interface {
	EnrichFile(ctx context.Context, in *EnrichFileRequest, opts ...client.CallOption) (*EnrichFileResponse, error)
	EnrichDatasource(ctx context.Context, in *EnrichDatasourceRequest, opts ...client.CallOption) (*EnrichDatasourceResponse, error)
	Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error)
}

type serviceClient struct {
	c           client.Client
	serviceName string
}

func NewServiceClient(serviceName string, c client.Client) ServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "proto.audio"
	}
	return &serviceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *serviceClient) EnrichFile(ctx context.Context, in *EnrichFileRequest, opts ...client.CallOption) (*EnrichFileResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.EnrichFile", in)
	out := new(EnrichFileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) EnrichDatasource(ctx context.Context, in *EnrichDatasourceRequest, opts ...client.CallOption) (*EnrichDatasourceResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.EnrichDatasource", in)
	out := new(EnrichDatasourceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Health", in)
	out := new(HealthResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceHandler interface {
	EnrichFile(context.Context, *EnrichFileRequest, *EnrichFileResponse) error
	EnrichDatasource(context.Context, *EnrichDatasourceRequest, *EnrichDatasourceResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Service{hdlr}, opts...))
}

type Service struct {
	ServiceHandler
}

func (h *Service) EnrichFile(ctx context.Context, in *EnrichFileRequest, out *EnrichFileResponse) error {
	return h.ServiceHandler.EnrichFile(ctx, in, out)
}

func (h *Service) EnrichDatasource(ctx context.Context, in *EnrichDatasourceRequest, out *EnrichDatasourceResponse) error {
	return h.ServiceHandler.EnrichDatasource(ctx, in, out)
}

func (h *Service) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.ServiceHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/audio/srv/proto/audio/audio.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x50, 0xcb, 0x4a, 0xf3, 0x40,
	0x14, 0xfe, 0x9b, 0x5f, 0x23, 0x1e, 0xb1, 0xe8, 0x41, 0x34, 0x44, 0x50, 0x19, 0x15, 0xea, 0x66,
	0x02, 0xba, 0x12, 0xba, 0xf3, 0x82, 0x3b, 0x21, 0x3e, 0xc1, 0x34, 0x99, 0x9a, 0xc1, 0x34, 0x13,
	0xe7, 0x52, 0xd4, 0xa5, 0x4f, 0x6e, 0x9d, 0x24, 0xbd, 0x98, 0x56, 0x37, 0xc9, 0xb9, 0x7c, 0x97,
	0x39, 0x1f, 0xf4, 0x9f, 0x85, 0xc9, 0xec, 0x80, 0x26, 0x72, 0x14, 0xbd, 0xb0, 0x0f, 0x69, 0xcb,
	0xa8, 0xcc, 0x99, 0x19, 0x4a, 0x35, 0x8a, 0x98, 0x4d, 0x85, 0x8c, 0xb4, 0x1a, 0x47, 0xa5, 0x92,
	0x46, 0xd6, 0xbd, 0xfb, 0x52, 0x37, 0xc1, 0x2d, 0xf7, 0xa3, 0x6e, 0x44, 0xae, 0x61, 0xf7, 0xae,
	0x50, 0x22, 0xc9, 0xee, 0x45, 0xce, 0x63, 0xfe, 0x6a, 0xb9, 0x36, 0xd8, 0x05, 0x4f, 0xa4, 0x41,
	0xe7, 0xa4, 0xd3, 0xdb, 0x8c, 0x27, 0x15, 0xee, 0xc1, 0xba, 0x28, 0x52, 0xfe, 0x16, 0x78, 0x6e,
	0x54, 0x35, 0xa4, 0x07, 0x38, 0x4f, 0xd5, 0xa5, 0x2c, 0x34, 0x47, 0x84, 0x35, 0x51, 0x0c, 0x65,
	0xcd, 0x76, 0x35, 0xb9, 0x80, 0x83, 0x0a, 0x79, 0xcb, 0x0c, 0xd3, 0xd2, 0xaa, 0x64, 0x95, 0x15,
	0xa1, 0x10, 0xb4, 0xa1, 0xbf, 0x48, 0x9f, 0xc2, 0xf6, 0x03, 0x67, 0xb9, 0xc9, 0x1a, 0xc1, 0x09,
	0xc8, 0xbc, 0x97, 0xbc, 0x01, 0x7d, 0xd7, 0xa4, 0x0f, 0xdd, 0x06, 0x54, 0x4b, 0xed, 0x83, 0xaf,
	0x0d, 0x33, 0x56, 0x3b, 0xdc, 0xff, 0xb8, 0xee, 0xa6, 0x16, 0xde, 0xcc, 0xe2, 0xf2, 0xd3, 0x83,
	0x8d, 0x27, 0xae, 0xc6, 0x22, 0xe1, 0xf8, 0x08, 0x30, 0xbb, 0x19, 0x8f, 0xe8, 0x5c, 0x94, 0xb4,
	0x95, 0x63, 0x78, 0xbc, 0x72, 0x5f, 0x3d, 0x83, 0xfc, 0x43, 0x06, 0x3b, 0x3f, 0xef, 0xc5, 0xb3,
	0x25, 0xb4, 0x56, 0x72, 0xe1, 0xf9, 0x1f, 0xa8, 0xa9, 0xc5, 0x0d, 0xf8, 0xd5, 0xf5, 0x18, 0x2e,
	0x50, 0x16, 0x72, 0x0b, 0x0f, 0x97, 0xee, 0x1a, 0x91, 0x81, 0xef, 0xb6, 0x57, 0x5f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x73, 0x53, 0xa8, 0xa5, 0x7b, 0x02, 0x00, 0x00,
}