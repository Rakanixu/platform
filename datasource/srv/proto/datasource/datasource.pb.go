// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/datasource/srv/proto/datasource/datasource.proto
// DO NOT EDIT!

/*
Package proto_datasource is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/datasource/srv/proto/datasource/datasource.proto

It has these top-level messages:
	Endpoint
	CreateRequest
	CreateResponse
	ReadRequest
	ReadResponse
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
*/
package proto_datasource

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

func (m *Endpoint) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Endpoint) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Endpoint) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Endpoint) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *Endpoint) GetLastScan() int64 {
	if m != nil {
		return m.LastScan
	}
	return 0
}

func (m *Endpoint) GetLastScanStarted() int64 {
	if m != nil {
		return m.LastScanStarted
	}
	return 0
}

func (m *Endpoint) GetCrawlerRunning() bool {
	if m != nil {
		return m.CrawlerRunning
	}
	return false
}

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

func (m *CreateResponse) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

type ReadRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ReadRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ReadResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *ReadResponse) Reset()                    { *m = ReadResponse{} }
func (m *ReadResponse) String() string            { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()               {}
func (*ReadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ReadResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

type DeleteRequest struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Index string `protobuf:"bytes,2,opt,name=index" json:"index,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DeleteRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

type DeleteResponse struct {
}

func (m *DeleteResponse) Reset()                    { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()               {}
func (*DeleteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *SearchRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *SearchRequest) GetTerm() string {
	if m != nil {
		return m.Term
	}
	return ""
}

func (m *SearchRequest) GetFrom() int64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *SearchRequest) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *SearchRequest) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *SearchRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *SearchRequest) GetDepth() int64 {
	if m != nil {
		return m.Depth
	}
	return 0
}

func (m *SearchRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *SearchResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func (m *SearchResponse) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type ScanRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *ScanRequest) Reset()                    { *m = ScanRequest{} }
func (m *ScanRequest) String() string            { return proto.CompactTextString(m) }
func (*ScanRequest) ProtoMessage()               {}
func (*ScanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ScanRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ScanResponse struct {
}

func (m *ScanResponse) Reset()                    { *m = ScanResponse{} }
func (m *ScanResponse) String() string            { return proto.CompactTextString(m) }
func (*ScanResponse) ProtoMessage()               {}
func (*ScanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type ScanAllRequest struct {
	DatasourcesId []string `protobuf:"bytes,2,rep,name=datasources_id,json=datasourcesId" json:"datasources_id,omitempty"`
}

func (m *ScanAllRequest) Reset()                    { *m = ScanAllRequest{} }
func (m *ScanAllRequest) String() string            { return proto.CompactTextString(m) }
func (*ScanAllRequest) ProtoMessage()               {}
func (*ScanAllRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *ScanAllRequest) GetDatasourcesId() []string {
	if m != nil {
		return m.DatasourcesId
	}
	return nil
}

type ScanAllResponse struct {
}

func (m *ScanAllResponse) Reset()                    { *m = ScanAllResponse{} }
func (m *ScanAllResponse) String() string            { return proto.CompactTextString(m) }
func (*ScanAllResponse) ProtoMessage()               {}
func (*ScanAllResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type HealthRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *HealthRequest) Reset()                    { *m = HealthRequest{} }
func (m *HealthRequest) String() string            { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()               {}
func (*HealthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

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
func (*HealthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

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

type Token struct {
	AccessToken  string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
	TokenType    string `protobuf:"bytes,2,opt,name=token_type,json=tokenType" json:"token_type,omitempty"`
	RefreshToken string `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
	Expiry       int64  `protobuf:"varint,4,opt,name=expiry" json:"expiry,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *Token) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *Token) GetTokenType() string {
	if m != nil {
		return m.TokenType
	}
	return ""
}

func (m *Token) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

func (m *Token) GetExpiry() int64 {
	if m != nil {
		return m.Expiry
	}
	return 0
}

func init() {
	proto.RegisterType((*Endpoint)(nil), "proto.datasource.Endpoint")
	proto.RegisterType((*CreateRequest)(nil), "proto.datasource.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "proto.datasource.CreateResponse")
	proto.RegisterType((*ReadRequest)(nil), "proto.datasource.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "proto.datasource.ReadResponse")
	proto.RegisterType((*DeleteRequest)(nil), "proto.datasource.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "proto.datasource.DeleteResponse")
	proto.RegisterType((*SearchRequest)(nil), "proto.datasource.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "proto.datasource.SearchResponse")
	proto.RegisterType((*ScanRequest)(nil), "proto.datasource.ScanRequest")
	proto.RegisterType((*ScanResponse)(nil), "proto.datasource.ScanResponse")
	proto.RegisterType((*ScanAllRequest)(nil), "proto.datasource.ScanAllRequest")
	proto.RegisterType((*ScanAllResponse)(nil), "proto.datasource.ScanAllResponse")
	proto.RegisterType((*HealthRequest)(nil), "proto.datasource.HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "proto.datasource.HealthResponse")
	proto.RegisterType((*Token)(nil), "proto.datasource.Token")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Service service

type ServiceClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...client.CallOption) (*ScanResponse, error)
	ScanAll(ctx context.Context, in *ScanAllRequest, opts ...client.CallOption) (*ScanAllResponse, error)
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
		serviceName = "proto.datasource"
	}
	return &serviceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *serviceClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Scan(ctx context.Context, in *ScanRequest, opts ...client.CallOption) (*ScanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.Scan", in)
	out := new(ScanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ScanAll(ctx context.Context, in *ScanAllRequest, opts ...client.CallOption) (*ScanAllResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Service.ScanAll", in)
	out := new(ScanAllResponse)
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
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Read(context.Context, *ReadRequest, *ReadResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	Search(context.Context, *SearchRequest, *SearchResponse) error
	Scan(context.Context, *ScanRequest, *ScanResponse) error
	ScanAll(context.Context, *ScanAllRequest, *ScanAllResponse) error
	Health(context.Context, *HealthRequest, *HealthResponse) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Service{hdlr}, opts...))
}

type Service struct {
	ServiceHandler
}

func (h *Service) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.ServiceHandler.Create(ctx, in, out)
}

func (h *Service) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.ServiceHandler.Read(ctx, in, out)
}

func (h *Service) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.ServiceHandler.Delete(ctx, in, out)
}

func (h *Service) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.ServiceHandler.Search(ctx, in, out)
}

func (h *Service) Scan(ctx context.Context, in *ScanRequest, out *ScanResponse) error {
	return h.ServiceHandler.Scan(ctx, in, out)
}

func (h *Service) ScanAll(ctx context.Context, in *ScanAllRequest, out *ScanAllResponse) error {
	return h.ServiceHandler.ScanAll(ctx, in, out)
}

func (h *Service) Health(ctx context.Context, in *HealthRequest, out *HealthResponse) error {
	return h.ServiceHandler.Health(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/datasource/srv/proto/datasource/datasource.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 731 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xcd, 0x6e, 0xdb, 0x46,
	0x10, 0x36, 0x45, 0x89, 0x92, 0x46, 0x12, 0x6d, 0x2f, 0x0a, 0x9b, 0x50, 0xe1, 0x56, 0xa6, 0xd1,
	0x56, 0x28, 0x5a, 0x09, 0x70, 0xd1, 0xf6, 0x92, 0x4b, 0x90, 0x04, 0x8e, 0x11, 0x18, 0x08, 0x28,
	0xdf, 0x85, 0x35, 0x39, 0xb2, 0x08, 0x53, 0x24, 0xb3, 0xbb, 0x74, 0x2c, 0x3f, 0x40, 0xde, 0x25,
	0x8f, 0x90, 0x77, 0xcb, 0x21, 0xe0, 0xee, 0x92, 0x12, 0x65, 0xda, 0x27, 0xcd, 0x7c, 0x33, 0xfb,
	0x71, 0x7e, 0xbe, 0x11, 0x7c, 0xb8, 0x0d, 0xc5, 0x32, 0xbb, 0x99, 0xf8, 0xc9, 0x6a, 0x7a, 0x47,
	0x1f, 0x93, 0x2c, 0x9d, 0xa6, 0x11, 0x15, 0x8b, 0x84, 0xad, 0xa6, 0x01, 0x15, 0x94, 0x27, 0x19,
	0xf3, 0x71, 0xca, 0xd9, 0xfd, 0x34, 0x65, 0x89, 0x48, 0xb6, 0xc1, 0x8d, 0x39, 0x91, 0x31, 0x72,
	0x20, 0x7f, 0x26, 0x1b, 0xdc, 0xfd, 0x6e, 0x40, 0xe7, 0x5d, 0x1c, 0xa4, 0x49, 0x18, 0x0b, 0x62,
	0x43, 0x23, 0x0c, 0x1c, 0x63, 0x64, 0x8c, 0xbb, 0x5e, 0x23, 0x0c, 0xc8, 0x31, 0xb4, 0x33, 0x8e,
	0x6c, 0x1e, 0x06, 0x4e, 0x43, 0x82, 0x56, 0xee, 0x5e, 0x06, 0xe4, 0x00, 0xcc, 0x8c, 0x45, 0x8e,
	0x29, 0xc1, 0xdc, 0x24, 0x3f, 0x41, 0x2b, 0x8c, 0x03, 0x7c, 0x70, 0x9a, 0x12, 0x53, 0x0e, 0xf9,
	0x19, 0xba, 0x11, 0xe5, 0x62, 0xce, 0x7d, 0x1a, 0x3b, 0xad, 0x91, 0x31, 0x36, 0xbd, 0x4e, 0x0e,
	0xcc, 0x7c, 0x1a, 0x93, 0x3f, 0xe1, 0xb0, 0x0c, 0xce, 0xb9, 0xa0, 0x4c, 0x60, 0xe0, 0x58, 0x32,
	0x69, 0xbf, 0x48, 0x9a, 0x29, 0x98, 0xfc, 0x01, 0xfb, 0x3e, 0xa3, 0x9f, 0x23, 0x64, 0x73, 0x96,
	0xc5, 0x71, 0x18, 0xdf, 0x3a, 0xed, 0x91, 0x31, 0xee, 0x78, 0xb6, 0x86, 0x3d, 0x85, 0x92, 0xbf,
	0xa1, 0x25, 0x92, 0x3b, 0x8c, 0x9d, 0xce, 0xc8, 0x18, 0xf7, 0xce, 0x8f, 0x27, 0xbb, 0x1d, 0x4f,
	0xae, 0xf3, 0xb0, 0xa7, 0xb2, 0xdc, 0x0b, 0x18, 0xbc, 0x61, 0x48, 0x05, 0x7a, 0xf8, 0x29, 0x43,
	0x2e, 0xc8, 0x7f, 0xd0, 0x41, 0x3d, 0x0e, 0x39, 0x88, 0xde, 0xf9, 0xf0, 0x29, 0x45, 0x31, 0x30,
	0xaf, 0xcc, 0x75, 0xff, 0x02, 0xbb, 0x20, 0xe2, 0x69, 0x12, 0x73, 0x24, 0x43, 0xe8, 0x30, 0x6d,
	0xeb, 0x91, 0x96, 0xbe, 0x7b, 0x02, 0x3d, 0x0f, 0x69, 0x50, 0x7c, 0x74, 0x67, 0xee, 0xee, 0xef,
	0xd0, 0x57, 0x61, 0x4d, 0x75, 0x04, 0x16, 0x43, 0x9e, 0x45, 0x42, 0xe7, 0x68, 0xcf, 0xfd, 0x17,
	0x06, 0x6f, 0x31, 0xc2, 0x4d, 0xf5, 0xbb, 0x0b, 0x2c, 0xb7, 0xd2, 0xd8, 0xda, 0x8a, 0x7b, 0x00,
	0x76, 0xf1, 0x4c, 0xd7, 0xf3, 0xcd, 0x80, 0xc1, 0x0c, 0x29, 0xf3, 0x97, 0x05, 0x53, 0xf9, 0xd2,
	0xd8, 0xde, 0x27, 0x81, 0xa6, 0x40, 0xb6, 0xd2, 0x74, 0xd2, 0xce, 0xb1, 0x05, 0x4b, 0x56, 0x52,
	0x0c, 0xa6, 0x27, 0xed, 0x1c, 0xe3, 0xe1, 0x23, 0x4a, 0x31, 0x98, 0x9e, 0xb4, 0xf3, 0x79, 0xf8,
	0x54, 0xe0, 0x6d, 0xc2, 0xd6, 0x52, 0x0a, 0x5d, 0xaf, 0xf4, 0x0b, 0x3d, 0x59, 0x15, 0x3d, 0x05,
	0x98, 0x8a, 0xa5, 0x5c, 0xb3, 0xe9, 0x29, 0x47, 0x7e, 0x7f, 0x9d, 0xa2, 0x5c, 0x6e, 0xfe, 0xfd,
	0x75, 0x8a, 0xee, 0x2b, 0xb0, 0x8b, 0xd2, 0x5f, 0x1e, 0x57, 0xfe, 0x3a, 0x8c, 0x17, 0x49, 0x51,
	0x7d, 0x6e, 0xe7, 0x9b, 0xc8, 0x75, 0xf6, 0xdc, 0x26, 0x6c, 0xe8, 0xab, 0xb0, 0x1e, 0xd4, 0xff,
	0x60, 0xe7, 0xfe, 0xeb, 0x28, 0x2a, 0x5e, 0xfc, 0x06, 0xf6, 0x46, 0x19, 0x5c, 0x9d, 0x8a, 0x39,
	0xee, 0x7a, 0x83, 0x2d, 0xf4, 0x32, 0x70, 0x0f, 0x61, 0xbf, 0x7c, 0xa8, 0xb9, 0xce, 0x60, 0xf0,
	0x1e, 0x69, 0x24, 0xca, 0x99, 0x17, 0xdd, 0x19, 0xd5, 0xee, 0x8a, 0xa4, 0x4d, 0x77, 0x5c, 0x50,
	0x91, 0x71, 0x99, 0x67, 0x7a, 0xda, 0xab, 0xed, 0xee, 0x8b, 0x01, 0x2d, 0xa9, 0x77, 0x72, 0x0a,
	0x7d, 0xea, 0xfb, 0xc8, 0xf9, 0x5c, 0x9d, 0x87, 0xfa, 0x46, 0x4f, 0x61, 0x2a, 0xe5, 0x04, 0x40,
	0xc6, 0xe6, 0xb2, 0x08, 0x45, 0xd3, 0x95, 0xc8, 0xf5, 0x3a, 0x45, 0x72, 0x06, 0x03, 0x86, 0x0b,
	0x86, 0x7c, 0xa9, 0x29, 0xd4, 0xf5, 0xf7, 0x35, 0xa8, 0x38, 0x8e, 0xc0, 0xc2, 0x87, 0x34, 0x64,
	0x6b, 0xbd, 0x7a, 0xed, 0x9d, 0x7f, 0x6d, 0x42, 0x7b, 0x86, 0xec, 0x3e, 0xf4, 0x91, 0x5c, 0x81,
	0xa5, 0x4e, 0x85, 0xfc, 0xfa, 0xf4, 0xb4, 0x2a, 0xd7, 0x38, 0x1c, 0x3d, 0x9f, 0xa0, 0x87, 0xb8,
	0x47, 0x2e, 0xa0, 0x99, 0x1f, 0x0b, 0x39, 0x79, 0x9a, 0xbb, 0x75, 0x63, 0xc3, 0x5f, 0x9e, 0x0b,
	0x97, 0x44, 0x57, 0x60, 0xa9, 0xb3, 0xa8, 0xab, 0xab, 0x72, 0x67, 0x75, 0x75, 0xed, 0x5c, 0x94,
	0xa4, 0x53, 0xba, 0xac, 0xa3, 0xab, 0x1c, 0x5b, 0x1d, 0x5d, 0x55, 0xd2, 0xaa, 0x4d, 0xf9, 0xaf,
	0x59, 0xd3, 0xe6, 0x96, 0x80, 0xeb, 0xda, 0xac, 0x08, 0x78, 0x8f, 0x7c, 0x84, 0xb6, 0x56, 0x22,
	0x19, 0xd5, 0x27, 0x6f, 0xd4, 0x3d, 0x3c, 0x7d, 0x21, 0x63, 0xbb, 0x53, 0xa5, 0xd1, 0xba, 0x4e,
	0x2b, 0x12, 0xaf, 0xeb, 0xb4, 0x2a, 0x6f, 0x77, 0xef, 0xc6, 0x92, 0x29, 0xff, 0xfc, 0x08, 0x00,
	0x00, 0xff, 0xff, 0xd1, 0x23, 0xd8, 0xf5, 0xfa, 0x06, 0x00, 0x00,
}
