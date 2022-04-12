//
//Copyright 2022 quarkcm Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.8
// source: pkg/grpc/quarkcmsvc.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TestRequestMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientName string `protobuf:"bytes,1,opt,name=client_name,json=clientName,proto3" json:"client_name,omitempty"`
}

func (x *TestRequestMessage) Reset() {
	*x = TestRequestMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_quarkcmsvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRequestMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequestMessage) ProtoMessage() {}

func (x *TestRequestMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_quarkcmsvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequestMessage.ProtoReflect.Descriptor instead.
func (*TestRequestMessage) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_quarkcmsvc_proto_rawDescGZIP(), []int{0}
}

func (x *TestRequestMessage) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

type TestResponseMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerName string `protobuf:"bytes,1,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`
}

func (x *TestResponseMessage) Reset() {
	*x = TestResponseMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_quarkcmsvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResponseMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponseMessage) ProtoMessage() {}

func (x *TestResponseMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_quarkcmsvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponseMessage.ProtoReflect.Descriptor instead.
func (*TestResponseMessage) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_quarkcmsvc_proto_rawDescGZIP(), []int{1}
}

func (x *TestResponseMessage) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

var File_pkg_grpc_quarkcmsvc_proto protoreflect.FileDescriptor

var file_pkg_grpc_quarkcmsvc_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x71, 0x75, 0x61, 0x72, 0x6b,
	0x63, 0x6d, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x71, 0x75, 0x61,
	0x72, 0x6b, 0x63, 0x6d, 0x73, 0x76, 0x63, 0x22, 0x35, 0x0a, 0x12, 0x54, 0x65, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x36,
	0x0a, 0x13, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0x5f, 0x0a, 0x0e, 0x51, 0x75, 0x61, 0x72, 0x6b, 0x43,
	0x4d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74,
	0x50, 0x69, 0x6e, 0x67, 0x12, 0x1e, 0x2e, 0x71, 0x75, 0x61, 0x72, 0x6b, 0x63, 0x6d, 0x73, 0x76,
	0x63, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x1f, 0x2e, 0x71, 0x75, 0x61, 0x72, 0x6b, 0x63, 0x6d, 0x73, 0x76,
	0x63, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x6b, 0x67, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_grpc_quarkcmsvc_proto_rawDescOnce sync.Once
	file_pkg_grpc_quarkcmsvc_proto_rawDescData = file_pkg_grpc_quarkcmsvc_proto_rawDesc
)

func file_pkg_grpc_quarkcmsvc_proto_rawDescGZIP() []byte {
	file_pkg_grpc_quarkcmsvc_proto_rawDescOnce.Do(func() {
		file_pkg_grpc_quarkcmsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_grpc_quarkcmsvc_proto_rawDescData)
	})
	return file_pkg_grpc_quarkcmsvc_proto_rawDescData
}

var file_pkg_grpc_quarkcmsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_grpc_quarkcmsvc_proto_goTypes = []interface{}{
	(*TestRequestMessage)(nil),  // 0: quarkcmsvc.TestRequestMessage
	(*TestResponseMessage)(nil), // 1: quarkcmsvc.TestResponseMessage
}
var file_pkg_grpc_quarkcmsvc_proto_depIdxs = []int32{
	0, // 0: quarkcmsvc.QuarkCMService.TestPing:input_type -> quarkcmsvc.TestRequestMessage
	1, // 1: quarkcmsvc.QuarkCMService.TestPing:output_type -> quarkcmsvc.TestResponseMessage
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_grpc_quarkcmsvc_proto_init() }
func file_pkg_grpc_quarkcmsvc_proto_init() {
	if File_pkg_grpc_quarkcmsvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_grpc_quarkcmsvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRequestMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_grpc_quarkcmsvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResponseMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_grpc_quarkcmsvc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_grpc_quarkcmsvc_proto_goTypes,
		DependencyIndexes: file_pkg_grpc_quarkcmsvc_proto_depIdxs,
		MessageInfos:      file_pkg_grpc_quarkcmsvc_proto_msgTypes,
	}.Build()
	File_pkg_grpc_quarkcmsvc_proto = out.File
	file_pkg_grpc_quarkcmsvc_proto_rawDesc = nil
	file_pkg_grpc_quarkcmsvc_proto_goTypes = nil
	file_pkg_grpc_quarkcmsvc_proto_depIdxs = nil
}
