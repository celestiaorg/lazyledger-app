// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: celestia/core/v1/tx/tx.proto

package tx

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

// TxStatusRequest is the request type for the TxStatus gRPC method.
type TxStatusRequest struct {
	// this is the hex encoded transaction hash (should be 64 bytes long)
	TxId string `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
}

func (m *TxStatusRequest) Reset()         { *m = TxStatusRequest{} }
func (m *TxStatusRequest) String() string { return proto.CompactTextString(m) }
func (*TxStatusRequest) ProtoMessage()    {}
func (*TxStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d8b070565b0dcb6, []int{0}
}
func (m *TxStatusRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxStatusRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxStatusRequest.Merge(m, src)
}
func (m *TxStatusRequest) XXX_Size() int {
	return m.Size()
}
func (m *TxStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TxStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TxStatusRequest proto.InternalMessageInfo

func (m *TxStatusRequest) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

// TxStatusResponse is the response type for the TxStatus gRPC method.
type TxStatusResponse struct {
	Height int64  `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Index  uint32 `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	// execution_code is returned when the transaction has been committed
	// and returns whether it was successful or errored. A non zero
	// execution code indicated an error.
	ExecutionCode uint32 `protobuf:"varint,3,opt,name=execution_code,json=executionCode,proto3" json:"execution_code,omitempty"`
	// error log for failed transactions.
	Error string `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	// status is the status of the transaction.
	Status string `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (m *TxStatusResponse) Reset()         { *m = TxStatusResponse{} }
func (m *TxStatusResponse) String() string { return proto.CompactTextString(m) }
func (*TxStatusResponse) ProtoMessage()    {}
func (*TxStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d8b070565b0dcb6, []int{1}
}
func (m *TxStatusResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxStatusResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxStatusResponse.Merge(m, src)
}
func (m *TxStatusResponse) XXX_Size() int {
	return m.Size()
}
func (m *TxStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TxStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TxStatusResponse proto.InternalMessageInfo

func (m *TxStatusResponse) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *TxStatusResponse) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *TxStatusResponse) GetExecutionCode() uint32 {
	if m != nil {
		return m.ExecutionCode
	}
	return 0
}

func (m *TxStatusResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *TxStatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*TxStatusRequest)(nil), "celestia.core.v1.tx.TxStatusRequest")
	proto.RegisterType((*TxStatusResponse)(nil), "celestia.core.v1.tx.TxStatusResponse")
}

func init() { proto.RegisterFile("celestia/core/v1/tx/tx.proto", fileDescriptor_7d8b070565b0dcb6) }

var fileDescriptor_7d8b070565b0dcb6 = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0xcd, 0x4a, 0x33, 0x31,
	0x14, 0x6d, 0xfa, 0xc7, 0xf7, 0x05, 0xfa, 0x7d, 0x92, 0x8a, 0x0c, 0xa5, 0x0c, 0xa5, 0xb4, 0xd2,
	0x8d, 0x13, 0xaa, 0x6f, 0xa0, 0xab, 0x6e, 0xc7, 0xae, 0xdc, 0x94, 0xe9, 0xcc, 0x65, 0x1a, 0xa8,
	0x93, 0x31, 0xb9, 0x53, 0x02, 0xd2, 0x8d, 0xbe, 0x80, 0x20, 0xbe, 0x93, 0xcb, 0x82, 0x1b, 0x97,
	0xd2, 0xfa, 0x20, 0x32, 0x99, 0xb6, 0x82, 0x14, 0x5c, 0x04, 0x72, 0xee, 0x39, 0x37, 0x27, 0xf7,
	0x5c, 0xda, 0x0e, 0x61, 0x0e, 0x1a, 0x45, 0xc0, 0x43, 0xa9, 0x80, 0x2f, 0x86, 0x1c, 0x0d, 0x47,
	0xe3, 0xa5, 0x4a, 0xa2, 0x64, 0xcd, 0x1d, 0xeb, 0xe5, 0xac, 0xb7, 0x18, 0x7a, 0x68, 0x5a, 0xed,
	0x58, 0xca, 0x78, 0x0e, 0x3c, 0x48, 0x05, 0x0f, 0x92, 0x44, 0x62, 0x80, 0x42, 0x26, 0xba, 0x68,
	0xe9, 0x9e, 0xd2, 0xff, 0x63, 0x73, 0x8d, 0x01, 0x66, 0xda, 0x87, 0xbb, 0x0c, 0x34, 0xb2, 0x26,
	0xad, 0xa1, 0x99, 0x88, 0xc8, 0x21, 0x1d, 0x32, 0xf8, 0xeb, 0x57, 0xd1, 0x8c, 0xa2, 0xee, 0x0b,
	0xa1, 0x47, 0xdf, 0x42, 0x9d, 0xca, 0x44, 0x03, 0x3b, 0xa1, 0xf5, 0x19, 0x88, 0x78, 0x86, 0x56,
	0x5a, 0xf1, 0xb7, 0x88, 0x1d, 0xd3, 0x9a, 0x48, 0x22, 0x30, 0x4e, 0xb9, 0x43, 0x06, 0x0d, 0xbf,
	0x00, 0xac, 0x4f, 0xff, 0x81, 0x81, 0x30, 0xcb, 0xed, 0x27, 0xa1, 0x8c, 0xc0, 0xa9, 0x58, 0xba,
	0xb1, 0xaf, 0x5e, 0xc9, 0x08, 0xf2, 0x66, 0x50, 0x4a, 0x2a, 0xa7, 0x6a, 0xed, 0x0b, 0x90, 0x5b,
	0x69, 0x6b, 0xee, 0xd4, 0x6c, 0x79, 0x8b, 0xce, 0x1f, 0x09, 0x2d, 0x8f, 0x0d, 0x5b, 0xd2, 0x3f,
	0xbb, 0xdf, 0xb1, 0x9e, 0x77, 0x20, 0x06, 0xef, 0xc7, 0x94, 0xad, 0xfe, 0x2f, 0xaa, 0x62, 0xc4,
	0x6e, 0xef, 0xe1, 0xed, 0xf3, 0xb9, 0xec, 0xb2, 0x36, 0x3f, 0x94, 0xfc, 0xbd, 0x0d, 0x6a, 0x79,
	0x39, 0x7a, 0x5d, 0xbb, 0x64, 0xb5, 0x76, 0xc9, 0xc7, 0xda, 0x25, 0x4f, 0x1b, 0xb7, 0xb4, 0xda,
	0xb8, 0xa5, 0xf7, 0x8d, 0x5b, 0xba, 0xe1, 0xb1, 0xc0, 0x59, 0x36, 0xf5, 0x42, 0x79, 0xbb, 0x7f,
	0x41, 0xaa, 0x78, 0x7f, 0x3f, 0x0b, 0xd2, 0x94, 0xe7, 0x27, 0x56, 0x69, 0xc8, 0xd1, 0x4c, 0xeb,
	0x76, 0x2f, 0x17, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x5a, 0x2f, 0x1b, 0xd1, 0xea, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TxClient is the client API for Tx service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TxClient interface {
	// TxStatus returns the status of a transaction. There are four possible states:
	// - Committed
	// - Pending
	// - Evicted
	// - Unknown
	TxStatus(ctx context.Context, in *TxStatusRequest, opts ...grpc.CallOption) (*TxStatusResponse, error)
}

type txClient struct {
	cc grpc1.ClientConn
}

func NewTxClient(cc grpc1.ClientConn) TxClient {
	return &txClient{cc}
}

func (c *txClient) TxStatus(ctx context.Context, in *TxStatusRequest, opts ...grpc.CallOption) (*TxStatusResponse, error) {
	out := new(TxStatusResponse)
	err := c.cc.Invoke(ctx, "/celestia.core.v1.tx.Tx/TxStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TxServer is the server API for Tx service.
type TxServer interface {
	// TxStatus returns the status of a transaction. There are four possible states:
	// - Committed
	// - Pending
	// - Evicted
	// - Unknown
	TxStatus(context.Context, *TxStatusRequest) (*TxStatusResponse, error)
}

// UnimplementedTxServer can be embedded to have forward compatible implementations.
type UnimplementedTxServer struct {
}

func (*UnimplementedTxServer) TxStatus(ctx context.Context, req *TxStatusRequest) (*TxStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TxStatus not implemented")
}

func RegisterTxServer(s grpc1.Server, srv TxServer) {
	s.RegisterService(&_Tx_serviceDesc, srv)
}

func _Tx_TxStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TxServer).TxStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/celestia.core.v1.tx.Tx/TxStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TxServer).TxStatus(ctx, req.(*TxStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tx_serviceDesc = grpc.ServiceDesc{
	ServiceName: "celestia.core.v1.tx.Tx",
	HandlerType: (*TxServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TxStatus",
			Handler:    _Tx_TxStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "celestia/core/v1/tx/tx.proto",
}

func (m *TxStatusRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxStatusRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxStatusRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxId) > 0 {
		i -= len(m.TxId)
		copy(dAtA[i:], m.TxId)
		i = encodeVarintTx(dAtA, i, uint64(len(m.TxId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TxStatusResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxStatusResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxStatusResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Error) > 0 {
		i -= len(m.Error)
		copy(dAtA[i:], m.Error)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Error)))
		i--
		dAtA[i] = 0x22
	}
	if m.ExecutionCode != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.ExecutionCode))
		i--
		dAtA[i] = 0x18
	}
	if m.Index != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxStatusRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TxId)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *TxStatusResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovTx(uint64(m.Height))
	}
	if m.Index != 0 {
		n += 1 + sovTx(uint64(m.Index))
	}
	if m.ExecutionCode != 0 {
		n += 1 + sovTx(uint64(m.ExecutionCode))
	}
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxStatusRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: TxStatusRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxStatusRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *TxStatusResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: TxStatusResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxStatusResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecutionCode", wireType)
			}
			m.ExecutionCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExecutionCode |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
