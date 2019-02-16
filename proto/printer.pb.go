// Code generated by protoc-gen-go. DO NOT EDIT.
// source: printer.proto

package printer

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
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

type Config struct {
	Requests             []*PrintRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f3006b867863bc3, []int{0}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetRequests() []*PrintRequest {
	if m != nil {
		return m.Requests
	}
	return nil
}

type PrintRequest struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Lines                []string `protobuf:"bytes,2,rep,name=lines,proto3" json:"lines,omitempty"`
	Origin               string   `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrintRequest) Reset()         { *m = PrintRequest{} }
func (m *PrintRequest) String() string { return proto.CompactTextString(m) }
func (*PrintRequest) ProtoMessage()    {}
func (*PrintRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f3006b867863bc3, []int{1}
}

func (m *PrintRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrintRequest.Unmarshal(m, b)
}
func (m *PrintRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrintRequest.Marshal(b, m, deterministic)
}
func (m *PrintRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrintRequest.Merge(m, src)
}
func (m *PrintRequest) XXX_Size() int {
	return xxx_messageInfo_PrintRequest.Size(m)
}
func (m *PrintRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PrintRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PrintRequest proto.InternalMessageInfo

func (m *PrintRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *PrintRequest) GetLines() []string {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *PrintRequest) GetOrigin() string {
	if m != nil {
		return m.Origin
	}
	return ""
}

type PrintResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrintResponse) Reset()         { *m = PrintResponse{} }
func (m *PrintResponse) String() string { return proto.CompactTextString(m) }
func (*PrintResponse) ProtoMessage()    {}
func (*PrintResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f3006b867863bc3, []int{2}
}

func (m *PrintResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrintResponse.Unmarshal(m, b)
}
func (m *PrintResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrintResponse.Marshal(b, m, deterministic)
}
func (m *PrintResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrintResponse.Merge(m, src)
}
func (m *PrintResponse) XXX_Size() int {
	return xxx_messageInfo_PrintResponse.Size(m)
}
func (m *PrintResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PrintResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PrintResponse proto.InternalMessageInfo

type ClearRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearRequest) Reset()         { *m = ClearRequest{} }
func (m *ClearRequest) String() string { return proto.CompactTextString(m) }
func (*ClearRequest) ProtoMessage()    {}
func (*ClearRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f3006b867863bc3, []int{3}
}

func (m *ClearRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearRequest.Unmarshal(m, b)
}
func (m *ClearRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearRequest.Marshal(b, m, deterministic)
}
func (m *ClearRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearRequest.Merge(m, src)
}
func (m *ClearRequest) XXX_Size() int {
	return xxx_messageInfo_ClearRequest.Size(m)
}
func (m *ClearRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClearRequest proto.InternalMessageInfo

type ClearResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearResponse) Reset()         { *m = ClearResponse{} }
func (m *ClearResponse) String() string { return proto.CompactTextString(m) }
func (*ClearResponse) ProtoMessage()    {}
func (*ClearResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f3006b867863bc3, []int{4}
}

func (m *ClearResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearResponse.Unmarshal(m, b)
}
func (m *ClearResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearResponse.Marshal(b, m, deterministic)
}
func (m *ClearResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearResponse.Merge(m, src)
}
func (m *ClearResponse) XXX_Size() int {
	return xxx_messageInfo_ClearResponse.Size(m)
}
func (m *ClearResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ClearResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Config)(nil), "printer.Config")
	proto.RegisterType((*PrintRequest)(nil), "printer.PrintRequest")
	proto.RegisterType((*PrintResponse)(nil), "printer.PrintResponse")
	proto.RegisterType((*ClearRequest)(nil), "printer.ClearRequest")
	proto.RegisterType((*ClearResponse)(nil), "printer.ClearResponse")
}

func init() { proto.RegisterFile("printer.proto", fileDescriptor_2f3006b867863bc3) }

var fileDescriptor_2f3006b867863bc3 = []byte{
	// 217 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4a, 0xc6, 0x30,
	0x10, 0x84, 0xff, 0x58, 0x5b, 0xed, 0xda, 0x2a, 0x2c, 0x5a, 0x82, 0xa7, 0x92, 0x53, 0x4f, 0x05,
	0xeb, 0x45, 0xf0, 0xd8, 0x17, 0x28, 0xf1, 0x09, 0x54, 0xd6, 0x12, 0x28, 0x49, 0x4d, 0xa2, 0x78,
	0xf6, 0xc9, 0xa5, 0x69, 0x94, 0x28, 0x9e, 0x92, 0x2f, 0x3b, 0xb3, 0x33, 0x04, 0xea, 0xd5, 0x2a,
	0xed, 0xc9, 0xf6, 0xab, 0x35, 0xde, 0xe0, 0x49, 0x44, 0x71, 0x0f, 0xc5, 0x68, 0xf4, 0x8b, 0x9a,
	0xf1, 0x06, 0x4e, 0x2d, 0xbd, 0xbe, 0x91, 0xf3, 0x8e, 0xb3, 0x36, 0xeb, 0xce, 0x86, 0xab, 0xfe,
	0xdb, 0x34, 0x6d, 0xa7, 0xdc, 0xa7, 0xf2, 0x47, 0x26, 0x26, 0xa8, 0xd2, 0x09, 0x22, 0x1c, 0x7b,
	0xfa, 0xf0, 0x9c, 0xb5, 0xac, 0x2b, 0x65, 0xb8, 0xe3, 0x25, 0xe4, 0x8b, 0xd2, 0xe4, 0xf8, 0x51,
	0x9b, 0x75, 0xa5, 0xdc, 0x01, 0x1b, 0x28, 0x8c, 0x55, 0xb3, 0xd2, 0x3c, 0x0b, 0xda, 0x48, 0xe2,
	0x02, 0xea, 0xb8, 0xd1, 0xad, 0x46, 0x3b, 0x12, 0xe7, 0x50, 0x8d, 0x0b, 0x3d, 0xda, 0x18, 0xb1,
	0x09, 0x22, 0xef, 0x82, 0xe1, 0x93, 0xc5, 0x12, 0x0f, 0x64, 0xdf, 0xd5, 0x33, 0xe1, 0x1d, 0xe4,
	0x81, 0xf1, 0xff, 0xfa, 0xd7, 0xcd, 0xdf, 0xe7, 0x98, 0x74, 0xd8, 0x9c, 0x61, 0x77, 0xe2, 0x4c,
	0xb3, 0x13, 0xe7, 0xaf, 0x0a, 0xe2, 0xf0, 0x54, 0x84, 0x5f, 0xbd, 0xfd, 0x0a, 0x00, 0x00, 0xff,
	0xff, 0xa6, 0x58, 0xdb, 0x34, 0x66, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PrintServiceClient is the client API for PrintService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PrintServiceClient interface {
	Print(ctx context.Context, in *PrintRequest, opts ...grpc.CallOption) (*PrintResponse, error)
	Clear(ctx context.Context, in *ClearRequest, opts ...grpc.CallOption) (*ClearResponse, error)
}

type printServiceClient struct {
	cc *grpc.ClientConn
}

func NewPrintServiceClient(cc *grpc.ClientConn) PrintServiceClient {
	return &printServiceClient{cc}
}

func (c *printServiceClient) Print(ctx context.Context, in *PrintRequest, opts ...grpc.CallOption) (*PrintResponse, error) {
	out := new(PrintResponse)
	err := c.cc.Invoke(ctx, "/printer.PrintService/Print", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *printServiceClient) Clear(ctx context.Context, in *ClearRequest, opts ...grpc.CallOption) (*ClearResponse, error) {
	out := new(ClearResponse)
	err := c.cc.Invoke(ctx, "/printer.PrintService/Clear", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrintServiceServer is the server API for PrintService service.
type PrintServiceServer interface {
	Print(context.Context, *PrintRequest) (*PrintResponse, error)
	Clear(context.Context, *ClearRequest) (*ClearResponse, error)
}

func RegisterPrintServiceServer(s *grpc.Server, srv PrintServiceServer) {
	s.RegisterService(&_PrintService_serviceDesc, srv)
}

func _PrintService_Print_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrintServiceServer).Print(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/printer.PrintService/Print",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrintServiceServer).Print(ctx, req.(*PrintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrintService_Clear_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrintServiceServer).Clear(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/printer.PrintService/Clear",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrintServiceServer).Clear(ctx, req.(*ClearRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PrintService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "printer.PrintService",
	HandlerType: (*PrintServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Print",
			Handler:    _PrintService_Print_Handler,
		},
		{
			MethodName: "Clear",
			Handler:    _PrintService_Clear_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "printer.proto",
}
