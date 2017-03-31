// Code generated by protoc-gen-go.
// source: github.com/kazoup/platform/lib/protomsg/announce/announce.proto
// DO NOT EDIT!

/*
Package message is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/platform/lib/protomsg/announce/announce.proto
	github.com/kazoup/platform/lib/protomsg/crawler/crawler.proto
	github.com/kazoup/platform/lib/protomsg/deletebucket/deletebucket.proto
	github.com/kazoup/platform/lib/protomsg/deletefileinbucket/deletefileinbucket.proto
	github.com/kazoup/platform/lib/protomsg/enrich/enrich.proto

It has these top-level messages:
	AnnounceMessage
	CrawlerFinishedMessage
	DeleteBucketMsg
	DeleteFileInBucketMsg
	EnrichMessage
*/
package message

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

type AnnounceMessage struct {
	Handler string `protobuf:"bytes,2,opt,name=handler" json:"handler,omitempty"`
	Data    string `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *AnnounceMessage) Reset()                    { *m = AnnounceMessage{} }
func (m *AnnounceMessage) String() string            { return proto.CompactTextString(m) }
func (*AnnounceMessage) ProtoMessage()               {}
func (*AnnounceMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*AnnounceMessage)(nil), "message.AnnounceMessage")
}

func init() {
	proto.RegisterFile("github.com/kazoup/platform/lib/protomsg/announce/announce.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xb2, 0x4f, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0x4e, 0xac, 0xca, 0x2f, 0x2d, 0xd0, 0x2f, 0xc8,
	0x49, 0x2c, 0x49, 0xcb, 0x2f, 0xca, 0xd5, 0xcf, 0xc9, 0x4c, 0xd2, 0x2f, 0x28, 0xca, 0x2f, 0xc9,
	0xcf, 0x2d, 0x4e, 0xd7, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0xcd, 0x4b, 0x4e, 0x85, 0x33, 0xf4, 0xc0,
	0x52, 0x42, 0xec, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x4a, 0xf6, 0x5c, 0xfc, 0x8e, 0x50,
	0x29, 0x5f, 0x88, 0x90, 0x90, 0x04, 0x17, 0x7b, 0x46, 0x62, 0x5e, 0x4a, 0x4e, 0x6a, 0x91, 0x04,
	0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8c, 0x2b, 0x24, 0xc4, 0xc5, 0x92, 0x92, 0x58, 0x92, 0x28,
	0xc1, 0x0c, 0x16, 0x06, 0xb3, 0x93, 0xd8, 0xc0, 0x06, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x9e, 0xaa, 0xb5, 0x90, 0x93, 0x00, 0x00, 0x00,
}