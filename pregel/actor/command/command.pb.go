// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: command.proto

package command

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
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

type InitVertex struct {
	VertexId string `protobuf:"bytes,1,opt,name=vertex_id,json=vertexId,proto3" json:"vertex_id,omitempty"`
}

func (m *InitVertex) Reset()      { *m = InitVertex{} }
func (*InitVertex) ProtoMessage() {}
func (*InitVertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{0}
}
func (m *InitVertex) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InitVertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InitVertex.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InitVertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitVertex.Merge(m, src)
}
func (m *InitVertex) XXX_Size() int {
	return m.Size()
}
func (m *InitVertex) XXX_DiscardUnknown() {
	xxx_messageInfo_InitVertex.DiscardUnknown(m)
}

var xxx_messageInfo_InitVertex proto.InternalMessageInfo

func (m *InitVertex) GetVertexId() string {
	if m != nil {
		return m.VertexId
	}
	return ""
}

type LoadVertex struct {
}

func (m *LoadVertex) Reset()      { *m = LoadVertex{} }
func (*LoadVertex) ProtoMessage() {}
func (*LoadVertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{1}
}
func (m *LoadVertex) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LoadVertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LoadVertex.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LoadVertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadVertex.Merge(m, src)
}
func (m *LoadVertex) XXX_Size() int {
	return m.Size()
}
func (m *LoadVertex) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadVertex.DiscardUnknown(m)
}

var xxx_messageInfo_LoadVertex proto.InternalMessageInfo

type Compute struct {
	SuperStep uint64 `protobuf:"varint,1,opt,name=super_step,json=superStep,proto3" json:"super_step,omitempty"`
}

func (m *Compute) Reset()      { *m = Compute{} }
func (*Compute) ProtoMessage() {}
func (*Compute) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{2}
}
func (m *Compute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Compute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Compute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Compute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Compute.Merge(m, src)
}
func (m *Compute) XXX_Size() int {
	return m.Size()
}
func (m *Compute) XXX_DiscardUnknown() {
	xxx_messageInfo_Compute.DiscardUnknown(m)
}

var xxx_messageInfo_Compute proto.InternalMessageInfo

func (m *Compute) GetSuperStep() uint64 {
	if m != nil {
		return m.SuperStep
	}
	return 0
}

type CustomVertexMessage struct {
	Message *types.Any `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *CustomVertexMessage) Reset()      { *m = CustomVertexMessage{} }
func (*CustomVertexMessage) ProtoMessage() {}
func (*CustomVertexMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{3}
}
func (m *CustomVertexMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CustomVertexMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CustomVertexMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CustomVertexMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CustomVertexMessage.Merge(m, src)
}
func (m *CustomVertexMessage) XXX_Size() int {
	return m.Size()
}
func (m *CustomVertexMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_CustomVertexMessage.DiscardUnknown(m)
}

var xxx_messageInfo_CustomVertexMessage proto.InternalMessageInfo

func (m *CustomVertexMessage) GetMessage() *types.Any {
	if m != nil {
		return m.Message
	}
	return nil
}

func init() {
	proto.RegisterType((*InitVertex)(nil), "command.InitVertex")
	proto.RegisterType((*LoadVertex)(nil), "command.LoadVertex")
	proto.RegisterType((*Compute)(nil), "command.Compute")
	proto.RegisterType((*CustomVertexMessage)(nil), "command.CustomVertexMessage")
}

func init() { proto.RegisterFile("command.proto", fileDescriptor_213c0bb044472049) }

var fileDescriptor_213c0bb044472049 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0xce, 0xcf, 0xcd,
	0x4d, 0xcc, 0x4b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0xa5, 0x24, 0xd3,
	0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0xc2, 0x49, 0xa5, 0x69, 0xfa, 0x89, 0x79, 0x95, 0x10,
	0x35, 0x4a, 0x9a, 0x5c, 0x5c, 0x9e, 0x79, 0x99, 0x25, 0x61, 0xa9, 0x45, 0x25, 0xa9, 0x15, 0x42,
	0xd2, 0x5c, 0x9c, 0x65, 0x60, 0x56, 0x7c, 0x66, 0x8a, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10,
	0x07, 0x44, 0xc0, 0x33, 0x45, 0x89, 0x87, 0x8b, 0xcb, 0x27, 0x3f, 0x31, 0x05, 0xa2, 0x54, 0x49,
	0x83, 0x8b, 0xdd, 0x39, 0x3f, 0xb7, 0xa0, 0xb4, 0x24, 0x55, 0x48, 0x96, 0x8b, 0xab, 0xb8, 0xb4,
	0x20, 0xb5, 0x28, 0xbe, 0xb8, 0x24, 0xb5, 0x00, 0xac, 0x8d, 0x25, 0x88, 0x13, 0x2c, 0x12, 0x5c,
	0x92, 0x5a, 0xa0, 0xe4, 0xca, 0x25, 0xec, 0x5c, 0x5a, 0x5c, 0x92, 0x9f, 0x0b, 0xd1, 0xe9, 0x9b,
	0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0x2a, 0xa4, 0xc7, 0xc5, 0x9e, 0x0b, 0x61, 0x82, 0xb5, 0x70, 0x1b,
	0x89, 0xe8, 0x41, 0x9c, 0xa9, 0x07, 0x73, 0xa6, 0x9e, 0x63, 0x5e, 0x65, 0x10, 0x4c, 0x91, 0x93,
	0xc9, 0x85, 0x87, 0x72, 0x0c, 0x37, 0x1e, 0xca, 0x31, 0x7c, 0x78, 0x28, 0xc7, 0xd8, 0xf0, 0x48,
	0x8e, 0x71, 0xc5, 0x23, 0x39, 0xc6, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0,
	0x48, 0x8e, 0xf1, 0xc5, 0x23, 0x39, 0x86, 0x0f, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8,
	0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x24, 0x36, 0xb0, 0x61, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x5e, 0xb8, 0xce, 0x40, 0x1b, 0x01, 0x00, 0x00,
}

func (this *InitVertex) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*InitVertex)
	if !ok {
		that2, ok := that.(InitVertex)
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
	if this.VertexId != that1.VertexId {
		return false
	}
	return true
}
func (this *LoadVertex) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*LoadVertex)
	if !ok {
		that2, ok := that.(LoadVertex)
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
	return true
}
func (this *Compute) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Compute)
	if !ok {
		that2, ok := that.(Compute)
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
	if this.SuperStep != that1.SuperStep {
		return false
	}
	return true
}
func (this *CustomVertexMessage) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CustomVertexMessage)
	if !ok {
		that2, ok := that.(CustomVertexMessage)
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
	if !this.Message.Equal(that1.Message) {
		return false
	}
	return true
}
func (this *InitVertex) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&command.InitVertex{")
	s = append(s, "VertexId: "+fmt.Sprintf("%#v", this.VertexId)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *LoadVertex) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&command.LoadVertex{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Compute) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&command.Compute{")
	s = append(s, "SuperStep: "+fmt.Sprintf("%#v", this.SuperStep)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *CustomVertexMessage) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&command.CustomVertexMessage{")
	if this.Message != nil {
		s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringCommand(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *InitVertex) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InitVertex) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.VertexId) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCommand(dAtA, i, uint64(len(m.VertexId)))
		i += copy(dAtA[i:], m.VertexId)
	}
	return i, nil
}

func (m *LoadVertex) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LoadVertex) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Compute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Compute) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SuperStep != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.SuperStep))
	}
	return i, nil
}

func (m *CustomVertexMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CustomVertexMessage) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Message != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.Message.Size()))
		n1, err1 := m.Message.MarshalTo(dAtA[i:])
		if err1 != nil {
			return 0, err1
		}
		i += n1
	}
	return i, nil
}

func encodeVarintCommand(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *InitVertex) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.VertexId)
	if l > 0 {
		n += 1 + l + sovCommand(uint64(l))
	}
	return n
}

func (m *LoadVertex) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Compute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SuperStep != 0 {
		n += 1 + sovCommand(uint64(m.SuperStep))
	}
	return n
}

func (m *CustomVertexMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Message != nil {
		l = m.Message.Size()
		n += 1 + l + sovCommand(uint64(l))
	}
	return n
}

func sovCommand(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommand(x uint64) (n int) {
	return sovCommand(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *InitVertex) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&InitVertex{`,
		`VertexId:` + fmt.Sprintf("%v", this.VertexId) + `,`,
		`}`,
	}, "")
	return s
}
func (this *LoadVertex) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&LoadVertex{`,
		`}`,
	}, "")
	return s
}
func (this *Compute) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Compute{`,
		`SuperStep:` + fmt.Sprintf("%v", this.SuperStep) + `,`,
		`}`,
	}, "")
	return s
}
func (this *CustomVertexMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&CustomVertexMessage{`,
		`Message:` + strings.Replace(fmt.Sprintf("%v", this.Message), "Any", "types.Any", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringCommand(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *InitVertex) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
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
			return fmt.Errorf("proto: InitVertex: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InitVertex: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VertexId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
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
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VertexId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthCommand
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
func (m *LoadVertex) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
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
			return fmt.Errorf("proto: LoadVertex: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LoadVertex: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthCommand
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
func (m *Compute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
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
			return fmt.Errorf("proto: Compute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Compute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuperStep", wireType)
			}
			m.SuperStep = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SuperStep |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthCommand
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
func (m *CustomVertexMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
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
			return fmt.Errorf("proto: CustomVertexMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CustomVertexMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
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
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Message == nil {
				m.Message = &types.Any{}
			}
			if err := m.Message.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthCommand
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
func skipCommand(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommand
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
					return 0, ErrIntOverflowCommand
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
					return 0, ErrIntOverflowCommand
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
				return 0, ErrInvalidLengthCommand
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthCommand
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCommand
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
				next, err := skipCommand(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthCommand
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
	ErrInvalidLengthCommand = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommand   = fmt.Errorf("proto: integer overflow")
)
