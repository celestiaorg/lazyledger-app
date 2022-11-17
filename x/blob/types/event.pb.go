// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: blob/event.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type EventPayForBlob struct {
	Signer   string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	BlobSize uint64 `protobuf:"varint,2,opt,name=blob_size,json=blobSize,proto3" json:"blob_size,omitempty"`
}

func (m *EventPayForBlob) Reset()         { *m = EventPayForBlob{} }
func (m *EventPayForBlob) String() string { return proto.CompactTextString(m) }
func (*EventPayForBlob) ProtoMessage()    {}
func (*EventPayForBlob) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c2ee49ac7fcb08f, []int{0}
}
func (m *EventPayForBlob) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventPayForBlob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventPayForBlob.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventPayForBlob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventPayForBlob.Merge(m, src)
}
func (m *EventPayForBlob) XXX_Size() int {
	return m.Size()
}
func (m *EventPayForBlob) XXX_DiscardUnknown() {
	xxx_messageInfo_EventPayForBlob.DiscardUnknown(m)
}

var xxx_messageInfo_EventPayForBlob proto.InternalMessageInfo

func (m *EventPayForBlob) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *EventPayForBlob) GetBlobSize() uint64 {
	if m != nil {
		return m.BlobSize
	}
	return 0
}

func init() {
	proto.RegisterType((*EventPayForBlob)(nil), "blob.EventPayForBlob")
}

func init() { proto.RegisterFile("blob/event.proto", fileDescriptor_6c2ee49ac7fcb08f) }

var fileDescriptor_6c2ee49ac7fcb08f = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0xca, 0xc9, 0x4f,
	0xd2, 0x4f, 0x2d, 0x4b, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0x89,
	0x48, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0x05, 0xf4, 0x41, 0x2c, 0x88, 0x9c, 0x92, 0x1b, 0x17,
	0xbf, 0x2b, 0x48, 0x69, 0x40, 0x62, 0xa5, 0x5b, 0x7e, 0x91, 0x53, 0x4e, 0x7e, 0x92, 0x90, 0x18,
	0x17, 0x5b, 0x71, 0x66, 0x7a, 0x5e, 0x6a, 0x91, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x94,
	0x27, 0x24, 0xcd, 0xc5, 0x09, 0x32, 0x28, 0xbe, 0x38, 0xb3, 0x2a, 0x55, 0x82, 0x49, 0x81, 0x51,
	0x83, 0x25, 0x88, 0x03, 0x24, 0x10, 0x9c, 0x59, 0x95, 0xea, 0xe4, 0x75, 0xe2, 0x91, 0x1c, 0xe3,
	0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c,
	0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x06, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9,
	0xb9, 0xfa, 0xc9, 0xa9, 0x39, 0xa9, 0xc5, 0x25, 0x99, 0x89, 0xf9, 0x45, 0xe9, 0x70, 0xb6, 0x6e,
	0x62, 0x41, 0x81, 0x7e, 0x85, 0x3e, 0xd8, 0xd5, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x49, 0x6c, 0x60,
	0xa7, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x13, 0x3e, 0xc3, 0x7a, 0xca, 0x00, 0x00, 0x00,
}

func (m *EventPayForBlob) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventPayForBlob) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventPayForBlob) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BlobSize != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.BlobSize))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventPayForBlob) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.BlobSize != 0 {
		n += 1 + sovEvent(uint64(m.BlobSize))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventPayForBlob) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventPayForBlob: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventPayForBlob: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlobSize", wireType)
			}
			m.BlobSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlobSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
