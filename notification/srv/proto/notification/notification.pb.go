// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/notification/srv/proto/notification/notification.proto
// DO NOT EDIT!

/*
Package proto_notification is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/notification/srv/proto/notification/notification.proto

It has these top-level messages:
	NotificationMessage
	StreamRequest
	StreamResponse
	HealthRequest
	HealthResponse
*/
package proto_notification

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

type NotificationMessage struct {
	Info   string `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	Method string `protobuf:"bytes,2,opt,name=method" json:"method,omitempty"`
	Data   string `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
	UserId string `protobuf:"bytes,4,opt,name=user_id,json=userId" json:"user_id,omitempty"`
}

func (m *NotificationMessage) Reset()                    { *m = NotificationMessage{} }
func (m *NotificationMessage) String() string            { return proto.CompactTextString(m) }
func (*NotificationMessage) ProtoMessage()               {}
func (*NotificationMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NotificationMessage) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *NotificationMessage) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *NotificationMessage) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *NotificationMessage) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type StreamRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *StreamRequest) Reset()                    { *m = StreamRequest{} }
func (m *StreamRequest) String() string            { return proto.CompactTextString(m) }
func (*StreamRequest) ProtoMessage()               {}
func (*StreamRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StreamRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *StreamRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type StreamResponse struct {
	Message *NotificationMessage `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *StreamResponse) Reset()                    { *m = StreamResponse{} }
func (m *StreamResponse) String() string            { return proto.CompactTextString(m) }
func (*StreamResponse) ProtoMessage()               {}
func (*StreamResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StreamResponse) GetMessage() *NotificationMessage {
	if m != nil {
		return m.Message
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
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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
	proto.RegisterType((*NotificationMessage)(nil), "proto.notification.NotificationMessage")
	proto.RegisterType((*StreamRequest)(nil), "proto.notification.StreamRequest")
	proto.RegisterType((*StreamResponse)(nil), "proto.notification.StreamResponse")
	proto.RegisterType((*HealthRequest)(nil), "proto.notification.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "proto.notification.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Service service

type ServiceClient interface {
	Stream(ctx context.Context, in *StreamRequest, opts ...client.CallOption) (Service_StreamClient, error)
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
		serviceName = "proto.notification"
	}
	return &serviceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *serviceClient) Stream(ctx context.Context, in *StreamRequest, opts ...client.CallOption) (Service_StreamClient, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Stream", &StreamRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &serviceStreamClient{stream}, nil
}

type Service_StreamClient interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamResponse, error)
}

type serviceStreamClient struct {
	stream client.Streamer
}

func (x *serviceStreamClient) Close() error {
	return x.stream.Close()
}

func (x *serviceStreamClient) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *serviceStreamClient) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *serviceStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
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
	Stream(context.Context, *StreamRequest, Service_StreamStream) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Service{hdlr}, opts...))
}

type Service struct {
	ServiceHandler
}

func (h *Service) Stream(ctx context.Context, stream server.Streamer) error {
	m := new(StreamRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.ServiceHandler.Stream(ctx, m, &serviceStreamStream{stream})
}

type Service_StreamStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamResponse) error
}

type serviceStreamStream struct {
	stream server.Streamer
}

func (x *serviceStreamStream) Close() error {
	return x.stream.Close()
}

func (x *serviceStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *serviceStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *serviceStreamStream) Send(m *StreamResponse) error {
	return x.stream.Send(m)
}

func (h *Service) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.ServiceHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/notification/srv/proto/notification/notification.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4e, 0xf3, 0x30,
	0x10, 0x85, 0x9b, 0xb6, 0x7f, 0xaa, 0x7f, 0x50, 0xbb, 0x30, 0xa8, 0x44, 0xac, 0xc0, 0x2c, 0x60,
	0x95, 0xa0, 0xb2, 0x45, 0x48, 0xec, 0x60, 0x01, 0x52, 0x93, 0x03, 0x20, 0xb7, 0x99, 0x36, 0x56,
	0x9b, 0x38, 0xc4, 0x93, 0x4a, 0x70, 0x2b, 0x6e, 0x88, 0x62, 0x27, 0x90, 0x88, 0x88, 0x95, 0xe7,
	0xd9, 0x9f, 0x67, 0xde, 0x1b, 0x58, 0x6e, 0x25, 0x25, 0xe5, 0xca, 0x5f, 0xab, 0x34, 0xd8, 0x89,
	0x0f, 0x55, 0xe6, 0x41, 0xbe, 0x17, 0xb4, 0x51, 0x45, 0x1a, 0x64, 0x8a, 0xe4, 0x46, 0xae, 0x05,
	0x49, 0x95, 0x05, 0xba, 0x38, 0x04, 0x79, 0xa1, 0x48, 0x75, 0xaf, 0xdb, 0xc2, 0x37, 0xef, 0x8c,
	0x99, 0xc3, 0x6f, 0xbf, 0xf0, 0x0c, 0x8e, 0x5f, 0x5a, 0xfa, 0x19, 0xb5, 0x16, 0x5b, 0x64, 0x0c,
	0xc6, 0x32, 0xdb, 0x28, 0xcf, 0x39, 0x77, 0xae, 0xff, 0x87, 0xa6, 0x66, 0x73, 0x70, 0x53, 0xa4,
	0x44, 0xc5, 0xde, 0xd0, 0xdc, 0xd6, 0xaa, 0x62, 0x63, 0x41, 0xc2, 0x1b, 0x59, 0xb6, 0xaa, 0xd9,
	0x29, 0x4c, 0x4a, 0x8d, 0xc5, 0xab, 0x8c, 0xbd, 0xb1, 0x85, 0x2b, 0xf9, 0x14, 0xf3, 0x7b, 0x98,
	0x46, 0x54, 0xa0, 0x48, 0x43, 0x7c, 0x2b, 0x51, 0x53, 0x9b, 0x74, 0xda, 0x24, 0x3b, 0x81, 0x7f,
	0xa4, 0x76, 0x98, 0xd5, 0xd3, 0xac, 0xe0, 0x11, 0xcc, 0x9a, 0xff, 0x3a, 0x57, 0x99, 0x46, 0xf6,
	0x00, 0x93, 0xd4, 0xba, 0x36, 0x0d, 0x8e, 0x16, 0x57, 0xfe, 0xef, 0x9c, 0x7e, 0x4f, 0xc8, 0xb0,
	0xf9, 0xc7, 0x2f, 0x61, 0xfa, 0x88, 0x62, 0x4f, 0x49, 0x63, 0x8a, 0xc1, 0x98, 0xde, 0x73, 0x6c,
	0xe2, 0x57, 0x35, 0xbf, 0x83, 0x59, 0x03, 0xd5, 0x93, 0xe7, 0xe0, 0x6a, 0x12, 0x54, 0x6a, 0xc3,
	0x8d, 0xc2, 0x5a, 0x7d, 0x2f, 0x6f, 0xf8, 0xb3, 0xbc, 0xc5, 0xa7, 0x03, 0x93, 0x08, 0x8b, 0x83,
	0x5c, 0x23, 0x8b, 0xc0, 0xb5, 0x19, 0xd8, 0x45, 0x9f, 0xd5, 0xce, 0x7e, 0xce, 0xf8, 0x5f, 0x88,
	0x35, 0xc2, 0x07, 0x37, 0x0e, 0x5b, 0x82, 0x6b, 0xed, 0xf5, 0x37, 0xed, 0xe4, 0xeb, 0x6f, 0xda,
	0x4d, 0xc7, 0x07, 0x2b, 0xd7, 0x40, 0xb7, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x15, 0xb4, 0x70,
	0x3f, 0x8b, 0x02, 0x00, 0x00,
}
