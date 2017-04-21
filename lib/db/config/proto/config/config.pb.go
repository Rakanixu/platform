// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/lib/db/config/proto/config/config.proto
// DO NOT EDIT!

/*
Package proto_config is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/lib/db/config/proto/config/config.proto

It has these top-level messages:
	CreateIndexRequest
	CreateIndexResponse
	StatusRequest
	StatusResponse
	AddAliasRequest
	AddAliasResponse
	DeleteIndexRequest
	DeleteIndexResponse
	DeleteAliasRequest
	DeleteAliasResponse
	RenameAliasRequest
	RenameAliasResponse
*/
package proto_config

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CreateIndexRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
}

func (m *CreateIndexRequest) Reset()                    { *m = CreateIndexRequest{} }
func (m *CreateIndexRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateIndexRequest) ProtoMessage()               {}
func (*CreateIndexRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CreateIndexRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

type CreateIndexResponse struct {
}

func (m *CreateIndexResponse) Reset()                    { *m = CreateIndexResponse{} }
func (m *CreateIndexResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateIndexResponse) ProtoMessage()               {}
func (*CreateIndexResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type StatusResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *StatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type AddAliasRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Alias string `protobuf:"bytes,2,opt,name=alias" json:"alias,omitempty"`
}

func (m *AddAliasRequest) Reset()                    { *m = AddAliasRequest{} }
func (m *AddAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*AddAliasRequest) ProtoMessage()               {}
func (*AddAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *AddAliasRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *AddAliasRequest) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

type AddAliasResponse struct {
}

func (m *AddAliasResponse) Reset()                    { *m = AddAliasResponse{} }
func (m *AddAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*AddAliasResponse) ProtoMessage()               {}
func (*AddAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type DeleteIndexRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
}

func (m *DeleteIndexRequest) Reset()                    { *m = DeleteIndexRequest{} }
func (m *DeleteIndexRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexRequest) ProtoMessage()               {}
func (*DeleteIndexRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *DeleteIndexRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

type DeleteIndexResponse struct {
}

func (m *DeleteIndexResponse) Reset()                    { *m = DeleteIndexResponse{} }
func (m *DeleteIndexResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexResponse) ProtoMessage()               {}
func (*DeleteIndexResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type DeleteAliasRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Alias string `protobuf:"bytes,2,opt,name=alias" json:"alias,omitempty"`
}

func (m *DeleteAliasRequest) Reset()                    { *m = DeleteAliasRequest{} }
func (m *DeleteAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteAliasRequest) ProtoMessage()               {}
func (*DeleteAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *DeleteAliasRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *DeleteAliasRequest) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

type DeleteAliasResponse struct {
}

func (m *DeleteAliasResponse) Reset()                    { *m = DeleteAliasResponse{} }
func (m *DeleteAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteAliasResponse) ProtoMessage()               {}
func (*DeleteAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type RenameAliasRequest struct {
	Index    string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	OldAlias string `protobuf:"bytes,2,opt,name=old_alias,json=oldAlias" json:"old_alias,omitempty"`
	NewAlias string `protobuf:"bytes,3,opt,name=new_alias,json=newAlias" json:"new_alias,omitempty"`
}

func (m *RenameAliasRequest) Reset()                    { *m = RenameAliasRequest{} }
func (m *RenameAliasRequest) String() string            { return proto.CompactTextString(m) }
func (*RenameAliasRequest) ProtoMessage()               {}
func (*RenameAliasRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *RenameAliasRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *RenameAliasRequest) GetOldAlias() string {
	if m != nil {
		return m.OldAlias
	}
	return ""
}

func (m *RenameAliasRequest) GetNewAlias() string {
	if m != nil {
		return m.NewAlias
	}
	return ""
}

type RenameAliasResponse struct {
}

func (m *RenameAliasResponse) Reset()                    { *m = RenameAliasResponse{} }
func (m *RenameAliasResponse) String() string            { return proto.CompactTextString(m) }
func (*RenameAliasResponse) ProtoMessage()               {}
func (*RenameAliasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func init() {
	proto.RegisterType((*CreateIndexRequest)(nil), "proto.config.CreateIndexRequest")
	proto.RegisterType((*CreateIndexResponse)(nil), "proto.config.CreateIndexResponse")
	proto.RegisterType((*StatusRequest)(nil), "proto.config.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "proto.config.StatusResponse")
	proto.RegisterType((*AddAliasRequest)(nil), "proto.config.AddAliasRequest")
	proto.RegisterType((*AddAliasResponse)(nil), "proto.config.AddAliasResponse")
	proto.RegisterType((*DeleteIndexRequest)(nil), "proto.config.DeleteIndexRequest")
	proto.RegisterType((*DeleteIndexResponse)(nil), "proto.config.DeleteIndexResponse")
	proto.RegisterType((*DeleteAliasRequest)(nil), "proto.config.DeleteAliasRequest")
	proto.RegisterType((*DeleteAliasResponse)(nil), "proto.config.DeleteAliasResponse")
	proto.RegisterType((*RenameAliasRequest)(nil), "proto.config.RenameAliasRequest")
	proto.RegisterType((*RenameAliasResponse)(nil), "proto.config.RenameAliasResponse")
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/lib/db/config/proto/config/config.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x51, 0x4d, 0x4b, 0xc4, 0x30,
	0x10, 0x65, 0x15, 0x17, 0x77, 0x50, 0x57, 0xe2, 0x07, 0x0b, 0x5e, 0x24, 0xa7, 0x65, 0x0f, 0xcd,
	0xc1, 0xb3, 0x60, 0xd5, 0x8b, 0xd7, 0xfa, 0x03, 0x24, 0xdd, 0xcc, 0xae, 0xc1, 0x34, 0x53, 0x9b,
	0x94, 0x15, 0x7f, 0xbd, 0x34, 0x49, 0xd5, 0x5e, 0xa4, 0x78, 0x0a, 0xef, 0xcd, 0x9b, 0xf7, 0xde,
	0x10, 0xb8, 0xdf, 0x6a, 0xff, 0xda, 0x96, 0xd9, 0x9a, 0x2a, 0xf1, 0x26, 0x3f, 0xa9, 0xad, 0x45,
	0x6d, 0xa4, 0xdf, 0x50, 0x53, 0x09, 0xa3, 0x4b, 0xa1, 0x4a, 0xb1, 0x26, 0xbb, 0xd1, 0x5b, 0x51,
	0x37, 0xe4, 0xa9, 0x07, 0xf1, 0xc9, 0x02, 0xc7, 0x8e, 0xc2, 0x93, 0x45, 0x8e, 0xaf, 0x80, 0x3d,
	0x34, 0x28, 0x3d, 0x3e, 0x59, 0x85, 0x1f, 0x05, 0xbe, 0xb7, 0xe8, 0x3c, 0x3b, 0x87, 0x03, 0xdd,
	0xe1, 0xc5, 0xe4, 0x7a, 0xb2, 0x9c, 0x15, 0x11, 0xf0, 0x0b, 0x38, 0x1b, 0x68, 0x5d, 0x4d, 0xd6,
	0x21, 0x9f, 0xc3, 0xf1, 0xb3, 0x97, 0xbe, 0x75, 0x69, 0x9b, 0x2f, 0xe1, 0xa4, 0x27, 0xa2, 0x84,
	0x5d, 0xc2, 0xd4, 0x05, 0x26, 0x19, 0x26, 0xc4, 0x6f, 0x61, 0x9e, 0x2b, 0x95, 0x1b, 0x2d, 0xdd,
	0x9f, 0xd1, 0x1d, 0x2b, 0x3b, 0xd5, 0x62, 0x2f, 0xb2, 0x01, 0x70, 0x06, 0xa7, 0x3f, 0xeb, 0xa9,
	0xcd, 0x0a, 0xd8, 0x23, 0x1a, 0x1c, 0x7b, 0xd0, 0x40, 0x9b, 0x2c, 0xee, 0x7a, 0x8b, 0x7f, 0x17,
	0xfb, 0x36, 0x1e, 0x76, 0x53, 0xc0, 0x0a, 0xb4, 0xb2, 0x1a, 0x63, 0x7c, 0x05, 0x33, 0x32, 0xea,
	0xe5, 0xb7, 0xf9, 0x21, 0x99, 0x78, 0x6c, 0x37, 0xb4, 0xb8, 0x4b, 0xc3, 0xfd, 0x38, 0xb4, 0xb8,
	0xcb, 0xfb, 0xf0, 0x41, 0x4a, 0x0c, 0x2f, 0xa7, 0xe1, 0xdf, 0x6f, 0xbe, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x7d, 0x10, 0xbf, 0xcc, 0x44, 0x02, 0x00, 0x00,
}