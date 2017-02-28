// Code generated by protoc-gen-go.
// source: server.proto
// DO NOT EDIT!

/*
Package grpc is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	Package
*/
package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc1 "google.golang.org/grpc"
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

type Package struct {
	Type     int32  `protobuf:"varint,1,opt,name=Type" json:"Type,omitempty"`
	Encoding int32  `protobuf:"varint,2,opt,name=Encoding" json:"Encoding,omitempty"`
	ID       int32  `protobuf:"varint,3,opt,name=ID" json:"ID,omitempty"`
	Content  []byte `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (m *Package) Reset()                    { *m = Package{} }
func (m *Package) String() string            { return proto.CompactTextString(m) }
func (*Package) ProtoMessage()               {}
func (*Package) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Package) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Package) GetEncoding() int32 {
	if m != nil {
		return m.Encoding
	}
	return 0
}

func (m *Package) GetID() int32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Package) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*Package)(nil), "grpc.Package")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc1.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc1.SupportPackageIsVersion4

// Client API for Grpc service

type GrpcClient interface {
	Stream(ctx context.Context, opts ...grpc1.CallOption) (Grpc_StreamClient, error)
}

type grpcClient struct {
	cc *grpc1.ClientConn
}

func NewGrpcClient(cc *grpc1.ClientConn) GrpcClient {
	return &grpcClient{cc}
}

func (c *grpcClient) Stream(ctx context.Context, opts ...grpc1.CallOption) (Grpc_StreamClient, error) {
	stream, err := grpc1.NewClientStream(ctx, &_Grpc_serviceDesc.Streams[0], c.cc, "/grpc.Grpc/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcStreamClient{stream}
	return x, nil
}

type Grpc_StreamClient interface {
	Send(*Package) error
	Recv() (*Package, error)
	grpc1.ClientStream
}

type grpcStreamClient struct {
	grpc1.ClientStream
}

func (x *grpcStreamClient) Send(m *Package) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcStreamClient) Recv() (*Package, error) {
	m := new(Package)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Grpc service

type GrpcServer interface {
	Stream(Grpc_StreamServer) error
}

func RegisterGrpcServer(s *grpc1.Server, srv GrpcServer) {
	s.RegisterService(&_Grpc_serviceDesc, srv)
}

func _Grpc_Stream_Handler(srv interface{}, stream grpc1.ServerStream) error {
	return srv.(GrpcServer).Stream(&grpcStreamServer{stream})
}

type Grpc_StreamServer interface {
	Send(*Package) error
	Recv() (*Package, error)
	grpc1.ServerStream
}

type grpcStreamServer struct {
	grpc1.ServerStream
}

func (x *grpcStreamServer) Send(m *Package) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcStreamServer) Recv() (*Package, error) {
	m := new(Package)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Grpc_serviceDesc = grpc1.ServiceDesc{
	ServiceName: "grpc.Grpc",
	HandlerType: (*GrpcServer)(nil),
	Methods:     []grpc1.MethodDesc{},
	Streams: []grpc1.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Grpc_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "server.proto",
}

func init() { proto.RegisterFile("server.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2f, 0x2a, 0x48, 0x56, 0x4a,
	0xe6, 0x62, 0x0f, 0x48, 0x4c, 0xce, 0x4e, 0x4c, 0x4f, 0x15, 0x12, 0xe2, 0x62, 0x09, 0xa9, 0x2c,
	0x48, 0x95, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0d, 0x02, 0xb3, 0x85, 0xa4, 0xb8, 0x38, 0x5c, 0xf3,
	0x92, 0xf3, 0x53, 0x32, 0xf3, 0xd2, 0x25, 0x98, 0xc0, 0xe2, 0x70, 0xbe, 0x10, 0x1f, 0x17, 0x93,
	0xa7, 0x8b, 0x04, 0x33, 0x58, 0x94, 0xc9, 0xd3, 0x45, 0x48, 0x82, 0x8b, 0xdd, 0x39, 0x3f, 0xaf,
	0x24, 0x35, 0xaf, 0x44, 0x82, 0x45, 0x81, 0x51, 0x83, 0x27, 0x08, 0xc6, 0x35, 0x32, 0xe1, 0x62,
	0x71, 0x2f, 0x2a, 0x48, 0x16, 0xd2, 0xe1, 0x62, 0x0b, 0x2e, 0x29, 0x4a, 0x4d, 0xcc, 0x15, 0xe2,
	0xd5, 0x03, 0xd9, 0xae, 0x07, 0xb5, 0x5a, 0x0a, 0x95, 0xab, 0xc4, 0xa0, 0xc1, 0x68, 0xc0, 0x98,
	0xc4, 0x06, 0x76, 0xa7, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x02, 0xeb, 0xb3, 0x63, 0xb7, 0x00,
	0x00, 0x00,
}
