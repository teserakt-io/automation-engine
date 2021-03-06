// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// List of supported ActionType
type ActionType int32

const (
	ActionType_UNDEFINED_ACTION ActionType = 0
	ActionType_KEY_ROTATION     ActionType = 1
)

var ActionType_name = map[int32]string{
	0: "UNDEFINED_ACTION",
	1: "KEY_ROTATION",
}

var ActionType_value = map[string]int32{
	"UNDEFINED_ACTION": 0,
	"KEY_ROTATION":     1,
}

func (x ActionType) String() string {
	return proto.EnumName(ActionType_name, int32(x))
}

func (ActionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

// List of supported TargetType
type TargetType int32

const (
	TargetType_ANY    TargetType = 0
	TargetType_TOPIC  TargetType = 1
	TargetType_CLIENT TargetType = 2
)

var TargetType_name = map[int32]string{
	0: "ANY",
	1: "TOPIC",
	2: "CLIENT",
}

var TargetType_value = map[string]int32{
	"ANY":    0,
	"TOPIC":  1,
	"CLIENT": 2,
}

func (x TargetType) String() string {
	return proto.EnumName(TargetType_name, int32(x))
}

func (TargetType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

// List of supported TriggerType
type TriggerType int32

const (
	TriggerType_UNDEFINED_TRIGGER TriggerType = 0
	TriggerType_TIME_INTERVAL     TriggerType = 1
	TriggerType_EVENT             TriggerType = 2
)

var TriggerType_name = map[int32]string{
	0: "UNDEFINED_TRIGGER",
	1: "TIME_INTERVAL",
	2: "EVENT",
}

var TriggerType_value = map[string]int32{
	"UNDEFINED_TRIGGER": 0,
	"TIME_INTERVAL":     1,
	"EVENT":             2,
}

func (x TriggerType) String() string {
	return proto.EnumName(TriggerType_name, int32(x))
}

func (TriggerType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

type Rule struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Description          string               `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Action               ActionType           `protobuf:"varint,3,opt,name=action,proto3,enum=pb.ActionType" json:"action,omitempty"`
	LastExecuted         *timestamp.Timestamp `protobuf:"bytes,4,opt,name=lastExecuted,proto3" json:"lastExecuted,omitempty"`
	Triggers             []*Trigger           `protobuf:"bytes,5,rep,name=triggers,proto3" json:"triggers,omitempty"`
	Targets              []*Target            `protobuf:"bytes,6,rep,name=targets,proto3" json:"targets,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Rule) Reset()         { *m = Rule{} }
func (m *Rule) String() string { return proto.CompactTextString(m) }
func (*Rule) ProtoMessage()    {}
func (*Rule) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *Rule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Rule.Unmarshal(m, b)
}
func (m *Rule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Rule.Marshal(b, m, deterministic)
}
func (m *Rule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Rule.Merge(m, src)
}
func (m *Rule) XXX_Size() int {
	return xxx_messageInfo_Rule.Size(m)
}
func (m *Rule) XXX_DiscardUnknown() {
	xxx_messageInfo_Rule.DiscardUnknown(m)
}

var xxx_messageInfo_Rule proto.InternalMessageInfo

func (m *Rule) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Rule) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Rule) GetAction() ActionType {
	if m != nil {
		return m.Action
	}
	return ActionType_UNDEFINED_ACTION
}

func (m *Rule) GetLastExecuted() *timestamp.Timestamp {
	if m != nil {
		return m.LastExecuted
	}
	return nil
}

func (m *Rule) GetTriggers() []*Trigger {
	if m != nil {
		return m.Triggers
	}
	return nil
}

func (m *Rule) GetTargets() []*Target {
	if m != nil {
		return m.Targets
	}
	return nil
}

type Target struct {
	Id                   int32      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 TargetType `protobuf:"varint,2,opt,name=type,proto3,enum=pb.TargetType" json:"type,omitempty"`
	Expr                 string     `protobuf:"bytes,3,opt,name=expr,proto3" json:"expr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Target) Reset()         { *m = Target{} }
func (m *Target) String() string { return proto.CompactTextString(m) }
func (*Target) ProtoMessage()    {}
func (*Target) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *Target) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Target.Unmarshal(m, b)
}
func (m *Target) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Target.Marshal(b, m, deterministic)
}
func (m *Target) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Target.Merge(m, src)
}
func (m *Target) XXX_Size() int {
	return xxx_messageInfo_Target.Size(m)
}
func (m *Target) XXX_DiscardUnknown() {
	xxx_messageInfo_Target.DiscardUnknown(m)
}

var xxx_messageInfo_Target proto.InternalMessageInfo

func (m *Target) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Target) GetType() TargetType {
	if m != nil {
		return m.Type
	}
	return TargetType_ANY
}

func (m *Target) GetExpr() string {
	if m != nil {
		return m.Expr
	}
	return ""
}

type Trigger struct {
	Id                   int32       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 TriggerType `protobuf:"varint,2,opt,name=type,proto3,enum=pb.TriggerType" json:"type,omitempty"`
	Settings             []byte      `protobuf:"bytes,3,opt,name=settings,proto3" json:"settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Trigger) Reset()         { *m = Trigger{} }
func (m *Trigger) String() string { return proto.CompactTextString(m) }
func (*Trigger) ProtoMessage()    {}
func (*Trigger) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *Trigger) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Trigger.Unmarshal(m, b)
}
func (m *Trigger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Trigger.Marshal(b, m, deterministic)
}
func (m *Trigger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Trigger.Merge(m, src)
}
func (m *Trigger) XXX_Size() int {
	return xxx_messageInfo_Trigger.Size(m)
}
func (m *Trigger) XXX_DiscardUnknown() {
	xxx_messageInfo_Trigger.DiscardUnknown(m)
}

var xxx_messageInfo_Trigger proto.InternalMessageInfo

func (m *Trigger) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Trigger) GetType() TriggerType {
	if m != nil {
		return m.Type
	}
	return TriggerType_UNDEFINED_TRIGGER
}

func (m *Trigger) GetSettings() []byte {
	if m != nil {
		return m.Settings
	}
	return nil
}

type RulesResponse struct {
	Rules                []*Rule  `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RulesResponse) Reset()         { *m = RulesResponse{} }
func (m *RulesResponse) String() string { return proto.CompactTextString(m) }
func (*RulesResponse) ProtoMessage()    {}
func (*RulesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *RulesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RulesResponse.Unmarshal(m, b)
}
func (m *RulesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RulesResponse.Marshal(b, m, deterministic)
}
func (m *RulesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RulesResponse.Merge(m, src)
}
func (m *RulesResponse) XXX_Size() int {
	return xxx_messageInfo_RulesResponse.Size(m)
}
func (m *RulesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RulesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RulesResponse proto.InternalMessageInfo

func (m *RulesResponse) GetRules() []*Rule {
	if m != nil {
		return m.Rules
	}
	return nil
}

type RuleResponse struct {
	Rule                 *Rule    `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RuleResponse) Reset()         { *m = RuleResponse{} }
func (m *RuleResponse) String() string { return proto.CompactTextString(m) }
func (*RuleResponse) ProtoMessage()    {}
func (*RuleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *RuleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuleResponse.Unmarshal(m, b)
}
func (m *RuleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuleResponse.Marshal(b, m, deterministic)
}
func (m *RuleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuleResponse.Merge(m, src)
}
func (m *RuleResponse) XXX_Size() int {
	return xxx_messageInfo_RuleResponse.Size(m)
}
func (m *RuleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RuleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RuleResponse proto.InternalMessageInfo

func (m *RuleResponse) GetRule() *Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

type ListRulesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRulesRequest) Reset()         { *m = ListRulesRequest{} }
func (m *ListRulesRequest) String() string { return proto.CompactTextString(m) }
func (*ListRulesRequest) ProtoMessage()    {}
func (*ListRulesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *ListRulesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRulesRequest.Unmarshal(m, b)
}
func (m *ListRulesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRulesRequest.Marshal(b, m, deterministic)
}
func (m *ListRulesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRulesRequest.Merge(m, src)
}
func (m *ListRulesRequest) XXX_Size() int {
	return xxx_messageInfo_ListRulesRequest.Size(m)
}
func (m *ListRulesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRulesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRulesRequest proto.InternalMessageInfo

type GetRuleRequest struct {
	RuleId               int32    `protobuf:"varint,1,opt,name=ruleId,proto3" json:"ruleId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRuleRequest) Reset()         { *m = GetRuleRequest{} }
func (m *GetRuleRequest) String() string { return proto.CompactTextString(m) }
func (*GetRuleRequest) ProtoMessage()    {}
func (*GetRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *GetRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRuleRequest.Unmarshal(m, b)
}
func (m *GetRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRuleRequest.Marshal(b, m, deterministic)
}
func (m *GetRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRuleRequest.Merge(m, src)
}
func (m *GetRuleRequest) XXX_Size() int {
	return xxx_messageInfo_GetRuleRequest.Size(m)
}
func (m *GetRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRuleRequest proto.InternalMessageInfo

func (m *GetRuleRequest) GetRuleId() int32 {
	if m != nil {
		return m.RuleId
	}
	return 0
}

type AddRuleRequest struct {
	Description          string     `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Action               ActionType `protobuf:"varint,2,opt,name=action,proto3,enum=pb.ActionType" json:"action,omitempty"`
	Triggers             []*Trigger `protobuf:"bytes,3,rep,name=triggers,proto3" json:"triggers,omitempty"`
	Targets              []*Target  `protobuf:"bytes,4,rep,name=targets,proto3" json:"targets,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AddRuleRequest) Reset()         { *m = AddRuleRequest{} }
func (m *AddRuleRequest) String() string { return proto.CompactTextString(m) }
func (*AddRuleRequest) ProtoMessage()    {}
func (*AddRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{7}
}

func (m *AddRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRuleRequest.Unmarshal(m, b)
}
func (m *AddRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRuleRequest.Marshal(b, m, deterministic)
}
func (m *AddRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRuleRequest.Merge(m, src)
}
func (m *AddRuleRequest) XXX_Size() int {
	return xxx_messageInfo_AddRuleRequest.Size(m)
}
func (m *AddRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRuleRequest proto.InternalMessageInfo

func (m *AddRuleRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *AddRuleRequest) GetAction() ActionType {
	if m != nil {
		return m.Action
	}
	return ActionType_UNDEFINED_ACTION
}

func (m *AddRuleRequest) GetTriggers() []*Trigger {
	if m != nil {
		return m.Triggers
	}
	return nil
}

func (m *AddRuleRequest) GetTargets() []*Target {
	if m != nil {
		return m.Targets
	}
	return nil
}

// UpdateRuleRequest will fetch the rule identified by ruleId,
// and override its description, action, triggers and targets values
// with those provided.
type UpdateRuleRequest struct {
	RuleId               int32      `protobuf:"varint,1,opt,name=ruleId,proto3" json:"ruleId,omitempty"`
	Description          string     `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Action               ActionType `protobuf:"varint,3,opt,name=action,proto3,enum=pb.ActionType" json:"action,omitempty"`
	Triggers             []*Trigger `protobuf:"bytes,4,rep,name=triggers,proto3" json:"triggers,omitempty"`
	Targets              []*Target  `protobuf:"bytes,5,rep,name=targets,proto3" json:"targets,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *UpdateRuleRequest) Reset()         { *m = UpdateRuleRequest{} }
func (m *UpdateRuleRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRuleRequest) ProtoMessage()    {}
func (*UpdateRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{8}
}

func (m *UpdateRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRuleRequest.Unmarshal(m, b)
}
func (m *UpdateRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRuleRequest.Marshal(b, m, deterministic)
}
func (m *UpdateRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRuleRequest.Merge(m, src)
}
func (m *UpdateRuleRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRuleRequest.Size(m)
}
func (m *UpdateRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRuleRequest proto.InternalMessageInfo

func (m *UpdateRuleRequest) GetRuleId() int32 {
	if m != nil {
		return m.RuleId
	}
	return 0
}

func (m *UpdateRuleRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *UpdateRuleRequest) GetAction() ActionType {
	if m != nil {
		return m.Action
	}
	return ActionType_UNDEFINED_ACTION
}

func (m *UpdateRuleRequest) GetTriggers() []*Trigger {
	if m != nil {
		return m.Triggers
	}
	return nil
}

func (m *UpdateRuleRequest) GetTargets() []*Target {
	if m != nil {
		return m.Targets
	}
	return nil
}

type DeleteRuleRequest struct {
	RuleId               int32    `protobuf:"varint,1,opt,name=ruleId,proto3" json:"ruleId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteRuleRequest) Reset()         { *m = DeleteRuleRequest{} }
func (m *DeleteRuleRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRuleRequest) ProtoMessage()    {}
func (*DeleteRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{9}
}

func (m *DeleteRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRuleRequest.Unmarshal(m, b)
}
func (m *DeleteRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRuleRequest.Marshal(b, m, deterministic)
}
func (m *DeleteRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRuleRequest.Merge(m, src)
}
func (m *DeleteRuleRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteRuleRequest.Size(m)
}
func (m *DeleteRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRuleRequest proto.InternalMessageInfo

func (m *DeleteRuleRequest) GetRuleId() int32 {
	if m != nil {
		return m.RuleId
	}
	return 0
}

type DeleteRuleResponse struct {
	RuleId               int32    `protobuf:"varint,1,opt,name=ruleId,proto3" json:"ruleId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteRuleResponse) Reset()         { *m = DeleteRuleResponse{} }
func (m *DeleteRuleResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteRuleResponse) ProtoMessage()    {}
func (*DeleteRuleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{10}
}

func (m *DeleteRuleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRuleResponse.Unmarshal(m, b)
}
func (m *DeleteRuleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRuleResponse.Marshal(b, m, deterministic)
}
func (m *DeleteRuleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRuleResponse.Merge(m, src)
}
func (m *DeleteRuleResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteRuleResponse.Size(m)
}
func (m *DeleteRuleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRuleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRuleResponse proto.InternalMessageInfo

func (m *DeleteRuleResponse) GetRuleId() int32 {
	if m != nil {
		return m.RuleId
	}
	return 0
}

type HealthCheckRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthCheckRequest) Reset()         { *m = HealthCheckRequest{} }
func (m *HealthCheckRequest) String() string { return proto.CompactTextString(m) }
func (*HealthCheckRequest) ProtoMessage()    {}
func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{11}
}

func (m *HealthCheckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheckRequest.Unmarshal(m, b)
}
func (m *HealthCheckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheckRequest.Marshal(b, m, deterministic)
}
func (m *HealthCheckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckRequest.Merge(m, src)
}
func (m *HealthCheckRequest) XXX_Size() int {
	return xxx_messageInfo_HealthCheckRequest.Size(m)
}
func (m *HealthCheckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckRequest proto.InternalMessageInfo

type HealthCheckResponse struct {
	Code                 int64    `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthCheckResponse) Reset()         { *m = HealthCheckResponse{} }
func (m *HealthCheckResponse) String() string { return proto.CompactTextString(m) }
func (*HealthCheckResponse) ProtoMessage()    {}
func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{12}
}

func (m *HealthCheckResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheckResponse.Unmarshal(m, b)
}
func (m *HealthCheckResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheckResponse.Marshal(b, m, deterministic)
}
func (m *HealthCheckResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckResponse.Merge(m, src)
}
func (m *HealthCheckResponse) XXX_Size() int {
	return xxx_messageInfo_HealthCheckResponse.Size(m)
}
func (m *HealthCheckResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckResponse proto.InternalMessageInfo

func (m *HealthCheckResponse) GetCode() int64 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *HealthCheckResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterEnum("pb.ActionType", ActionType_name, ActionType_value)
	proto.RegisterEnum("pb.TargetType", TargetType_name, TargetType_value)
	proto.RegisterEnum("pb.TriggerType", TriggerType_name, TriggerType_value)
	proto.RegisterType((*Rule)(nil), "pb.Rule")
	proto.RegisterType((*Target)(nil), "pb.Target")
	proto.RegisterType((*Trigger)(nil), "pb.Trigger")
	proto.RegisterType((*RulesResponse)(nil), "pb.RulesResponse")
	proto.RegisterType((*RuleResponse)(nil), "pb.RuleResponse")
	proto.RegisterType((*ListRulesRequest)(nil), "pb.ListRulesRequest")
	proto.RegisterType((*GetRuleRequest)(nil), "pb.GetRuleRequest")
	proto.RegisterType((*AddRuleRequest)(nil), "pb.AddRuleRequest")
	proto.RegisterType((*UpdateRuleRequest)(nil), "pb.UpdateRuleRequest")
	proto.RegisterType((*DeleteRuleRequest)(nil), "pb.DeleteRuleRequest")
	proto.RegisterType((*DeleteRuleResponse)(nil), "pb.DeleteRuleResponse")
	proto.RegisterType((*HealthCheckRequest)(nil), "pb.HealthCheckRequest")
	proto.RegisterType((*HealthCheckResponse)(nil), "pb.HealthCheckResponse")
}

func init() {
	proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c)
}

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 850 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xdd, 0x6e, 0xe3, 0x44,
	0x18, 0xad, 0xf3, 0x9f, 0x2f, 0x6d, 0xea, 0x7c, 0xb4, 0xbb, 0x91, 0xb5, 0x82, 0xc8, 0x20, 0x88,
	0x42, 0x6b, 0xef, 0x06, 0x04, 0x55, 0x2f, 0x90, 0xdc, 0xd4, 0xdb, 0x46, 0x94, 0x74, 0x35, 0xb8,
	0x2b, 0x2d, 0x37, 0x95, 0x9b, 0x0c, 0xa9, 0x21, 0xb5, 0x4d, 0x66, 0x22, 0x76, 0x85, 0xb8, 0xe1,
	0x11, 0xe0, 0x29, 0x78, 0x0f, 0xde, 0x80, 0x07, 0xe0, 0x86, 0x7b, 0x5e, 0x01, 0xf9, 0xb3, 0x9d,
	0x38, 0x35, 0x5d, 0xe5, 0x82, 0xab, 0x7a, 0xce, 0x77, 0x7c, 0x66, 0x8e, 0xe7, 0x9c, 0x06, 0xea,
	0x6e, 0xe8, 0x19, 0xe1, 0x3c, 0x90, 0x01, 0x16, 0xc2, 0x1b, 0xed, 0xbd, 0x69, 0x10, 0x4c, 0x67,
	0xdc, 0x24, 0xe4, 0x66, 0xf1, 0xad, 0x29, 0xbd, 0x3b, 0x2e, 0xa4, 0x7b, 0x17, 0xc6, 0x24, 0xed,
	0x49, 0x42, 0x70, 0x43, 0xcf, 0x74, 0x7d, 0x3f, 0x90, 0xae, 0xf4, 0x02, 0x5f, 0x24, 0xd3, 0x03,
	0xfa, 0x33, 0x3e, 0x9c, 0x72, 0xff, 0x50, 0xfc, 0xe8, 0x4e, 0xa7, 0x7c, 0x6e, 0x06, 0x21, 0x31,
	0xf2, 0x6c, 0xfd, 0x1f, 0x05, 0x4a, 0x6c, 0x31, 0xe3, 0xd8, 0x84, 0x82, 0x37, 0x69, 0x2b, 0x1d,
	0xa5, 0x5b, 0x66, 0x05, 0x6f, 0x82, 0x1d, 0x68, 0x4c, 0xb8, 0x18, 0xcf, 0x3d, 0x7a, 0xb5, 0x5d,
	0xe8, 0x28, 0xdd, 0x3a, 0xcb, 0x42, 0xf8, 0x21, 0x54, 0xdc, 0x31, 0x0d, 0x8b, 0x1d, 0xa5, 0xdb,
	0xec, 0x37, 0x8d, 0xf0, 0xc6, 0xb0, 0x08, 0x71, 0xde, 0x84, 0x9c, 0x25, 0x53, 0xfc, 0x02, 0xb6,
	0x67, 0xae, 0x90, 0xf6, 0x6b, 0x3e, 0x5e, 0x48, 0x3e, 0x69, 0x97, 0x3a, 0x4a, 0xb7, 0xd1, 0xd7,
	0x8c, 0xd8, 0x85, 0x91, 0xda, 0x34, 0x9c, 0xd4, 0x26, 0x5b, 0xe3, 0xe3, 0x47, 0x50, 0x93, 0x73,
	0x2f, 0xf2, 0x21, 0xda, 0xe5, 0x4e, 0xb1, 0xdb, 0xe8, 0x37, 0xa2, 0x9d, 0x9c, 0x18, 0x63, 0xcb,
	0x21, 0x7e, 0x00, 0x55, 0xe9, 0xce, 0xa7, 0x5c, 0x8a, 0x76, 0x85, 0x78, 0x40, 0x3c, 0x82, 0x58,
	0x3a, 0xd2, 0x5f, 0x40, 0x25, 0x86, 0x72, 0x96, 0x75, 0x28, 0xc9, 0x37, 0x21, 0x27, 0xaf, 0x89,
	0x9d, 0x98, 0x49, 0x76, 0x68, 0x86, 0x08, 0x25, 0xfe, 0x3a, 0x9c, 0x93, 0xe5, 0x3a, 0xa3, 0x67,
	0xfd, 0x1b, 0xa8, 0x26, 0x87, 0xc9, 0x49, 0xbe, 0xbf, 0x26, 0xb9, 0x9b, 0x39, 0x77, 0x46, 0x53,
	0x83, 0x9a, 0xe0, 0x52, 0x7a, 0xfe, 0x54, 0x90, 0xee, 0x36, 0x5b, 0xae, 0x75, 0x13, 0x76, 0xa2,
	0xeb, 0x11, 0x8c, 0x8b, 0x30, 0xf0, 0x05, 0xc7, 0x77, 0xa1, 0x3c, 0x8f, 0x80, 0xb6, 0x42, 0x16,
	0x6b, 0x91, 0x64, 0xc4, 0x60, 0x31, 0xac, 0x1f, 0xc0, 0x36, 0x2d, 0x53, 0xfe, 0x13, 0x28, 0x45,
	0x03, 0x3a, 0x53, 0x96, 0x4e, 0xa8, 0x8e, 0xa0, 0x5e, 0x78, 0x42, 0x26, 0x5b, 0xfc, 0xb0, 0xe0,
	0x42, 0xea, 0x5d, 0x68, 0x9e, 0x71, 0x19, 0x8b, 0x10, 0x82, 0x8f, 0xa0, 0x12, 0xb1, 0x87, 0xa9,
	0xb3, 0x64, 0xa5, 0xff, 0xae, 0x40, 0xd3, 0x9a, 0x4c, 0xb2, 0xd4, 0x7b, 0xb1, 0x51, 0xde, 0x16,
	0x9b, 0xc2, 0x5b, 0x63, 0x93, 0xbd, 0xf6, 0xe2, 0x86, 0xd7, 0x5e, 0x7a, 0xf8, 0xda, 0xff, 0x50,
	0xa0, 0x75, 0x15, 0x4e, 0x5c, 0xc9, 0x37, 0x70, 0xf6, 0x3f, 0xa6, 0x3f, 0x6b, 0xa3, 0xb4, 0xa1,
	0x8d, 0xf2, 0xc3, 0x36, 0x3e, 0x86, 0xd6, 0x29, 0x9f, 0xf1, 0x8d, 0x5c, 0xe8, 0x07, 0x80, 0x59,
	0x72, 0x92, 0x88, 0x87, 0xd8, 0x7b, 0x80, 0xe7, 0xdc, 0x9d, 0xc9, 0xdb, 0xc1, 0x2d, 0x1f, 0x7f,
	0x9f, 0xa6, 0xc1, 0x82, 0x77, 0xd6, 0xd0, 0x44, 0x04, 0xa1, 0x34, 0x08, 0x26, 0x71, 0xac, 0x8a,
	0x8c, 0x9e, 0x23, 0xe1, 0xaf, 0xa5, 0x2b, 0x17, 0x22, 0xf9, 0x5e, 0xc9, 0xaa, 0xf7, 0x29, 0xc0,
	0xea, 0xc3, 0xe0, 0x1e, 0xa8, 0x57, 0xa3, 0x53, 0xfb, 0xf9, 0x70, 0x64, 0x9f, 0x5e, 0x5b, 0x03,
	0x67, 0x78, 0x39, 0x52, 0xb7, 0x50, 0x85, 0xed, 0x2f, 0xed, 0x57, 0xd7, 0xec, 0xd2, 0xb1, 0x08,
	0x51, 0x7a, 0x07, 0x00, 0xab, 0xf6, 0x61, 0x15, 0x8a, 0xd6, 0xe8, 0x95, 0xba, 0x85, 0x75, 0x28,
	0x3b, 0x97, 0x2f, 0x86, 0x03, 0x55, 0x41, 0x80, 0xca, 0xe0, 0x62, 0x68, 0x8f, 0x1c, 0xb5, 0xd0,
	0x3b, 0x81, 0x46, 0xa6, 0x58, 0xb8, 0x0f, 0xad, 0xd5, 0x26, 0x0e, 0x1b, 0x9e, 0x9d, 0xd9, 0x4c,
	0xdd, 0xc2, 0x16, 0xec, 0x38, 0xc3, 0xaf, 0xec, 0xeb, 0xe1, 0xc8, 0xb1, 0xd9, 0x4b, 0xeb, 0x42,
	0x55, 0x22, 0x3d, 0xfb, 0x25, 0x69, 0xf4, 0xff, 0x2a, 0x02, 0x0e, 0xfa, 0xd6, 0x42, 0x06, 0x77,
	0xf4, 0x3f, 0xd2, 0xf6, 0xa7, 0x9e, 0xcf, 0xf1, 0x14, 0xea, 0xcb, 0x8e, 0xe0, 0x5e, 0x74, 0x29,
	0xf7, 0x2b, 0xa3, 0xb5, 0xd2, 0x5a, 0x2d, 0x7b, 0xaa, 0x37, 0x7f, 0xf9, 0xf3, 0xef, 0xdf, 0x0a,
	0x35, 0xac, 0x98, 0xd4, 0x4b, 0x3c, 0x87, 0x6a, 0xd2, 0x2a, 0xc4, 0x88, 0xbd, 0x5e, 0x31, 0x4d,
	0x5d, 0x16, 0x33, 0x15, 0x78, 0x4c, 0x02, 0x2d, 0xdc, 0x8d, 0x05, 0xcc, 0x9f, 0xe2, 0x6b, 0xfa,
	0x19, 0x4f, 0xa0, 0x9a, 0x94, 0x2e, 0x56, 0x5a, 0x6f, 0xe0, 0x7f, 0x28, 0xb5, 0x48, 0xa9, 0xa1,
	0x27, 0x47, 0x39, 0x56, 0x7a, 0x78, 0x0e, 0xb0, 0x2a, 0x03, 0xee, 0x47, 0xaf, 0xe4, 0xca, 0xf1,
	0xb0, 0x92, 0x96, 0x51, 0x72, 0x00, 0x56, 0x19, 0x8b, 0x95, 0x72, 0x01, 0xd5, 0x1e, 0xdd, 0x87,
	0xd7, 0x3d, 0xf6, 0x72, 0x1e, 0xaf, 0xa0, 0x91, 0x49, 0x1d, 0xd2, 0xfb, 0xf9, 0x70, 0x6a, 0x8f,
	0x73, 0x78, 0x22, 0xbc, 0x4f, 0xc2, 0xbb, 0xb8, 0x63, 0xde, 0xd2, 0xf4, 0x70, 0x1c, 0x8d, 0x4f,
	0x9e, 0xff, 0x6a, 0x0d, 0x10, 0xa0, 0x36, 0xee, 0xbb, 0xfc, 0xd0, 0x0d, 0x3d, 0xad, 0xf9, 0xac,
	0xff, 0xb9, 0xf1, 0xd4, 0x78, 0x6a, 0x3c, 0x3b, 0x3e, 0x3a, 0x3a, 0xfa, 0xac, 0xa7, 0x14, 0xfa,
	0xaa, 0x1b, 0x86, 0x33, 0x6f, 0x4c, 0x09, 0x30, 0xbf, 0x13, 0x81, 0x7f, 0x9c, 0x43, 0x6e, 0x2a,
	0xf4, 0xa3, 0xf5, 0xc9, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x50, 0x53, 0x58, 0xba, 0x07,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// C2AutomationEngineClient is the client API for C2AutomationEngine service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type C2AutomationEngineClient interface {
	// Retrieve list of existing rules
	ListRules(ctx context.Context, in *ListRulesRequest, opts ...grpc.CallOption) (*RulesResponse, error)
	// Retrieve a single rule, by its ID
	GetRule(ctx context.Context, in *GetRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error)
	// Create a new rule
	AddRule(ctx context.Context, in *AddRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error)
	// Update an existing rule
	UpdateRule(ctx context.Context, in *UpdateRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error)
	// Remove a rule
	DeleteRule(ctx context.Context, in *DeleteRuleRequest, opts ...grpc.CallOption) (*DeleteRuleResponse, error)
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type c2AutomationEngineClient struct {
	cc grpc.ClientConnInterface
}

func NewC2AutomationEngineClient(cc grpc.ClientConnInterface) C2AutomationEngineClient {
	return &c2AutomationEngineClient{cc}
}

func (c *c2AutomationEngineClient) ListRules(ctx context.Context, in *ListRulesRequest, opts ...grpc.CallOption) (*RulesResponse, error) {
	out := new(RulesResponse)
	err := c.cc.Invoke(ctx, "/pb.C2AutomationEngine/ListRules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *c2AutomationEngineClient) GetRule(ctx context.Context, in *GetRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error) {
	out := new(RuleResponse)
	err := c.cc.Invoke(ctx, "/pb.C2AutomationEngine/GetRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *c2AutomationEngineClient) AddRule(ctx context.Context, in *AddRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error) {
	out := new(RuleResponse)
	err := c.cc.Invoke(ctx, "/pb.C2AutomationEngine/AddRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *c2AutomationEngineClient) UpdateRule(ctx context.Context, in *UpdateRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error) {
	out := new(RuleResponse)
	err := c.cc.Invoke(ctx, "/pb.C2AutomationEngine/UpdateRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *c2AutomationEngineClient) DeleteRule(ctx context.Context, in *DeleteRuleRequest, opts ...grpc.CallOption) (*DeleteRuleResponse, error) {
	out := new(DeleteRuleResponse)
	err := c.cc.Invoke(ctx, "/pb.C2AutomationEngine/DeleteRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *c2AutomationEngineClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/pb.C2AutomationEngine/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// C2AutomationEngineServer is the server API for C2AutomationEngine service.
type C2AutomationEngineServer interface {
	// Retrieve list of existing rules
	ListRules(context.Context, *ListRulesRequest) (*RulesResponse, error)
	// Retrieve a single rule, by its ID
	GetRule(context.Context, *GetRuleRequest) (*RuleResponse, error)
	// Create a new rule
	AddRule(context.Context, *AddRuleRequest) (*RuleResponse, error)
	// Update an existing rule
	UpdateRule(context.Context, *UpdateRuleRequest) (*RuleResponse, error)
	// Remove a rule
	DeleteRule(context.Context, *DeleteRuleRequest) (*DeleteRuleResponse, error)
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

// UnimplementedC2AutomationEngineServer can be embedded to have forward compatible implementations.
type UnimplementedC2AutomationEngineServer struct {
}

func (*UnimplementedC2AutomationEngineServer) ListRules(ctx context.Context, req *ListRulesRequest) (*RulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRules not implemented")
}
func (*UnimplementedC2AutomationEngineServer) GetRule(ctx context.Context, req *GetRuleRequest) (*RuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRule not implemented")
}
func (*UnimplementedC2AutomationEngineServer) AddRule(ctx context.Context, req *AddRuleRequest) (*RuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRule not implemented")
}
func (*UnimplementedC2AutomationEngineServer) UpdateRule(ctx context.Context, req *UpdateRuleRequest) (*RuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRule not implemented")
}
func (*UnimplementedC2AutomationEngineServer) DeleteRule(ctx context.Context, req *DeleteRuleRequest) (*DeleteRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRule not implemented")
}
func (*UnimplementedC2AutomationEngineServer) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}

func RegisterC2AutomationEngineServer(s *grpc.Server, srv C2AutomationEngineServer) {
	s.RegisterService(&_C2AutomationEngine_serviceDesc, srv)
}

func _C2AutomationEngine_ListRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(C2AutomationEngineServer).ListRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.C2AutomationEngine/ListRules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(C2AutomationEngineServer).ListRules(ctx, req.(*ListRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _C2AutomationEngine_GetRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(C2AutomationEngineServer).GetRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.C2AutomationEngine/GetRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(C2AutomationEngineServer).GetRule(ctx, req.(*GetRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _C2AutomationEngine_AddRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(C2AutomationEngineServer).AddRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.C2AutomationEngine/AddRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(C2AutomationEngineServer).AddRule(ctx, req.(*AddRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _C2AutomationEngine_UpdateRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(C2AutomationEngineServer).UpdateRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.C2AutomationEngine/UpdateRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(C2AutomationEngineServer).UpdateRule(ctx, req.(*UpdateRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _C2AutomationEngine_DeleteRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(C2AutomationEngineServer).DeleteRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.C2AutomationEngine/DeleteRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(C2AutomationEngineServer).DeleteRule(ctx, req.(*DeleteRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _C2AutomationEngine_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(C2AutomationEngineServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.C2AutomationEngine/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(C2AutomationEngineServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _C2AutomationEngine_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.C2AutomationEngine",
	HandlerType: (*C2AutomationEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRules",
			Handler:    _C2AutomationEngine_ListRules_Handler,
		},
		{
			MethodName: "GetRule",
			Handler:    _C2AutomationEngine_GetRule_Handler,
		},
		{
			MethodName: "AddRule",
			Handler:    _C2AutomationEngine_AddRule_Handler,
		},
		{
			MethodName: "UpdateRule",
			Handler:    _C2AutomationEngine_UpdateRule_Handler,
		},
		{
			MethodName: "DeleteRule",
			Handler:    _C2AutomationEngine_DeleteRule_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _C2AutomationEngine_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
