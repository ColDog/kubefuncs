// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package kubefuncs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Event struct {
	// ID will be a UUID configured that should be unique to every message. It is
	// generated by the NewEvent() function from each client library.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Topic is the topic for the given event. It should be provided by the user
	// this represents the function to call.
	Topic string `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	// Return is a return topic, this will be set by the client libraries Call()
	// method.
	Return  string               `protobuf:"bytes,3,opt,name=return" json:"return,omitempty"`
	Payload *google_protobuf.Any `protobuf:"bytes,4,opt,name=payload" json:"payload,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *Event) GetReturn() string {
	if m != nil {
		return m.Return
	}
	return ""
}

func (m *Event) GetPayload() *google_protobuf.Any {
	if m != nil {
		return m.Payload
	}
	return nil
}

type HTTPRequest struct {
	Url     string            `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Headers map[string]string `protobuf:"bytes,2,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Body    []byte            `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *HTTPRequest) Reset()                    { *m = HTTPRequest{} }
func (m *HTTPRequest) String() string            { return proto.CompactTextString(m) }
func (*HTTPRequest) ProtoMessage()               {}
func (*HTTPRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HTTPRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *HTTPRequest) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *HTTPRequest) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type HTTPResponse struct {
	Status  int32             `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Headers map[string]string `protobuf:"bytes,2,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Body    []byte            `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *HTTPResponse) Reset()                    { *m = HTTPResponse{} }
func (m *HTTPResponse) String() string            { return proto.CompactTextString(m) }
func (*HTTPResponse) ProtoMessage()               {}
func (*HTTPResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *HTTPResponse) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *HTTPResponse) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *HTTPResponse) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type Error struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Error) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*Event)(nil), "kubefuncs.Event")
	proto.RegisterType((*HTTPRequest)(nil), "kubefuncs.HTTPRequest")
	proto.RegisterType((*HTTPResponse)(nil), "kubefuncs.HTTPResponse")
	proto.RegisterType((*Empty)(nil), "kubefuncs.Empty")
	proto.RegisterType((*Error)(nil), "kubefuncs.Error")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x51, 0x4f, 0x6b, 0xfa, 0x40,
	0x14, 0x24, 0xd1, 0x28, 0x3e, 0xfd, 0xfd, 0x28, 0x8b, 0x94, 0x54, 0x28, 0x48, 0xda, 0x83, 0xa7,
	0x15, 0xec, 0xa5, 0x08, 0x2d, 0xf4, 0x10, 0xf0, 0x58, 0x82, 0x5f, 0x60, 0x63, 0x9e, 0x56, 0x12,
	0x77, 0xd3, 0xfd, 0x23, 0xec, 0x97, 0xea, 0xb1, 0x9f, 0xaf, 0x64, 0x37, 0x29, 0xd2, 0xde, 0x7b,
	0x7b, 0xf3, 0x76, 0x66, 0x67, 0x66, 0x17, 0xfe, 0x9d, 0x50, 0x29, 0x76, 0x40, 0x5a, 0x4b, 0xa1,
	0x05, 0x19, 0x95, 0x26, 0xc7, 0xbd, 0xe1, 0x3b, 0x35, 0xbb, 0x39, 0x08, 0x71, 0xa8, 0x70, 0xe9,
	0x0e, 0x72, 0xb3, 0x5f, 0x32, 0x6e, 0x3d, 0x2b, 0x31, 0x10, 0xa5, 0x67, 0xe4, 0x9a, 0xfc, 0x87,
	0xf0, 0x58, 0xc4, 0xc1, 0x3c, 0x58, 0x8c, 0xb2, 0xf0, 0x58, 0x90, 0x29, 0x44, 0x5a, 0xd4, 0xc7,
	0x5d, 0x1c, 0xba, 0x95, 0x07, 0xe4, 0x1a, 0x06, 0x12, 0xb5, 0x91, 0x3c, 0xee, 0xb9, 0x75, 0x8b,
	0x08, 0x85, 0x61, 0xcd, 0x6c, 0x25, 0x58, 0x11, 0xf7, 0xe7, 0xc1, 0x62, 0xbc, 0x9a, 0x52, 0xef,
	0x49, 0x3b, 0x4f, 0xfa, 0xc2, 0x6d, 0xd6, 0x91, 0x92, 0x8f, 0x00, 0xc6, 0x9b, 0xed, 0xf6, 0x35,
	0xc3, 0x77, 0x83, 0x4a, 0x93, 0x2b, 0xe8, 0x19, 0x59, 0xb5, 0xf6, 0xcd, 0x48, 0x9e, 0x60, 0xf8,
	0x86, 0xac, 0x40, 0xa9, 0xe2, 0x70, 0xde, 0x5b, 0x8c, 0x57, 0x77, 0xf4, 0xbb, 0x10, 0xbd, 0x90,
	0xd2, 0x8d, 0x67, 0xa5, 0x5c, 0x4b, 0x9b, 0x75, 0x1a, 0x42, 0xa0, 0x9f, 0x8b, 0xc2, 0xba, 0x98,
	0x93, 0xcc, 0xcd, 0xb3, 0x35, 0x4c, 0x2e, 0xc9, 0x8d, 0x69, 0x89, 0xb6, 0x33, 0x2d, 0xd1, 0x36,
	0xa5, 0xcf, 0xac, 0x32, 0xd8, 0x95, 0x76, 0x60, 0x1d, 0x3e, 0x06, 0xc9, 0x67, 0x00, 0x13, 0xef,
	0xaa, 0x6a, 0xc1, 0x15, 0x36, 0x2f, 0xa1, 0x34, 0xd3, 0x46, 0x39, 0x7d, 0x94, 0xb5, 0x88, 0x3c,
	0xff, 0xcc, 0x7d, 0xff, 0x2b, 0xb7, 0xbf, 0xe1, 0x8f, 0x82, 0x0f, 0x21, 0x4a, 0x4f, 0xb5, 0xb6,
	0xc9, 0x2d, 0x44, 0xa9, 0x94, 0x42, 0x36, 0x5c, 0x6c, 0x86, 0x56, 0xef, 0x41, 0x3e, 0x70, 0x1f,
	0xf5, 0xf0, 0x15, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x52, 0x5c, 0xb0, 0x46, 0x02, 0x00, 0x00,
}
