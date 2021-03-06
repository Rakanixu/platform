// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/lib/db/bulk/proto/bulk/bulk.proto
// DO NOT EDIT!

/*
Package proto_bulk is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/lib/db/bulk/proto/bulk/bulk.proto

It has these top-level messages:
	DocRef
	BulkCreateRequest
	BulkCreateResponse
*/
package proto_bulk

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

type DocRef struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *DocRef) Reset()                    { *m = DocRef{} }
func (m *DocRef) String() string            { return proto.CompactTextString(m) }
func (*DocRef) ProtoMessage()               {}
func (*DocRef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DocRef) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *DocRef) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *DocRef) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type BulkCreateRequest struct {
	Index string `protobuf:"bytes,1,opt,name=index" json:"index,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Data  string `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
}

func (m *BulkCreateRequest) Reset()                    { *m = BulkCreateRequest{} }
func (m *BulkCreateRequest) String() string            { return proto.CompactTextString(m) }
func (*BulkCreateRequest) ProtoMessage()               {}
func (*BulkCreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BulkCreateRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *BulkCreateRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *BulkCreateRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BulkCreateRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type BulkCreateResponse struct {
}

func (m *BulkCreateResponse) Reset()                    { *m = BulkCreateResponse{} }
func (m *BulkCreateResponse) String() string            { return proto.CompactTextString(m) }
func (*BulkCreateResponse) ProtoMessage()               {}
func (*BulkCreateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*DocRef)(nil), "proto.bulk.DocRef")
	proto.RegisterType((*BulkCreateRequest)(nil), "proto.bulk.BulkCreateRequest")
	proto.RegisterType((*BulkCreateResponse)(nil), "proto.bulk.BulkCreateResponse")
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/lib/db/bulk/proto/bulk/bulk.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x8f, 0x31, 0x8e, 0x83, 0x30,
	0x10, 0x45, 0x05, 0xcb, 0x22, 0xed, 0x14, 0x2b, 0xc5, 0xa2, 0x70, 0x19, 0x51, 0xa5, 0xc2, 0x45,
	0xda, 0x54, 0x24, 0x27, 0xe0, 0x06, 0x36, 0x1e, 0x12, 0x0b, 0x83, 0x1d, 0x18, 0x4b, 0x49, 0x4e,
	0x1f, 0xd9, 0x34, 0xa9, 0xd3, 0xd8, 0xef, 0x3f, 0x7d, 0xcd, 0x68, 0xe0, 0x74, 0x35, 0x74, 0x0b,
	0xaa, 0xe9, 0xdd, 0x24, 0x46, 0xf9, 0x72, 0xc1, 0x0b, 0x6f, 0x25, 0x0d, 0x6e, 0x99, 0x84, 0x35,
	0x4a, 0x68, 0x25, 0x54, 0xb0, 0xa3, 0xf0, 0x8b, 0x23, 0xb7, 0x61, 0x7c, 0x9a, 0x94, 0x19, 0xa4,
	0xaf, 0x89, 0xa6, 0x6e, 0xa1, 0xbc, 0xb8, 0xbe, 0xc3, 0x81, 0x55, 0xf0, 0x6b, 0x66, 0x8d, 0x0f,
	0x9e, 0xed, 0xb3, 0xc3, 0x5f, 0xb7, 0x05, 0xc6, 0xa0, 0xa0, 0xa7, 0x47, 0x9e, 0x27, 0x99, 0x98,
	0xfd, 0x43, 0x6e, 0x34, 0xff, 0x49, 0x26, 0x37, 0xba, 0x96, 0xb0, 0x6b, 0x83, 0x1d, 0xcf, 0x0b,
	0x4a, 0xc2, 0x0e, 0xef, 0x01, 0x57, 0xfa, 0x7e, 0x5c, 0xec, 0x68, 0x49, 0x92, 0x17, 0x5b, 0x27,
	0x72, 0x5d, 0x01, 0xfb, 0x5c, 0xb1, 0x7a, 0x37, 0xaf, 0xa8, 0xca, 0x74, 0xc8, 0xf1, 0x1d, 0x00,
	0x00, 0xff, 0xff, 0xf6, 0x5e, 0x58, 0x45, 0x0f, 0x01, 0x00, 0x00,
}
