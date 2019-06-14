// Code generated by protoc-gen-go. DO NOT EDIT.
// source: command.proto

package command

import (
	fmt "fmt"
	actor "github.com/AsynkronIT/protoactor-go/actor"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LoadVertex struct {
	PartitionId          uint64   `protobuf:"varint,1,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
	VertexId             string   `protobuf:"bytes,2,opt,name=vertex_id,json=vertexId,proto3" json:"vertex_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoadVertex) Reset()         { *m = LoadVertex{} }
func (m *LoadVertex) String() string { return proto.CompactTextString(m) }
func (*LoadVertex) ProtoMessage()    {}
func (*LoadVertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{0}
}

func (m *LoadVertex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadVertex.Unmarshal(m, b)
}
func (m *LoadVertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadVertex.Marshal(b, m, deterministic)
}
func (m *LoadVertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadVertex.Merge(m, src)
}
func (m *LoadVertex) XXX_Size() int {
	return xxx_messageInfo_LoadVertex.Size(m)
}
func (m *LoadVertex) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadVertex.DiscardUnknown(m)
}

var xxx_messageInfo_LoadVertex proto.InternalMessageInfo

func (m *LoadVertex) GetPartitionId() uint64 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

func (m *LoadVertex) GetVertexId() string {
	if m != nil {
		return m.VertexId
	}
	return ""
}

type LoadVertexAck struct {
	PartitionId          uint64   `protobuf:"varint,1,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
	VertexId             string   `protobuf:"bytes,2,opt,name=vertex_id,json=vertexId,proto3" json:"vertex_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoadVertexAck) Reset()         { *m = LoadVertexAck{} }
func (m *LoadVertexAck) String() string { return proto.CompactTextString(m) }
func (*LoadVertexAck) ProtoMessage()    {}
func (*LoadVertexAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{1}
}

func (m *LoadVertexAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadVertexAck.Unmarshal(m, b)
}
func (m *LoadVertexAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadVertexAck.Marshal(b, m, deterministic)
}
func (m *LoadVertexAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadVertexAck.Merge(m, src)
}
func (m *LoadVertexAck) XXX_Size() int {
	return xxx_messageInfo_LoadVertexAck.Size(m)
}
func (m *LoadVertexAck) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadVertexAck.DiscardUnknown(m)
}

var xxx_messageInfo_LoadVertexAck proto.InternalMessageInfo

func (m *LoadVertexAck) GetPartitionId() uint64 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

func (m *LoadVertexAck) GetVertexId() string {
	if m != nil {
		return m.VertexId
	}
	return ""
}

type SuperStepBarrier struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SuperStepBarrier) Reset()         { *m = SuperStepBarrier{} }
func (m *SuperStepBarrier) String() string { return proto.CompactTextString(m) }
func (*SuperStepBarrier) ProtoMessage()    {}
func (*SuperStepBarrier) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{2}
}

func (m *SuperStepBarrier) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperStepBarrier.Unmarshal(m, b)
}
func (m *SuperStepBarrier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperStepBarrier.Marshal(b, m, deterministic)
}
func (m *SuperStepBarrier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperStepBarrier.Merge(m, src)
}
func (m *SuperStepBarrier) XXX_Size() int {
	return xxx_messageInfo_SuperStepBarrier.Size(m)
}
func (m *SuperStepBarrier) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperStepBarrier.DiscardUnknown(m)
}

var xxx_messageInfo_SuperStepBarrier proto.InternalMessageInfo

type SuperStepBarrierAck struct {
	VertexId             string   `protobuf:"bytes,1,opt,name=vertex_id,json=vertexId,proto3" json:"vertex_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SuperStepBarrierAck) Reset()         { *m = SuperStepBarrierAck{} }
func (m *SuperStepBarrierAck) String() string { return proto.CompactTextString(m) }
func (*SuperStepBarrierAck) ProtoMessage()    {}
func (*SuperStepBarrierAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{3}
}

func (m *SuperStepBarrierAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperStepBarrierAck.Unmarshal(m, b)
}
func (m *SuperStepBarrierAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperStepBarrierAck.Marshal(b, m, deterministic)
}
func (m *SuperStepBarrierAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperStepBarrierAck.Merge(m, src)
}
func (m *SuperStepBarrierAck) XXX_Size() int {
	return xxx_messageInfo_SuperStepBarrierAck.Size(m)
}
func (m *SuperStepBarrierAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperStepBarrierAck.DiscardUnknown(m)
}

var xxx_messageInfo_SuperStepBarrierAck proto.InternalMessageInfo

func (m *SuperStepBarrierAck) GetVertexId() string {
	if m != nil {
		return m.VertexId
	}
	return ""
}

type SuperStepBarrierPartitionAck struct {
	PartitionId          uint64   `protobuf:"varint,1,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SuperStepBarrierPartitionAck) Reset()         { *m = SuperStepBarrierPartitionAck{} }
func (m *SuperStepBarrierPartitionAck) String() string { return proto.CompactTextString(m) }
func (*SuperStepBarrierPartitionAck) ProtoMessage()    {}
func (*SuperStepBarrierPartitionAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{4}
}

func (m *SuperStepBarrierPartitionAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperStepBarrierPartitionAck.Unmarshal(m, b)
}
func (m *SuperStepBarrierPartitionAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperStepBarrierPartitionAck.Marshal(b, m, deterministic)
}
func (m *SuperStepBarrierPartitionAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperStepBarrierPartitionAck.Merge(m, src)
}
func (m *SuperStepBarrierPartitionAck) XXX_Size() int {
	return xxx_messageInfo_SuperStepBarrierPartitionAck.Size(m)
}
func (m *SuperStepBarrierPartitionAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperStepBarrierPartitionAck.DiscardUnknown(m)
}

var xxx_messageInfo_SuperStepBarrierPartitionAck proto.InternalMessageInfo

func (m *SuperStepBarrierPartitionAck) GetPartitionId() uint64 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

type SuperStepBarrierWorkerAck struct {
	WorkerPid            *actor.PID `protobuf:"bytes,1,opt,name=worker_pid,json=workerPid,proto3" json:"worker_pid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SuperStepBarrierWorkerAck) Reset()         { *m = SuperStepBarrierWorkerAck{} }
func (m *SuperStepBarrierWorkerAck) String() string { return proto.CompactTextString(m) }
func (*SuperStepBarrierWorkerAck) ProtoMessage()    {}
func (*SuperStepBarrierWorkerAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{5}
}

func (m *SuperStepBarrierWorkerAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperStepBarrierWorkerAck.Unmarshal(m, b)
}
func (m *SuperStepBarrierWorkerAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperStepBarrierWorkerAck.Marshal(b, m, deterministic)
}
func (m *SuperStepBarrierWorkerAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperStepBarrierWorkerAck.Merge(m, src)
}
func (m *SuperStepBarrierWorkerAck) XXX_Size() int {
	return xxx_messageInfo_SuperStepBarrierWorkerAck.Size(m)
}
func (m *SuperStepBarrierWorkerAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperStepBarrierWorkerAck.DiscardUnknown(m)
}

var xxx_messageInfo_SuperStepBarrierWorkerAck proto.InternalMessageInfo

func (m *SuperStepBarrierWorkerAck) GetWorkerPid() *actor.PID {
	if m != nil {
		return m.WorkerPid
	}
	return nil
}

type Compute struct {
	SuperStep            uint64   `protobuf:"varint,1,opt,name=super_step,json=superStep,proto3" json:"super_step,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Compute) Reset()         { *m = Compute{} }
func (m *Compute) String() string { return proto.CompactTextString(m) }
func (*Compute) ProtoMessage()    {}
func (*Compute) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{6}
}

func (m *Compute) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Compute.Unmarshal(m, b)
}
func (m *Compute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Compute.Marshal(b, m, deterministic)
}
func (m *Compute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Compute.Merge(m, src)
}
func (m *Compute) XXX_Size() int {
	return xxx_messageInfo_Compute.Size(m)
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

type ComputeAck struct {
	VertexId             string   `protobuf:"bytes,1,opt,name=vertex_id,json=vertexId,proto3" json:"vertex_id,omitempty"`
	Halted               bool     `protobuf:"varint,2,opt,name=halted,proto3" json:"halted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComputeAck) Reset()         { *m = ComputeAck{} }
func (m *ComputeAck) String() string { return proto.CompactTextString(m) }
func (*ComputeAck) ProtoMessage()    {}
func (*ComputeAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{7}
}

func (m *ComputeAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComputeAck.Unmarshal(m, b)
}
func (m *ComputeAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComputeAck.Marshal(b, m, deterministic)
}
func (m *ComputeAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComputeAck.Merge(m, src)
}
func (m *ComputeAck) XXX_Size() int {
	return xxx_messageInfo_ComputeAck.Size(m)
}
func (m *ComputeAck) XXX_DiscardUnknown() {
	xxx_messageInfo_ComputeAck.DiscardUnknown(m)
}

var xxx_messageInfo_ComputeAck proto.InternalMessageInfo

func (m *ComputeAck) GetVertexId() string {
	if m != nil {
		return m.VertexId
	}
	return ""
}

func (m *ComputeAck) GetHalted() bool {
	if m != nil {
		return m.Halted
	}
	return false
}

type ComputePartitionAck struct {
	PartitionId          uint64   `protobuf:"varint,1,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComputePartitionAck) Reset()         { *m = ComputePartitionAck{} }
func (m *ComputePartitionAck) String() string { return proto.CompactTextString(m) }
func (*ComputePartitionAck) ProtoMessage()    {}
func (*ComputePartitionAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{8}
}

func (m *ComputePartitionAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComputePartitionAck.Unmarshal(m, b)
}
func (m *ComputePartitionAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComputePartitionAck.Marshal(b, m, deterministic)
}
func (m *ComputePartitionAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComputePartitionAck.Merge(m, src)
}
func (m *ComputePartitionAck) XXX_Size() int {
	return xxx_messageInfo_ComputePartitionAck.Size(m)
}
func (m *ComputePartitionAck) XXX_DiscardUnknown() {
	xxx_messageInfo_ComputePartitionAck.DiscardUnknown(m)
}

var xxx_messageInfo_ComputePartitionAck proto.InternalMessageInfo

func (m *ComputePartitionAck) GetPartitionId() uint64 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

type ComputeWorkerAck struct {
	WorkerPid            *actor.PID `protobuf:"bytes,1,opt,name=worker_pid,json=workerPid,proto3" json:"worker_pid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ComputeWorkerAck) Reset()         { *m = ComputeWorkerAck{} }
func (m *ComputeWorkerAck) String() string { return proto.CompactTextString(m) }
func (*ComputeWorkerAck) ProtoMessage()    {}
func (*ComputeWorkerAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{9}
}

func (m *ComputeWorkerAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComputeWorkerAck.Unmarshal(m, b)
}
func (m *ComputeWorkerAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComputeWorkerAck.Marshal(b, m, deterministic)
}
func (m *ComputeWorkerAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComputeWorkerAck.Merge(m, src)
}
func (m *ComputeWorkerAck) XXX_Size() int {
	return xxx_messageInfo_ComputeWorkerAck.Size(m)
}
func (m *ComputeWorkerAck) XXX_DiscardUnknown() {
	xxx_messageInfo_ComputeWorkerAck.DiscardUnknown(m)
}

var xxx_messageInfo_ComputeWorkerAck proto.InternalMessageInfo

func (m *ComputeWorkerAck) GetWorkerPid() *actor.PID {
	if m != nil {
		return m.WorkerPid
	}
	return nil
}

type SuperStepMessage struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	SuperStep            uint64   `protobuf:"varint,2,opt,name=super_step,json=superStep,proto3" json:"super_step,omitempty"`
	SrcVertexId          string   `protobuf:"bytes,3,opt,name=src_vertex_id,json=srcVertexId,proto3" json:"src_vertex_id,omitempty"`
	DestVertexId         string   `protobuf:"bytes,4,opt,name=dest_vertex_id,json=destVertexId,proto3" json:"dest_vertex_id,omitempty"`
	Message              *any.Any `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SuperStepMessage) Reset()         { *m = SuperStepMessage{} }
func (m *SuperStepMessage) String() string { return proto.CompactTextString(m) }
func (*SuperStepMessage) ProtoMessage()    {}
func (*SuperStepMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{10}
}

func (m *SuperStepMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperStepMessage.Unmarshal(m, b)
}
func (m *SuperStepMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperStepMessage.Marshal(b, m, deterministic)
}
func (m *SuperStepMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperStepMessage.Merge(m, src)
}
func (m *SuperStepMessage) XXX_Size() int {
	return xxx_messageInfo_SuperStepMessage.Size(m)
}
func (m *SuperStepMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperStepMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SuperStepMessage proto.InternalMessageInfo

func (m *SuperStepMessage) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *SuperStepMessage) GetSuperStep() uint64 {
	if m != nil {
		return m.SuperStep
	}
	return 0
}

func (m *SuperStepMessage) GetSrcVertexId() string {
	if m != nil {
		return m.SrcVertexId
	}
	return ""
}

func (m *SuperStepMessage) GetDestVertexId() string {
	if m != nil {
		return m.DestVertexId
	}
	return ""
}

func (m *SuperStepMessage) GetMessage() *any.Any {
	if m != nil {
		return m.Message
	}
	return nil
}

// SuperStepMessageAck is used to guarantee receipt messages at each end of super step
type SuperStepMessageAck struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SuperStepMessageAck) Reset()         { *m = SuperStepMessageAck{} }
func (m *SuperStepMessageAck) String() string { return proto.CompactTextString(m) }
func (*SuperStepMessageAck) ProtoMessage()    {}
func (*SuperStepMessageAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{11}
}

func (m *SuperStepMessageAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperStepMessageAck.Unmarshal(m, b)
}
func (m *SuperStepMessageAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperStepMessageAck.Marshal(b, m, deterministic)
}
func (m *SuperStepMessageAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperStepMessageAck.Merge(m, src)
}
func (m *SuperStepMessageAck) XXX_Size() int {
	return xxx_messageInfo_SuperStepMessageAck.Size(m)
}
func (m *SuperStepMessageAck) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperStepMessageAck.DiscardUnknown(m)
}

var xxx_messageInfo_SuperStepMessageAck proto.InternalMessageInfo

func (m *SuperStepMessageAck) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type InitPartition struct {
	PartitionId          uint64   `protobuf:"varint,1,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitPartition) Reset()         { *m = InitPartition{} }
func (m *InitPartition) String() string { return proto.CompactTextString(m) }
func (*InitPartition) ProtoMessage()    {}
func (*InitPartition) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{12}
}

func (m *InitPartition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitPartition.Unmarshal(m, b)
}
func (m *InitPartition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitPartition.Marshal(b, m, deterministic)
}
func (m *InitPartition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitPartition.Merge(m, src)
}
func (m *InitPartition) XXX_Size() int {
	return xxx_messageInfo_InitPartition.Size(m)
}
func (m *InitPartition) XXX_DiscardUnknown() {
	xxx_messageInfo_InitPartition.DiscardUnknown(m)
}

var xxx_messageInfo_InitPartition proto.InternalMessageInfo

func (m *InitPartition) GetPartitionId() uint64 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

type InitPartitionAck struct {
	PartitionId          uint64   `protobuf:"varint,1,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitPartitionAck) Reset()         { *m = InitPartitionAck{} }
func (m *InitPartitionAck) String() string { return proto.CompactTextString(m) }
func (*InitPartitionAck) ProtoMessage()    {}
func (*InitPartitionAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{13}
}

func (m *InitPartitionAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitPartitionAck.Unmarshal(m, b)
}
func (m *InitPartitionAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitPartitionAck.Marshal(b, m, deterministic)
}
func (m *InitPartitionAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitPartitionAck.Merge(m, src)
}
func (m *InitPartitionAck) XXX_Size() int {
	return xxx_messageInfo_InitPartitionAck.Size(m)
}
func (m *InitPartitionAck) XXX_DiscardUnknown() {
	xxx_messageInfo_InitPartitionAck.DiscardUnknown(m)
}

var xxx_messageInfo_InitPartitionAck proto.InternalMessageInfo

func (m *InitPartitionAck) GetPartitionId() uint64 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

type ClusterInfo struct {
	WorkerInfo           []*WorkerInfo `protobuf:"bytes,1,rep,name=worker_info,json=workerInfo,proto3" json:"worker_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ClusterInfo) Reset()         { *m = ClusterInfo{} }
func (m *ClusterInfo) String() string { return proto.CompactTextString(m) }
func (*ClusterInfo) ProtoMessage()    {}
func (*ClusterInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{14}
}

func (m *ClusterInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClusterInfo.Unmarshal(m, b)
}
func (m *ClusterInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClusterInfo.Marshal(b, m, deterministic)
}
func (m *ClusterInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterInfo.Merge(m, src)
}
func (m *ClusterInfo) XXX_Size() int {
	return xxx_messageInfo_ClusterInfo.Size(m)
}
func (m *ClusterInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterInfo proto.InternalMessageInfo

func (m *ClusterInfo) GetWorkerInfo() []*WorkerInfo {
	if m != nil {
		return m.WorkerInfo
	}
	return nil
}

type WorkerInfo struct {
	WorkerPid            *actor.PID `protobuf:"bytes,1,opt,name=worker_pid,json=workerPid,proto3" json:"worker_pid,omitempty"`
	Partitions           []uint64   `protobuf:"varint,2,rep,packed,name=partitions,proto3" json:"partitions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *WorkerInfo) Reset()         { *m = WorkerInfo{} }
func (m *WorkerInfo) String() string { return proto.CompactTextString(m) }
func (*WorkerInfo) ProtoMessage()    {}
func (*WorkerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{15}
}

func (m *WorkerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkerInfo.Unmarshal(m, b)
}
func (m *WorkerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkerInfo.Marshal(b, m, deterministic)
}
func (m *WorkerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkerInfo.Merge(m, src)
}
func (m *WorkerInfo) XXX_Size() int {
	return xxx_messageInfo_WorkerInfo.Size(m)
}
func (m *WorkerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_WorkerInfo proto.InternalMessageInfo

func (m *WorkerInfo) GetWorkerPid() *actor.PID {
	if m != nil {
		return m.WorkerPid
	}
	return nil
}

func (m *WorkerInfo) GetPartitions() []uint64 {
	if m != nil {
		return m.Partitions
	}
	return nil
}

type InitWorker struct {
	ClusterInfo          *ClusterInfo `protobuf:"bytes,1,opt,name=clusterInfo,proto3" json:"clusterInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *InitWorker) Reset()         { *m = InitWorker{} }
func (m *InitWorker) String() string { return proto.CompactTextString(m) }
func (*InitWorker) ProtoMessage()    {}
func (*InitWorker) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{16}
}

func (m *InitWorker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitWorker.Unmarshal(m, b)
}
func (m *InitWorker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitWorker.Marshal(b, m, deterministic)
}
func (m *InitWorker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitWorker.Merge(m, src)
}
func (m *InitWorker) XXX_Size() int {
	return xxx_messageInfo_InitWorker.Size(m)
}
func (m *InitWorker) XXX_DiscardUnknown() {
	xxx_messageInfo_InitWorker.DiscardUnknown(m)
}

var xxx_messageInfo_InitWorker proto.InternalMessageInfo

func (m *InitWorker) GetClusterInfo() *ClusterInfo {
	if m != nil {
		return m.ClusterInfo
	}
	return nil
}

type InitWorkerAck struct {
	WorkerPid            *actor.PID `protobuf:"bytes,1,opt,name=worker_pid,json=workerPid,proto3" json:"worker_pid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *InitWorkerAck) Reset()         { *m = InitWorkerAck{} }
func (m *InitWorkerAck) String() string { return proto.CompactTextString(m) }
func (*InitWorkerAck) ProtoMessage()    {}
func (*InitWorkerAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_213c0bb044472049, []int{17}
}

func (m *InitWorkerAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitWorkerAck.Unmarshal(m, b)
}
func (m *InitWorkerAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitWorkerAck.Marshal(b, m, deterministic)
}
func (m *InitWorkerAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitWorkerAck.Merge(m, src)
}
func (m *InitWorkerAck) XXX_Size() int {
	return xxx_messageInfo_InitWorkerAck.Size(m)
}
func (m *InitWorkerAck) XXX_DiscardUnknown() {
	xxx_messageInfo_InitWorkerAck.DiscardUnknown(m)
}

var xxx_messageInfo_InitWorkerAck proto.InternalMessageInfo

func (m *InitWorkerAck) GetWorkerPid() *actor.PID {
	if m != nil {
		return m.WorkerPid
	}
	return nil
}

func init() {
	proto.RegisterType((*LoadVertex)(nil), "LoadVertex")
	proto.RegisterType((*LoadVertexAck)(nil), "LoadVertexAck")
	proto.RegisterType((*SuperStepBarrier)(nil), "SuperStepBarrier")
	proto.RegisterType((*SuperStepBarrierAck)(nil), "SuperStepBarrierAck")
	proto.RegisterType((*SuperStepBarrierPartitionAck)(nil), "SuperStepBarrierPartitionAck")
	proto.RegisterType((*SuperStepBarrierWorkerAck)(nil), "SuperStepBarrierWorkerAck")
	proto.RegisterType((*Compute)(nil), "Compute")
	proto.RegisterType((*ComputeAck)(nil), "ComputeAck")
	proto.RegisterType((*ComputePartitionAck)(nil), "ComputePartitionAck")
	proto.RegisterType((*ComputeWorkerAck)(nil), "ComputeWorkerAck")
	proto.RegisterType((*SuperStepMessage)(nil), "SuperStepMessage")
	proto.RegisterType((*SuperStepMessageAck)(nil), "SuperStepMessageAck")
	proto.RegisterType((*InitPartition)(nil), "InitPartition")
	proto.RegisterType((*InitPartitionAck)(nil), "InitPartitionAck")
	proto.RegisterType((*ClusterInfo)(nil), "ClusterInfo")
	proto.RegisterType((*WorkerInfo)(nil), "WorkerInfo")
	proto.RegisterType((*InitWorker)(nil), "InitWorker")
	proto.RegisterType((*InitWorkerAck)(nil), "InitWorkerAck")
}

func init() { proto.RegisterFile("command.proto", fileDescriptor_213c0bb044472049) }

var fileDescriptor_213c0bb044472049 = []byte{
	// 513 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x5f, 0x6b, 0xdb, 0x3e,
	0x14, 0xc5, 0x49, 0x7e, 0x6d, 0x73, 0x9d, 0xfc, 0x08, 0xea, 0x18, 0xe9, 0xfe, 0x91, 0x89, 0x3d,
	0xb8, 0xb0, 0x29, 0x90, 0xb1, 0x31, 0xf6, 0xe7, 0xc1, 0xeb, 0x18, 0x18, 0x3a, 0x16, 0xdc, 0xd1,
	0x3e, 0x06, 0xc7, 0x56, 0x5c, 0x93, 0xd8, 0x32, 0x92, 0xbc, 0x2d, 0x9f, 0x6e, 0x5f, 0x6d, 0x58,
	0x92, 0x1d, 0xd7, 0x0c, 0x96, 0x76, 0x6f, 0xd2, 0xd5, 0xb9, 0xe7, 0x9e, 0x73, 0xae, 0x0d, 0xc3,
	0x90, 0xa5, 0x69, 0x90, 0x45, 0x24, 0xe7, 0x4c, 0xb2, 0x07, 0x27, 0x31, 0x63, 0xf1, 0x86, 0x4e,
	0xd5, 0x6d, 0x59, 0xac, 0xa6, 0x41, 0xb6, 0x35, 0x4f, 0xaf, 0xe3, 0x44, 0x5e, 0x17, 0x4b, 0x12,
	0xb2, 0x74, 0xea, 0x8a, 0x6d, 0xb6, 0xe6, 0x2c, 0xf3, 0xbe, 0x69, 0x64, 0x10, 0x4a, 0xc6, 0x5f,
	0xc4, 0x6c, 0xaa, 0x0e, 0xba, 0x26, 0x74, 0x1f, 0x3e, 0x07, 0x38, 0x67, 0x41, 0x74, 0x49, 0xb9,
	0xa4, 0x3f, 0xd1, 0x53, 0x18, 0xe4, 0x01, 0x97, 0x89, 0x4c, 0x58, 0xb6, 0x48, 0xa2, 0xb1, 0x35,
	0xb1, 0x9c, 0x9e, 0x6f, 0xd7, 0x35, 0x2f, 0x42, 0x0f, 0xa1, 0xff, 0x5d, 0x81, 0xcb, 0xf7, 0xce,
	0xc4, 0x72, 0xfa, 0xfe, 0x91, 0x2e, 0x78, 0x11, 0xfe, 0x0a, 0xc3, 0x1d, 0x9b, 0x1b, 0xae, 0xff,
	0x99, 0x10, 0xc1, 0xe8, 0xa2, 0xc8, 0x29, 0xbf, 0x90, 0x34, 0xff, 0x18, 0x70, 0x9e, 0x50, 0x8e,
	0x67, 0x70, 0xdc, 0xae, 0x95, 0xa3, 0x6e, 0xf0, 0x58, 0x2d, 0x1e, 0x17, 0x1e, 0xb5, 0x7b, 0xe6,
	0x95, 0x86, 0xfd, 0x74, 0xe2, 0xcf, 0x70, 0xd2, 0xa6, 0xb8, 0x62, 0x7c, 0xad, 0x87, 0x9f, 0x02,
	0xfc, 0x50, 0x97, 0x45, 0x6e, 0xba, 0xed, 0x19, 0x10, 0x95, 0x37, 0x99, 0x7b, 0x9f, 0xfc, 0xbe,
	0x7e, 0x9d, 0x27, 0x11, 0x76, 0xe0, 0xf0, 0x8c, 0xa5, 0x79, 0x21, 0x29, 0x7a, 0x0c, 0x20, 0x4a,
	0xca, 0x85, 0x90, 0x34, 0x37, 0x33, 0xfb, 0xa2, 0x1a, 0x82, 0x5d, 0x00, 0x83, 0xfc, 0x9b, 0x3f,
	0x74, 0x1f, 0x0e, 0xae, 0x83, 0x8d, 0xa4, 0x3a, 0xc1, 0x23, 0xdf, 0xdc, 0xf0, 0x1b, 0x38, 0x36,
	0x14, 0xb7, 0xb5, 0xfb, 0x01, 0x46, 0xa6, 0xf3, 0x4e, 0x2e, 0x7f, 0x59, 0x8d, 0xcd, 0x7d, 0xa1,
	0x42, 0x04, 0x31, 0x45, 0x08, 0x7a, 0x45, 0x51, 0xab, 0x57, 0xe7, 0x56, 0x06, 0x9d, 0x56, 0x06,
	0x08, 0xc3, 0x50, 0xf0, 0x70, 0xb1, 0x73, 0xde, 0x55, 0xbd, 0xb6, 0xe0, 0xe1, 0x65, 0x65, 0xfe,
	0x19, 0xfc, 0x1f, 0x51, 0x21, 0x1b, 0xa0, 0x9e, 0x02, 0x0d, 0xca, 0x6a, 0x8d, 0x22, 0x70, 0x98,
	0x6a, 0x1d, 0xe3, 0xff, 0x94, 0xf2, 0x7b, 0x44, 0xff, 0x4e, 0xa4, 0xfa, 0x9d, 0x88, 0x9b, 0x6d,
	0xfd, 0x0a, 0x84, 0x4f, 0x1b, 0x9f, 0x99, 0x31, 0x50, 0x66, 0xf0, 0x07, 0x0f, 0x78, 0x06, 0x43,
	0x2f, 0x4b, 0x64, 0x1d, 0xf1, 0x3e, 0xf9, 0xbe, 0x82, 0xd1, 0x8d, 0x9e, 0x3d, 0xd7, 0xf2, 0x0e,
	0xec, 0xb3, 0x4d, 0x21, 0x24, 0xe5, 0x5e, 0xb6, 0x62, 0xe8, 0x39, 0xd8, 0x66, 0x23, 0x49, 0xb6,
	0x62, 0x63, 0x6b, 0xd2, 0x75, 0xec, 0x99, 0x4d, 0xf4, 0xca, 0x4a, 0x84, 0x6f, 0x36, 0x56, 0x9e,
	0xf1, 0x15, 0xc0, 0xee, 0xe5, 0x16, 0xdb, 0x44, 0x4f, 0x00, 0x6a, 0x11, 0x62, 0xdc, 0x99, 0x74,
	0x9d, 0x9e, 0xdf, 0xa8, 0xe0, 0xf7, 0x00, 0xa5, 0x19, 0x4d, 0x8e, 0x08, 0xd8, 0xe1, 0x4e, 0xa3,
	0x61, 0x1e, 0x90, 0x86, 0x6e, 0xbf, 0x09, 0xc0, 0x6f, 0x75, 0x7c, 0x77, 0xf9, 0xce, 0x96, 0x07,
	0x6a, 0x79, 0x2f, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xad, 0xdf, 0xd1, 0x9a, 0x2a, 0x05, 0x00,
	0x00,
}
