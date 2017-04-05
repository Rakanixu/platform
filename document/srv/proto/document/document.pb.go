// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/document/srv/proto/document/document.proto
// DO NOT EDIT!

/*
Package proto_document is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/document/srv/proto/document/document.proto

It has these top-level messages:
	EnrichFileRequest
	EnrichFileResponse
	EnrichDatasourceRequest
	EnrichDatasourceResponse
	HealthRequest
	HealthResponse
*/
package proto_document

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

func (m *EnrichFileRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EnrichFileRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

type EnrichFileResponse struct {
	Info string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *EnrichFileResponse) Reset()                    { *m = EnrichFileResponse{} }
func (m *EnrichFileResponse) String() string            { return proto.CompactTextString(m) }
func (*EnrichFileResponse) ProtoMessage()               {}
func (*EnrichFileResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EnrichFileResponse) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type EnrichDatasourceRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *EnrichDatasourceRequest) Reset()                    { *m = EnrichDatasourceRequest{} }
func (m *EnrichDatasourceRequest) String() string            { return proto.CompactTextString(m) }
func (*EnrichDatasourceRequest) ProtoMessage()               {}
func (*EnrichDatasourceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *EnrichDatasourceRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type EnrichDatasourceResponse struct {
	Info string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *EnrichDatasourceResponse) Reset()                    { *m = EnrichDatasourceResponse{} }
func (m *EnrichDatasourceResponse) String() string            { return proto.CompactTextString(m) }
func (*EnrichDatasourceResponse) ProtoMessage()               {}
func (*EnrichDatasourceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *EnrichDatasourceResponse) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *HealthRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type HealthResponse struct {
	Status int64  `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *HealthResponse) Reset()                    { *m = HealthResponse{} }
func (m *HealthResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()               {}
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *HealthResponse) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *HealthResponse) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func init() {
	proto.RegisterType((*EnrichFileRequest)(nil), "proto.document.EnrichFileRequest")
	proto.RegisterType((*EnrichFileResponse)(nil), "proto.document.EnrichFileResponse")
	proto.RegisterType((*EnrichDatasourceRequest)(nil), "proto.document.EnrichDatasourceRequest")
	proto.RegisterType((*EnrichDatasourceResponse)(nil), "proto.document.EnrichDatasourceResponse")
	proto.RegisterType((*HealthRequest)(nil), "proto.document.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "proto.document.HealthResponse")
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
		serviceName = "proto.document"
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
	proto.RegisterFile("github.com/kazoup/platform/document/srv/proto/document/document.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 304 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xcf, 0x4e, 0xf3, 0x30,
	0x10, 0xc4, 0xbf, 0xe6, 0x2b, 0x45, 0xac, 0x44, 0x05, 0x2b, 0x04, 0x55, 0x25, 0x10, 0x98, 0x03,
	0xe5, 0xe2, 0x48, 0x70, 0x42, 0xe2, 0x48, 0x11, 0x5c, 0x83, 0x78, 0x00, 0x37, 0x71, 0x1b, 0x8b,
	0x24, 0x0e, 0xfe, 0x53, 0x01, 0xcf, 0xc1, 0x03, 0x23, 0x9c, 0xb8, 0xa1, 0x0d, 0x2d, 0xa7, 0xac,
	0x37, 0xbf, 0x99, 0xb1, 0x07, 0xc6, 0x33, 0x61, 0x52, 0x3b, 0xa1, 0xb1, 0xcc, 0xc3, 0x17, 0xf6,
	0x21, 0x6d, 0x19, 0x96, 0x19, 0x33, 0x53, 0xa9, 0xf2, 0x30, 0x91, 0xb1, 0xcd, 0x79, 0x61, 0x42,
	0xad, 0xe6, 0x61, 0xa9, 0xa4, 0x91, 0xcd, 0xca, 0x0f, 0xd4, 0xed, 0xb1, 0xef, 0x3e, 0xd4, 0x6f,
	0xc9, 0x0d, 0xec, 0x8f, 0x0b, 0x25, 0xe2, 0xf4, 0x5e, 0x64, 0x3c, 0xe2, 0xaf, 0x96, 0x6b, 0x83,
	0x7d, 0x08, 0x44, 0x32, 0xe8, 0x9c, 0x76, 0x46, 0x3b, 0x51, 0x20, 0x12, 0x3c, 0x80, 0x2d, 0x51,
	0x24, 0xfc, 0x6d, 0x10, 0xb8, 0x55, 0x75, 0x20, 0x23, 0xc0, 0x9f, 0x52, 0x5d, 0xca, 0x42, 0x73,
	0x44, 0xe8, 0x8a, 0x62, 0x2a, 0x6b, 0xb5, 0x9b, 0xc9, 0x25, 0x1c, 0x55, 0xe4, 0x1d, 0x33, 0x4c,
	0x4b, 0xab, 0xe2, 0x75, 0x51, 0x84, 0xc2, 0xa0, 0x8d, 0x6e, 0xb0, 0x3e, 0x87, 0xdd, 0x07, 0xce,
	0x32, 0x93, 0x7a, 0x43, 0x84, 0xae, 0x79, 0x2f, 0xb9, 0x87, 0xbe, 0x67, 0x72, 0x0b, 0x7d, 0x0f,
	0xd5, 0x56, 0x87, 0xd0, 0xd3, 0x86, 0x19, 0xab, 0x1d, 0xf7, 0x3f, 0xaa, 0x4f, 0x8b, 0x88, 0xa0,
	0x89, 0xb8, 0xfa, 0x0c, 0x60, 0xfb, 0x89, 0xab, 0xb9, 0x88, 0x39, 0x3e, 0x03, 0x34, 0x6f, 0xc6,
	0x33, 0xba, 0xdc, 0x26, 0x6d, 0x55, 0x39, 0x24, 0x9b, 0x90, 0xea, 0x32, 0xe4, 0x1f, 0xce, 0x60,
	0x6f, 0xf5, 0xd5, 0x78, 0xf1, 0xbb, 0xb2, 0x55, 0xe1, 0x70, 0xf4, 0x37, 0xb8, 0x08, 0x7a, 0x84,
	0x5e, 0xd5, 0x04, 0x1e, 0xaf, 0xaa, 0x96, 0x6a, 0x1c, 0x9e, 0xac, 0xfb, 0xed, 0xad, 0x26, 0x3d,
	0x07, 0x5c, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x8d, 0xde, 0x23, 0x08, 0x99, 0x02, 0x00, 0x00,
}
