// Code generated by protoc-gen-go.
// source: cluster.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	cluster.proto
	rpc.proto

It has these top-level messages:
	ServerStatus
	HeartRequest
	HeartRespose
	RpcHeart
	RpcNil
	RpcHandlers
	RpcHandler
*/
package pb

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

// 3秒上传一次到center服务器
type ServerStatus struct {
	Addr             *string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	Load             *int32  `protobuf:"varint,2,opt,name=load" json:"load,omitempty"`
	Sid              *int32  `protobuf:"varint,3,opt,name=sid" json:"sid,omitempty"`
	Stype            *string `protobuf:"bytes,4,opt,name=stype" json:"stype,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ServerStatus) Reset()                    { *m = ServerStatus{} }
func (m *ServerStatus) String() string            { return proto.CompactTextString(m) }
func (*ServerStatus) ProtoMessage()               {}
func (*ServerStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServerStatus) GetAddr() string {
	if m != nil && m.Addr != nil {
		return *m.Addr
	}
	return ""
}

func (m *ServerStatus) GetLoad() int32 {
	if m != nil && m.Load != nil {
		return *m.Load
	}
	return 0
}

func (m *ServerStatus) GetSid() int32 {
	if m != nil && m.Sid != nil {
		return *m.Sid
	}
	return 0
}

func (m *ServerStatus) GetStype() string {
	if m != nil && m.Stype != nil {
		return *m.Stype
	}
	return ""
}

type HeartRequest struct {
	Version          *uint32       `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Status           *ServerStatus `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *HeartRequest) Reset()                    { *m = HeartRequest{} }
func (m *HeartRequest) String() string            { return proto.CompactTextString(m) }
func (*HeartRequest) ProtoMessage()               {}
func (*HeartRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HeartRequest) GetVersion() uint32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *HeartRequest) GetStatus() *ServerStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

type HeartRespose struct {
	Version          *uint32         `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	StatusList       []*ServerStatus `protobuf:"bytes,2,rep,name=statusList" json:"statusList,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *HeartRespose) Reset()                    { *m = HeartRespose{} }
func (m *HeartRespose) String() string            { return proto.CompactTextString(m) }
func (*HeartRespose) ProtoMessage()               {}
func (*HeartRespose) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *HeartRespose) GetVersion() uint32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *HeartRespose) GetStatusList() []*ServerStatus {
	if m != nil {
		return m.StatusList
	}
	return nil
}

func init() {
	proto.RegisterType((*ServerStatus)(nil), "pb.ServerStatus")
	proto.RegisterType((*HeartRequest)(nil), "pb.HeartRequest")
	proto.RegisterType((*HeartRespose)(nil), "pb.HeartRespose")
}

func init() { proto.RegisterFile("cluster.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0xcf, 0xbf, 0x6a, 0xc3, 0x30,
	0x10, 0xc7, 0x71, 0xfc, 0xaf, 0xc5, 0x67, 0x9b, 0x16, 0x4d, 0xa2, 0x93, 0x31, 0x1d, 0x3a, 0x79,
	0xe8, 0x1b, 0xb4, 0x90, 0x90, 0x21, 0x53, 0xbc, 0x65, 0x93, 0xad, 0x1b, 0x0c, 0x26, 0x52, 0x74,
	0xe7, 0x40, 0xde, 0x3e, 0xb2, 0x0c, 0x21, 0x90, 0x8c, 0x5f, 0x38, 0x3e, 0x3f, 0x0e, 0xaa, 0x61,
	0x9a, 0x89, 0xd1, 0xb5, 0xd6, 0x19, 0x36, 0x22, 0xb6, 0xfd, 0x57, 0xee, 0xec, 0xb0, 0x66, 0xb3,
	0x85, 0xb2, 0x43, 0x77, 0x41, 0xd7, 0xb1, 0xe2, 0x99, 0x44, 0x09, 0xa9, 0xd2, 0xda, 0xc9, 0xa8,
	0x8e, 0x7e, 0xf2, 0xa5, 0x26, 0xa3, 0xb4, 0x8c, 0x7d, 0x65, 0xa2, 0x80, 0x84, 0x46, 0x2d, 0x93,
	0x10, 0x15, 0x64, 0xc4, 0x57, 0x8b, 0x32, 0x5d, 0x2e, 0x9b, 0x3f, 0x28, 0x77, 0xa8, 0x1c, 0x1f,
	0xf0, 0x3c, 0x23, 0xb1, 0xf8, 0x80, 0x77, 0x8f, 0xd2, 0x68, 0x4e, 0x81, 0xaa, 0x44, 0x0d, 0x6f,
	0x14, 0x26, 0x02, 0x56, 0xfc, 0x7e, 0xb6, 0xb6, 0x6f, 0x1f, 0xa7, 0x9b, 0xcd, 0x9d, 0x20, 0x6b,
	0x08, 0x9f, 0x89, 0x6f, 0x80, 0x95, 0xd8, 0x8f, 0xc4, 0x9e, 0x49, 0x5e, 0x31, 0xff, 0xe9, 0xd1,
	0xbf, 0x78, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x19, 0x89, 0x42, 0xe0, 0xf6, 0x00, 0x00, 0x00,
}
