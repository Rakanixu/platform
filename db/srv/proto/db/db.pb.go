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

type CreateIndexWithSettingsRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Settings string `protobuf:"bytes,2,opt,name=settings" json:"settings,omitempty"`
}

func (m *CreateIndexWithSettingsRequest) Reset()                    { *m = CreateIndexWithSettingsRequest{} }
func (m *CreateIndexWithSettingsRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateIndexWithSettingsRequest) ProtoMessage()               {}
func (*CreateIndexWithSettingsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type CreateIndexWithSettingsResponse struct {
}

func (m *CreateIndexWithSettingsResponse) Reset()         { *m = CreateIndexWithSettingsResponse{} }
func (m *CreateIndexWithSettingsResponse) String() string { return proto.CompactTextString(m) }
func (*CreateIndexWithSettingsResponse) ProtoMessage()    {}
func (*CreateIndexWithSettingsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{12}
}

type PutMappingFromJSONRequest struct {
	Index   string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type    string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Mapping string `protobuf:"bytes,3,opt,name=mapping" json:"mapping,omitempty"`
}

func (m *PutMappingFromJSONRequest) Reset()                    { *m = PutMappingFromJSONRequest{} }
func (m *PutMappingFromJSONRequest) String() string            { return proto.CompactTextString(m) }
func (*PutMappingFromJSONRequest) ProtoMessage()               {}
func (*PutMappingFromJSONRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type PutMappingFromJSONResponse struct {
}

func (m *PutMappingFromJSONResponse) Reset()                    { *m = PutMappingFromJSONResponse{} }
func (m *PutMappingFromJSONResponse) String() string            { return proto.CompactTextString(m) }
func (*PutMappingFromJSONResponse) ProtoMessage()               {}
func (*PutMappingFromJSONResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type StatusResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

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
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

type SearchByIdRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *SearchByIdRequest) Reset()                    { *m = SearchByIdRequest{} }
func (m *SearchByIdRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchByIdRequest) ProtoMessage()               {}
func (*SearchByIdRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

type SearchByIdResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *SearchByIdResponse) Reset()                    { *m = SearchByIdResponse{} }
func (m *SearchByIdResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchByIdResponse) ProtoMessage()               {}
func (*SearchByIdResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

type AddAliasRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Alias string `protobuf:"bytes,2,opt,name=alias" json:"alias,omitempty"`
}

func (m *AddAliasRequest) Reset()                    { *m = AddAliasRequest{} }
func (m *AddAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*AddAliasRequest) ProtoMessage()               {}
func (*AddAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

type AddAliasResponse struct {
}

func (m *AddAliasResponse) Reset()                    { *m = AddAliasResponse{} }
func (m *AddAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*AddAliasResponse) ProtoMessage()               {}
func (*AddAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{22} }

type DeleteIndexRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
}

func (m *DeleteIndexRequest) Reset()                    { *m = DeleteIndexRequest{} }
func (m *DeleteIndexRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexRequest) ProtoMessage()               {}
func (*DeleteIndexRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{23} }

type DeleteIndexResponse struct {
}

func (m *DeleteIndexResponse) Reset()                    { *m = DeleteIndexResponse{} }
func (m *DeleteIndexResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexResponse) ProtoMessage()               {}
func (*DeleteIndexResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{24} }

type DeleteAliasRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Alias string `protobuf:"bytes,2,opt,name=alias" json:"alias,omitempty"`
}

func (m *DeleteAliasRequest) Reset()                    { *m = DeleteAliasRequest{} }
func (m *DeleteAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteAliasRequest) ProtoMessage()               {}
func (*DeleteAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{25} }

type DeleteAliasResponse struct {
}

func (m *DeleteAliasResponse) Reset()                    { *m = DeleteAliasResponse{} }
func (m *DeleteAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteAliasResponse) ProtoMessage()               {}
func (*DeleteAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{26} }

type RenameAliasRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	OldAlias string `protobuf:"bytes,2,opt,name=old_alias,json=oldAlias" json:"old_alias,omitempty"`
	NewAlias string `protobuf:"bytes,3,opt,name=new_alias,json=newAlias" json:"new_alias,omitempty"`
}

func (m *RenameAliasRequest) Reset()                    { *m = RenameAliasRequest{} }
func (m *RenameAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*RenameAliasRequest) ProtoMessage()               {}
func (*RenameAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{27} }

type RenameAliasResponse struct {
}

func (m *RenameAliasResponse) Reset()                    { *m = RenameAliasResponse{} }
func (m *RenameAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*RenameAliasResponse) ProtoMessage()               {}
func (*RenameAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{28} }

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
	// 790 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x96, 0x5b, 0x4f, 0xd4, 0x40,
	0x14, 0xc7, 0xb9, 0x2e, 0xcb, 0x41, 0x6e, 0x23, 0x6a, 0xad, 0xc8, 0x65, 0x20, 0x86, 0x10, 0xb3,
	0x6b, 0xe4, 0x55, 0x13, 0x59, 0x51, 0x83, 0x06, 0x2f, 0x45, 0x43, 0x62, 0x62, 0x48, 0x77, 0x3b,
	0xbb, 0x34, 0xf4, 0x66, 0x3b, 0x55, 0xe1, 0xc1, 0x8f, 0xe9, 0xa7, 0xf1, 0xc1, 0xb9, 0xb5, 0xb4,
	0xbb, 0x6d, 0xb9, 0x08, 0x6f, 0xe7, 0x9c, 0x39, 0xf3, 0x3b, 0x67, 0x4e, 0x67, 0xfe, 0xbb, 0xb0,
	0xd5, 0xb3, 0xe9, 0x51, 0xdc, 0x6e, 0x74, 0x7c, 0xb7, 0x79, 0x6c, 0x9e, 0xfa, 0x71, 0xd0, 0x0c,
	0x1c, 0x93, 0x76, 0xfd, 0xd0, 0x6d, 0x5a, 0xed, 0x66, 0x14, 0xfe, 0x68, 0x06, 0xa1, 0x4f, 0x7d,
	0xee, 0x58, 0xed, 0x86, 0xb0, 0xd1, 0x6c, 0xcf, 0x6f, 0xb8, 0x76, 0x27, 0xf4, 0x1b, 0x6c, 0xbd,
	0x61, 0xb5, 0x71, 0x0b, 0x6a, 0x3b, 0x7e, 0xc7, 0x20, 0x5d, 0xb4, 0x00, 0xe3, 0xb6, 0x67, 0x91,
	0x5f, 0xda, 0xf0, 0xca, 0xf0, 0xc6, 0xa4, 0x21, 0x1d, 0x84, 0x60, 0x8c, 0x9e, 0x04, 0x44, 0x1b,
	0x11, 0x41, 0x61, 0xa3, 0x19, 0x18, 0xb1, 0x2d, 0x6d, 0x54, 0x44, 0x98, 0x85, 0xbf, 0xc1, 0xf4,
	0xcb, 0x90, 0x98, 0x94, 0x18, 0xe4, 0x7b, 0x4c, 0x22, 0x7a, 0x75, 0x14, 0xcf, 0xb1, 0x4c, 0x6a,
	0x6a, 0x63, 0x32, 0x87, 0xdb, 0x78, 0x0e, 0x66, 0x12, 0x7c, 0x14, 0xf8, 0x5e, 0x44, 0xb0, 0x09,
	0xf3, 0xad, 0xd8, 0x39, 0xbe, 0xc9, 0xa2, 0x0b, 0x80, 0xb2, 0x25, 0x54, 0xe1, 0x37, 0x30, 0x65,
	0x10, 0xd3, 0xfa, 0xef, 0x92, 0xf8, 0x11, 0xdc, 0x92, 0x20, 0x09, 0x46, 0x77, 0xa1, 0x16, 0x92,
	0x28, 0x76, 0xa8, 0x42, 0x29, 0x8f, 0x8f, 0xf6, 0x4b, 0x60, 0xdd, 0xe4, 0x68, 0x13, 0xbc, 0x3a,
	0xe1, 0x2e, 0x4c, 0xef, 0x10, 0x87, 0x5c, 0x43, 0x41, 0x0e, 0x4f, 0x50, 0x0a, 0x6e, 0xc0, 0x92,
	0x1c, 0xe8, 0x2e, 0x87, 0x1c, 0xb0, 0xfb, 0xbb, 0x4f, 0x28, 0xb5, 0xbd, 0x5e, 0x54, 0x5d, 0x4d,
	0x87, 0x7a, 0xa4, 0x12, 0x55, 0xc5, 0xd4, 0xc7, 0xab, 0xb0, 0x5c, 0xca, 0x54, 0x65, 0x0f, 0xe1,
	0xfe, 0xc7, 0x98, 0xee, 0x99, 0x41, 0xc0, 0xc2, 0xaf, 0x43, 0xdf, 0x7d, 0xbb, 0xff, 0xe1, 0xfd,
	0xe5, 0xcf, 0xa7, 0xc1, 0x84, 0x2b, 0x19, 0xea, 0x90, 0x89, 0x8b, 0x17, 0x41, 0x2f, 0x2a, 0xa0,
	0xca, 0xcf, 0xc2, 0xf4, 0x3e, 0x35, 0x69, 0x9c, 0x1c, 0x12, 0x6f, 0xc0, 0x4c, 0x12, 0x38, 0xfb,
	0xfc, 0x91, 0x88, 0x24, 0x9f, 0x5f, 0x7a, 0xf8, 0xcf, 0x30, 0xdb, 0x4b, 0xcc, 0xb0, 0x73, 0x74,
	0x7e, 0xbb, 0x24, 0x74, 0xd3, 0x76, 0x99, 0xcd, 0x63, 0x5d, 0xd6, 0x8a, 0xe8, 0x75, 0xd4, 0x10,
	0x36, 0x8f, 0x45, 0xf6, 0x29, 0x11, 0x77, 0x80, 0xc5, 0xb8, 0xcd, 0x87, 0xdb, 0x61, 0xe3, 0xeb,
	0xf9, 0xe1, 0x89, 0x36, 0x2e, 0x87, 0x9b, 0xf8, 0x68, 0x0e, 0x46, 0xe3, 0xd0, 0xd1, 0x6a, 0x22,
	0xcc, 0x4d, 0x5e, 0xdf, 0x22, 0x01, 0x3d, 0xd2, 0x26, 0x04, 0x42, 0x3a, 0xe9, 0xb8, 0xea, 0x99,
	0x71, 0x3d, 0x80, 0xc9, 0xae, 0xed, 0x90, 0x43, 0xb1, 0x30, 0x29, 0xc1, 0x3c, 0xf0, 0x99, 0xf9,
	0xf8, 0x19, 0x1b, 0x81, 0x3a, 0x57, 0xf5, 0x0b, 0xe0, 0x68, 0xdb, 0xeb, 0xfa, 0xc9, 0xd1, 0xb8,
	0x8d, 0xf7, 0x60, 0x5e, 0xee, 0x6e, 0x9d, 0xec, 0x5e, 0xc3, 0x63, 0x7c, 0x0c, 0x28, 0x8b, 0x3b,
	0xe7, 0x49, 0x3e, 0x87, 0xd9, 0x6d, 0xcb, 0xda, 0x76, 0x6c, 0xf3, 0x9c, 0x5b, 0xcb, 0xa2, 0x26,
	0xcf, 0x52, 0xb5, 0xa5, 0x83, 0x11, 0xcc, 0x9d, 0x6d, 0x57, 0x37, 0x64, 0x13, 0x90, 0x7c, 0x29,
	0xe2, 0x0e, 0x57, 0x52, 0xf1, 0x1d, 0xb8, 0x9d, 0xcb, 0x55, 0x88, 0x17, 0x09, 0xe2, 0xca, 0x8d,
	0xa5, 0xe0, 0x7c, 0x6f, 0x16, 0x20, 0x83, 0x78, 0xa6, 0x7b, 0x11, 0x30, 0xfb, 0xe4, 0xbe, 0x63,
	0x1d, 0x66, 0xe1, 0x75, 0x16, 0x10, 0x3b, 0xf9, 0xa2, 0x47, 0x7e, 0xaa, 0x45, 0x39, 0xfc, 0x3a,
	0x0b, 0x6c, 0x27, 0xc5, 0x73, 0x55, 0x64, 0xf1, 0xa7, 0x7f, 0xeb, 0x30, 0xb2, 0xd3, 0x42, 0xef,
	0xa0, 0x26, 0xdf, 0x38, 0x5a, 0x6a, 0xf4, 0xfd, 0x80, 0x35, 0x72, 0x3f, 0x02, 0xfa, 0x72, 0xe9,
	0xba, 0x3a, 0xce, 0x10, 0x7a, 0x05, 0x63, 0x5c, 0x7a, 0xd1, 0xe2, 0x40, 0x6a, 0x46, 0xda, 0xf5,
	0x87, 0x25, 0xab, 0x29, 0x86, 0xf5, 0x24, 0xa5, 0xb3, 0xa0, 0xa7, 0x9c, 0x64, 0x17, 0xf4, 0xd4,
	0xa7, 0xb9, 0x02, 0x26, 0x67, 0x5f, 0x00, 0xcb, 0xc9, 0x71, 0x01, 0xac, 0x4f, 0x63, 0x87, 0xd0,
	0x6f, 0xb8, 0x57, 0xa2, 0x88, 0xa8, 0x59, 0x32, 0x9e, 0x32, 0x3d, 0xd6, 0x9f, 0x5c, 0x7c, 0x43,
	0x5a, 0xdf, 0x03, 0x34, 0xa8, 0x86, 0x68, 0x73, 0x80, 0x54, 0xaa, 0xc9, 0xfa, 0x25, 0x72, 0xe5,
	0xf0, 0xa4, 0x9c, 0x16, 0x0c, 0x2f, 0x27, 0xbc, 0x05, 0xc3, 0xcb, 0xeb, 0xb0, 0x82, 0x09, 0x2d,
	0x28, 0x82, 0x65, 0x95, 0xb8, 0x08, 0x96, 0x53, 0x34, 0x06, 0x3b, 0x00, 0x38, 0x13, 0x16, 0x84,
	0x4b, 0x36, 0x64, 0x44, 0x4c, 0x5f, 0xab, 0xcc, 0x49, 0xc1, 0x9f, 0xa0, 0x9e, 0x88, 0x08, 0x5a,
	0x19, 0xd8, 0xd2, 0x27, 0x4f, 0xfa, 0x6a, 0x45, 0x46, 0x8a, 0xfc, 0x0a, 0x53, 0x19, 0x5d, 0x41,
	0x6b, 0x25, 0xf7, 0x2c, 0xab, 0x50, 0xfa, 0x7a, 0x75, 0xd2, 0x20, 0x5b, 0x76, 0x5c, 0xc6, 0xce,
	0x35, 0xbd, 0x5e, 0x9d, 0x94, 0x65, 0x67, 0x94, 0xa3, 0x80, 0x3d, 0xa8, 0x5e, 0x05, 0xec, 0x02,
	0xf1, 0xc1, 0x43, 0xed, 0x9a, 0xf8, 0xd3, 0xbc, 0xf5, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xf0, 0xe3,
	0xda, 0x00, 0x6b, 0x0b, 0x00, 0x00,
}
