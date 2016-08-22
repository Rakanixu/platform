// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/search/srv/proto/search/messages.proto
// DO NOT EDIT!

/*
Package go_micro_srv_search_search is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/search/srv/proto/search/messages.proto
	github.com/kazoup/platform/search/srv/proto/search/search.proto

It has these top-level messages:
	SearchRequest
	SearchResponse
*/
package go_micro_srv_search_search

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
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SearchResponse struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   string `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*SearchRequest)(nil), "go.micro.srv.search.search.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "go.micro.srv.search.search.SearchResponse")
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/search/srv/proto/search/messages.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x3c, 0x90, 0x4d, 0x4e, 0xc3, 0x30,
	0x10, 0x85, 0x15, 0xd2, 0x86, 0x62, 0x09, 0x84, 0x2c, 0x84, 0xac, 0xae, 0x50, 0x57, 0xac, 0x92,
	0x05, 0x5b, 0x36, 0x5c, 0x21, 0x9c, 0x20, 0x4d, 0x27, 0x3f, 0x22, 0xce, 0x98, 0x19, 0x1b, 0xd1,
	0xde, 0x8c, 0xdb, 0x11, 0x4f, 0x92, 0xae, 0xf2, 0xbd, 0x4f, 0x33, 0x99, 0x27, 0xab, 0x8f, 0xb6,
	0xf7, 0x5d, 0x38, 0xe6, 0x35, 0xda, 0xe2, 0xab, 0xba, 0x60, 0x70, 0x85, 0x1b, 0x2a, 0xdf, 0x20,
	0xd9, 0x82, 0xa1, 0xa2, 0xba, 0x2b, 0x98, 0x7e, 0x0a, 0x47, 0xe8, 0x71, 0x15, 0x16, 0x98, 0xab,
	0x16, 0x38, 0x17, 0xab, 0xf7, 0x2d, 0xe6, 0xb6, 0xaf, 0x09, 0xf3, 0x69, 0x32, 0x9f, 0x67, 0x96,
	0xcf, 0xe1, 0x2f, 0x51, 0xf7, 0x9f, 0x82, 0x25, 0x7c, 0x07, 0x60, 0xaf, 0x9f, 0xd4, 0xb6, 0x1f,
	0x4f, 0xf0, 0x6b, 0x92, 0x97, 0xe4, 0xf5, 0xae, 0x9c, 0x83, 0xd6, 0x6a, 0xe3, 0x81, 0xac, 0xb9,
	0x11, 0x29, 0x1c, 0x5d, 0x43, 0x68, 0x4d, 0x3a, 0xb9, 0xb4, 0x14, 0x8e, 0x8e, 0xfb, 0x0b, 0x98,
	0xcd, 0xec, 0x22, 0xeb, 0xbd, 0xda, 0xd5, 0x95, 0x87, 0x16, 0xe9, 0x6c, 0xb6, 0xb2, 0x7f, 0xcd,
	0xfa, 0x51, 0xa5, 0x81, 0x06, 0x93, 0x89, 0x8e, 0x18, 0xef, 0x9f, 0xc0, 0xf9, 0xce, 0xdc, 0xca,
	0x2f, 0xe6, 0x20, 0xf7, 0xcf, 0x0e, 0xcc, 0x6e, 0xb9, 0x3f, 0xf1, 0xe1, 0x5d, 0x3d, 0xac, 0xd5,
	0xd9, 0xe1, 0xc8, 0xa0, 0x9f, 0x55, 0x46, 0xc0, 0x61, 0xf0, 0x4b, 0xf9, 0x25, 0xc5, 0xed, 0x7e,
	0x6c, 0x70, 0x6d, 0x1f, 0xf9, 0x98, 0xc9, 0xe3, 0xbc, 0xfd, 0x07, 0x00, 0x00, 0xff, 0xff, 0x15,
	0xc5, 0xdb, 0x7f, 0x61, 0x01, 0x00, 0x00,
}
