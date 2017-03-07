// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/quota/srv/proto/quota/quota.proto
// DO NOT EDIT!

/*
Package go_micro_srv_quota is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/quota/srv/proto/quota/quota.proto

It has these top-level messages:
	Quota
	ReadRequest
	ReadResponse
	HealthRequest
	HealthResponse
*/
package go_micro_srv_quota

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

type Quota struct {
	Name           string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Rate           int64  `protobuf:"varint,2,opt,name=rate" json:"rate,omitempty"`
	ResetTimestamp int64  `protobuf:"varint,3,opt,name=reset_timestamp,json=resetTimestamp" json:"reset_timestamp,omitempty"`
	Quota          int64  `protobuf:"varint,4,opt,name=quota" json:"quota,omitempty"`
}

func (m *Quota) Reset()                    { *m = Quota{} }
func (m *Quota) String() string            { return proto.CompactTextString(m) }
func (*Quota) ProtoMessage()               {}
func (*Quota) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ReadRequest struct {
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ReadResponse struct {
	Quota []*Quota `protobuf:"bytes,1,rep,name=quota" json:"quota,omitempty"`
}

func (m *ReadResponse) Reset()                    { *m = ReadResponse{} }
func (m *ReadResponse) String() string            { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()               {}
func (*ReadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ReadResponse) GetQuota() []*Quota {
	if m != nil {
		return m.Quota
	}
	return nil
}

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type HealthResponse struct {
	Status int64  `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *HealthResponse) Reset()                    { *m = HealthResponse{} }
func (m *HealthResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()               {}
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*Quota)(nil), "go.micro.srv.quota.Quota")
	proto.RegisterType((*ReadRequest)(nil), "go.micro.srv.quota.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "go.micro.srv.quota.ReadResponse")
	proto.RegisterType((*HealthRequest)(nil), "go.micro.srv.quota.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "go.micro.srv.quota.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for File service

type FileClient interface {
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error)
}

type fileClient struct {
	c           client.Client
	serviceName string
}

func NewFileClient(serviceName string, c client.Client) FileClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.quota"
	}
	return &fileClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *fileClient) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.serviceName, "File.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error) {
	req := c.c.NewRequest(c.serviceName, "File.Health", in)
	out := new(HealthResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for File service

type FileHandler interface {
	Read(context.Context, *ReadRequest, *ReadResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterFileHandler(s server.Server, hdlr FileHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&File{hdlr}, opts...))
}

type File struct {
	FileHandler
}

func (h *File) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.FileHandler.Read(ctx, in, out)
}

func (h *File) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.FileHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/quota/srv/proto/quota/quota.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x51, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x75, 0xdd, 0x6d, 0xc1, 0xa9, 0xad, 0x10, 0x44, 0xd6, 0x5e, 0xac, 0xf1, 0x60, 0x4f, 0x09,
	0xd4, 0x6b, 0xc1, 0x9b, 0x08, 0x9e, 0xba, 0x78, 0x97, 0xb4, 0xa6, 0xed, 0xe2, 0x66, 0x93, 0x26,
	0x59, 0x41, 0x3f, 0xc8, 0xef, 0x74, 0x33, 0xbb, 0x2b, 0x8a, 0xd5, 0x4b, 0x78, 0xf3, 0xe6, 0xcd,
	0xcc, 0x7b, 0x04, 0xe6, 0x9b, 0xdc, 0x6f, 0xab, 0x25, 0x5b, 0x69, 0xc5, 0x5f, 0xc4, 0xbb, 0xae,
	0x0c, 0x37, 0x85, 0xf0, 0x6b, 0x6d, 0x15, 0xdf, 0x55, 0xda, 0x0b, 0xee, 0xec, 0x2b, 0x37, 0x56,
	0x7b, 0xdd, 0xd6, 0xf8, 0x32, 0x64, 0x08, 0xd9, 0x68, 0xa6, 0xf2, 0x95, 0xd5, 0xac, 0x56, 0x31,
	0xec, 0xd0, 0x12, 0x7a, 0x8b, 0x00, 0x08, 0x81, 0xa4, 0x14, 0x4a, 0xa6, 0xd1, 0x24, 0x9a, 0x1e,
	0x65, 0x88, 0x03, 0x67, 0x85, 0x97, 0xe9, 0x61, 0xcd, 0xc5, 0x19, 0x62, 0x72, 0x0d, 0x27, 0x56,
	0x3a, 0xe9, 0x9f, 0x7c, 0xae, 0xa4, 0xf3, 0x42, 0x99, 0x34, 0xc6, 0xf6, 0x08, 0xe9, 0xc7, 0x8e,
	0x25, 0xa7, 0xd0, 0xc3, 0x13, 0x69, 0x82, 0xed, 0xa6, 0xa0, 0x43, 0x18, 0x64, 0x52, 0x3c, 0x67,
	0x72, 0x57, 0xd5, 0x3a, 0x7a, 0x0b, 0xc7, 0x4d, 0xe9, 0x8c, 0x2e, 0x9d, 0x24, 0xbc, 0x1b, 0x8a,
	0x26, 0xf1, 0x74, 0x30, 0x3b, 0x67, 0xbf, 0x2d, 0x33, 0xf4, 0xdb, 0xed, 0xbb, 0x82, 0xe1, 0xbd,
	0x14, 0x85, 0xdf, 0xb6, 0x1b, 0x83, 0x67, 0xff, 0x66, 0xbe, 0x72, 0x04, 0x4c, 0xe7, 0x30, 0xea,
	0x44, 0xed, 0x9d, 0x33, 0xe8, 0xd7, 0x2e, 0x7d, 0xe5, 0x50, 0x17, 0x67, 0x6d, 0x15, 0xa6, 0xf3,
	0x72, 0xad, 0x31, 0x71, 0x3d, 0x1d, 0xf0, 0xec, 0x23, 0x82, 0xe4, 0x2e, 0x2f, 0x24, 0x79, 0x80,
	0x24, 0x98, 0x25, 0x17, 0xfb, 0x5c, 0x7d, 0x4b, 0x35, 0x9e, 0xfc, 0x2d, 0x68, 0xee, 0xd3, 0x03,
	0xb2, 0x80, 0x7e, 0xe3, 0x89, 0x5c, 0xee, 0x53, 0xff, 0x08, 0x35, 0xa6, 0xff, 0x49, 0xba, 0x95,
	0xcb, 0x3e, 0x7e, 0xf3, 0xcd, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe7, 0xd9, 0xd4, 0x7c, 0x26,
	0x02, 0x00, 0x00,
}
