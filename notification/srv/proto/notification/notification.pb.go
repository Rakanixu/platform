// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/notification/srv/proto/notification/notification.proto
// DO NOT EDIT!

/*
Package go_micro_srv_notification is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/notification/srv/proto/notification/notification.proto

It has these top-level messages:
	NotificationMessage
	StreamRequest
	StreamResponse
*/
package go_micro_srv_notification

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

type StreamRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
}

func (m *StreamRequest) Reset()                    { *m = StreamRequest{} }
func (m *StreamRequest) String() string            { return proto.CompactTextString(m) }
func (*StreamRequest) ProtoMessage()               {}
func (*StreamRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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

func init() {
	proto.RegisterType((*NotificationMessage)(nil), "go.micro.srv.notification.NotificationMessage")
	proto.RegisterType((*StreamRequest)(nil), "go.micro.srv.notification.StreamRequest")
	proto.RegisterType((*StreamResponse)(nil), "go.micro.srv.notification.StreamResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Notification service

type NotificationClient interface {
	Stream(ctx context.Context, in *StreamRequest, opts ...client.CallOption) (Notification_StreamClient, error)
}

type notificationClient struct {
	c           client.Client
	serviceName string
}

func NewNotificationClient(serviceName string, c client.Client) NotificationClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.notification"
	}
	return &notificationClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *notificationClient) Stream(ctx context.Context, in *StreamRequest, opts ...client.CallOption) (Notification_StreamClient, error) {
	req := c.c.NewRequest(c.serviceName, "Notification.Stream", &StreamRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &notificationStreamClient{stream}, nil
}

type Notification_StreamClient interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamResponse, error)
}

type notificationStreamClient struct {
	stream client.Streamer
}

func (x *notificationStreamClient) Close() error {
	return x.stream.Close()
}

func (x *notificationStreamClient) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *notificationStreamClient) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *notificationStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Notification service

type NotificationHandler interface {
	Stream(context.Context, *StreamRequest, Notification_StreamStream) error
}

func RegisterNotificationHandler(s server.Server, hdlr NotificationHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Notification{hdlr}, opts...))
}

type Notification struct {
	NotificationHandler
}

func (h *Notification) Stream(ctx context.Context, stream server.Streamer) error {
	m := new(StreamRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.NotificationHandler.Stream(ctx, m, &notificationStreamStream{stream})
}

type Notification_StreamStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamResponse) error
}

type notificationStreamStream struct {
	stream server.Streamer
}

func (x *notificationStreamStream) Close() error {
	return x.stream.Close()
}

func (x *notificationStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *notificationStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *notificationStreamStream) Send(m *StreamResponse) error {
	return x.stream.Send(m)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/notification/srv/proto/notification/notification.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x91, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x09, 0x54, 0xad, 0x38, 0xfe, 0x0c, 0x46, 0x82, 0xc0, 0x84, 0x32, 0x85, 0xc5, 0x46,
	0xe5, 0x4b, 0xc0, 0x00, 0x12, 0x61, 0x63, 0x41, 0x6e, 0xe2, 0xa4, 0x16, 0xd8, 0x97, 0xda, 0x17,
	0x06, 0x3e, 0x3d, 0xc4, 0xa6, 0x92, 0x23, 0x01, 0x62, 0xbb, 0x77, 0xf7, 0x3b, 0xfb, 0x3d, 0x1d,
	0x3c, 0x76, 0x9a, 0xd6, 0xc3, 0x8a, 0xd7, 0x68, 0xc4, 0xab, 0xfc, 0xc0, 0xa1, 0x17, 0xfd, 0x9b,
	0xa4, 0x16, 0x9d, 0x11, 0x16, 0x49, 0xb7, 0xba, 0x96, 0xa4, 0xd1, 0x0a, 0xef, 0xde, 0x45, 0xef,
	0x90, 0x70, 0xda, 0x4e, 0x05, 0x0f, 0x73, 0x76, 0xde, 0x21, 0x37, 0xba, 0x76, 0xc8, 0xbf, 0x76,
	0x78, 0x0a, 0x14, 0x16, 0x4e, 0x1e, 0x12, 0x7d, 0xaf, 0xbc, 0x97, 0x9d, 0x62, 0x0c, 0x66, 0xda,
	0xb6, 0x98, 0x67, 0x97, 0x59, 0xb9, 0x5f, 0x85, 0x9a, 0x9d, 0xc2, 0xdc, 0x28, 0x5a, 0x63, 0x93,
	0xef, 0x86, 0xee, 0xb7, 0x1a, 0xd9, 0x46, 0x92, 0xcc, 0xf7, 0x22, 0x3b, 0xd6, 0xec, 0x0c, 0x16,
	0x83, 0x57, 0xee, 0x45, 0x37, 0xf9, 0x2c, 0xc2, 0xa3, 0xbc, 0x6b, 0x8a, 0x12, 0x8e, 0x9e, 0xc8,
	0x29, 0x69, 0x2a, 0xb5, 0x19, 0x94, 0xa7, 0x94, 0xcc, 0x26, 0xe4, 0x33, 0x1c, 0x6f, 0x49, 0xdf,
	0xa3, 0xf5, 0x8a, 0xdd, 0xc2, 0xc2, 0x44, 0x7f, 0x01, 0x3d, 0x58, 0x72, 0xfe, 0x6b, 0x30, 0xfe,
	0x43, 0xaa, 0x6a, 0xbb, 0xbe, 0xdc, 0xc0, 0x61, 0x3a, 0x67, 0x12, 0xe6, 0xf1, 0x2f, 0x56, 0xfe,
	0xf1, 0xe4, 0xc4, 0xf8, 0xc5, 0xd5, 0x3f, 0xc8, 0x68, 0xbc, 0xd8, 0xb9, 0xce, 0x56, 0xf3, 0x70,
	0x8a, 0x9b, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x16, 0xe1, 0xd0, 0x5f, 0xdf, 0x01, 0x00, 0x00,
}