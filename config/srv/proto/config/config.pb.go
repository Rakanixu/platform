// Code generated by protoc-gen-go.
// source: github.com/kazoup/config/srv/proto/config/config.proto
// DO NOT EDIT!

/*
Package go_micro_srv_config is a generated protocol buffer package.

It is generated from these files:
	github.com/kazoup/config/srv/proto/config/config.proto

It has these top-level messages:
	StatusRequest
	StatusResponse
	SetElasticSettingsRequest
	SetElasticSettingsResponse
	SetElasticMappingRequest
	SetElasticMappingResponse
	SetFlagsRequest
	SetFlagsResponse
*/
package go_micro_srv_config

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

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StatusResponse struct {
	APPLIANCE_IS_CONFIGURED bool   `protobuf:"varint,1,opt,name=APPLIANCE_IS_CONFIGURED,json=aPPLIANCEISCONFIGURED" json:"APPLIANCE_IS_CONFIGURED,omitempty"`
	APPLIANCE_IS_DEMO       bool   `protobuf:"varint,2,opt,name=APPLIANCE_IS_DEMO,json=aPPLIANCEISDEMO" json:"APPLIANCE_IS_DEMO,omitempty"`
	APPLIANCE_IS_REGISTERED bool   `protobuf:"varint,3,opt,name=APPLIANCE_IS_REGISTERED,json=aPPLIANCEISREGISTERED" json:"APPLIANCE_IS_REGISTERED,omitempty"`
	GIT_COMMIT_STRING       string `protobuf:"bytes,4,opt,name=GIT_COMMIT_STRING,json=gITCOMMITSTRING" json:"GIT_COMMIT_STRING,omitempty"`
	SMB_USER_EXISTS         bool   `protobuf:"varint,5,opt,name=SMB_USER_EXISTS,json=sMBUSEREXISTS" json:"SMB_USER_EXISTS,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type SetElasticSettingsRequest struct {
}

func (m *SetElasticSettingsRequest) Reset()                    { *m = SetElasticSettingsRequest{} }
func (m *SetElasticSettingsRequest) String() string            { return proto.CompactTextString(m) }
func (*SetElasticSettingsRequest) ProtoMessage()               {}
func (*SetElasticSettingsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type SetElasticSettingsResponse struct {
}

func (m *SetElasticSettingsResponse) Reset()                    { *m = SetElasticSettingsResponse{} }
func (m *SetElasticSettingsResponse) String() string            { return proto.CompactTextString(m) }
func (*SetElasticSettingsResponse) ProtoMessage()               {}
func (*SetElasticSettingsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type SetElasticMappingRequest struct {
}

func (m *SetElasticMappingRequest) Reset()                    { *m = SetElasticMappingRequest{} }
func (m *SetElasticMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*SetElasticMappingRequest) ProtoMessage()               {}
func (*SetElasticMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type SetElasticMappingResponse struct {
}

func (m *SetElasticMappingResponse) Reset()                    { *m = SetElasticMappingResponse{} }
func (m *SetElasticMappingResponse) String() string            { return proto.CompactTextString(m) }
func (*SetElasticMappingResponse) ProtoMessage()               {}
func (*SetElasticMappingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type SetFlagsRequest struct {
}

func (m *SetFlagsRequest) Reset()                    { *m = SetFlagsRequest{} }
func (m *SetFlagsRequest) String() string            { return proto.CompactTextString(m) }
func (*SetFlagsRequest) ProtoMessage()               {}
func (*SetFlagsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type SetFlagsResponse struct {
}

func (m *SetFlagsResponse) Reset()                    { *m = SetFlagsResponse{} }
func (m *SetFlagsResponse) String() string            { return proto.CompactTextString(m) }
func (*SetFlagsResponse) ProtoMessage()               {}
func (*SetFlagsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*StatusRequest)(nil), "go.micro.srv.config.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "go.micro.srv.config.StatusResponse")
	proto.RegisterType((*SetElasticSettingsRequest)(nil), "go.micro.srv.config.SetElasticSettingsRequest")
	proto.RegisterType((*SetElasticSettingsResponse)(nil), "go.micro.srv.config.SetElasticSettingsResponse")
	proto.RegisterType((*SetElasticMappingRequest)(nil), "go.micro.srv.config.SetElasticMappingRequest")
	proto.RegisterType((*SetElasticMappingResponse)(nil), "go.micro.srv.config.SetElasticMappingResponse")
	proto.RegisterType((*SetFlagsRequest)(nil), "go.micro.srv.config.SetFlagsRequest")
	proto.RegisterType((*SetFlagsResponse)(nil), "go.micro.srv.config.SetFlagsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Config service

type ConfigClient interface {
	Status(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error)
	SetElasticSettings(ctx context.Context, in *SetElasticSettingsRequest, opts ...client.CallOption) (*SetElasticSettingsResponse, error)
	SetElasticMapping(ctx context.Context, in *SetElasticMappingRequest, opts ...client.CallOption) (*SetElasticMappingResponse, error)
	SetFlags(ctx context.Context, in *SetFlagsRequest, opts ...client.CallOption) (*SetFlagsResponse, error)
}

type configClient struct {
	c           client.Client
	serviceName string
}

func NewConfigClient(serviceName string, c client.Client) ConfigClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.config"
	}
	return &configClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *configClient) Status(ctx context.Context, in *StatusRequest, opts ...client.CallOption) (*StatusResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Config.Status", in)
	out := new(StatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) SetElasticSettings(ctx context.Context, in *SetElasticSettingsRequest, opts ...client.CallOption) (*SetElasticSettingsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Config.SetElasticSettings", in)
	out := new(SetElasticSettingsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) SetElasticMapping(ctx context.Context, in *SetElasticMappingRequest, opts ...client.CallOption) (*SetElasticMappingResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Config.SetElasticMapping", in)
	out := new(SetElasticMappingResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) SetFlags(ctx context.Context, in *SetFlagsRequest, opts ...client.CallOption) (*SetFlagsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Config.SetFlags", in)
	out := new(SetFlagsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Config service

type ConfigHandler interface {
	Status(context.Context, *StatusRequest, *StatusResponse) error
	SetElasticSettings(context.Context, *SetElasticSettingsRequest, *SetElasticSettingsResponse) error
	SetElasticMapping(context.Context, *SetElasticMappingRequest, *SetElasticMappingResponse) error
	SetFlags(context.Context, *SetFlagsRequest, *SetFlagsResponse) error
}

func RegisterConfigHandler(s server.Server, hdlr ConfigHandler) {
	s.Handle(s.NewHandler(&Config{hdlr}))
}

type Config struct {
	ConfigHandler
}

func (h *Config) Status(ctx context.Context, in *StatusRequest, out *StatusResponse) error {
	return h.ConfigHandler.Status(ctx, in, out)
}

func (h *Config) SetElasticSettings(ctx context.Context, in *SetElasticSettingsRequest, out *SetElasticSettingsResponse) error {
	return h.ConfigHandler.SetElasticSettings(ctx, in, out)
}

func (h *Config) SetElasticMapping(ctx context.Context, in *SetElasticMappingRequest, out *SetElasticMappingResponse) error {
	return h.ConfigHandler.SetElasticMapping(ctx, in, out)
}

func (h *Config) SetFlags(ctx context.Context, in *SetFlagsRequest, out *SetFlagsResponse) error {
	return h.ConfigHandler.SetFlags(ctx, in, out)
}

var fileDescriptor0 = []byte{
	// 400 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x53, 0x5d, 0xab, 0xd3, 0x30,
	0x18, 0x76, 0x9b, 0x8e, 0x19, 0x98, 0x75, 0x11, 0xb1, 0x56, 0x2f, 0xa4, 0x7e, 0x20, 0x82, 0x29,
	0x28, 0xec, 0x7e, 0x1f, 0x59, 0x09, 0xd8, 0x6d, 0x34, 0x1d, 0xe8, 0x55, 0xe9, 0x4a, 0xad, 0xc5,
	0x6d, 0xa9, 0x6b, 0xaa, 0xe0, 0xd5, 0xf9, 0x37, 0xe7, 0x6f, 0x9e, 0xb4, 0x59, 0xdb, 0x95, 0x6e,
	0x67, 0xbb, 0xca, 0x78, 0x3e, 0xde, 0xe7, 0xdd, 0xf3, 0x52, 0x30, 0x0c, 0x23, 0xfe, 0x2b, 0x5d,
	0x23, 0x9f, 0x6d, 0x8d, 0xdf, 0xde, 0x7f, 0x96, 0xc6, 0x86, 0xcf, 0x76, 0x3f, 0xa3, 0xd0, 0x48,
	0xf6, 0x7f, 0x8d, 0x78, 0xcf, 0x38, 0x2b, 0x00, 0xf9, 0xa0, 0x1c, 0x83, 0xcf, 0x42, 0x86, 0xb6,
	0x91, 0xbf, 0x67, 0x48, 0xe8, 0x90, 0xa4, 0x74, 0x05, 0xf4, 0x29, 0xf7, 0x78, 0x9a, 0xd8, 0xc1,
	0x9f, 0x34, 0x48, 0xb8, 0x7e, 0xd3, 0x06, 0x4f, 0x0a, 0x24, 0x89, 0xd9, 0x2e, 0x09, 0xe0, 0x10,
	0xbc, 0x18, 0x2d, 0x97, 0xdf, 0xc8, 0x68, 0x3e, 0xc1, 0x2e, 0xa1, 0xee, 0x64, 0x31, 0x9f, 0x11,
	0x73, 0x65, 0xe3, 0xa9, 0xda, 0x7a, 0xd3, 0xfa, 0xd8, 0xb3, 0x9f, 0x7b, 0x05, 0x4d, 0x68, 0x45,
	0xc2, 0x4f, 0x60, 0x50, 0xf3, 0x4d, 0xb1, 0xb5, 0x50, 0xdb, 0xb9, 0x43, 0x39, 0x72, 0x64, 0x70,
	0x23, 0xc3, 0xc6, 0x26, 0xa1, 0x0e, 0xce, 0x32, 0x3a, 0x8d, 0x8c, 0x8a, 0xcc, 0x32, 0x4c, 0xe2,
	0x88, 0x95, 0x2c, 0x4b, 0x3c, 0xd4, 0xb1, 0xc9, 0xdc, 0x54, 0x1f, 0x0a, 0xc7, 0x63, 0x5b, 0x09,
	0x89, 0x23, 0x71, 0x09, 0xc3, 0x0f, 0x40, 0xa1, 0xd6, 0xd8, 0x5d, 0x51, 0x6c, 0xbb, 0xf8, 0xbb,
	0x98, 0x40, 0xd5, 0x47, 0xf9, 0xec, 0x7e, 0x62, 0x8d, 0x33, 0x54, 0x82, 0xfa, 0x2b, 0xf0, 0x92,
	0x06, 0x1c, 0x6f, 0xbc, 0x84, 0x47, 0xbe, 0xf8, 0xc5, 0xa3, 0x5d, 0x58, 0xf6, 0xf3, 0x1a, 0x68,
	0xa7, 0x48, 0x59, 0x95, 0xae, 0x01, 0xb5, 0x62, 0x2d, 0x2f, 0x8e, 0x05, 0x5b, 0x38, 0x6b, 0x63,
	0x4b, 0xee, 0x60, 0x1c, 0x88, 0xdd, 0x02, 0x3e, 0xdb, 0x78, 0x55, 0x12, 0x04, 0x4f, 0x2b, 0x48,
	0xca, 0xbe, 0xdc, 0x76, 0x40, 0x77, 0x92, 0x5f, 0x0e, 0x52, 0xd0, 0x95, 0x77, 0x82, 0x3a, 0x3a,
	0x71, 0x59, 0x54, 0x3b, 0xab, 0xf6, 0xf6, 0x5e, 0xcd, 0x61, 0x89, 0x07, 0xf0, 0x1f, 0x80, 0xcd,
	0x7f, 0x07, 0xd1, 0x69, 0xf3, 0xb9, 0x8e, 0x34, 0xe3, 0x6a, 0x7d, 0x19, 0xcc, 0xc1, 0xa0, 0x51,
	0x0e, 0xfc, 0x7c, 0x61, 0x4e, 0xbd, 0x60, 0x0d, 0x5d, 0x2b, 0x2f, 0x53, 0x7f, 0x80, 0x5e, 0x51,
	0x31, 0x7c, 0x77, 0xce, 0x7d, 0x7c, 0x14, 0xed, 0xfd, 0x05, 0x55, 0x31, 0x7a, 0xdd, 0xcd, 0x3f,
	0xba, 0xaf, 0x77, 0x01, 0x00, 0x00, 0xff, 0xff, 0x33, 0x0a, 0x9f, 0x6b, 0xae, 0x03, 0x00, 0x00,
}
