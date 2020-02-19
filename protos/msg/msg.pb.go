// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.protos

package Tmsg

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the protos package it is being compiled against.
// A compilation error at this line likely means your copy of the
// protos package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the protos package

//测试请求接口
type ReqTestMessage struct {
	Uid                  uint32   `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqTestMessage) Reset()         { *m = ReqTestMessage{} }
func (m *ReqTestMessage) String() string { return proto.CompactTextString(m) }
func (*ReqTestMessage) ProtoMessage()    {}
func (*ReqTestMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *ReqTestMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqTestMessage.Unmarshal(m, b)
}
func (m *ReqTestMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqTestMessage.Marshal(b, m, deterministic)
}
func (m *ReqTestMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqTestMessage.Merge(m, src)
}
func (m *ReqTestMessage) XXX_Size() int {
	return xxx_messageInfo_ReqTestMessage.Size(m)
}
func (m *ReqTestMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqTestMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ReqTestMessage proto.InternalMessageInfo

func (m *ReqTestMessage) GetUid() uint32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *ReqTestMessage) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

//测试请求回应接口
type AckTestMessage struct {
	Code                 uint32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AckTestMessage) Reset()         { *m = AckTestMessage{} }
func (m *AckTestMessage) String() string { return proto.CompactTextString(m) }
func (*AckTestMessage) ProtoMessage()    {}
func (*AckTestMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}

func (m *AckTestMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckTestMessage.Unmarshal(m, b)
}
func (m *AckTestMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckTestMessage.Marshal(b, m, deterministic)
}
func (m *AckTestMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckTestMessage.Merge(m, src)
}
func (m *AckTestMessage) XXX_Size() int {
	return xxx_messageInfo_AckTestMessage.Size(m)
}
func (m *AckTestMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_AckTestMessage.DiscardUnknown(m)
}

var xxx_messageInfo_AckTestMessage proto.InternalMessageInfo

func (m *AckTestMessage) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *AckTestMessage) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*ReqTestMessage)(nil), "Tmsg.ReqTestMessage")
	proto.RegisterType((*AckTestMessage)(nil), "Tmsg.AckTestMessage")
}

func init() { proto.RegisterFile("msg.protos", fileDescriptor_c06e4cca6c2cc899) }

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 116 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0x2d, 0x4e, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x09, 0xc9, 0x2d, 0x4e, 0x57, 0x32, 0xe1, 0xe2, 0x0b,
	0x4a, 0x2d, 0x0c, 0x49, 0x2d, 0x2e, 0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x12, 0xe0,
	0x62, 0x2e, 0xcd, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0x02, 0x31, 0x41, 0x22, 0xb9,
	0xc5, 0xe9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x20, 0xa6, 0x92, 0x19, 0x17, 0x9f, 0x63,
	0x72, 0x36, 0xb2, 0x2e, 0x21, 0x2e, 0x96, 0xe4, 0xfc, 0x94, 0x54, 0xa8, 0x36, 0x30, 0x1b, 0x53,
	0x5f, 0x12, 0x1b, 0xd8, 0x6a, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc0, 0x2d, 0x8c, 0x78,
	0x87, 0x00, 0x00, 0x00,
}
