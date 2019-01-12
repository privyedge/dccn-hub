// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/taskmgr/taskmgr.proto

package go_micro_srv_taskmgr

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

type Task_Status int32

const (
	Task_CREATING  Task_Status = 0
	Task_CREATED   Task_Status = 1
	Task_RUNNING   Task_Status = 2
	Task_CANCELING Task_Status = 3
	Task_CANCELED  Task_Status = 4
	Task_UPDATING  Task_Status = 5
	Task_UPDATED   Task_Status = 6
	Task_FAILURE   Task_Status = 7
	Task_SUCCESS   Task_Status = 8
)

var Task_Status_name = map[int32]string{
	0: "CREATING",
	1: "CREATED",
	2: "RUNNING",
	3: "CANCELING",
	4: "CANCELED",
	5: "UPDATING",
	6: "UPDATED",
	7: "FAILURE",
	8: "SUCCESS",
}
var Task_Status_value = map[string]int32{
	"CREATING":  0,
	"CREATED":   1,
	"RUNNING":   2,
	"CANCELING": 3,
	"CANCELED":  4,
	"UPDATING":  5,
	"UPDATED":   6,
	"FAILURE":   7,
	"SUCCESS":   8,
}

func (x Task_Status) String() string {
	return proto.EnumName(Task_Status_name, int32(x))
}
func (Task_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_taskmgr_f2b1a0b0af520394, []int{1, 0}
}

type Event_OpCode int32

const (
	Event_CREATE  Event_OpCode = 0
	Event_CANCEL  Event_OpCode = 1
	Event_UPDATE  Event_OpCode = 2
	Event_RETURN  Event_OpCode = 3
	Event_REFUSED Event_OpCode = 4
)

var Event_OpCode_name = map[int32]string{
	0: "CREATE",
	1: "CANCEL",
	2: "UPDATE",
	3: "RETURN",
	4: "REFUSED",
}
var Event_OpCode_value = map[string]int32{
	"CREATE":  0,
	"CANCEL":  1,
	"UPDATE":  2,
	"RETURN":  3,
	"REFUSED": 4,
}

func (x Event_OpCode) String() string {
	return proto.EnumName(Event_OpCode_name, int32(x))
}
func (Event_OpCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_taskmgr_f2b1a0b0af520394, []int{3, 0}
}

type ID struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ID) Reset()         { *m = ID{} }
func (m *ID) String() string { return proto.CompactTextString(m) }
func (*ID) ProtoMessage()    {}
func (*ID) Descriptor() ([]byte, []int) {
	return fileDescriptor_taskmgr_f2b1a0b0af520394, []int{0}
}
func (m *ID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ID.Unmarshal(m, b)
}
func (m *ID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ID.Marshal(b, m, deterministic)
}
func (dst *ID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ID.Merge(dst, src)
}
func (m *ID) XXX_Size() int {
	return xxx_messageInfo_ID.Size(m)
}
func (m *ID) XXX_DiscardUnknown() {
	xxx_messageInfo_ID.DiscardUnknown(m)
}

var xxx_messageInfo_ID proto.InternalMessageInfo

func (m *ID) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Task struct {
	// id task id, unique.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// event_id the flag of updating task info.
	// filter log with event id
	EventId string `protobuf:"bytes,2,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	// user_id the task belongs.
	UserId int64 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// name task name
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// task type
	Type string `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	// status [CREATING,RUNNING,CANCELING,CANCELED,UPDATING,UPDATED,FAILURE,SUCCESS]
	Status Task_Status `protobuf:"varint,6,opt,name=status,proto3,enum=go.micro.srv.taskmgr.Task_Status" json:"status,omitempty"`
	// startup time of the this task
	StartupTime uint32 `protobuf:"varint,7,opt,name=startup_time,json=startupTime,proto3" json:"startup_time,omitempty"`
	// mission data center id
	DataCenterId int64 `protobuf:"varint,8,opt,name=data_center_id,json=dataCenterId,proto3" json:"data_center_id,omitempty"`
	// task creation date
	CreateTime uint64 `protobuf:"varint,9,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// task cancel date
	CancelTime uint64 `protobuf:"varint,10,opt,name=cancel_time,json=cancelTime,proto3" json:"cancel_time,omitempty"`
	// task update date
	UpdateTime uint64 `protobuf:"varint,11,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	// task returns date
	ReturnTime uint64 `protobuf:"varint,12,opt,name=return_time,json=returnTime,proto3" json:"return_time,omitempty"`
	// extra for other arguments
	Extra []byte `protobuf:"bytes,13,opt,name=extra,proto3" json:"extra,omitempty"`
	// result of the task
	Result               []byte   `protobuf:"bytes,14,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_taskmgr_f2b1a0b0af520394, []int{1}
}
func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (dst *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(dst, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Task) GetEventId() string {
	if m != nil {
		return m.EventId
	}
	return ""
}

func (m *Task) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *Task) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Task) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Task) GetStatus() Task_Status {
	if m != nil {
		return m.Status
	}
	return Task_CREATING
}

func (m *Task) GetStartupTime() uint32 {
	if m != nil {
		return m.StartupTime
	}
	return 0
}

func (m *Task) GetDataCenterId() int64 {
	if m != nil {
		return m.DataCenterId
	}
	return 0
}

func (m *Task) GetCreateTime() uint64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *Task) GetCancelTime() uint64 {
	if m != nil {
		return m.CancelTime
	}
	return 0
}

func (m *Task) GetUpdateTime() uint64 {
	if m != nil {
		return m.UpdateTime
	}
	return 0
}

func (m *Task) GetReturnTime() uint64 {
	if m != nil {
		return m.ReturnTime
	}
	return 0
}

func (m *Task) GetExtra() []byte {
	if m != nil {
		return m.Extra
	}
	return nil
}

func (m *Task) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

type Response struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_taskmgr_f2b1a0b0af520394, []int{2}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

type Event struct {
	// unique id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// event based on task_id.
	TaskId int64 `protobuf:"varint,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	// unix timestamp
	Timestamp int64 `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// operation code []
	Operation Event_OpCode `protobuf:"varint,4,opt,name=operation,proto3,enum=go.micro.srv.taskmgr.Event_OpCode" json:"operation,omitempty"`
	// data_center_id task's executer dc
	DataCenterId int64 `protobuf:"varint,5,opt,name=data_center_id,json=dataCenterId,proto3" json:"data_center_id,omitempty"`
	// modify_time status's change time
	ModifyTime int64 `protobuf:"varint,6,opt,name=modify_time,json=modifyTime,proto3" json:"modify_time,omitempty"`
	// result task result
	Result               []byte   `protobuf:"bytes,7,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_taskmgr_f2b1a0b0af520394, []int{3}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetTaskId() int64 {
	if m != nil {
		return m.TaskId
	}
	return 0
}

func (m *Event) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Event) GetOperation() Event_OpCode {
	if m != nil {
		return m.Operation
	}
	return Event_CREATE
}

func (m *Event) GetDataCenterId() int64 {
	if m != nil {
		return m.DataCenterId
	}
	return 0
}

func (m *Event) GetModifyTime() int64 {
	if m != nil {
		return m.ModifyTime
	}
	return 0
}

func (m *Event) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*ID)(nil), "go.micro.srv.taskmgr.ID")
	proto.RegisterType((*Task)(nil), "go.micro.srv.taskmgr.Task")
	proto.RegisterType((*Response)(nil), "go.micro.srv.taskmgr.Response")
	proto.RegisterType((*Event)(nil), "go.micro.srv.taskmgr.Event")
	proto.RegisterEnum("go.micro.srv.taskmgr.Task_Status", Task_Status_name, Task_Status_value)
	proto.RegisterEnum("go.micro.srv.taskmgr.Event_OpCode", Event_OpCode_name, Event_OpCode_value)
}

func init() {
	proto.RegisterFile("proto/taskmgr/taskmgr.proto", fileDescriptor_taskmgr_f2b1a0b0af520394)
}

var fileDescriptor_taskmgr_f2b1a0b0af520394 = []byte{
	// 611 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x4f, 0x6f, 0xd3, 0x4e,
	0x10, 0x8d, 0x9d, 0x64, 0x9d, 0x4c, 0xd2, 0xc8, 0x5a, 0x55, 0xbf, 0xfa, 0x57, 0x10, 0xa4, 0x16,
	0x87, 0x9c, 0x8c, 0x54, 0x4e, 0x1c, 0x90, 0xa8, 0x6c, 0xa7, 0xb2, 0x54, 0x02, 0xda, 0xc4, 0xe7,
	0x6a, 0x89, 0x97, 0xca, 0x6a, 0xfd, 0x47, 0xeb, 0x4d, 0x45, 0xaf, 0x7c, 0x10, 0xbe, 0x01, 0x57,
	0x3e, 0x1f, 0xda, 0x59, 0xa7, 0x05, 0x64, 0xca, 0x81, 0x53, 0xf6, 0xbd, 0x79, 0x3b, 0x3b, 0x99,
	0xf7, 0x12, 0x78, 0x52, 0xcb, 0x4a, 0x55, 0x2f, 0x15, 0x6f, 0xae, 0x8b, 0x2b, 0xb9, 0xff, 0x0c,
	0x90, 0xa5, 0x87, 0x57, 0x55, 0x50, 0xe4, 0x5b, 0x59, 0x05, 0x8d, 0xbc, 0x0d, 0xda, 0x9a, 0x7f,
	0x08, 0x76, 0x12, 0xd1, 0x19, 0xd8, 0x79, 0xe6, 0x59, 0x73, 0x6b, 0xd1, 0x67, 0x76, 0x9e, 0xf9,
	0xdf, 0x06, 0x30, 0xd8, 0xf0, 0xe6, 0xfa, 0xf7, 0x02, 0xfd, 0x1f, 0x46, 0xe2, 0x56, 0x94, 0xea,
	0x32, 0xcf, 0x3c, 0x7b, 0x6e, 0x2d, 0xc6, 0xcc, 0x41, 0x9c, 0x64, 0xf4, 0x08, 0x9c, 0x5d, 0x23,
	0xa4, 0xae, 0xf4, 0x51, 0x4f, 0x34, 0x4c, 0x32, 0x4a, 0x61, 0x50, 0xf2, 0x42, 0x78, 0x03, 0xd4,
	0xe3, 0x59, 0x73, 0xea, 0xae, 0x16, 0xde, 0xd0, 0x70, 0xfa, 0x4c, 0x5f, 0x03, 0x69, 0x14, 0x57,
	0xbb, 0xc6, 0x23, 0x73, 0x6b, 0x31, 0x3b, 0x3d, 0x09, 0xba, 0x26, 0x0e, 0xf4, 0x5c, 0xc1, 0x1a,
	0x85, 0xac, 0xbd, 0x40, 0x4f, 0x60, 0xda, 0x28, 0x2e, 0xd5, 0xae, 0xbe, 0x54, 0x79, 0x21, 0x3c,
	0x67, 0x6e, 0x2d, 0x0e, 0xd8, 0xa4, 0xe5, 0x36, 0x79, 0x21, 0xe8, 0x0b, 0x98, 0x65, 0x5c, 0xf1,
	0xcb, 0xad, 0x28, 0x95, 0x99, 0x72, 0x84, 0x53, 0x4e, 0x35, 0x1b, 0x22, 0x99, 0x64, 0xf4, 0x39,
	0x4c, 0xb6, 0x52, 0x70, 0x25, 0x4c, 0x9f, 0xf1, 0xdc, 0x5a, 0x0c, 0x18, 0x18, 0x0a, 0xdb, 0x68,
	0x01, 0x2f, 0xb7, 0xe2, 0xc6, 0x08, 0xa0, 0x15, 0x20, 0xb5, 0x17, 0xec, 0xea, 0xec, 0xbe, 0xc3,
	0xc4, 0x08, 0x0c, 0xb5, 0x17, 0x48, 0xa1, 0x76, 0xb2, 0x34, 0x82, 0xa9, 0x11, 0x18, 0x0a, 0x05,
	0x87, 0x30, 0x14, 0x9f, 0x95, 0xe4, 0xde, 0xc1, 0xdc, 0x5a, 0x4c, 0x99, 0x01, 0xf4, 0x3f, 0x20,
	0x52, 0x34, 0xbb, 0x1b, 0xe5, 0xcd, 0x90, 0x6e, 0x91, 0xff, 0xc5, 0x02, 0x62, 0xb6, 0x41, 0xa7,
	0x30, 0x0a, 0x59, 0x7c, 0xb6, 0x49, 0x56, 0xe7, 0x6e, 0x8f, 0x4e, 0xc0, 0x41, 0x14, 0x47, 0xae,
	0xa5, 0x01, 0x4b, 0x57, 0x2b, 0x5d, 0xb1, 0xe9, 0x01, 0x8c, 0xc3, 0xb3, 0x55, 0x18, 0x5f, 0x68,
	0xd8, 0xc7, 0x6b, 0x08, 0xe3, 0xc8, 0x1d, 0x68, 0x94, 0x7e, 0x88, 0x4c, 0x93, 0xa1, 0xbe, 0x87,
	0x28, 0x8e, 0x5c, 0xa2, 0xc1, 0xf2, 0x2c, 0xb9, 0x48, 0x59, 0xec, 0x3a, 0x1a, 0xac, 0xd3, 0x30,
	0x8c, 0xd7, 0x6b, 0x77, 0xe4, 0x03, 0x8c, 0x98, 0x68, 0xea, 0xaa, 0x6c, 0x84, 0xff, 0xdd, 0x86,
	0x61, 0xac, 0x33, 0xf1, 0x53, 0x78, 0xc6, 0x18, 0x9e, 0x23, 0x70, 0xb4, 0x89, 0xfb, 0xec, 0xf4,
	0x19, 0xd1, 0x30, 0xc9, 0xe8, 0x53, 0x18, 0xeb, 0x5d, 0x34, 0x8a, 0x17, 0x75, 0x1b, 0x9e, 0x07,
	0x82, 0xbe, 0x85, 0x71, 0x55, 0x0b, 0xc9, 0x55, 0x5e, 0x95, 0x18, 0xa2, 0xd9, 0xa9, 0xdf, 0x1d,
	0x0d, 0x7c, 0x36, 0x78, 0x5f, 0x87, 0x55, 0x26, 0xd8, 0xc3, 0xa5, 0x0e, 0xef, 0x87, 0xdd, 0xde,
	0x17, 0x55, 0x96, 0x7f, 0xba, 0x33, 0xc6, 0x10, 0x94, 0x80, 0xa1, 0xd0, 0x98, 0x07, 0x0b, 0x9c,
	0x5f, 0x2c, 0x88, 0x81, 0x98, 0x37, 0x29, 0x00, 0x31, 0x3b, 0x77, 0x7b, 0x78, 0xc6, 0xb5, 0xba,
	0x96, 0x3e, 0x9b, 0x35, 0xba, 0xb6, 0x3e, 0xb3, 0x78, 0x93, 0xb2, 0x95, 0xdb, 0x47, 0x5b, 0xe2,
	0x65, 0xba, 0xd6, 0x9b, 0x3f, 0xfd, 0x6a, 0x83, 0xa3, 0xc3, 0xfd, 0xee, 0x4a, 0xd2, 0x37, 0xd0,
	0x3f, 0x17, 0x8a, 0x7a, 0xdd, 0xdf, 0x33, 0x89, 0x8e, 0x8f, 0xff, 0xfc, 0xe3, 0xf0, 0x7b, 0x74,
	0x09, 0x24, 0xc4, 0xcc, 0xd2, 0x47, 0x74, 0xc7, 0xcf, 0xba, 0x6b, 0xf7, 0x4e, 0xf6, 0x68, 0x04,
	0x24, 0xc4, 0x68, 0x3f, 0x32, 0xc9, 0xdf, 0xbb, 0x2c, 0x81, 0xa4, 0x98, 0xff, 0x7f, 0x9b, 0xe6,
	0x23, 0xc1, 0x3f, 0xb2, 0x57, 0x3f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x13, 0x20, 0x12, 0xe7,
	0x04, 0x00, 0x00,
}