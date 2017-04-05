// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/sentiment/srv/proto/sentiment/sentiment.proto
// DO NOT EDIT!

/*
Package proto_sentiment is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/sentiment/srv/proto/sentiment/sentiment.proto

It has these top-level messages:
	AnalyzeFileRequest
	AnalyzeFileResponse
	HealthRequest
	HealthResponse
*/
package proto_sentiment

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

type AnalyzeFileRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Id    string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *AnalyzeFileRequest) Reset()                    { *m = AnalyzeFileRequest{} }
func (m *AnalyzeFileRequest) String() string            { return proto.CompactTextString(m) }
func (*AnalyzeFileRequest) ProtoMessage()               {}
func (*AnalyzeFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AnalyzeFileRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *AnalyzeFileRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type AnalyzeFileResponse struct {
	Info string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *AnalyzeFileResponse) Reset()                    { *m = AnalyzeFileResponse{} }
func (m *AnalyzeFileResponse) String() string            { return proto.CompactTextString(m) }
func (*AnalyzeFileResponse) ProtoMessage()               {}
func (*AnalyzeFileResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AnalyzeFileResponse) GetInfo() string {
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
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

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
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
	proto.RegisterType((*AnalyzeFileRequest)(nil), "proto.sentiment.AnalyzeFileRequest")
	proto.RegisterType((*AnalyzeFileResponse)(nil), "proto.sentiment.AnalyzeFileResponse")
	proto.RegisterType((*HealthRequest)(nil), "proto.sentiment.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "proto.sentiment.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Service service

type ServiceClient interface {
	AnalyzeFile(ctx context.Context, in *AnalyzeFileRequest, opts ...client.CallOption) (*AnalyzeFileResponse, error)
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
		serviceName = "proto.sentiment"
	}
	return &serviceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *serviceClient) AnalyzeFile(ctx context.Context, in *AnalyzeFileRequest, opts ...client.CallOption) (*AnalyzeFileResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.AnalyzeFile", in)
	out := new(AnalyzeFileResponse)
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
	AnalyzeFile(context.Context, *AnalyzeFileRequest, *AnalyzeFileResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Service{hdlr}, opts...))
}

type Service struct {
	ServiceHandler
}

func (h *Service) AnalyzeFile(ctx context.Context, in *AnalyzeFileRequest, out *AnalyzeFileResponse) error {
	return h.ServiceHandler.AnalyzeFile(ctx, in, out)
}

func (h *Service) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.ServiceHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/sentiment/srv/proto/sentiment/sentiment.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x4f, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x25, 0xa1, 0x04, 0x71, 0x88, 0x22, 0x1d, 0x08, 0x55, 0x1d, 0x00, 0xa5, 0x0c, 0xb0, 0x24,
	0x12, 0x6c, 0x88, 0x85, 0x05, 0x55, 0x62, 0x0b, 0x1b, 0x9b, 0xdb, 0x5c, 0xa9, 0x45, 0x62, 0x9b,
	0xf8, 0x52, 0xd1, 0x7e, 0x16, 0x5f, 0x88, 0x70, 0xe2, 0xd2, 0x52, 0xd1, 0xc9, 0xe7, 0x77, 0xef,
	0xdd, 0x7b, 0x0f, 0x86, 0x6f, 0x92, 0xa7, 0xf5, 0x28, 0x19, 0xeb, 0x32, 0x7d, 0x17, 0x0b, 0x5d,
	0x9b, 0xd4, 0x14, 0x82, 0x27, 0xba, 0x2a, 0x53, 0x4b, 0x8a, 0x65, 0x49, 0x8a, 0x53, 0x5b, 0xcd,
	0x52, 0x53, 0x69, 0xd6, 0xab, 0x98, 0x9f, 0x12, 0xb7, 0xc1, 0x63, 0xf7, 0x24, 0x4b, 0x38, 0xbe,
	0x07, 0x7c, 0x54, 0xa2, 0x98, 0x2f, 0xe8, 0x49, 0x16, 0x94, 0xd1, 0x47, 0x4d, 0x96, 0xf1, 0x14,
	0xf6, 0xa4, 0xca, 0xe9, 0xb3, 0x17, 0x5c, 0x06, 0xd7, 0x07, 0x59, 0xf3, 0xc1, 0x2e, 0x84, 0x32,
	0xef, 0x85, 0x0e, 0x0a, 0x65, 0x1e, 0xdf, 0xc0, 0xc9, 0x9a, 0xd6, 0x1a, 0xad, 0x2c, 0x21, 0x42,
	0x47, 0xaa, 0x89, 0x6e, 0xb5, 0x6e, 0x8e, 0x07, 0x70, 0x34, 0x24, 0x51, 0xf0, 0xd4, 0x3b, 0x20,
	0x74, 0x78, 0x6e, 0xc8, 0x93, 0x7e, 0xe6, 0xf8, 0x01, 0xba, 0x9e, 0xd4, 0x9e, 0x3a, 0x83, 0xc8,
	0xb2, 0xe0, 0xda, 0x3a, 0xde, 0x6e, 0xd6, 0xfe, 0x96, 0x16, 0xe1, 0xaf, 0xc5, 0xed, 0x57, 0x00,
	0xfb, 0x2f, 0x54, 0xcd, 0xe4, 0x98, 0xf0, 0x15, 0x0e, 0x57, 0x92, 0xe1, 0x20, 0xf9, 0x53, 0x3b,
	0xd9, 0xec, 0xdc, 0xbf, 0xda, 0x4e, 0x6a, 0x12, 0xc5, 0x3b, 0xf8, 0x0c, 0x51, 0x93, 0x12, 0xcf,
	0x37, 0x14, 0x6b, 0x1d, 0xfb, 0x17, 0xff, 0xee, 0xfd, 0xb1, 0x51, 0xe4, 0x18, 0x77, 0xdf, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xf6, 0xc9, 0x50, 0x77, 0xe2, 0x01, 0x00, 0x00,
}
