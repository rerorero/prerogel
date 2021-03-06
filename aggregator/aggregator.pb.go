// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: aggregator.proto

package aggregator

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type VertexStats struct {
	ActiveVertices uint64 `protobuf:"varint,1,opt,name=active_vertices,json=activeVertices,proto3" json:"active_vertices,omitempty"`
	TotalVertices  uint64 `protobuf:"varint,2,opt,name=total_vertices,json=totalVertices,proto3" json:"total_vertices,omitempty"`
	MessagesSent   uint64 `protobuf:"varint,3,opt,name=messages_sent,json=messagesSent,proto3" json:"messages_sent,omitempty"`
}

func (m *VertexStats) Reset()      { *m = VertexStats{} }
func (*VertexStats) ProtoMessage() {}
func (*VertexStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{0}
}
func (m *VertexStats) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VertexStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VertexStats.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VertexStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VertexStats.Merge(m, src)
}
func (m *VertexStats) XXX_Size() int {
	return m.Size()
}
func (m *VertexStats) XXX_DiscardUnknown() {
	xxx_messageInfo_VertexStats.DiscardUnknown(m)
}

var xxx_messageInfo_VertexStats proto.InternalMessageInfo

func (m *VertexStats) GetActiveVertices() uint64 {
	if m != nil {
		return m.ActiveVertices
	}
	return 0
}

func (m *VertexStats) GetTotalVertices() uint64 {
	if m != nil {
		return m.TotalVertices
	}
	return 0
}

func (m *VertexStats) GetMessagesSent() uint64 {
	if m != nil {
		return m.MessagesSent
	}
	return 0
}

func init() {
	proto.RegisterType((*VertexStats)(nil), "VertexStats")
}

func init() { proto.RegisterFile("aggregator.proto", fileDescriptor_60785b04c84bec7e) }

var fileDescriptor_60785b04c84bec7e = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x4c, 0x4f, 0x2f,
	0x4a, 0x4d, 0x4f, 0x2c, 0xc9, 0x2f, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x6a, 0x62, 0xe4,
	0xe2, 0x0e, 0x4b, 0x2d, 0x2a, 0x49, 0xad, 0x08, 0x2e, 0x49, 0x2c, 0x29, 0x16, 0x52, 0xe7, 0xe2,
	0x4f, 0x4c, 0x2e, 0xc9, 0x2c, 0x4b, 0x8d, 0x2f, 0x4b, 0x2d, 0x2a, 0xc9, 0x4c, 0x4e, 0x2d, 0x96,
	0x60, 0x54, 0x60, 0xd4, 0x60, 0x09, 0xe2, 0x83, 0x08, 0x87, 0x41, 0x45, 0x85, 0x54, 0xb9, 0xf8,
	0x4a, 0xf2, 0x4b, 0x12, 0x73, 0x10, 0xea, 0x98, 0xc0, 0xea, 0x78, 0xc1, 0xa2, 0x70, 0x65, 0xca,
	0x5c, 0xbc, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0xc5, 0xf1, 0xc5, 0xa9, 0x79, 0x25, 0x12,
	0xcc, 0x60, 0x55, 0x3c, 0x30, 0xc1, 0xe0, 0xd4, 0xbc, 0x12, 0x27, 0x93, 0x0b, 0x0f, 0xe5, 0x18,
	0x6e, 0x3c, 0x94, 0x63, 0xf8, 0xf0, 0x50, 0x8e, 0xb1, 0xe1, 0x91, 0x1c, 0xe3, 0x8a, 0x47, 0x72,
	0x8c, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x8b, 0x47,
	0x72, 0x0c, 0x1f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d,
	0xc7, 0x72, 0x0c, 0x49, 0x6c, 0x60, 0x1f, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x09, 0x3a,
	0x13, 0x87, 0xd5, 0x00, 0x00, 0x00,
}

func (this *VertexStats) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*VertexStats)
	if !ok {
		that2, ok := that.(VertexStats)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ActiveVertices != that1.ActiveVertices {
		return false
	}
	if this.TotalVertices != that1.TotalVertices {
		return false
	}
	if this.MessagesSent != that1.MessagesSent {
		return false
	}
	return true
}
func (this *VertexStats) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&aggregator.VertexStats{")
	s = append(s, "ActiveVertices: "+fmt.Sprintf("%#v", this.ActiveVertices)+",\n")
	s = append(s, "TotalVertices: "+fmt.Sprintf("%#v", this.TotalVertices)+",\n")
	s = append(s, "MessagesSent: "+fmt.Sprintf("%#v", this.MessagesSent)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringAggregator(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *VertexStats) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VertexStats) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ActiveVertices != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintAggregator(dAtA, i, uint64(m.ActiveVertices))
	}
	if m.TotalVertices != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintAggregator(dAtA, i, uint64(m.TotalVertices))
	}
	if m.MessagesSent != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintAggregator(dAtA, i, uint64(m.MessagesSent))
	}
	return i, nil
}

func encodeVarintAggregator(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *VertexStats) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActiveVertices != 0 {
		n += 1 + sovAggregator(uint64(m.ActiveVertices))
	}
	if m.TotalVertices != 0 {
		n += 1 + sovAggregator(uint64(m.TotalVertices))
	}
	if m.MessagesSent != 0 {
		n += 1 + sovAggregator(uint64(m.MessagesSent))
	}
	return n
}

func sovAggregator(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozAggregator(x uint64) (n int) {
	return sovAggregator(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *VertexStats) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&VertexStats{`,
		`ActiveVertices:` + fmt.Sprintf("%v", this.ActiveVertices) + `,`,
		`TotalVertices:` + fmt.Sprintf("%v", this.TotalVertices) + `,`,
		`MessagesSent:` + fmt.Sprintf("%v", this.MessagesSent) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringAggregator(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *VertexStats) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAggregator
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
			return fmt.Errorf("proto: VertexStats: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VertexStats: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActiveVertices", wireType)
			}
			m.ActiveVertices = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAggregator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ActiveVertices |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalVertices", wireType)
			}
			m.TotalVertices = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAggregator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalVertices |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessagesSent", wireType)
			}
			m.MessagesSent = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAggregator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MessagesSent |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAggregator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAggregator
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAggregator
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
func skipAggregator(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAggregator
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
					return 0, ErrIntOverflowAggregator
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAggregator
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
				return 0, ErrInvalidLengthAggregator
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthAggregator
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAggregator
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipAggregator(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthAggregator
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthAggregator = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAggregator   = fmt.Errorf("proto: integer overflow")
)
