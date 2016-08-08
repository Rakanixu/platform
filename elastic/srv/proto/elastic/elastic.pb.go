// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/elastic/srv/proto/elastic/elastic.proto
// DO NOT EDIT!

/*
Package go_micro_srv_elastic is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/elastic/srv/proto/elastic/elastic.proto

It has these top-level messages:
	DocRef
	CreateRequest
	CreateResponse
	BulkCreateRequest
	BulkCreateResponse
	ReadRequest
	ReadResponse
	UpdateRequest
	UpdateResponse
	DeleteRequest
	DeleteResponse
	SearchRequest
	SearchResponse
	QueryRequest
	QueryResponse
	CreateIndexWithSettingsRequest
	CreateIndexWithSettingsResponse
	PutMappingFromJSONRequest
	PutMappingFromJSONResponse
	StatusRequest
	StatusResponse
*/
package go_micro_srv_elastic

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
const _ = proto.ProtoPackageIsVersion1

type DocRef struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *DocRef) Reset()                    { *m = DocRef{} }
func (m *DocRef) String() string            { return proto.CompactTextString(m) }
func (*DocRef) ProtoMessage()               {}
func (*DocRef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CreateRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Data  string `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type CreateResponse struct {
}

func (m *CreateResponse) Reset()                    { *m = CreateResponse{} }
func (m *CreateResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()               {}
func (*CreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type BulkCreateRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Data  string `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
}

func (m *BulkCreateRequest) Reset()                    { *m = BulkCreateRequest{} }
func (m *BulkCreateRequest) String() string            { return proto.CompactTextString(m) }
func (*BulkCreateRequest) ProtoMessage()               {}
func (*BulkCreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type BulkCreateResponse struct {
}

func (m *BulkCreateResponse) Reset()                    { *m = BulkCreateResponse{} }
func (m *BulkCreateResponse) String() string            { return proto.CompactTextString(m) }
func (*BulkCreateResponse) ProtoMessage()               {}
func (*BulkCreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type ReadRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type ReadResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *ReadResponse) Reset()                    { *m = ReadResponse{} }
func (m *ReadResponse) String() string            { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()               {}
func (*ReadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type UpdateRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Data  string `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
}

func (m *UpdateRequest) Reset()                    { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()               {}
func (*UpdateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type UpdateResponse struct {
}

func (m *UpdateResponse) Reset()                    { *m = UpdateResponse{} }
func (m *UpdateResponse) String() string            { return proto.CompactTextString(m) }
func (*UpdateResponse) ProtoMessage()               {}
func (*UpdateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type DeleteRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type DeleteResponse struct {
}

func (m *DeleteResponse) Reset()                    { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()               {}
func (*DeleteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type SearchRequest struct {
	Index  string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type   string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Query  string `protobuf:"bytes,3,opt,name=query" json:"query,omitempty"`
	Limit  int64  `protobuf:"varint,4,opt,name=limit" json:"limit,omitempty"`
	Offset int64  `protobuf:"varint,5,opt,name=offset" json:"offset,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type QueryRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Query string `protobuf:"bytes,3,opt,name=query" json:"query,omitempty"`
}

func (m *QueryRequest) Reset()                    { *m = QueryRequest{} }
func (m *QueryRequest) String() string            { return proto.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()               {}
func (*QueryRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type QueryResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *QueryResponse) Reset()                    { *m = QueryResponse{} }
func (m *QueryResponse) String() string            { return proto.CompactTextString(m) }
func (*QueryResponse) ProtoMessage()               {}
func (*QueryResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

type CreateIndexWithSettingsRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Settings string `protobuf:"bytes,2,opt,name=settings" json:"settings,omitempty"`
}

func (m *CreateIndexWithSettingsRequest) Reset()                    { *m = CreateIndexWithSettingsRequest{} }
func (m *CreateIndexWithSettingsRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateIndexWithSettingsRequest) ProtoMessage()               {}
func (*CreateIndexWithSettingsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type CreateIndexWithSettingsResponse struct {
}

func (m *CreateIndexWithSettingsResponse) Reset()         { *m = CreateIndexWithSettingsResponse{} }
func (m *CreateIndexWithSettingsResponse) String() string { return proto.CompactTextString(m) }
func (*CreateIndexWithSettingsResponse) ProtoMessage()    {}
func (*CreateIndexWithSettingsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{16}
}

type PutMappingFromJSONRequest struct {
	Index   string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type    string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Mapping string `protobuf:"bytes,3,opt,name=mapping" json:"mapping,omitempty"`
}

func (m *PutMappingFromJSONRequest) Reset()                    { *m = PutMappingFromJSONRequest{} }
func (m *PutMappingFromJSONRequest) String() string            { return proto.CompactTextString(m) }
func (*PutMappingFromJSONRequest) ProtoMessage()               {}
func (*PutMappingFromJSONRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

type PutMappingFromJSONResponse struct {
}

func (m *PutMappingFromJSONResponse) Reset()                    { *m = PutMappingFromJSONResponse{} }
func (m *PutMappingFromJSONResponse) String() string            { return proto.CompactTextString(m) }
func (*PutMappingFromJSONResponse) ProtoMessage()               {}
func (*PutMappingFromJSONResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

type StatusResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func init() {
	proto.RegisterType((*DocRef)(nil), "go.micro.srv.elastic.DocRef")
	proto.RegisterType((*CreateRequest)(nil), "go.micro.srv.elastic.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "go.micro.srv.elastic.CreateResponse")
	proto.RegisterType((*BulkCreateRequest)(nil), "go.micro.srv.elastic.BulkCreateRequest")
	proto.RegisterType((*BulkCreateResponse)(nil), "go.micro.srv.elastic.BulkCreateResponse")
	proto.RegisterType((*ReadRequest)(nil), "go.micro.srv.elastic.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "go.micro.srv.elastic.ReadResponse")
	proto.RegisterType((*UpdateRequest)(nil), "go.micro.srv.elastic.UpdateRequest")
	proto.RegisterType((*UpdateResponse)(nil), "go.micro.srv.elastic.UpdateResponse")
	proto.RegisterType((*DeleteRequest)(nil), "go.micro.srv.elastic.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "go.micro.srv.elastic.DeleteResponse")
	proto.RegisterType((*SearchRequest)(nil), "go.micro.srv.elastic.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "go.micro.srv.elastic.SearchResponse")
	proto.RegisterType((*QueryRequest)(nil), "go.micro.srv.elastic.QueryRequest")
	proto.RegisterType((*QueryResponse)(nil), "go.micro.srv.elastic.QueryResponse")
	proto.RegisterType((*CreateIndexWithSettingsRequest)(nil), "go.micro.srv.elastic.CreateIndexWithSettingsRequest")
	proto.RegisterType((*CreateIndexWithSettingsResponse)(nil), "go.micro.srv.elastic.CreateIndexWithSettingsResponse")
	proto.RegisterType((*PutMappingFromJSONRequest)(nil), "go.micro.srv.elastic.PutMappingFromJSONRequest")
	proto.RegisterType((*PutMappingFromJSONResponse)(nil), "go.micro.srv.elastic.PutMappingFromJSONResponse")
	proto.RegisterType((*StatusRequest)(nil), "go.micro.srv.elastic.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "go.micro.srv.elastic.StatusResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Elastic service

type ElasticClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	BulkCreate(ctx context.Context, in *BulkCreateRequest, opts ...client.CallOption) (*BulkCreateResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	Query(ctx context.Context, in *QueryRequest, opts ...client.CallOption) (*QueryResponse, error)
	CreateIndexWithSettings(ctx context.Context, in *CreateIndexWithSettingsRequest, opts ...client.CallOption) (*CreateIndexWithSettingsResponse, error)
	PutMappingFromJSON(ctx context.Context, in *PutMappingFromJSONRequest, opts ...client.CallOption) (*PutMappingFromJSONRequest, error)
	Status(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error)
}

type elasticClient struct {
	c           client.Client
	serviceName string
}

func NewElasticClient(serviceName string, c client.Client) ElasticClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.elastic"
	}
	return &elasticClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *elasticClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) BulkCreate(ctx context.Context, in *BulkCreateRequest, opts ...client.CallOption) (*BulkCreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.BulkCreate", in)
	out := new(BulkCreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Update", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) Query(ctx context.Context, in *QueryRequest, opts ...client.CallOption) (*QueryResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Query", in)
	out := new(QueryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) CreateIndexWithSettings(ctx context.Context, in *CreateIndexWithSettingsRequest, opts ...client.CallOption) (*CreateIndexWithSettingsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.CreateIndexWithSettings", in)
	out := new(CreateIndexWithSettingsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) PutMappingFromJSON(ctx context.Context, in *PutMappingFromJSONRequest, opts ...client.CallOption) (*PutMappingFromJSONRequest, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.PutMappingFromJSON", in)
	out := new(PutMappingFromJSONRequest)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elasticClient) Status(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Elastic.Status", in)
	out := new(StatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Elastic service

type ElasticHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	BulkCreate(context.Context, *BulkCreateRequest, *BulkCreateResponse) error
	Read(context.Context, *ReadRequest, *ReadResponse) error
	Update(context.Context, *UpdateRequest, *UpdateResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	Search(context.Context, *SearchRequest, *SearchResponse) error
	Query(context.Context, *QueryRequest, *QueryResponse) error
	CreateIndexWithSettings(context.Context, *CreateIndexWithSettingsRequest, *CreateIndexWithSettingsResponse) error
	PutMappingFromJSON(context.Context, *PutMappingFromJSONRequest, *PutMappingFromJSONRequest) error
	Status(context.Context, *StatusRequest, *StatusResponse) error
}

func RegisterElasticHandler(s server.Server, hdlr ElasticHandler) {
	s.Handle(s.NewHandler(&Elastic{hdlr}))
}

type Elastic struct {
	ElasticHandler
}

func (h *Elastic) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.ElasticHandler.Create(ctx, in, out)
}

func (h *Elastic) BulkCreate(ctx context.Context, in *BulkCreateRequest, out *BulkCreateResponse) error {
	return h.ElasticHandler.BulkCreate(ctx, in, out)
}

func (h *Elastic) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.ElasticHandler.Read(ctx, in, out)
}

func (h *Elastic) Update(ctx context.Context, in *UpdateRequest, out *UpdateResponse) error {
	return h.ElasticHandler.Update(ctx, in, out)
}

func (h *Elastic) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.ElasticHandler.Delete(ctx, in, out)
}

func (h *Elastic) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.ElasticHandler.Search(ctx, in, out)
}

func (h *Elastic) Query(ctx context.Context, in *QueryRequest, out *QueryResponse) error {
	return h.ElasticHandler.Query(ctx, in, out)
}

func (h *Elastic) CreateIndexWithSettings(ctx context.Context, in *CreateIndexWithSettingsRequest, out *CreateIndexWithSettingsResponse) error {
	return h.ElasticHandler.CreateIndexWithSettings(ctx, in, out)
}

func (h *Elastic) PutMappingFromJSON(ctx context.Context, in *PutMappingFromJSONRequest, out *PutMappingFromJSONRequest) error {
	return h.ElasticHandler.PutMappingFromJSON(ctx, in, out)
}

func (h *Elastic) Status(ctx context.Context, in *StatusRequest, out *StatusResponse) error {
	return h.ElasticHandler.Status(ctx, in, out)
}

var fileDescriptor0 = []byte{
	// 603 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x96, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0xd7, 0xad, 0x1f, 0x70, 0x58, 0x0b, 0x58, 0x15, 0x84, 0x08, 0x01, 0xf3, 0x10, 0xeb,
	0x55, 0x22, 0xf1, 0xf1, 0x02, 0x65, 0x80, 0x86, 0xc4, 0x06, 0xa9, 0x26, 0xae, 0x10, 0xf2, 0x5a,
	0xb7, 0x8b, 0x96, 0xd4, 0x21, 0x76, 0xd0, 0x86, 0xc4, 0x3d, 0xaf, 0xc6, 0x5b, 0x91, 0xd8, 0x4e,
	0xd6, 0x6c, 0x71, 0xaa, 0xb2, 0xed, 0x6a, 0x39, 0x67, 0xc7, 0xbf, 0x63, 0x1f, 0xfb, 0xff, 0x57,
	0x61, 0x38, 0xf3, 0xc5, 0x71, 0x72, 0xe4, 0x8c, 0x59, 0xe8, 0x9e, 0x90, 0x5f, 0x2c, 0x89, 0xdc,
	0x28, 0x20, 0x62, 0xca, 0xe2, 0xd0, 0xa5, 0x01, 0xe1, 0xc2, 0x1f, 0xbb, 0x3c, 0xfe, 0xe9, 0x46,
	0x31, 0x13, 0xac, 0xc8, 0xe8, 0xbf, 0x8e, 0xcc, 0xa2, 0xfe, 0x8c, 0x39, 0xa1, 0x3f, 0x8e, 0x99,
	0x93, 0x56, 0x3a, 0xfa, 0x7f, 0x78, 0x08, 0xed, 0x5d, 0x36, 0xf6, 0xe8, 0x14, 0xf5, 0xa1, 0xe5,
	0xcf, 0x27, 0xf4, 0xd4, 0x6a, 0x3c, 0x6b, 0x0c, 0x6e, 0x7b, 0x2a, 0x40, 0x08, 0x9a, 0xe2, 0x2c,
	0xa2, 0xd6, 0xba, 0x4c, 0xca, 0x6f, 0xd4, 0x83, 0x75, 0x7f, 0x62, 0x6d, 0xc8, 0x4c, 0xfa, 0x85,
	0xbf, 0x41, 0xf7, 0x6d, 0x4c, 0x89, 0xa0, 0x1e, 0xfd, 0x91, 0x50, 0x2e, 0xfe, 0x1f, 0x95, 0xd5,
	0x4c, 0x88, 0x20, 0x56, 0x53, 0xd5, 0x64, 0xdf, 0xf8, 0x1e, 0xf4, 0x72, 0x3c, 0x8f, 0xd8, 0x9c,
	0x53, 0x4c, 0xe0, 0xfe, 0x30, 0x09, 0x4e, 0x6e, 0xb2, 0x69, 0x1f, 0xd0, 0x62, 0x0b, 0xdd, 0xf8,
	0x03, 0xdc, 0xf1, 0x28, 0x99, 0x5c, 0xb9, 0x25, 0x7e, 0x01, 0x9b, 0x0a, 0xa4, 0xc0, 0xe8, 0x01,
	0xb4, 0x63, 0xca, 0x93, 0x40, 0x68, 0x94, 0x8e, 0xb2, 0xd1, 0x1e, 0x46, 0x93, 0x9b, 0x1c, 0x6d,
	0x8e, 0xd7, 0x27, 0xdc, 0x83, 0xee, 0x2e, 0x0d, 0xe8, 0x35, 0x34, 0xcc, 0xe0, 0x39, 0x4a, 0xc3,
	0x7f, 0x43, 0x77, 0x44, 0x49, 0x3c, 0x3e, 0x5e, 0x1d, 0x9e, 0x56, 0xa6, 0x4b, 0xe2, 0x33, 0xcd,
	0x57, 0x41, 0x96, 0x0d, 0xfc, 0xd0, 0x17, 0xf2, 0x50, 0x1b, 0x9e, 0x0a, 0xb2, 0x61, 0xb2, 0xe9,
	0x94, 0x53, 0x61, 0xb5, 0x64, 0x5a, 0x47, 0x78, 0x00, 0xbd, 0xbc, 0xfd, 0x92, 0xb1, 0xef, 0xc3,
	0xe6, 0x97, 0xac, 0xc1, 0x35, 0xed, 0x13, 0xef, 0x40, 0x57, 0xf3, 0x96, 0x34, 0xf6, 0xe0, 0x89,
	0x7a, 0x72, 0x7b, 0x59, 0x87, 0xaf, 0xa9, 0xea, 0x47, 0x54, 0x08, 0x7f, 0x3e, 0xe3, 0xf5, 0x5b,
	0xb1, 0xe1, 0x16, 0xd7, 0x85, 0x7a, 0x3b, 0x45, 0x8c, 0xb7, 0xe0, 0xa9, 0x91, 0xa9, 0x2f, 0xe6,
	0x3b, 0x3c, 0xfa, 0x9c, 0x88, 0x4f, 0x24, 0x8a, 0xd2, 0xf4, 0xfb, 0x98, 0x85, 0x1f, 0x47, 0x07,
	0xfb, 0xab, 0x1f, 0xde, 0x82, 0x4e, 0xa8, 0x18, 0xfa, 0xf8, 0x79, 0x88, 0x1f, 0x83, 0x5d, 0xd5,
	0x40, 0xb7, 0xbf, 0x9b, 0xbe, 0x0b, 0x41, 0x44, 0x92, 0x1f, 0x52, 0xde, 0x94, 0x4e, 0x9c, 0x0f,
	0x8c, 0xcb, 0x4c, 0x3e, 0x30, 0x15, 0xbd, 0xfc, 0xdb, 0x81, 0xce, 0x3b, 0xe5, 0x65, 0xe8, 0x10,
	0xda, 0xea, 0xa0, 0x68, 0xdb, 0xa9, 0x32, 0x3b, 0xa7, 0x64, 0x18, 0xf6, 0xf3, 0xfa, 0x22, 0xbd,
	0xb7, 0x35, 0x44, 0x00, 0xce, 0xad, 0x00, 0xed, 0x54, 0xaf, 0xba, 0xe4, 0x47, 0xf6, 0x60, 0x79,
	0x61, 0xd1, 0xe2, 0x00, 0x9a, 0x99, 0x1d, 0xa0, 0xad, 0xea, 0x35, 0x0b, 0x9e, 0x63, 0xe3, 0xba,
	0x92, 0x02, 0x98, 0x8e, 0x42, 0x09, 0xdb, 0x34, 0x8a, 0x92, 0xab, 0x98, 0x46, 0x71, 0xc1, 0x1b,
	0x24, 0x56, 0x49, 0xda, 0x84, 0x2d, 0x79, 0x87, 0x09, 0x7b, 0xc1, 0x15, 0x24, 0x56, 0x09, 0xd3,
	0x84, 0x2d, 0xb9, 0x86, 0x09, 0x5b, 0xd6, 0x76, 0x8a, 0xf5, 0xa0, 0x25, 0x55, 0x87, 0x0c, 0x33,
	0x5b, 0x94, 0xb8, 0xbd, 0x5d, 0x5b, 0x53, 0x30, 0xff, 0x34, 0xe0, 0xa1, 0x41, 0x4d, 0xe8, 0x75,
	0xdd, 0x83, 0x32, 0x09, 0xda, 0x7e, 0xb3, 0xe2, 0xaa, 0x62, 0x2b, 0xa7, 0x80, 0x2e, 0x6b, 0x0a,
	0xb9, 0xd5, 0x38, 0xa3, 0xbc, 0xed, 0x55, 0x17, 0xe8, 0xfb, 0x92, 0xf2, 0x33, 0xde, 0xd7, 0xa2,
	0x9a, 0x8d, 0xf7, 0x55, 0x52, 0x38, 0x5e, 0x3b, 0x6a, 0xcb, 0x1f, 0x2a, 0xaf, 0xfe, 0x05, 0x00,
	0x00, 0xff, 0xff, 0x79, 0x82, 0x13, 0xae, 0xee, 0x08, 0x00, 0x00,
}
