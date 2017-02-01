// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/datasource/srv/proto/datasource/datasource.proto
// DO NOT EDIT!

/*
Package go_micro_srv_datasource is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/datasource/srv/proto/datasource/datasource.proto

It has these top-level messages:
	Endpoint
	CreateRequest
	CreateResponse
	DeleteRequest
	DeleteResponse
	SearchRequest
	SearchResponse
	ScanRequest
	ScanResponse
	ScanAllRequest
	ScanAllResponse
	HealthRequest
	HealthResponse
	Token
	DeleteBucketMessage
	DeleteFileInBucketMessage
*/
package go_micro_srv_datasource

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

type Endpoint struct {
	Id              string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	UserId          string `protobuf:"bytes,2,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	Url             string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	Index           string `protobuf:"bytes,4,opt,name=index" json:"index,omitempty"`
	LastScan        int64  `protobuf:"varint,5,opt,name=last_scan,json=lastScan" json:"last_scan,omitempty"`
	LastScanStarted int64  `protobuf:"varint,6,opt,name=last_scan_started,json=lastScanStarted" json:"last_scan_started,omitempty"`
	CrawlerRunning  bool   `protobuf:"varint,7,opt,name=crawler_running,json=crawlerRunning" json:"crawler_running,omitempty"`
	Token           *Token `protobuf:"bytes,8,opt,name=token" json:"token,omitempty"`
}

func (m *Endpoint) Reset()                    { *m = Endpoint{} }
func (m *Endpoint) String() string            { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()               {}
func (*Endpoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Endpoint) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

type CreateRequest struct {
	Endpoint *Endpoint `protobuf:"bytes,1,opt,name=endpoint" json:"endpoint,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateRequest) GetEndpoint() *Endpoint {
	if m != nil {
		return m.Endpoint
	}
	return nil
}

type CreateResponse struct {
	Response string `protobuf:"bytes,1,opt,name=response" json:"response,omitempty"`
}

func (m *CreateResponse) Reset()                    { *m = CreateResponse{} }
func (m *CreateResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()               {}
func (*CreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type DeleteRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type DeleteResponse struct {
}

func (m *DeleteResponse) Reset()                    { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()               {}
func (*DeleteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type SearchRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Term     string `protobuf:"bytes,2,opt,name=term" json:"term,omitempty"`
	From     int64  `protobuf:"varint,3,opt,name=from" json:"from,omitempty"`
	Size     int64  `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
	Category string `protobuf:"bytes,5,opt,name=category" json:"category,omitempty"`
	Url      string `protobuf:"bytes,6,opt,name=url" json:"url,omitempty"`
	Depth    int64  `protobuf:"varint,7,opt,name=depth" json:"depth,omitempty"`
	Type     string `protobuf:"bytes,8,opt,name=type" json:"type,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type ScanRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *ScanRequest) Reset()                    { *m = ScanRequest{} }
func (m *ScanRequest) String() string            { return proto.CompactTextString(m) }
func (*ScanRequest) ProtoMessage()               {}
func (*ScanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type ScanResponse struct {
}

func (m *ScanResponse) Reset()                    { *m = ScanResponse{} }
func (m *ScanResponse) String() string            { return proto.CompactTextString(m) }
func (*ScanResponse) ProtoMessage()               {}
func (*ScanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type ScanAllRequest struct {
	UserId        string   `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	DatasourcesId []string `protobuf:"bytes,2,rep,name=datasources_id,json=datasourcesId" json:"datasources_id,omitempty"`
}

func (m *ScanAllRequest) Reset()                    { *m = ScanAllRequest{} }
func (m *ScanAllRequest) String() string            { return proto.CompactTextString(m) }
func (*ScanAllRequest) ProtoMessage()               {}
func (*ScanAllRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type ScanAllResponse struct {
}

func (m *ScanAllResponse) Reset()                    { *m = ScanAllResponse{} }
func (m *ScanAllResponse) String() string            { return proto.CompactTextString(m) }
func (*ScanAllResponse) ProtoMessage()               {}
func (*ScanAllResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type HealthResponse struct {
	Status int64  `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *HealthResponse) Reset()                    { *m = HealthResponse{} }
func (m *HealthResponse) String() string            { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()               {}
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type Token struct {
	AccessToken  string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
	TokenType    string `protobuf:"bytes,2,opt,name=token_type,json=tokenType" json:"token_type,omitempty"`
	RefreshToken string `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
	Expiry       int64  `protobuf:"varint,4,opt,name=expiry" json:"expiry,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type DeleteBucketMessage struct {
	Endpoint *Endpoint `protobuf:"bytes,1,opt,name=endpoint" json:"endpoint,omitempty"`
}

func (m *DeleteBucketMessage) Reset()                    { *m = DeleteBucketMessage{} }
func (m *DeleteBucketMessage) String() string            { return proto.CompactTextString(m) }
func (*DeleteBucketMessage) ProtoMessage()               {}
func (*DeleteBucketMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *DeleteBucketMessage) GetEndpoint() *Endpoint {
	if m != nil {
		return m.Endpoint
	}
	return nil
}

type DeleteFileInBucketMessage struct {
	FileId string `protobuf:"bytes,1,opt,name=file_id,json=fileId" json:"file_id,omitempty"`
	Index  string `protobuf:"bytes,2,opt,name=index" json:"index,omitempty"`
}

func (m *DeleteFileInBucketMessage) Reset()                    { *m = DeleteFileInBucketMessage{} }
func (m *DeleteFileInBucketMessage) String() string            { return proto.CompactTextString(m) }
func (*DeleteFileInBucketMessage) ProtoMessage()               {}
func (*DeleteFileInBucketMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func init() {
	proto.RegisterType((*Endpoint)(nil), "go.micro.srv.datasource.Endpoint")
	proto.RegisterType((*CreateRequest)(nil), "go.micro.srv.datasource.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "go.micro.srv.datasource.CreateResponse")
	proto.RegisterType((*DeleteRequest)(nil), "go.micro.srv.datasource.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "go.micro.srv.datasource.DeleteResponse")
	proto.RegisterType((*SearchRequest)(nil), "go.micro.srv.datasource.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "go.micro.srv.datasource.SearchResponse")
	proto.RegisterType((*ScanRequest)(nil), "go.micro.srv.datasource.ScanRequest")
	proto.RegisterType((*ScanResponse)(nil), "go.micro.srv.datasource.ScanResponse")
	proto.RegisterType((*ScanAllRequest)(nil), "go.micro.srv.datasource.ScanAllRequest")
	proto.RegisterType((*ScanAllResponse)(nil), "go.micro.srv.datasource.ScanAllResponse")
	proto.RegisterType((*HealthRequest)(nil), "go.micro.srv.datasource.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "go.micro.srv.datasource.HealthResponse")
	proto.RegisterType((*Token)(nil), "go.micro.srv.datasource.Token")
	proto.RegisterType((*DeleteBucketMessage)(nil), "go.micro.srv.datasource.DeleteBucketMessage")
	proto.RegisterType((*DeleteFileInBucketMessage)(nil), "go.micro.srv.datasource.DeleteFileInBucketMessage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for DataSource service

type DataSourceClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...client.CallOption) (*ScanResponse, error)
	ScanAll(ctx context.Context, in *ScanAllRequest, opts ...client.CallOption) (*ScanAllResponse, error)
	Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error)
}

type dataSourceClient struct {
	c           client.Client
	serviceName string
}

func NewDataSourceClient(serviceName string, c client.Client) DataSourceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.datasource"
	}
	return &dataSourceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *dataSourceClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DataSource.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DataSource.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DataSource.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Scan(ctx context.Context, in *ScanRequest, opts ...client.CallOption) (*ScanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DataSource.Scan", in)
	out := new(ScanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) ScanAll(ctx context.Context, in *ScanAllRequest, opts ...client.CallOption) (*ScanAllResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DataSource.ScanAll", in)
	out := new(ScanAllResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceClient) Health(ctx context.Context, in *HealthRequest, opts ...client.CallOption) (*HealthResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DataSource.Health", in)
	out := new(HealthResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DataSource service

type DataSourceHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	Search(context.Context, *SearchRequest, *SearchResponse) error
	Scan(context.Context, *ScanRequest, *ScanResponse) error
	ScanAll(context.Context, *ScanAllRequest, *ScanAllResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterDataSourceHandler(s server.Server, hdlr DataSourceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&DataSource{hdlr}, opts...))
}

type DataSource struct {
	DataSourceHandler
}

func (h *DataSource) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.DataSourceHandler.Create(ctx, in, out)
}

func (h *DataSource) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.DataSourceHandler.Delete(ctx, in, out)
}

func (h *DataSource) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.DataSourceHandler.Search(ctx, in, out)
}

func (h *DataSource) Scan(ctx context.Context, in *ScanRequest, out *ScanResponse) error {
	return h.DataSourceHandler.Scan(ctx, in, out)
}

func (h *DataSource) ScanAll(ctx context.Context, in *ScanAllRequest, out *ScanAllResponse) error {
	return h.DataSourceHandler.ScanAll(ctx, in, out)
}

func (h *DataSource) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.DataSourceHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/datasource/srv/proto/datasource/datasource.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 748 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x55, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0x6e, 0x3e, 0x9b, 0x4e, 0x9a, 0xb4, 0xdd, 0xf7, 0x15, 0x0d, 0x41, 0x05, 0x6a, 0x28, 0xad,
	0x10, 0x4a, 0xa4, 0xc2, 0x11, 0x0e, 0x40, 0x41, 0x14, 0x04, 0x42, 0x4e, 0x39, 0x21, 0x11, 0xb9,
	0xf6, 0x26, 0xb1, 0xea, 0x78, 0xcd, 0x7a, 0x0d, 0x6d, 0xaf, 0x48, 0xfc, 0x25, 0xc4, 0xbf, 0x63,
	0xf6, 0xc3, 0x8e, 0x5d, 0x11, 0xf7, 0xc0, 0x6d, 0xe6, 0xd9, 0xd9, 0x67, 0x66, 0xe7, 0x99, 0xb1,
	0xe1, 0xdd, 0xd4, 0x17, 0xb3, 0xe4, 0x74, 0xe0, 0xb2, 0xf9, 0xf0, 0xcc, 0xb9, 0x64, 0x49, 0x34,
	0x8c, 0x02, 0x47, 0x4c, 0x18, 0x9f, 0x0f, 0x3d, 0x47, 0x38, 0x31, 0x4b, 0xb8, 0x4b, 0x87, 0x31,
	0xff, 0x36, 0x8c, 0x38, 0x13, 0x2c, 0x0f, 0x2e, 0xcc, 0x81, 0x3a, 0x23, 0xdb, 0x53, 0x36, 0x98,
	0xfb, 0x2e, 0x67, 0x03, 0x8c, 0x1f, 0x2c, 0x8e, 0xad, 0x1f, 0x55, 0x68, 0xbd, 0x0a, 0xbd, 0x88,
	0xf9, 0xa1, 0x20, 0x5d, 0xa8, 0xfa, 0x5e, 0xaf, 0x72, 0xb7, 0x72, 0xb0, 0x66, 0xa3, 0x45, 0xb6,
	0x61, 0x35, 0x89, 0x29, 0x1f, 0x23, 0x58, 0x55, 0x60, 0x53, 0xba, 0xc7, 0x1e, 0xd9, 0x84, 0x5a,
	0xc2, 0x83, 0x5e, 0x4d, 0x81, 0xd2, 0x24, 0xff, 0x43, 0xc3, 0x0f, 0x3d, 0x7a, 0xde, 0xab, 0x2b,
	0x4c, 0x3b, 0xe4, 0x16, 0xac, 0x05, 0x4e, 0x2c, 0xc6, 0xb1, 0xeb, 0x84, 0xbd, 0x06, 0x9e, 0xd4,
	0xec, 0x96, 0x04, 0x46, 0xe8, 0x93, 0x87, 0xb0, 0x95, 0x1d, 0x8e, 0x63, 0xe1, 0x70, 0x41, 0xbd,
	0x5e, 0x53, 0x05, 0x6d, 0xa4, 0x41, 0x23, 0x0d, 0x93, 0x7d, 0xd8, 0x70, 0xb9, 0xf3, 0x3d, 0xc0,
	0x62, 0x78, 0x12, 0x86, 0x7e, 0x38, 0xed, 0xad, 0x62, 0x64, 0xcb, 0xee, 0x1a, 0xd8, 0xd6, 0x28,
	0x79, 0x02, 0x0d, 0xc1, 0xce, 0x68, 0xd8, 0x6b, 0xe1, 0x71, 0xfb, 0xf0, 0xf6, 0x60, 0xc9, 0xc3,
	0x07, 0x27, 0x32, 0xca, 0xd6, 0xc1, 0xd6, 0x07, 0xe8, 0xbc, 0xe4, 0xd4, 0x11, 0xd4, 0xa6, 0x5f,
	0x13, 0x1a, 0x0b, 0xf2, 0x0c, 0x5a, 0xd4, 0x74, 0x45, 0xf5, 0xa3, 0x7d, 0xb8, 0xbb, 0x94, 0x29,
	0x6d, 0x9f, 0x9d, 0x5d, 0xb1, 0x1e, 0x41, 0x37, 0xe5, 0x8b, 0x23, 0x16, 0xc6, 0x94, 0xf4, 0xa1,
	0xc5, 0x8d, 0x6d, 0x1a, 0x9c, 0xf9, 0xd6, 0x1d, 0xe8, 0x1c, 0xd1, 0x80, 0x2e, 0xb2, 0x5f, 0xd1,
	0xc1, 0xda, 0x84, 0x6e, 0x1a, 0x60, 0xae, 0xfc, 0xae, 0x40, 0x67, 0x44, 0x1d, 0xee, 0xce, 0xd2,
	0x3b, 0x99, 0x00, 0x95, 0xbc, 0x00, 0x04, 0xea, 0x82, 0xf2, 0xb9, 0x91, 0x4f, 0xd9, 0x12, 0x9b,
	0x70, 0x36, 0x57, 0xea, 0xd5, 0x6c, 0x65, 0x4b, 0x2c, 0xf6, 0x2f, 0xa9, 0x52, 0x0f, 0x31, 0x69,
	0xcb, 0x92, 0x5d, 0x7c, 0xc2, 0x94, 0xf1, 0x0b, 0xa5, 0x1d, 0x96, 0x9c, 0xfa, 0xe9, 0x00, 0x34,
	0x0b, 0x03, 0xe0, 0xd1, 0x48, 0xcc, 0x94, 0x2e, 0x35, 0x5b, 0x3b, 0x2a, 0xff, 0x45, 0x44, 0x95,
	0x1a, 0x32, 0x3f, 0xda, 0xd6, 0x53, 0xe8, 0xa6, 0xa5, 0x9b, 0xe6, 0xdc, 0x80, 0x26, 0x36, 0x23,
	0x09, 0x84, 0x29, 0xde, 0x78, 0xf2, 0xb6, 0x1f, 0x4e, 0x58, 0x5a, 0xbd, 0xb4, 0xad, 0x1d, 0x68,
	0xcb, 0xc1, 0x58, 0xd6, 0xaa, 0x2e, 0xac, 0xeb, 0x63, 0xd3, 0xa8, 0x8f, 0x98, 0x0c, 0xfd, 0xe7,
	0x41, 0x90, 0xde, 0xc8, 0x0d, 0x75, 0xa5, 0x30, 0xd4, 0x7b, 0xd0, 0x5d, 0xa8, 0x1a, 0xeb, 0xa1,
	0xaf, 0xe1, 0x79, 0x27, 0x87, 0x1e, 0x7b, 0xd6, 0x16, 0x6c, 0x64, 0x8c, 0x26, 0xc9, 0x3d, 0xe8,
	0xbc, 0xa1, 0x4e, 0x20, 0x32, 0x31, 0xd2, 0x67, 0x57, 0x8a, 0xcf, 0x4e, 0x83, 0x16, 0xcf, 0xc6,
	0xb1, 0x17, 0x49, 0xac, 0xe2, 0x6a, 0xb6, 0xf1, 0xfe, 0xfa, 0xec, 0x9f, 0x15, 0x68, 0xa8, 0x91,
	0x25, 0xbb, 0xb0, 0xee, 0xb8, 0x58, 0x4a, 0x3c, 0xd6, 0x83, 0xae, 0x73, 0xb4, 0x35, 0xa6, 0x43,
	0x76, 0x00, 0xd4, 0xd9, 0x58, 0x15, 0xa1, 0x69, 0xd6, 0x14, 0x72, 0x82, 0x00, 0xc1, 0x72, 0x39,
	0x9d, 0x60, 0x8f, 0x67, 0x86, 0x42, 0xef, 0xf1, 0xba, 0x01, 0x35, 0x07, 0x16, 0x47, 0xcf, 0x23,
	0x1f, 0xb5, 0xd7, 0x33, 0x61, 0x3c, 0xeb, 0x04, 0xfe, 0xd3, 0xb3, 0xf8, 0x22, 0x71, 0xcf, 0xa8,
	0x78, 0x8f, 0x49, 0x9d, 0x29, 0xfd, 0xd7, 0x85, 0x79, 0x0b, 0x37, 0x35, 0xeb, 0x6b, 0x3f, 0xa0,
	0xc7, 0x61, 0x91, 0x1b, 0x15, 0x9b, 0x20, 0x9c, 0x53, 0x4c, 0xba, 0xa8, 0x58, 0x36, 0xf3, 0xd5,
	0xdc, 0xcc, 0x1f, 0xfe, 0xaa, 0x03, 0x1c, 0x61, 0xb6, 0x91, 0xca, 0x46, 0x3e, 0x43, 0x53, 0xef,
	0x22, 0x79, 0xb0, 0xb4, 0xa2, 0xc2, 0xf2, 0xf7, 0xf7, 0xaf, 0x8d, 0x33, 0xba, 0xaf, 0x48, 0x72,
	0x5d, 0x77, 0x09, 0x79, 0x61, 0xb7, 0x4b, 0xc8, 0xaf, 0xac, 0xb8, 0x22, 0xd7, 0x8b, 0x52, 0x42,
	0x5e, 0xf8, 0x08, 0x94, 0x90, 0x17, 0x37, 0x0e, 0xc9, 0x3f, 0x41, 0x5d, 0x7d, 0x85, 0xef, 0x2f,
	0xbf, 0xb2, 0x58, 0xb3, 0xfe, 0xde, 0x35, 0x51, 0x19, 0xed, 0x17, 0x58, 0x35, 0xdb, 0x41, 0xf6,
	0x4b, 0xef, 0x2c, 0x36, 0xb2, 0x7f, 0x70, 0x7d, 0x60, 0xbe, 0x27, 0x7a, 0x8b, 0x4a, 0x7a, 0x52,
	0xd8, 0xc5, 0x92, 0x9e, 0x14, 0xd7, 0xd1, 0x5a, 0x39, 0x6d, 0xaa, 0x9f, 0xe5, 0xe3, 0x3f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x1b, 0xe2, 0x55, 0x25, 0x7b, 0x07, 0x00, 0x00,
}
