// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/db/srv/proto/db/db.proto
// DO NOT EDIT!

/*
Package go_micro_srv_db is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/db/srv/proto/db/db.proto

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
	DeleteByQueryRequest
	DeleteByQueryResponse
	CreateIndexWithSettingsRequest
	CreateIndexWithSettingsResponse
	PutMappingFromJSONRequest
	PutMappingFromJSONResponse
	StatusRequest
	StatusResponse
	SearchRequest
	SearchResponse
	SearchByIdRequest
	SearchByIdResponse
	AddAliasRequest
	AddAliasResponse
	DeleteIndexRequest
	DeleteIndexResponse
	DeleteAliasRequest
	DeleteAliasResponse
	RenameAliasRequest
	RenameAliasResponse
*/
package go_micro_srv_db

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

type DeleteByQueryRequest struct {
	Indexes  []string `protobuf:"bytes,1,rep,name=indexes" json:"indexes,omitempty"`
	Types    []string `protobuf:"bytes,2,rep,name=types" json:"types,omitempty"`
	Term     string   `protobuf:"bytes,3,opt,name=term" json:"term,omitempty"`
	Category string   `protobuf:"bytes,4,opt,name=category" json:"category,omitempty"`
	Url      string   `protobuf:"bytes,5,opt,name=url" json:"url,omitempty"`
	Depth    int64    `protobuf:"varint,6,opt,name=depth" json:"depth,omitempty"`
	FileType string   `protobuf:"bytes,7,opt,name=file_type,json=fileType" json:"file_type,omitempty"`
	LastSeen int64    `protobuf:"varint,8,opt,name=last_seen,json=lastSeen" json:"last_seen,omitempty"`
}

func (m *DeleteByQueryRequest) Reset()                    { *m = DeleteByQueryRequest{} }
func (m *DeleteByQueryRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteByQueryRequest) ProtoMessage()               {}
func (*DeleteByQueryRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type DeleteByQueryResponse struct {
}

func (m *DeleteByQueryResponse) Reset()                    { *m = DeleteByQueryResponse{} }
func (m *DeleteByQueryResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteByQueryResponse) ProtoMessage()               {}
func (*DeleteByQueryResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type CreateIndexWithSettingsRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Settings string `protobuf:"bytes,2,opt,name=settings" json:"settings,omitempty"`
}

func (m *CreateIndexWithSettingsRequest) Reset()                    { *m = CreateIndexWithSettingsRequest{} }
func (m *CreateIndexWithSettingsRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateIndexWithSettingsRequest) ProtoMessage()               {}
func (*CreateIndexWithSettingsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type CreateIndexWithSettingsResponse struct {
}

func (m *CreateIndexWithSettingsResponse) Reset()         { *m = CreateIndexWithSettingsResponse{} }
func (m *CreateIndexWithSettingsResponse) String() string { return proto.CompactTextString(m) }
func (*CreateIndexWithSettingsResponse) ProtoMessage()    {}
func (*CreateIndexWithSettingsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{14}
}

type PutMappingFromJSONRequest struct {
	Index   string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type    string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Mapping string `protobuf:"bytes,3,opt,name=mapping" json:"mapping,omitempty"`
}

func (m *PutMappingFromJSONRequest) Reset()                    { *m = PutMappingFromJSONRequest{} }
func (m *PutMappingFromJSONRequest) String() string            { return proto.CompactTextString(m) }
func (*PutMappingFromJSONRequest) ProtoMessage()               {}
func (*PutMappingFromJSONRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type PutMappingFromJSONResponse struct {
}

func (m *PutMappingFromJSONResponse) Reset()                    { *m = PutMappingFromJSONResponse{} }
func (m *PutMappingFromJSONResponse) String() string            { return proto.CompactTextString(m) }
func (*PutMappingFromJSONResponse) ProtoMessage()               {}
func (*PutMappingFromJSONResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

type StatusResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

type SearchRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Term     string `protobuf:"bytes,2,opt,name=term" json:"term,omitempty"`
	From     int64  `protobuf:"varint,3,opt,name=from" json:"from,omitempty"`
	Size     int64  `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
	Category string `protobuf:"bytes,5,opt,name=category" json:"category,omitempty"`
	Url      string `protobuf:"bytes,6,opt,name=url" json:"url,omitempty"`
	Depth    int64  `protobuf:"varint,7,opt,name=depth" json:"depth,omitempty"`
	Type     string `protobuf:"bytes,8,opt,name=type" json:"type,omitempty"`
	FileType string `protobuf:"bytes,9,opt,name=file_type,json=fileType" json:"file_type,omitempty"`
	UserId   string `protobuf:"bytes,10,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	LastSeen int64  `protobuf:"varint,11,opt,name=last_seen,json=lastSeen" json:"last_seen,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

type SearchByIdRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *SearchByIdRequest) Reset()                    { *m = SearchByIdRequest{} }
func (m *SearchByIdRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchByIdRequest) ProtoMessage()               {}
func (*SearchByIdRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

type SearchByIdResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *SearchByIdResponse) Reset()                    { *m = SearchByIdResponse{} }
func (m *SearchByIdResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchByIdResponse) ProtoMessage()               {}
func (*SearchByIdResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{22} }

type AddAliasRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Alias string `protobuf:"bytes,2,opt,name=alias" json:"alias,omitempty"`
}

func (m *AddAliasRequest) Reset()                    { *m = AddAliasRequest{} }
func (m *AddAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*AddAliasRequest) ProtoMessage()               {}
func (*AddAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{23} }

type AddAliasResponse struct {
}

func (m *AddAliasResponse) Reset()                    { *m = AddAliasResponse{} }
func (m *AddAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*AddAliasResponse) ProtoMessage()               {}
func (*AddAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{24} }

type DeleteIndexRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
}

func (m *DeleteIndexRequest) Reset()                    { *m = DeleteIndexRequest{} }
func (m *DeleteIndexRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexRequest) ProtoMessage()               {}
func (*DeleteIndexRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{25} }

type DeleteIndexResponse struct {
}

func (m *DeleteIndexResponse) Reset()                    { *m = DeleteIndexResponse{} }
func (m *DeleteIndexResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexResponse) ProtoMessage()               {}
func (*DeleteIndexResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{26} }

type DeleteAliasRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Alias string `protobuf:"bytes,2,opt,name=alias" json:"alias,omitempty"`
}

func (m *DeleteAliasRequest) Reset()                    { *m = DeleteAliasRequest{} }
func (m *DeleteAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteAliasRequest) ProtoMessage()               {}
func (*DeleteAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{27} }

type DeleteAliasResponse struct {
}

func (m *DeleteAliasResponse) Reset()                    { *m = DeleteAliasResponse{} }
func (m *DeleteAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteAliasResponse) ProtoMessage()               {}
func (*DeleteAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{28} }

type RenameAliasRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	OldAlias string `protobuf:"bytes,2,opt,name=old_alias,json=oldAlias" json:"old_alias,omitempty"`
	NewAlias string `protobuf:"bytes,3,opt,name=new_alias,json=newAlias" json:"new_alias,omitempty"`
}

func (m *RenameAliasRequest) Reset()                    { *m = RenameAliasRequest{} }
func (m *RenameAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*RenameAliasRequest) ProtoMessage()               {}
func (*RenameAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{29} }

type RenameAliasResponse struct {
}

func (m *RenameAliasResponse) Reset()                    { *m = RenameAliasResponse{} }
func (m *RenameAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*RenameAliasResponse) ProtoMessage()               {}
func (*RenameAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{30} }

func init() {
	proto.RegisterType((*DocRef)(nil), "go.micro.srv.db.DocRef")
	proto.RegisterType((*CreateRequest)(nil), "go.micro.srv.db.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "go.micro.srv.db.CreateResponse")
	proto.RegisterType((*BulkCreateRequest)(nil), "go.micro.srv.db.BulkCreateRequest")
	proto.RegisterType((*BulkCreateResponse)(nil), "go.micro.srv.db.BulkCreateResponse")
	proto.RegisterType((*ReadRequest)(nil), "go.micro.srv.db.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "go.micro.srv.db.ReadResponse")
	proto.RegisterType((*UpdateRequest)(nil), "go.micro.srv.db.UpdateRequest")
	proto.RegisterType((*UpdateResponse)(nil), "go.micro.srv.db.UpdateResponse")
	proto.RegisterType((*DeleteRequest)(nil), "go.micro.srv.db.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "go.micro.srv.db.DeleteResponse")
	proto.RegisterType((*DeleteByQueryRequest)(nil), "go.micro.srv.db.DeleteByQueryRequest")
	proto.RegisterType((*DeleteByQueryResponse)(nil), "go.micro.srv.db.DeleteByQueryResponse")
	proto.RegisterType((*CreateIndexWithSettingsRequest)(nil), "go.micro.srv.db.CreateIndexWithSettingsRequest")
	proto.RegisterType((*CreateIndexWithSettingsResponse)(nil), "go.micro.srv.db.CreateIndexWithSettingsResponse")
	proto.RegisterType((*PutMappingFromJSONRequest)(nil), "go.micro.srv.db.PutMappingFromJSONRequest")
	proto.RegisterType((*PutMappingFromJSONResponse)(nil), "go.micro.srv.db.PutMappingFromJSONResponse")
	proto.RegisterType((*StatusRequest)(nil), "go.micro.srv.db.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "go.micro.srv.db.StatusResponse")
	proto.RegisterType((*SearchRequest)(nil), "go.micro.srv.db.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "go.micro.srv.db.SearchResponse")
	proto.RegisterType((*SearchByIdRequest)(nil), "go.micro.srv.db.SearchByIdRequest")
	proto.RegisterType((*SearchByIdResponse)(nil), "go.micro.srv.db.SearchByIdResponse")
	proto.RegisterType((*AddAliasRequest)(nil), "go.micro.srv.db.AddAliasRequest")
	proto.RegisterType((*AddAliasResponse)(nil), "go.micro.srv.db.AddAliasResponse")
	proto.RegisterType((*DeleteIndexRequest)(nil), "go.micro.srv.db.DeleteIndexRequest")
	proto.RegisterType((*DeleteIndexResponse)(nil), "go.micro.srv.db.DeleteIndexResponse")
	proto.RegisterType((*DeleteAliasRequest)(nil), "go.micro.srv.db.DeleteAliasRequest")
	proto.RegisterType((*DeleteAliasResponse)(nil), "go.micro.srv.db.DeleteAliasResponse")
	proto.RegisterType((*RenameAliasRequest)(nil), "go.micro.srv.db.RenameAliasRequest")
	proto.RegisterType((*RenameAliasResponse)(nil), "go.micro.srv.db.RenameAliasResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for DB service

type DBClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	DeleteByQuery(ctx context.Context, in *DeleteByQueryRequest, opts ...client.CallOption) (*DeleteByQueryResponse, error)
	CreateIndexWithSettings(ctx context.Context, in *CreateIndexWithSettingsRequest, opts ...client.CallOption) (*CreateIndexWithSettingsResponse, error)
	PutMappingFromJSON(ctx context.Context, in *PutMappingFromJSONRequest, opts ...client.CallOption) (*PutMappingFromJSONRequest, error)
	Status(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	SearchById(ctx context.Context, in *SearchByIdRequest, opts ...client.CallOption) (*SearchByIdResponse, error)
	AddAlias(ctx context.Context, in *AddAliasRequest, opts ...client.CallOption) (*AddAliasResponse, error)
	DeleteIndex(ctx context.Context, in *DeleteIndexRequest, opts ...client.CallOption) (*DeleteIndexResponse, error)
	DeleteAlias(ctx context.Context, in *DeleteAliasRequest, opts ...client.CallOption) (*DeleteAliasResponse, error)
	RenameAlias(ctx context.Context, in *RenameAliasRequest, opts ...client.CallOption) (*RenameAliasResponse, error)
}

type dBClient struct {
	c           client.Client
	serviceName string
}

func NewDBClient(serviceName string, c client.Client) DBClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.db"
	}
	return &dBClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *dBClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.Update", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) DeleteByQuery(ctx context.Context, in *DeleteByQueryRequest, opts ...client.CallOption) (*DeleteByQueryResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.DeleteByQuery", in)
	out := new(DeleteByQueryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) CreateIndexWithSettings(ctx context.Context, in *CreateIndexWithSettingsRequest, opts ...client.CallOption) (*CreateIndexWithSettingsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.CreateIndexWithSettings", in)
	out := new(CreateIndexWithSettingsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) PutMappingFromJSON(ctx context.Context, in *PutMappingFromJSONRequest, opts ...client.CallOption) (*PutMappingFromJSONRequest, error) {
	req := c.c.NewRequest(c.serviceName, "DB.PutMappingFromJSON", in)
	out := new(PutMappingFromJSONRequest)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) Status(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.Status", in)
	out := new(StatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) SearchById(ctx context.Context, in *SearchByIdRequest, opts ...client.CallOption) (*SearchByIdResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.SearchById", in)
	out := new(SearchByIdResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) AddAlias(ctx context.Context, in *AddAliasRequest, opts ...client.CallOption) (*AddAliasResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.AddAlias", in)
	out := new(AddAliasResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) DeleteIndex(ctx context.Context, in *DeleteIndexRequest, opts ...client.CallOption) (*DeleteIndexResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.DeleteIndex", in)
	out := new(DeleteIndexResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) DeleteAlias(ctx context.Context, in *DeleteAliasRequest, opts ...client.CallOption) (*DeleteAliasResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.DeleteAlias", in)
	out := new(DeleteAliasResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBClient) RenameAlias(ctx context.Context, in *RenameAliasRequest, opts ...client.CallOption) (*RenameAliasResponse, error) {
	req := c.c.NewRequest(c.serviceName, "DB.RenameAlias", in)
	out := new(RenameAliasResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DB service

type DBHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Read(context.Context, *ReadRequest, *ReadResponse) error
	Update(context.Context, *UpdateRequest, *UpdateResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	DeleteByQuery(context.Context, *DeleteByQueryRequest, *DeleteByQueryResponse) error
	CreateIndexWithSettings(context.Context, *CreateIndexWithSettingsRequest, *CreateIndexWithSettingsResponse) error
	PutMappingFromJSON(context.Context, *PutMappingFromJSONRequest, *PutMappingFromJSONRequest) error
	Status(context.Context, *StatusRequest, *StatusResponse) error
	Search(context.Context, *SearchRequest, *SearchResponse) error
	SearchById(context.Context, *SearchByIdRequest, *SearchByIdResponse) error
	AddAlias(context.Context, *AddAliasRequest, *AddAliasResponse) error
	DeleteIndex(context.Context, *DeleteIndexRequest, *DeleteIndexResponse) error
	DeleteAlias(context.Context, *DeleteAliasRequest, *DeleteAliasResponse) error
	RenameAlias(context.Context, *RenameAliasRequest, *RenameAliasResponse) error
}

func RegisterDBHandler(s server.Server, hdlr DBHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&DB{hdlr}, opts...))
}

type DB struct {
	DBHandler
}

func (h *DB) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.DBHandler.Create(ctx, in, out)
}

func (h *DB) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.DBHandler.Read(ctx, in, out)
}

func (h *DB) Update(ctx context.Context, in *UpdateRequest, out *UpdateResponse) error {
	return h.DBHandler.Update(ctx, in, out)
}

func (h *DB) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.DBHandler.Delete(ctx, in, out)
}

func (h *DB) DeleteByQuery(ctx context.Context, in *DeleteByQueryRequest, out *DeleteByQueryResponse) error {
	return h.DBHandler.DeleteByQuery(ctx, in, out)
}

func (h *DB) CreateIndexWithSettings(ctx context.Context, in *CreateIndexWithSettingsRequest, out *CreateIndexWithSettingsResponse) error {
	return h.DBHandler.CreateIndexWithSettings(ctx, in, out)
}

func (h *DB) PutMappingFromJSON(ctx context.Context, in *PutMappingFromJSONRequest, out *PutMappingFromJSONRequest) error {
	return h.DBHandler.PutMappingFromJSON(ctx, in, out)
}

func (h *DB) Status(ctx context.Context, in *StatusRequest, out *StatusResponse) error {
	return h.DBHandler.Status(ctx, in, out)
}

func (h *DB) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.DBHandler.Search(ctx, in, out)
}

func (h *DB) SearchById(ctx context.Context, in *SearchByIdRequest, out *SearchByIdResponse) error {
	return h.DBHandler.SearchById(ctx, in, out)
}

func (h *DB) AddAlias(ctx context.Context, in *AddAliasRequest, out *AddAliasResponse) error {
	return h.DBHandler.AddAlias(ctx, in, out)
}

func (h *DB) DeleteIndex(ctx context.Context, in *DeleteIndexRequest, out *DeleteIndexResponse) error {
	return h.DBHandler.DeleteIndex(ctx, in, out)
}

func (h *DB) DeleteAlias(ctx context.Context, in *DeleteAliasRequest, out *DeleteAliasResponse) error {
	return h.DBHandler.DeleteAlias(ctx, in, out)
}

func (h *DB) RenameAlias(ctx context.Context, in *RenameAliasRequest, out *RenameAliasResponse) error {
	return h.DBHandler.RenameAlias(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/db/srv/proto/db/db.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 921 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x57, 0x5b, 0x53, 0xd3, 0x5e,
	0x10, 0x87, 0x16, 0x7a, 0x59, 0xfe, 0xdc, 0xce, 0x1f, 0x24, 0x46, 0xe4, 0x12, 0x90, 0x61, 0x18,
	0xa7, 0x75, 0xe4, 0x55, 0x67, 0xa4, 0xa2, 0x4e, 0x75, 0x50, 0x49, 0x75, 0x98, 0x71, 0xc6, 0xa9,
	0x69, 0x73, 0x5a, 0x32, 0xa4, 0x49, 0xcc, 0x45, 0x2d, 0x0f, 0x7e, 0x02, 0x3f, 0xa0, 0x4f, 0x7e,
	0x16, 0xcf, 0x2d, 0x21, 0x69, 0x93, 0x16, 0x10, 0xde, 0xce, 0xee, 0xd9, 0xfd, 0xfd, 0x76, 0xf7,
	0x6c, 0x76, 0x27, 0xb0, 0xdf, 0x35, 0xfc, 0xd3, 0xa0, 0x55, 0x69, 0xdb, 0xbd, 0xea, 0x99, 0x76,
	0x6e, 0x07, 0x4e, 0xd5, 0x31, 0x35, 0xbf, 0x63, 0xbb, 0xbd, 0xaa, 0xde, 0xaa, 0x7a, 0xee, 0xb7,
	0xaa, 0xe3, 0xda, 0xbe, 0x4d, 0x05, 0xbd, 0x55, 0x61, 0x67, 0x34, 0xdf, 0xb5, 0x2b, 0x3d, 0xa3,
	0xed, 0xda, 0x15, 0x72, 0x5f, 0xd1, 0x5b, 0x4a, 0x0d, 0x0a, 0x87, 0x76, 0x5b, 0xc5, 0x1d, 0xb4,
	0x04, 0xd3, 0x86, 0xa5, 0xe3, 0x1f, 0xd2, 0xe4, 0xc6, 0xe4, 0x6e, 0x59, 0xe5, 0x02, 0x42, 0x30,
	0xe5, 0xf7, 0x1d, 0x2c, 0xe5, 0x98, 0x92, 0x9d, 0xd1, 0x1c, 0xe4, 0x0c, 0x5d, 0xca, 0x33, 0x0d,
	0x39, 0x29, 0x9f, 0x61, 0xf6, 0xb9, 0x8b, 0x35, 0x1f, 0xab, 0xf8, 0x6b, 0x80, 0x3d, 0xff, 0xfa,
	0x50, 0xd4, 0x46, 0xd7, 0x7c, 0x4d, 0x9a, 0xe2, 0x36, 0xf4, 0xac, 0x2c, 0xc0, 0x5c, 0x08, 0xef,
	0x39, 0xb6, 0xe5, 0x61, 0x45, 0x83, 0xc5, 0x5a, 0x60, 0x9e, 0xdd, 0x26, 0xe9, 0x12, 0xa0, 0x38,
	0x85, 0x20, 0x7e, 0x05, 0x33, 0x2a, 0xd6, 0xf4, 0x7f, 0xa6, 0x54, 0x76, 0xe0, 0x3f, 0x0e, 0xc4,
	0x81, 0xd1, 0x1d, 0x28, 0xb8, 0xd8, 0x0b, 0x4c, 0x5f, 0x40, 0x09, 0x89, 0x96, 0xf6, 0xa3, 0xa3,
	0xdf, 0x66, 0x69, 0x43, 0x78, 0x91, 0x61, 0x1d, 0x66, 0x0f, 0xb1, 0x89, 0x6f, 0x80, 0x90, 0x82,
	0x87, 0x50, 0x02, 0xfc, 0xf7, 0x24, 0x2c, 0x71, 0x55, 0xad, 0x7f, 0x1c, 0x60, 0xb7, 0x1f, 0x92,
	0x48, 0x50, 0x64, 0xb8, 0xd8, 0x23, 0x34, 0x79, 0xe2, 0x1f, 0x8a, 0x94, 0x9e, 0x82, 0x7b, 0x84,
	0x89, 0xea, 0xb9, 0xc0, 0xe8, 0xb1, 0xdb, 0x13, 0x64, 0xec, 0x8c, 0x64, 0x28, 0xb5, 0x49, 0x26,
	0x5d, 0xdb, 0xed, 0x8b, 0x1c, 0x23, 0x19, 0x2d, 0x40, 0x3e, 0x70, 0x4d, 0x69, 0x9a, 0xa9, 0xe9,
	0x91, 0xe2, 0xea, 0xd8, 0xf1, 0x4f, 0xa5, 0x02, 0xd1, 0xe5, 0x55, 0x2e, 0xa0, 0x7b, 0x50, 0xee,
	0x18, 0x26, 0x6e, 0xb2, 0xdc, 0x8a, 0x1c, 0x84, 0x2a, 0x3e, 0xd0, 0xfc, 0xc8, 0xa5, 0xa9, 0x79,
	0x7e, 0xd3, 0xc3, 0xd8, 0x92, 0x4a, 0xcc, 0xad, 0x44, 0x15, 0x0d, 0x22, 0x2b, 0x2b, 0xb0, 0x3c,
	0x90, 0x99, 0xc8, 0x59, 0x85, 0x35, 0xde, 0x44, 0x75, 0x9a, 0xd1, 0x09, 0xf9, 0x66, 0x1b, 0xd8,
	0xf7, 0x0d, 0xab, 0xeb, 0x8d, 0xae, 0x30, 0x49, 0xc7, 0x13, 0x86, 0xa2, 0xca, 0x91, 0xac, 0x6c,
	0xc2, 0x7a, 0x26, 0xa6, 0xa0, 0x6d, 0xc2, 0xdd, 0xf7, 0x81, 0x7f, 0xa4, 0x39, 0x0e, 0x51, 0xbf,
	0x74, 0xed, 0xde, 0xeb, 0xc6, 0xbb, 0xb7, 0x57, 0x7f, 0x53, 0xf2, 0x30, 0x3d, 0x8e, 0x21, 0x6a,
	0x1d, 0x8a, 0xca, 0x2a, 0xc8, 0x69, 0x04, 0x82, 0x7e, 0x1e, 0x66, 0x1b, 0xbe, 0xe6, 0x07, 0x61,
	0x92, 0xca, 0x2e, 0xcc, 0x85, 0x8a, 0x8b, 0x96, 0xf7, 0x98, 0x26, 0x6c, 0x79, 0x2e, 0x29, 0xbf,
	0x72, 0xc4, 0x17, 0x6b, 0x6e, 0xfb, 0x74, 0x7c, 0xb8, 0xb4, 0x07, 0x72, 0xb1, 0x1e, 0x20, 0xba,
	0x0e, 0x09, 0x85, 0xc5, 0x9a, 0x57, 0xd9, 0x99, 0xea, 0x3c, 0xe3, 0x1c, 0xb3, 0x9e, 0x20, 0x3a,
	0x7a, 0x4e, 0xf4, 0xca, 0x74, 0x7a, 0xaf, 0x14, 0x52, 0x7a, 0xa5, 0x18, 0xef, 0x95, 0xb0, 0x5c,
	0xa5, 0x58, 0xb9, 0x12, 0xfd, 0x53, 0x1e, 0xe8, 0x9f, 0x15, 0x28, 0x06, 0x1e, 0x76, 0x9b, 0xe4,
	0x23, 0x01, 0x9e, 0x31, 0x15, 0xeb, 0x7a, 0xb2, 0xb1, 0x66, 0x06, 0x1a, 0xeb, 0x09, 0x29, 0x9c,
	0xa8, 0xc6, 0xe8, 0x59, 0x41, 0x03, 0x32, 0xac, 0x8e, 0x1d, 0x16, 0x84, 0x9e, 0x95, 0x23, 0x58,
	0xe4, 0xde, 0xb5, 0x7e, 0xfd, 0x06, 0xc6, 0xd6, 0x43, 0x40, 0x71, 0xb8, 0x31, 0xc3, 0xeb, 0x29,
	0xcc, 0x1f, 0xe8, 0xfa, 0x81, 0x69, 0x68, 0x63, 0x7a, 0x9d, 0x68, 0x35, 0x6a, 0x25, 0xb8, 0xb9,
	0xa0, 0x20, 0x58, 0xb8, 0x70, 0x17, 0x7d, 0xb5, 0x07, 0x88, 0x7f, 0x66, 0xac, 0xf3, 0x47, 0xa2,
	0x2a, 0xcb, 0xf0, 0x7f, 0xc2, 0x56, 0x40, 0x3c, 0x0b, 0x21, 0xae, 0x1d, 0x58, 0x04, 0x9c, 0x8c,
	0x4d, 0x07, 0xa4, 0x62, 0x4b, 0xeb, 0x5d, 0x06, 0x98, 0x3c, 0xb9, 0x6d, 0xea, 0xcd, 0x38, 0x78,
	0x89, 0x28, 0x98, 0x27, 0xbd, 0xb4, 0xf0, 0x77, 0x71, 0xc9, 0x8b, 0x5f, 0x22, 0x8a, 0x83, 0x90,
	0x3c, 0xc1, 0xc2, 0xc9, 0x1f, 0xff, 0x29, 0x43, 0xee, 0xb0, 0x86, 0xde, 0x40, 0x81, 0x4f, 0x06,
	0xb4, 0x56, 0x19, 0x58, 0xf5, 0x95, 0xc4, 0xba, 0x94, 0xd7, 0x33, 0xef, 0x45, 0x3a, 0x13, 0xe8,
	0x05, 0x4c, 0xd1, 0x25, 0x85, 0x56, 0x87, 0x4c, 0x63, 0x4b, 0x50, 0xbe, 0x9f, 0x71, 0x1b, 0xc1,
	0x90, 0x98, 0xf8, 0x92, 0x49, 0x89, 0x29, 0xb1, 0xdc, 0x52, 0x62, 0x1a, 0xd8, 0x4e, 0x0c, 0x8c,
	0xd7, 0x3e, 0x05, 0x2c, 0xb1, 0xb8, 0x52, 0xc0, 0x06, 0xb6, 0xd1, 0x04, 0xfa, 0x12, 0x2e, 0x3b,
	0x31, 0xb4, 0xd1, 0x83, 0x0c, 0x9f, 0xe4, 0xba, 0x92, 0x77, 0xc6, 0x99, 0x45, 0x0c, 0x3f, 0x61,
	0x25, 0x63, 0x52, 0xa3, 0x6a, 0xc6, 0x03, 0x64, 0xed, 0x09, 0xf9, 0xd1, 0xe5, 0x1d, 0x22, 0x7e,
	0x0b, 0xd0, 0xf0, 0x94, 0x46, 0x7b, 0x43, 0x48, 0x99, 0xbb, 0x42, 0xbe, 0x82, 0x2d, 0x7f, 0x1e,
	0x3e, 0xe6, 0x53, 0x9e, 0x27, 0xb1, 0x10, 0x52, 0x9e, 0x27, 0xb9, 0x1f, 0x04, 0x18, 0x9b, 0x36,
	0x69, 0x60, 0xf1, 0x0d, 0x91, 0x06, 0x96, 0x98, 0x99, 0x04, 0xec, 0x04, 0xe0, 0x62, 0x74, 0x21,
	0x25, 0xc3, 0x21, 0x36, 0x26, 0xe5, 0xad, 0x91, 0x36, 0x11, 0xf0, 0x31, 0x94, 0xc2, 0x31, 0x85,
	0x36, 0x86, 0x5c, 0x06, 0x06, 0xa0, 0xbc, 0x39, 0xc2, 0x22, 0x82, 0xfc, 0x04, 0x33, 0xb1, 0xc9,
	0x85, 0xb6, 0x32, 0xda, 0x2d, 0x3e, 0x03, 0xe5, 0xed, 0xd1, 0x46, 0xc3, 0xd8, 0x3c, 0xe2, 0x2c,
	0xec, 0x44, 0xd0, 0xdb, 0xa3, 0x8d, 0xe2, 0xd8, 0xb1, 0xd9, 0x94, 0x82, 0x3d, 0x3c, 0x1f, 0x53,
	0xb0, 0x53, 0xc6, 0x9b, 0x32, 0xd1, 0x2a, 0xb0, 0x1f, 0x98, 0xfd, 0xbf, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xd9, 0x6a, 0x73, 0x43, 0xf7, 0x0c, 0x00, 0x00,
}