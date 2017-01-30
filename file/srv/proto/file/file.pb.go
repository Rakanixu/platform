// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/file/srv/proto/file/file.proto
// DO NOT EDIT!

/*
Package go_micro_srv_file is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/file/srv/proto/file/file.proto

It has these top-level messages:
	CreateRequest
	CreateResponse
	DeleteRequest
	DeleteResponse
	ShareRequest
	ShareResponse
	HealthRequest
	HealthResponse
*/
package go_micro_srv_file

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
	DatasourceId string `protobuf:"bytes,1,opt,name=datasource_id,json=datasourceId" json:"datasource_id,omitempty"`
	FileName     string `protobuf:"bytes,2,opt,name=file_name,json=fileName" json:"file_name,omitempty"`
	MimeType     string `protobuf:"bytes,3,opt,name=mime_type,json=mimeType" json:"mime_type,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CreateResponse struct {
	DocUrl string `protobuf:"bytes,1,opt,name=doc_url,json=docUrl" json:"doc_url,omitempty"`
	Data   string `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *CreateResponse) Reset()                    { *m = CreateResponse{} }
func (m *CreateResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()               {}
func (*CreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type DeleteRequest struct {
	DatasourceId     string `protobuf:"bytes,1,opt,name=datasource_id,json=datasourceId" json:"datasource_id,omitempty"`
	Index            string `protobuf:"bytes,2,opt,name=index" json:"index,omitempty"`
	FileId           string `protobuf:"bytes,3,opt,name=file_id,json=fileId" json:"file_id,omitempty"`
	OriginalId       string `protobuf:"bytes,4,opt,name=original_id,json=originalId" json:"original_id,omitempty"`
	OriginalFilePath string `protobuf:"bytes,5,opt,name=original_file_path,json=originalFilePath" json:"original_file_path,omitempty"`
	UserId           string `protobuf:"bytes,6,opt,name=user_id,json=userId" json:"user_id,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type DeleteResponse struct {
}

func (m *DeleteResponse) Reset()                    { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()               {}
func (*DeleteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ShareRequest struct {
	DatasourceId  string `protobuf:"bytes,1,opt,name=datasource_id,json=datasourceId" json:"datasource_id,omitempty"`
	Index         string `protobuf:"bytes,2,opt,name=index" json:"index,omitempty"`
	FileId        string `protobuf:"bytes,3,opt,name=file_id,json=fileId" json:"file_id,omitempty"`
	OriginalId    string `protobuf:"bytes,4,opt,name=original_id,json=originalId" json:"original_id,omitempty"`
	DestinationId string `protobuf:"bytes,5,opt,name=destination_id,json=destinationId" json:"destination_id,omitempty"`
	SharePublicly bool   `protobuf:"varint,6,opt,name=share_publicly,json=sharePublicly" json:"share_publicly,omitempty"`
}

func (m *ShareRequest) Reset()                    { *m = ShareRequest{} }
func (m *ShareRequest) String() string            { return proto.CompactTextString(m) }
func (*ShareRequest) ProtoMessage()               {}
func (*ShareRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type ShareResponse struct {
	PublicUrl     string `protobuf:"bytes,1,opt,name=public_url,json=publicUrl" json:"public_url,omitempty"`
	SharePublicly bool   `protobuf:"varint,2,opt,name=share_publicly,json=sharePublicly" json:"share_publicly,omitempty"`
}

func (m *ShareResponse) Reset()                    { *m = ShareResponse{} }
func (m *ShareResponse) String() string            { return proto.CompactTextString(m) }
func (*ShareResponse) ProtoMessage()               {}
func (*ShareResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type HealthResponse struct {
	Status int64  `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *HealthResponse) Reset()                    { *m = HealthResponse{} }
func (m *HealthResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()               {}
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*CreateRequest)(nil), "go.micro.srv.file.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "go.micro.srv.file.CreateResponse")
	proto.RegisterType((*DeleteRequest)(nil), "go.micro.srv.file.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "go.micro.srv.file.DeleteResponse")
	proto.RegisterType((*ShareRequest)(nil), "go.micro.srv.file.ShareRequest")
	proto.RegisterType((*ShareResponse)(nil), "go.micro.srv.file.ShareResponse")
	proto.RegisterType((*HealthRequest)(nil), "go.micro.srv.file.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "go.micro.srv.file.HealthResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for File service

type FileClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*CreateResponse, error)
	Share(ctx context.Context, in *ShareRequest, opts ...client.CallOption) (*ShareResponse, error)
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
		serviceName = "go.micro.srv.file"
	}
	return &fileClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *fileClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "File.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "File.Delete", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) Share(ctx context.Context, in *ShareRequest, opts ...client.CallOption) (*ShareResponse, error) {
	req := c.c.NewRequest(c.serviceName, "File.Share", in)
	out := new(ShareResponse)
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
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Delete(context.Context, *DeleteRequest, *CreateResponse) error
	Share(context.Context, *ShareRequest, *ShareResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterFileHandler(s server.Server, hdlr FileHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&File{hdlr}, opts...))
}

type File struct {
	FileHandler
}

func (h *File) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.FileHandler.Create(ctx, in, out)
}

func (h *File) Delete(ctx context.Context, in *DeleteRequest, out *CreateResponse) error {
	return h.FileHandler.Delete(ctx, in, out)
}

func (h *File) Share(ctx context.Context, in *ShareRequest, out *ShareResponse) error {
	return h.FileHandler.Share(ctx, in, out)
}

func (h *File) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.FileHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/file/srv/proto/file/file.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 500 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xc4, 0x54, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0x5e, 0xbb, 0x36, 0xac, 0x8f, 0xa5, 0x1a, 0x16, 0x82, 0x6a, 0x08, 0xad, 0x04, 0x21, 0x71,
	0x40, 0x89, 0x04, 0x27, 0x24, 0x38, 0x81, 0x10, 0x93, 0x10, 0x4c, 0x85, 0x9d, 0x23, 0x37, 0x71,
	0x5b, 0x0b, 0x27, 0x0e, 0xb6, 0x83, 0x18, 0xff, 0x20, 0x07, 0x4e, 0xfc, 0x47, 0xf8, 0xd9, 0x0e,
	0xb4, 0x5a, 0x40, 0x70, 0xda, 0xa5, 0xf2, 0xf7, 0xde, 0xf7, 0xbe, 0xf7, 0x33, 0x85, 0xa7, 0x6b,
	0x6e, 0x36, 0xed, 0x32, 0x2d, 0x64, 0x95, 0x7d, 0xa4, 0x5f, 0x65, 0xdb, 0x64, 0x8d, 0xa0, 0x66,
	0x25, 0x55, 0x95, 0xad, 0xb8, 0x60, 0x99, 0x56, 0x9f, 0xb3, 0x46, 0x49, 0x23, 0x3d, 0xc4, 0x9f,
	0xd4, 0x61, 0x72, 0x63, 0x2d, 0xd3, 0x8a, 0x17, 0x4a, 0xa6, 0x96, 0x93, 0xa2, 0x23, 0xa9, 0x21,
	0x7e, 0xa1, 0x18, 0x35, 0x6c, 0xc1, 0x3e, 0xb5, 0x4c, 0x1b, 0x72, 0x1f, 0xe2, 0x92, 0x1a, 0xaa,
	0x65, 0xab, 0x0a, 0x96, 0xf3, 0x72, 0x36, 0x98, 0x0f, 0x1e, 0x4e, 0x16, 0x87, 0xbf, 0x8d, 0xa7,
	0x25, 0xb9, 0x03, 0x13, 0x8c, 0xce, 0x6b, 0x5a, 0xb1, 0xd9, 0xd0, 0x11, 0x0e, 0xd0, 0xf0, 0xd6,
	0x62, 0x74, 0x56, 0xbc, 0x62, 0xb9, 0xb9, 0x68, 0xd8, 0x6c, 0xdf, 0x3b, 0xd1, 0xf0, 0xc1, 0xe2,
	0xe4, 0x39, 0x4c, 0xbb, 0x7c, 0xba, 0x91, 0xb5, 0x66, 0xe4, 0x36, 0x5c, 0x2b, 0x65, 0x91, 0xb7,
	0x4a, 0x84, 0x54, 0x91, 0x85, 0xe7, 0x4a, 0x10, 0x02, 0x23, 0x4c, 0x1a, 0xf4, 0xdd, 0x3b, 0xf9,
	0x3e, 0x80, 0xf8, 0x25, 0x13, 0xec, 0x3f, 0xeb, 0xbd, 0x09, 0x63, 0x5e, 0x97, 0xec, 0x4b, 0xd0,
	0xf2, 0x00, 0x33, 0xbb, 0x2e, 0x6c, 0x90, 0x2f, 0x33, 0x42, 0x68, 0xe9, 0x27, 0x70, 0x5d, 0x2a,
	0xbe, 0xe6, 0x35, 0x15, 0xe8, 0x1c, 0x39, 0x27, 0x74, 0x26, 0x4b, 0x78, 0x04, 0xe4, 0x17, 0xc1,
	0x49, 0x34, 0xd4, 0x6c, 0x66, 0x63, 0xc7, 0x3b, 0xea, 0x3c, 0xaf, 0xac, 0xe3, 0xcc, 0xda, 0x31,
	0x4f, 0xab, 0x99, 0x42, 0xa9, 0xc8, 0xe7, 0x41, 0x78, 0x5a, 0x26, 0x47, 0x30, 0xed, 0x9a, 0xf1,
	0xc3, 0x48, 0x7e, 0x0c, 0xe0, 0xf0, 0xfd, 0x86, 0xaa, 0xab, 0x6d, 0xef, 0x01, 0x4c, 0x4b, 0x9b,
	0xdc, 0x22, 0xc3, 0x65, 0x8d, 0x1c, 0xdf, 0x5a, 0xbc, 0x65, 0xf5, 0x34, 0x8d, 0xb5, 0xe6, 0x4d,
	0xbb, 0x14, 0xbc, 0x10, 0x17, 0xae, 0xbd, 0x83, 0x45, 0xec, 0xac, 0x67, 0xc1, 0x98, 0x9c, 0x43,
	0x1c, 0x5a, 0x0a, 0x1b, 0xbf, 0x0b, 0xe0, 0x23, 0xb6, 0x96, 0x3e, 0xf1, 0x16, 0xdc, 0xfb, 0x65,
	0xd9, 0x61, 0x9f, 0xac, 0x9d, 0xcc, 0x6b, 0x46, 0x85, 0xd9, 0x74, 0xa3, 0xb2, 0xf7, 0xe2, 0x4e,
	0xce, 0x0b, 0xba, 0x77, 0xf2, 0x0c, 0xa6, 0x1d, 0x29, 0x24, 0xbf, 0x05, 0x91, 0x36, 0xd4, 0xb4,
	0xda, 0xf1, 0xf6, 0x17, 0x01, 0x61, 0x34, 0xaf, 0x57, 0xb2, 0xbb, 0x36, 0x7c, 0x3f, 0xfe, 0x36,
	0x84, 0x11, 0x6e, 0x91, 0xbc, 0x83, 0xc8, 0x5f, 0x2d, 0x99, 0xa7, 0x97, 0xbe, 0xa1, 0x74, 0xe7,
	0x03, 0x3a, 0xbe, 0xf7, 0x17, 0x46, 0xd8, 0xf2, 0x1e, 0x0a, 0xfa, 0xcd, 0xf7, 0x0a, 0xee, 0x5c,
	0xf8, 0xbf, 0x09, 0xbe, 0x81, 0xb1, 0x1b, 0x32, 0x39, 0xe9, 0x61, 0x6f, 0x5f, 0xd4, 0xf1, 0xfc,
	0xcf, 0x84, 0xed, 0xf2, 0xfc, 0xd8, 0x7a, 0xcb, 0xdb, 0x19, 0x7b, 0x6f, 0x79, 0xbb, 0x33, 0x4f,
	0xf6, 0x96, 0x91, 0xfb, 0x03, 0x7a, 0xf2, 0x33, 0x00, 0x00, 0xff, 0xff, 0x19, 0xf5, 0x14, 0x1d,
	0xbd, 0x04, 0x00, 0x00,
}
