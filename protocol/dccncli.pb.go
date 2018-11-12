// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dccncli.proto

package dccncli

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

// The dccn client request message containing the user's token
type AddTaskRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Region               string   `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	Zone                 string   `protobuf:"bytes,3,opt,name=zone,proto3" json:"zone,omitempty"`
	Usertoken            string   `protobuf:"bytes,4,opt,name=usertoken,proto3" json:"usertoken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddTaskRequest) Reset()         { *m = AddTaskRequest{} }
func (m *AddTaskRequest) String() string { return proto.CompactTextString(m) }
func (*AddTaskRequest) ProtoMessage()    {}
func (*AddTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{0}
}
func (m *AddTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTaskRequest.Unmarshal(m, b)
}
func (m *AddTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTaskRequest.Marshal(b, m, deterministic)
}
func (dst *AddTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTaskRequest.Merge(dst, src)
}
func (m *AddTaskRequest) XXX_Size() int {
	return xxx_messageInfo_AddTaskRequest.Size(m)
}
func (m *AddTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddTaskRequest proto.InternalMessageInfo

func (m *AddTaskRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AddTaskRequest) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *AddTaskRequest) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

func (m *AddTaskRequest) GetUsertoken() string {
	if m != nil {
		return m.Usertoken
	}
	return ""
}

// The Ankr Hub response message containing the success or failure
type AddTaskResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Taskid               int64    `protobuf:"varint,2,opt,name=taskid,proto3" json:"taskid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddTaskResponse) Reset()         { *m = AddTaskResponse{} }
func (m *AddTaskResponse) String() string { return proto.CompactTextString(m) }
func (*AddTaskResponse) ProtoMessage()    {}
func (*AddTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{1}
}
func (m *AddTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTaskResponse.Unmarshal(m, b)
}
func (m *AddTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTaskResponse.Marshal(b, m, deterministic)
}
func (dst *AddTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTaskResponse.Merge(dst, src)
}
func (m *AddTaskResponse) XXX_Size() int {
	return xxx_messageInfo_AddTaskResponse.Size(m)
}
func (m *AddTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddTaskResponse proto.InternalMessageInfo

func (m *AddTaskResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *AddTaskResponse) GetTaskid() int64 {
	if m != nil {
		return m.Taskid
	}
	return 0
}

// The Client List request message
type TaskListRequest struct {
	Usertoken            string   `protobuf:"bytes,1,opt,name=usertoken,proto3" json:"usertoken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskListRequest) Reset()         { *m = TaskListRequest{} }
func (m *TaskListRequest) String() string { return proto.CompactTextString(m) }
func (*TaskListRequest) ProtoMessage()    {}
func (*TaskListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{2}
}
func (m *TaskListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskListRequest.Unmarshal(m, b)
}
func (m *TaskListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskListRequest.Marshal(b, m, deterministic)
}
func (dst *TaskListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskListRequest.Merge(dst, src)
}
func (m *TaskListRequest) XXX_Size() int {
	return xxx_messageInfo_TaskListRequest.Size(m)
}
func (m *TaskListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TaskListRequest proto.InternalMessageInfo

func (m *TaskListRequest) GetUsertoken() string {
	if m != nil {
		return m.Usertoken
	}
	return ""
}

type TaskInfo struct {
	Taskid               int64    `protobuf:"varint,1,opt,name=taskid,proto3" json:"taskid,omitempty"`
	Taskname             string   `protobuf:"bytes,2,opt,name=taskname,proto3" json:"taskname,omitempty"`
	Uptime               uint32   `protobuf:"varint,3,opt,name=uptime,proto3" json:"uptime,omitempty"`
	Creationdate         uint64   `protobuf:"varint,4,opt,name=creationdate,proto3" json:"creationdate,omitempty"`
	Status               string   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskInfo) Reset()         { *m = TaskInfo{} }
func (m *TaskInfo) String() string { return proto.CompactTextString(m) }
func (*TaskInfo) ProtoMessage()    {}
func (*TaskInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{3}
}
func (m *TaskInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskInfo.Unmarshal(m, b)
}
func (m *TaskInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskInfo.Marshal(b, m, deterministic)
}
func (dst *TaskInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskInfo.Merge(dst, src)
}
func (m *TaskInfo) XXX_Size() int {
	return xxx_messageInfo_TaskInfo.Size(m)
}
func (m *TaskInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskInfo.DiscardUnknown(m)
}

var xxx_messageInfo_TaskInfo proto.InternalMessageInfo

func (m *TaskInfo) GetTaskid() int64 {
	if m != nil {
		return m.Taskid
	}
	return 0
}

func (m *TaskInfo) GetTaskname() string {
	if m != nil {
		return m.Taskname
	}
	return ""
}

func (m *TaskInfo) GetUptime() uint32 {
	if m != nil {
		return m.Uptime
	}
	return 0
}

func (m *TaskInfo) GetCreationdate() uint64 {
	if m != nil {
		return m.Creationdate
	}
	return 0
}

func (m *TaskInfo) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type TaskListResponse struct {
	Tasksinfo            []*TaskInfo `protobuf:"bytes,1,rep,name=tasksinfo,proto3" json:"tasksinfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TaskListResponse) Reset()         { *m = TaskListResponse{} }
func (m *TaskListResponse) String() string { return proto.CompactTextString(m) }
func (*TaskListResponse) ProtoMessage()    {}
func (*TaskListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{4}
}
func (m *TaskListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskListResponse.Unmarshal(m, b)
}
func (m *TaskListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskListResponse.Marshal(b, m, deterministic)
}
func (dst *TaskListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskListResponse.Merge(dst, src)
}
func (m *TaskListResponse) XXX_Size() int {
	return xxx_messageInfo_TaskListResponse.Size(m)
}
func (m *TaskListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TaskListResponse proto.InternalMessageInfo

func (m *TaskListResponse) GetTasksinfo() []*TaskInfo {
	if m != nil {
		return m.Tasksinfo
	}
	return nil
}

type CancelTaskRequest struct {
	Usertoken            string   `protobuf:"bytes,1,opt,name=usertoken,proto3" json:"usertoken,omitempty"`
	Taskid               int64    `protobuf:"varint,2,opt,name=taskid,proto3" json:"taskid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelTaskRequest) Reset()         { *m = CancelTaskRequest{} }
func (m *CancelTaskRequest) String() string { return proto.CompactTextString(m) }
func (*CancelTaskRequest) ProtoMessage()    {}
func (*CancelTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{5}
}
func (m *CancelTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelTaskRequest.Unmarshal(m, b)
}
func (m *CancelTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelTaskRequest.Marshal(b, m, deterministic)
}
func (dst *CancelTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelTaskRequest.Merge(dst, src)
}
func (m *CancelTaskRequest) XXX_Size() int {
	return xxx_messageInfo_CancelTaskRequest.Size(m)
}
func (m *CancelTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CancelTaskRequest proto.InternalMessageInfo

func (m *CancelTaskRequest) GetUsertoken() string {
	if m != nil {
		return m.Usertoken
	}
	return ""
}

func (m *CancelTaskRequest) GetTaskid() int64 {
	if m != nil {
		return m.Taskid
	}
	return 0
}

type CancelTaskResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelTaskResponse) Reset()         { *m = CancelTaskResponse{} }
func (m *CancelTaskResponse) String() string { return proto.CompactTextString(m) }
func (*CancelTaskResponse) ProtoMessage()    {}
func (*CancelTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{6}
}
func (m *CancelTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelTaskResponse.Unmarshal(m, b)
}
func (m *CancelTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelTaskResponse.Marshal(b, m, deterministic)
}
func (dst *CancelTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelTaskResponse.Merge(dst, src)
}
func (m *CancelTaskResponse) XXX_Size() int {
	return xxx_messageInfo_CancelTaskResponse.Size(m)
}
func (m *CancelTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CancelTaskResponse proto.InternalMessageInfo

func (m *CancelTaskResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type ReportRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Report               string   `protobuf:"bytes,2,opt,name=report,proto3" json:"report,omitempty"`
	Host                 string   `protobuf:"bytes,3,opt,name=host,proto3" json:"host,omitempty"`
	Port                 int64    `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportRequest) Reset()         { *m = ReportRequest{} }
func (m *ReportRequest) String() string { return proto.CompactTextString(m) }
func (*ReportRequest) ProtoMessage()    {}
func (*ReportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{7}
}
func (m *ReportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportRequest.Unmarshal(m, b)
}
func (m *ReportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportRequest.Marshal(b, m, deterministic)
}
func (dst *ReportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportRequest.Merge(dst, src)
}
func (m *ReportRequest) XXX_Size() int {
	return xxx_messageInfo_ReportRequest.Size(m)
}
func (m *ReportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportRequest proto.InternalMessageInfo

func (m *ReportRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ReportRequest) GetReport() string {
	if m != nil {
		return m.Report
	}
	return ""
}

func (m *ReportRequest) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *ReportRequest) GetPort() int64 {
	if m != nil {
		return m.Port
	}
	return 0
}

type ReportResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportResponse) Reset()         { *m = ReportResponse{} }
func (m *ReportResponse) String() string { return proto.CompactTextString(m) }
func (*ReportResponse) ProtoMessage()    {}
func (*ReportResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{8}
}
func (m *ReportResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportResponse.Unmarshal(m, b)
}
func (m *ReportResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportResponse.Marshal(b, m, deterministic)
}
func (dst *ReportResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportResponse.Merge(dst, src)
}
func (m *ReportResponse) XXX_Size() int {
	return xxx_messageInfo_ReportResponse.Size(m)
}
func (m *ReportResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReportResponse proto.InternalMessageInfo

func (m *ReportResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type QueryTaskRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryTaskRequest) Reset()         { *m = QueryTaskRequest{} }
func (m *QueryTaskRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTaskRequest) ProtoMessage()    {}
func (*QueryTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{9}
}
func (m *QueryTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryTaskRequest.Unmarshal(m, b)
}
func (m *QueryTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryTaskRequest.Marshal(b, m, deterministic)
}
func (dst *QueryTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTaskRequest.Merge(dst, src)
}
func (m *QueryTaskRequest) XXX_Size() int {
	return xxx_messageInfo_QueryTaskRequest.Size(m)
}
func (m *QueryTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTaskRequest proto.InternalMessageInfo

func (m *QueryTaskRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type QueryTaskResponse struct {
	Taskid               int64    `protobuf:"varint,1,opt,name=taskid,proto3" json:"taskid,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Extra                string   `protobuf:"bytes,3,opt,name=extra,proto3" json:"extra,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryTaskResponse) Reset()         { *m = QueryTaskResponse{} }
func (m *QueryTaskResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTaskResponse) ProtoMessage()    {}
func (*QueryTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dccncli_2ff1a87dec470137, []int{10}
}
func (m *QueryTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryTaskResponse.Unmarshal(m, b)
}
func (m *QueryTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryTaskResponse.Marshal(b, m, deterministic)
}
func (dst *QueryTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTaskResponse.Merge(dst, src)
}
func (m *QueryTaskResponse) XXX_Size() int {
	return xxx_messageInfo_QueryTaskResponse.Size(m)
}
func (m *QueryTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTaskResponse proto.InternalMessageInfo

func (m *QueryTaskResponse) GetTaskid() int64 {
	if m != nil {
		return m.Taskid
	}
	return 0
}

func (m *QueryTaskResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *QueryTaskResponse) GetExtra() string {
	if m != nil {
		return m.Extra
	}
	return ""
}

func init() {
	proto.RegisterType((*AddTaskRequest)(nil), "dccncli.AddTaskRequest")
	proto.RegisterType((*AddTaskResponse)(nil), "dccncli.AddTaskResponse")
	proto.RegisterType((*TaskListRequest)(nil), "dccncli.TaskListRequest")
	proto.RegisterType((*TaskInfo)(nil), "dccncli.TaskInfo")
	proto.RegisterType((*TaskListResponse)(nil), "dccncli.TaskListResponse")
	proto.RegisterType((*CancelTaskRequest)(nil), "dccncli.CancelTaskRequest")
	proto.RegisterType((*CancelTaskResponse)(nil), "dccncli.CancelTaskResponse")
	proto.RegisterType((*ReportRequest)(nil), "dccncli.ReportRequest")
	proto.RegisterType((*ReportResponse)(nil), "dccncli.ReportResponse")
	proto.RegisterType((*QueryTaskRequest)(nil), "dccncli.QueryTaskRequest")
	proto.RegisterType((*QueryTaskResponse)(nil), "dccncli.QueryTaskResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DccncliClient is the client API for Dccncli service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DccncliClient interface {
	// Sends request to start a task and list task
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error)
	TaskList(ctx context.Context, in *TaskListRequest, opts ...grpc.CallOption) (*TaskListResponse, error)
	CancelTask(ctx context.Context, in *CancelTaskRequest, opts ...grpc.CallOption) (*CancelTaskResponse, error)
	K8ReportStatus(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error)
	K8QueryTask(ctx context.Context, in *QueryTaskRequest, opts ...grpc.CallOption) (*QueryTaskResponse, error)
}

type dccncliClient struct {
	cc *grpc.ClientConn
}

func NewDccncliClient(cc *grpc.ClientConn) DccncliClient {
	return &dccncliClient{cc}
}

func (c *dccncliClient) AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error) {
	out := new(AddTaskResponse)
	err := c.cc.Invoke(ctx, "/dccncli.dccncli/AddTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dccncliClient) TaskList(ctx context.Context, in *TaskListRequest, opts ...grpc.CallOption) (*TaskListResponse, error) {
	out := new(TaskListResponse)
	err := c.cc.Invoke(ctx, "/dccncli.dccncli/TaskList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dccncliClient) CancelTask(ctx context.Context, in *CancelTaskRequest, opts ...grpc.CallOption) (*CancelTaskResponse, error) {
	out := new(CancelTaskResponse)
	err := c.cc.Invoke(ctx, "/dccncli.dccncli/CancelTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dccncliClient) K8ReportStatus(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error) {
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, "/dccncli.dccncli/K8ReportStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dccncliClient) K8QueryTask(ctx context.Context, in *QueryTaskRequest, opts ...grpc.CallOption) (*QueryTaskResponse, error) {
	out := new(QueryTaskResponse)
	err := c.cc.Invoke(ctx, "/dccncli.dccncli/K8QueryTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DccncliServer is the server API for Dccncli service.
type DccncliServer interface {
	// Sends request to start a task and list task
	AddTask(context.Context, *AddTaskRequest) (*AddTaskResponse, error)
	TaskList(context.Context, *TaskListRequest) (*TaskListResponse, error)
	CancelTask(context.Context, *CancelTaskRequest) (*CancelTaskResponse, error)
	K8ReportStatus(context.Context, *ReportRequest) (*ReportResponse, error)
	K8QueryTask(context.Context, *QueryTaskRequest) (*QueryTaskResponse, error)
}

func RegisterDccncliServer(s *grpc.Server, srv DccncliServer) {
	s.RegisterService(&_Dccncli_serviceDesc, srv)
}

func _Dccncli_AddTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DccncliServer).AddTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dccncli.dccncli/AddTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DccncliServer).AddTask(ctx, req.(*AddTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dccncli_TaskList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DccncliServer).TaskList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dccncli.dccncli/TaskList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DccncliServer).TaskList(ctx, req.(*TaskListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dccncli_CancelTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DccncliServer).CancelTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dccncli.dccncli/CancelTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DccncliServer).CancelTask(ctx, req.(*CancelTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dccncli_K8ReportStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DccncliServer).K8ReportStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dccncli.dccncli/K8ReportStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DccncliServer).K8ReportStatus(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dccncli_K8QueryTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DccncliServer).K8QueryTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dccncli.dccncli/K8QueryTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DccncliServer).K8QueryTask(ctx, req.(*QueryTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Dccncli_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dccncli.dccncli",
	HandlerType: (*DccncliServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTask",
			Handler:    _Dccncli_AddTask_Handler,
		},
		{
			MethodName: "TaskList",
			Handler:    _Dccncli_TaskList_Handler,
		},
		{
			MethodName: "CancelTask",
			Handler:    _Dccncli_CancelTask_Handler,
		},
		{
			MethodName: "K8ReportStatus",
			Handler:    _Dccncli_K8ReportStatus_Handler,
		},
		{
			MethodName: "K8QueryTask",
			Handler:    _Dccncli_K8QueryTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dccncli.proto",
}

func init() { proto.RegisterFile("dccncli.proto", fileDescriptor_dccncli_2ff1a87dec470137) }

var fileDescriptor_dccncli_2ff1a87dec470137 = []byte{
	// 506 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xcd, 0x6e, 0x13, 0x31,
	0x10, 0x66, 0x93, 0xf4, 0x27, 0xd3, 0x26, 0x6d, 0x2c, 0x44, 0x97, 0x85, 0x43, 0xe5, 0x03, 0x8a,
	0x10, 0x4a, 0xa5, 0x72, 0xe9, 0x09, 0x29, 0x0d, 0x02, 0x55, 0xe5, 0x50, 0x16, 0x78, 0x00, 0xe3,
	0xb8, 0xc5, 0x4a, 0x62, 0x2f, 0x6b, 0xaf, 0x54, 0x78, 0x0c, 0x1e, 0x91, 0x27, 0x41, 0xf6, 0xda,
	0x5e, 0x6f, 0x13, 0x9a, 0xdb, 0xfc, 0xf9, 0x9b, 0x19, 0x7f, 0x9f, 0x0d, 0x83, 0x39, 0xa5, 0x82,
	0x2e, 0xf9, 0xa4, 0x28, 0xa5, 0x96, 0x68, 0xcf, 0xb9, 0x58, 0xc0, 0x70, 0x3a, 0x9f, 0x7f, 0x25,
	0x6a, 0x91, 0xb3, 0x9f, 0x15, 0x53, 0x1a, 0x21, 0xe8, 0x09, 0xb2, 0x62, 0x69, 0x72, 0x9a, 0x8c,
	0xfb, 0xb9, 0xb5, 0xd1, 0x33, 0xd8, 0x2d, 0xd9, 0x1d, 0x97, 0x22, 0xed, 0xd8, 0xa8, 0xf3, 0x4c,
	0xed, 0x6f, 0x29, 0x58, 0xda, 0xad, 0x6b, 0x8d, 0x8d, 0x5e, 0x42, 0xbf, 0x52, 0xac, 0xd4, 0x72,
	0xc1, 0x44, 0xda, 0xb3, 0x89, 0x26, 0x80, 0xa7, 0x70, 0x14, 0xfa, 0xa9, 0x42, 0x0a, 0x65, 0xc1,
	0x95, 0x26, 0xba, 0x52, 0xae, 0xa5, 0xf3, 0x4c, 0x5c, 0x13, 0xb5, 0xe0, 0x73, 0xdb, 0xb4, 0x9b,
	0x3b, 0x0f, 0x9f, 0xc1, 0x91, 0x39, 0xff, 0x89, 0x2b, 0xed, 0x67, 0x6e, 0xf5, 0x4c, 0x1e, 0xf6,
	0xfc, 0x93, 0xc0, 0xbe, 0x39, 0x71, 0x25, 0x6e, 0x65, 0x84, 0x9a, 0xc4, 0xa8, 0x28, 0x83, 0x7d,
	0x63, 0xd9, 0xd5, 0xeb, 0x25, 0x83, 0x6f, 0xce, 0x54, 0x85, 0xe6, 0xab, 0x7a, 0xd1, 0x41, 0xee,
	0x3c, 0x84, 0xe1, 0x90, 0x96, 0x8c, 0x68, 0x2e, 0xc5, 0x9c, 0x68, 0x66, 0xb7, 0xed, 0xe5, 0xad,
	0x58, 0xb4, 0xdd, 0x4e, 0xbc, 0x1d, 0x9e, 0xc1, 0x71, 0xb3, 0x85, 0xbb, 0x89, 0x33, 0xe8, 0x9b,
	0x9e, 0x8a, 0x8b, 0x5b, 0x99, 0x26, 0xa7, 0xdd, 0xf1, 0xc1, 0xf9, 0x68, 0xe2, 0x89, 0xf3, 0x1b,
	0xe4, 0x4d, 0x0d, 0xbe, 0x82, 0xd1, 0x8c, 0x08, 0xca, 0x96, 0x31, 0x81, 0x8f, 0x5e, 0xc6, 0x7f,
	0x6f, 0xf5, 0x0d, 0xa0, 0x18, 0xea, 0x71, 0x6e, 0x30, 0x85, 0x41, 0xce, 0x0a, 0x59, 0xea, 0xad,
	0xaa, 0x31, 0x45, 0x8d, 0x6a, 0x8c, 0x67, 0x6a, 0x7f, 0x48, 0xa5, 0xbd, 0x6a, 0x8c, 0x6d, 0x62,
	0xb6, 0xb2, 0x67, 0x87, 0xb2, 0x36, 0x1e, 0xc3, 0xd0, 0x37, 0xd9, 0x32, 0xce, 0x2b, 0x38, 0xfe,
	0x5c, 0xb1, 0xf2, 0xd7, 0x16, 0x1d, 0xe3, 0x6f, 0x30, 0x8a, 0xea, 0x1a, 0xd0, 0x8d, 0x8a, 0xf0,
	0x00, 0x9d, 0x68, 0xa5, 0xa7, 0xb0, 0xc3, 0xee, 0x75, 0x49, 0xdc, 0xec, 0xb5, 0x73, 0xfe, 0xb7,
	0x03, 0xfe, 0x41, 0xa1, 0x77, 0xb0, 0xe7, 0x04, 0x8e, 0x4e, 0x02, 0x77, 0xed, 0x27, 0x96, 0xa5,
	0xeb, 0x89, 0x7a, 0x16, 0xfc, 0x04, 0x4d, 0x6b, 0xad, 0x1a, 0x5d, 0xa0, 0xb4, 0x45, 0x7e, 0x24,
	0xf8, 0xec, 0xf9, 0x86, 0x4c, 0x80, 0xf8, 0x08, 0xd0, 0x50, 0x89, 0xb2, 0x50, 0xba, 0x26, 0x95,
	0xec, 0xc5, 0xc6, 0x5c, 0x00, 0x9a, 0xc1, 0xf0, 0xfa, 0xa2, 0xa6, 0xe0, 0x8b, 0x7b, 0x93, 0xe1,
	0x40, 0x8b, 0xfe, 0xec, 0x64, 0x2d, 0x1e, 0x40, 0x3e, 0xc0, 0xc1, 0xf5, 0x45, 0xb8, 0x75, 0xd4,
	0x4c, 0xfe, 0x90, 0xb1, 0x2c, 0xdb, 0x94, 0xf2, 0x38, 0x97, 0xaf, 0x21, 0xe5, 0x72, 0x72, 0x57,
	0x16, 0x74, 0xc2, 0xee, 0xc9, 0xaa, 0x58, 0x32, 0xe5, 0xeb, 0x2f, 0x0f, 0xdf, 0x53, 0x2a, 0x66,
	0x4b, 0x7e, 0x63, 0x3e, 0xb7, 0x9b, 0xe4, 0xfb, 0xae, 0xfd, 0xe5, 0xde, 0xfe, 0x0b, 0x00, 0x00,
	0xff, 0xff, 0x53, 0xc6, 0xe8, 0x65, 0xf6, 0x04, 0x00, 0x00,
}
