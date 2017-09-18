// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg.HandShake.proto

package pkg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type HostPort struct {
	Host string `protobuf:"bytes,1,opt,name=Host" json:"Host,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=Port" json:"Port,omitempty"`
}

func (m *HostPort) Reset()                    { *m = HostPort{} }
func (m *HostPort) String() string            { return proto.CompactTextString(m) }
func (*HostPort) ProtoMessage()               {}
func (*HostPort) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *HostPort) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *HostPort) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type HandShakeClientData struct {
	ClientType    string `protobuf:"bytes,1,opt,name=ClientType" json:"ClientType,omitempty"`
	ClientVersion string `protobuf:"bytes,2,opt,name=ClientVersion" json:"ClientVersion,omitempty"`
	DictMd5       string `protobuf:"bytes,3,opt,name=DictMd5" json:"DictMd5,omitempty"`
}

func (m *HandShakeClientData) Reset()                    { *m = HandShakeClientData{} }
func (m *HandShakeClientData) String() string            { return proto.CompactTextString(m) }
func (*HandShakeClientData) ProtoMessage()               {}
func (*HandShakeClientData) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *HandShakeClientData) GetClientType() string {
	if m != nil {
		return m.ClientType
	}
	return ""
}

func (m *HandShakeClientData) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

func (m *HandShakeClientData) GetDictMd5() string {
	if m != nil {
		return m.DictMd5
	}
	return ""
}

type HandShakeResponse struct {
	ServerVersion string            `protobuf:"bytes,1,opt,name=ServerVersion" json:"ServerVersion,omitempty"`
	Now           string            `protobuf:"bytes,2,opt,name=Now" json:"Now,omitempty"`
	HeartBeatRate int32             `protobuf:"varint,3,opt,name=HeartBeatRate" json:"HeartBeatRate,omitempty"`
	Routes        map[string]uint32 `protobuf:"bytes,4,rep,name=Routes" json:"Routes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	IsReconnect   bool              `protobuf:"varint,5,opt,name=IsReconnect" json:"IsReconnect,omitempty"`
	ReconnectTo   *HostPort         `protobuf:"bytes,6,opt,name=ReconnectTo" json:"ReconnectTo,omitempty"`
}

func (m *HandShakeResponse) Reset()                    { *m = HandShakeResponse{} }
func (m *HandShakeResponse) String() string            { return proto.CompactTextString(m) }
func (*HandShakeResponse) ProtoMessage()               {}
func (*HandShakeResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *HandShakeResponse) GetServerVersion() string {
	if m != nil {
		return m.ServerVersion
	}
	return ""
}

func (m *HandShakeResponse) GetNow() string {
	if m != nil {
		return m.Now
	}
	return ""
}

func (m *HandShakeResponse) GetHeartBeatRate() int32 {
	if m != nil {
		return m.HeartBeatRate
	}
	return 0
}

func (m *HandShakeResponse) GetRoutes() map[string]uint32 {
	if m != nil {
		return m.Routes
	}
	return nil
}

func (m *HandShakeResponse) GetIsReconnect() bool {
	if m != nil {
		return m.IsReconnect
	}
	return false
}

func (m *HandShakeResponse) GetReconnectTo() *HostPort {
	if m != nil {
		return m.ReconnectTo
	}
	return nil
}

func init() {
	proto.RegisterType((*HostPort)(nil), "pkg.HostPort")
	proto.RegisterType((*HandShakeClientData)(nil), "pkg.HandShakeClientData")
	proto.RegisterType((*HandShakeResponse)(nil), "pkg.HandShakeResponse")
}

func init() { proto.RegisterFile("pkg.HandShake.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x52, 0xc1, 0x4e, 0xeb, 0x30,
	0x10, 0x54, 0x92, 0xa6, 0xaf, 0xdd, 0xa8, 0x4f, 0xd4, 0xe5, 0x60, 0x71, 0x40, 0x51, 0xc4, 0x21,
	0xa7, 0x20, 0x15, 0x21, 0x41, 0x8f, 0xa5, 0x88, 0x72, 0x00, 0x55, 0x6e, 0xc5, 0xdd, 0xa4, 0xab,
	0x52, 0xa5, 0xb2, 0x23, 0xc7, 0x2d, 0xca, 0x2f, 0xf1, 0x0b, 0xfc, 0x1c, 0xb2, 0x93, 0x46, 0xc9,
	0x6d, 0x66, 0x76, 0x77, 0xc6, 0x6b, 0x1b, 0x26, 0x79, 0xb6, 0x4b, 0x96, 0x5c, 0x6c, 0xd7, 0x5f,
	0x3c, 0xc3, 0x24, 0x57, 0x52, 0x4b, 0xe2, 0xe5, 0xd9, 0x2e, 0x9a, 0xc2, 0x60, 0x29, 0x0b, 0xbd,
	0x92, 0x4a, 0x13, 0x02, 0x3d, 0x83, 0xa9, 0x13, 0x3a, 0xf1, 0x90, 0x59, 0x6c, 0x34, 0x53, 0xa3,
	0x6e, 0xe8, 0xc4, 0x3e, 0xb3, 0x38, 0x3a, 0xc2, 0xa4, 0xf1, 0x7a, 0x3a, 0xec, 0x51, 0xe8, 0x05,
	0xd7, 0x9c, 0x5c, 0x03, 0x54, 0x6c, 0x53, 0xe6, 0x58, 0x9b, 0xb4, 0x14, 0x72, 0x03, 0xa3, 0x8a,
	0x7d, 0xa0, 0x2a, 0xf6, 0x52, 0x58, 0xcf, 0x21, 0xeb, 0x8a, 0x84, 0xc2, 0xbf, 0xc5, 0x3e, 0xd5,
	0x6f, 0xdb, 0x7b, 0xea, 0xd9, 0xfa, 0x99, 0x46, 0xbf, 0x2e, 0x8c, 0x9b, 0x5c, 0x86, 0x45, 0x2e,
	0x45, 0x61, 0x5d, 0xd7, 0xa8, 0x4e, 0xa8, 0xce, 0xae, 0x55, 0x70, 0x57, 0x24, 0x17, 0xe0, 0xbd,
	0xcb, 0xef, 0x3a, 0xd1, 0x40, 0x33, 0xb7, 0x44, 0xae, 0xf4, 0x1c, 0xb9, 0x66, 0x5c, 0xa3, 0x4d,
	0xf3, 0x59, 0x57, 0x24, 0x33, 0xe8, 0x33, 0x79, 0xd4, 0x58, 0xd0, 0x5e, 0xe8, 0xc5, 0xc1, 0x34,
	0x4a, 0x3a, 0x37, 0x79, 0x3e, 0x45, 0x52, 0x35, 0x3d, 0x0b, 0xad, 0x4a, 0x56, 0x4f, 0x90, 0x10,
	0x82, 0xd7, 0x82, 0x61, 0x2a, 0x85, 0xc0, 0x54, 0x53, 0x3f, 0x74, 0xe2, 0x01, 0x6b, 0x4b, 0xe4,
	0x16, 0x82, 0x86, 0x6c, 0x24, 0xed, 0x87, 0x4e, 0x1c, 0x4c, 0x47, 0x55, 0x44, 0xfd, 0x28, 0xac,
	0xdd, 0x71, 0xf5, 0x08, 0x41, 0x2b, 0xc9, 0x6c, 0x95, 0x61, 0x59, 0x6f, 0x6c, 0x20, 0xb9, 0x04,
	0xff, 0xc4, 0x0f, 0x47, 0xb4, 0x9b, 0x8e, 0x58, 0x45, 0x66, 0xee, 0x83, 0x33, 0x1f, 0xff, 0xb8,
	0xff, 0x5f, 0xe4, 0xea, 0xc0, 0xcb, 0x64, 0xc5, 0xd3, 0x8c, 0xef, 0xf0, 0xb3, 0x6f, 0xff, 0xc1,
	0xdd, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x79, 0x8e, 0x13, 0x3f, 0x1e, 0x02, 0x00, 0x00,
}
