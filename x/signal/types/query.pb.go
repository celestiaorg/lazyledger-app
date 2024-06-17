// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: celestia/signal/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryVersionTallyRequest is the request type for the VersionTally query.
type QueryVersionTallyRequest struct {
	Version uint64 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (m *QueryVersionTallyRequest) Reset()         { *m = QueryVersionTallyRequest{} }
func (m *QueryVersionTallyRequest) String() string { return proto.CompactTextString(m) }
func (*QueryVersionTallyRequest) ProtoMessage()    {}
func (*QueryVersionTallyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7af24246367e432c, []int{0}
}
func (m *QueryVersionTallyRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryVersionTallyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryVersionTallyRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryVersionTallyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryVersionTallyRequest.Merge(m, src)
}
func (m *QueryVersionTallyRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryVersionTallyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryVersionTallyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryVersionTallyRequest proto.InternalMessageInfo

func (m *QueryVersionTallyRequest) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

// QueryVersionTallyResponse is the response type for the VersionTally query.
type QueryVersionTallyResponse struct {
	VotingPower      uint64 `protobuf:"varint,1,opt,name=voting_power,json=votingPower,proto3" json:"voting_power,omitempty"`
	ThresholdPower   uint64 `protobuf:"varint,2,opt,name=threshold_power,json=thresholdPower,proto3" json:"threshold_power,omitempty"`
	TotalVotingPower uint64 `protobuf:"varint,3,opt,name=total_voting_power,json=totalVotingPower,proto3" json:"total_voting_power,omitempty"`
}

func (m *QueryVersionTallyResponse) Reset()         { *m = QueryVersionTallyResponse{} }
func (m *QueryVersionTallyResponse) String() string { return proto.CompactTextString(m) }
func (*QueryVersionTallyResponse) ProtoMessage()    {}
func (*QueryVersionTallyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7af24246367e432c, []int{1}
}
func (m *QueryVersionTallyResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryVersionTallyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryVersionTallyResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryVersionTallyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryVersionTallyResponse.Merge(m, src)
}
func (m *QueryVersionTallyResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryVersionTallyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryVersionTallyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryVersionTallyResponse proto.InternalMessageInfo

func (m *QueryVersionTallyResponse) GetVotingPower() uint64 {
	if m != nil {
		return m.VotingPower
	}
	return 0
}

func (m *QueryVersionTallyResponse) GetThresholdPower() uint64 {
	if m != nil {
		return m.ThresholdPower
	}
	return 0
}

func (m *QueryVersionTallyResponse) GetTotalVotingPower() uint64 {
	if m != nil {
		return m.TotalVotingPower
	}
	return 0
}

// QueryGetUpgradeRequest is the request type for the GetUpgrade query.
type QueryGetUpgradeRequest struct {
}

func (m *QueryGetUpgradeRequest) Reset()         { *m = QueryGetUpgradeRequest{} }
func (m *QueryGetUpgradeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetUpgradeRequest) ProtoMessage()    {}
func (*QueryGetUpgradeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7af24246367e432c, []int{2}
}
func (m *QueryGetUpgradeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetUpgradeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetUpgradeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetUpgradeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetUpgradeRequest.Merge(m, src)
}
func (m *QueryGetUpgradeRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetUpgradeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetUpgradeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetUpgradeRequest proto.InternalMessageInfo

// QueryGetUpgradeResponse is the response type for the GetUpgrade query.
type QueryGetUpgradeResponse struct {
	Upgrade *Upgrade `protobuf:"bytes,1,opt,name=upgrade,proto3" json:"upgrade,omitempty"`
}

func (m *QueryGetUpgradeResponse) Reset()         { *m = QueryGetUpgradeResponse{} }
func (m *QueryGetUpgradeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetUpgradeResponse) ProtoMessage()    {}
func (*QueryGetUpgradeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7af24246367e432c, []int{3}
}
func (m *QueryGetUpgradeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetUpgradeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetUpgradeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetUpgradeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetUpgradeResponse.Merge(m, src)
}
func (m *QueryGetUpgradeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetUpgradeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetUpgradeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetUpgradeResponse proto.InternalMessageInfo

func (m *QueryGetUpgradeResponse) GetUpgrade() *Upgrade {
	if m != nil {
		return m.Upgrade
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryVersionTallyRequest)(nil), "celestia.signal.v1.QueryVersionTallyRequest")
	proto.RegisterType((*QueryVersionTallyResponse)(nil), "celestia.signal.v1.QueryVersionTallyResponse")
	proto.RegisterType((*QueryGetUpgradeRequest)(nil), "celestia.signal.v1.QueryGetUpgradeRequest")
	proto.RegisterType((*QueryGetUpgradeResponse)(nil), "celestia.signal.v1.QueryGetUpgradeResponse")
}

func init() { proto.RegisterFile("celestia/signal/v1/query.proto", fileDescriptor_7af24246367e432c) }

var fileDescriptor_7af24246367e432c = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcf, 0x6e, 0xda, 0x40,
	0x10, 0xc6, 0x59, 0xfa, 0x07, 0x69, 0x41, 0x6d, 0xb5, 0xaa, 0x5a, 0xd7, 0xb4, 0x16, 0x75, 0x0f,
	0xad, 0x5a, 0xf0, 0x0a, 0xda, 0xbe, 0x40, 0x2e, 0xb9, 0xe4, 0x40, 0x50, 0xc2, 0x21, 0x17, 0xb4,
	0xc0, 0xca, 0x58, 0x72, 0xbc, 0x8b, 0x77, 0xed, 0x04, 0x45, 0x39, 0x24, 0x4f, 0x80, 0x14, 0xe5,
	0x9a, 0xe7, 0xc9, 0x11, 0x29, 0x97, 0x1c, 0x23, 0xc8, 0x83, 0x44, 0xde, 0xb5, 0x21, 0x91, 0x41,
	0xe2, 0x66, 0xcf, 0xfc, 0xe6, 0x9b, 0x6f, 0x66, 0x16, 0x5a, 0x03, 0xea, 0x53, 0x21, 0x3d, 0x82,
	0x85, 0xe7, 0x06, 0xc4, 0xc7, 0x71, 0x13, 0x8f, 0x23, 0x1a, 0x4e, 0x1c, 0x1e, 0x32, 0xc9, 0x10,
	0xca, 0xf2, 0x8e, 0xce, 0x3b, 0x71, 0xd3, 0xfc, 0xea, 0x32, 0xe6, 0xfa, 0x14, 0x13, 0xee, 0x61,
	0x12, 0x04, 0x4c, 0x12, 0xe9, 0xb1, 0x40, 0xe8, 0x0a, 0xb3, 0xb6, 0x46, 0x31, 0xe2, 0x6e, 0x48,
	0x86, 0x54, 0x13, 0xf6, 0x3f, 0x68, 0xec, 0x27, 0x2d, 0xba, 0x34, 0x14, 0x1e, 0x0b, 0x0e, 0x88,
	0xef, 0x4f, 0x3a, 0x74, 0x1c, 0x51, 0x21, 0x91, 0x01, 0x4b, 0xb1, 0x0e, 0x1b, 0xa0, 0x06, 0x7e,
	0xbd, 0xee, 0x64, 0xbf, 0xf6, 0x35, 0x80, 0x5f, 0xd6, 0x94, 0x09, 0xce, 0x02, 0x41, 0xd1, 0x77,
	0x58, 0x89, 0x99, 0xf4, 0x02, 0xb7, 0xc7, 0xd9, 0x09, 0x0d, 0xd3, 0xe2, 0xb2, 0x8e, 0xb5, 0x93,
	0x10, 0xfa, 0x09, 0xdf, 0xcb, 0x51, 0x48, 0xc5, 0x88, 0xf9, 0xc3, 0x94, 0x2a, 0x2a, 0xea, 0xdd,
	0x32, 0xac, 0xc1, 0x3a, 0x44, 0x92, 0x49, 0xe2, 0xf7, 0x5e, 0x28, 0xbe, 0x52, 0xec, 0x07, 0x95,
	0xe9, 0xae, 0x64, 0x6d, 0x03, 0x7e, 0x52, 0xb6, 0x76, 0xa9, 0x3c, 0xd4, 0x63, 0xa6, 0xb3, 0xd8,
	0x6d, 0xf8, 0x39, 0x97, 0x49, 0xed, 0xfe, 0x87, 0xa5, 0x74, 0x27, 0xca, 0x69, 0xb9, 0x55, 0x75,
	0xf2, 0x8b, 0x76, 0xb2, 0xaa, 0x8c, 0x6d, 0xdd, 0x14, 0xe1, 0x1b, 0x25, 0x89, 0xa6, 0x00, 0x56,
	0x9e, 0x2f, 0x02, 0xd5, 0xd7, 0x09, 0x6c, 0x5a, 0xb3, 0xd9, 0xd8, 0x92, 0xd6, 0x76, 0xed, 0x1f,
	0x97, 0x77, 0x8f, 0x57, 0xc5, 0x6f, 0xa8, 0x9a, 0x5d, 0x32, 0x39, 0xaa, 0x4c, 0x10, 0x7c, 0x96,
	0xde, 0xe7, 0x1c, 0x5d, 0x00, 0x08, 0x57, 0xa3, 0xa2, 0xdf, 0x1b, 0x5b, 0xe4, 0x36, 0x65, 0xfe,
	0xd9, 0x8a, 0x4d, 0xcd, 0x98, 0xca, 0xcc, 0x47, 0x84, 0xf2, 0x0f, 0x6c, 0x67, 0xef, 0x76, 0x6e,
	0x81, 0xd9, 0xdc, 0x02, 0x0f, 0x73, 0x0b, 0x4c, 0x17, 0x56, 0x61, 0xb6, 0xb0, 0x0a, 0xf7, 0x0b,
	0xab, 0x70, 0xd4, 0x72, 0x3d, 0x39, 0x8a, 0xfa, 0xce, 0x80, 0x1d, 0xe3, 0xac, 0x19, 0x0b, 0xdd,
	0xe5, 0x77, 0x83, 0x70, 0x8e, 0x4f, 0x33, 0x49, 0x39, 0xe1, 0x54, 0xf4, 0xdf, 0xaa, 0xf7, 0xfa,
	0xf7, 0x29, 0x00, 0x00, 0xff, 0xff, 0x3d, 0x57, 0x0e, 0xab, 0x25, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// VersionTally enables a client to query for the tally of voting power has
	// signalled for a particular version.
	VersionTally(ctx context.Context, in *QueryVersionTallyRequest, opts ...grpc.CallOption) (*QueryVersionTallyResponse, error)
	// GetUpgrade enables a client to query for upgrade information if an upgrade is pending.
	// The response will be empty if no upgrade is pending.
	GetUpgrade(ctx context.Context, in *QueryGetUpgradeRequest, opts ...grpc.CallOption) (*QueryGetUpgradeResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) VersionTally(ctx context.Context, in *QueryVersionTallyRequest, opts ...grpc.CallOption) (*QueryVersionTallyResponse, error) {
	out := new(QueryVersionTallyResponse)
	err := c.cc.Invoke(ctx, "/celestia.signal.v1.Query/VersionTally", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetUpgrade(ctx context.Context, in *QueryGetUpgradeRequest, opts ...grpc.CallOption) (*QueryGetUpgradeResponse, error) {
	out := new(QueryGetUpgradeResponse)
	err := c.cc.Invoke(ctx, "/celestia.signal.v1.Query/GetUpgrade", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// VersionTally enables a client to query for the tally of voting power has
	// signalled for a particular version.
	VersionTally(context.Context, *QueryVersionTallyRequest) (*QueryVersionTallyResponse, error)
	// GetUpgrade enables a client to query for upgrade information if an upgrade is pending.
	// The response will be empty if no upgrade is pending.
	GetUpgrade(context.Context, *QueryGetUpgradeRequest) (*QueryGetUpgradeResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) VersionTally(ctx context.Context, req *QueryVersionTallyRequest) (*QueryVersionTallyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VersionTally not implemented")
}
func (*UnimplementedQueryServer) GetUpgrade(ctx context.Context, req *QueryGetUpgradeRequest) (*QueryGetUpgradeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUpgrade not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_VersionTally_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryVersionTallyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).VersionTally(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/celestia.signal.v1.Query/VersionTally",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).VersionTally(ctx, req.(*QueryVersionTallyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetUpgradeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/celestia.signal.v1.Query/GetUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetUpgrade(ctx, req.(*QueryGetUpgradeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "celestia.signal.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VersionTally",
			Handler:    _Query_VersionTally_Handler,
		},
		{
			MethodName: "GetUpgrade",
			Handler:    _Query_GetUpgrade_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "celestia/signal/v1/query.proto",
}

func (m *QueryVersionTallyRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryVersionTallyRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryVersionTallyRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Version != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryVersionTallyResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryVersionTallyResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryVersionTallyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TotalVotingPower != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.TotalVotingPower))
		i--
		dAtA[i] = 0x18
	}
	if m.ThresholdPower != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.ThresholdPower))
		i--
		dAtA[i] = 0x10
	}
	if m.VotingPower != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.VotingPower))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetUpgradeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetUpgradeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetUpgradeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryGetUpgradeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetUpgradeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetUpgradeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Upgrade != nil {
		{
			size, err := m.Upgrade.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryVersionTallyRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Version != 0 {
		n += 1 + sovQuery(uint64(m.Version))
	}
	return n
}

func (m *QueryVersionTallyResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.VotingPower != 0 {
		n += 1 + sovQuery(uint64(m.VotingPower))
	}
	if m.ThresholdPower != 0 {
		n += 1 + sovQuery(uint64(m.ThresholdPower))
	}
	if m.TotalVotingPower != 0 {
		n += 1 + sovQuery(uint64(m.TotalVotingPower))
	}
	return n
}

func (m *QueryGetUpgradeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryGetUpgradeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Upgrade != nil {
		l = m.Upgrade.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryVersionTallyRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryVersionTallyRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryVersionTallyRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryVersionTallyResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryVersionTallyResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryVersionTallyResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VotingPower", wireType)
			}
			m.VotingPower = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VotingPower |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ThresholdPower", wireType)
			}
			m.ThresholdPower = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ThresholdPower |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalVotingPower", wireType)
			}
			m.TotalVotingPower = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalVotingPower |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetUpgradeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetUpgradeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetUpgradeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetUpgradeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetUpgradeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetUpgradeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Upgrade", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Upgrade == nil {
				m.Upgrade = &Upgrade{}
			}
			if err := m.Upgrade.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
