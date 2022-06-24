// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: internal/testing/storage.proto

package testing

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	schema "namespacelabs.dev/foundation/schema"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PreStoredTestBundle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result    *schema.TestResult `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	TestLog   *Log               `protobuf:"bytes,2,opt,name=test_log,json=testLog,proto3" json:"test_log,omitempty"`
	ServerLog []*Log             `protobuf:"bytes,3,rep,name=server_log,json=serverLog,proto3" json:"server_log,omitempty"`
}

func (x *PreStoredTestBundle) Reset() {
	*x = PreStoredTestBundle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_testing_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreStoredTestBundle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreStoredTestBundle) ProtoMessage() {}

func (x *PreStoredTestBundle) ProtoReflect() protoreflect.Message {
	mi := &file_internal_testing_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreStoredTestBundle.ProtoReflect.Descriptor instead.
func (*PreStoredTestBundle) Descriptor() ([]byte, []int) {
	return file_internal_testing_storage_proto_rawDescGZIP(), []int{0}
}

func (x *PreStoredTestBundle) GetResult() *schema.TestResult {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *PreStoredTestBundle) GetTestLog() *Log {
	if x != nil {
		return x.TestLog
	}
	return nil
}

func (x *PreStoredTestBundle) GetServerLog() []*Log {
	if x != nil {
		return x.ServerLog
	}
	return nil
}

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PackageName   string               `protobuf:"bytes,1,opt,name=package_name,json=packageName,proto3" json:"package_name,omitempty"`
	ContainerName string               `protobuf:"bytes,3,opt,name=container_name,json=containerName,proto3" json:"container_name,omitempty"`
	ContainerKind schema.ContainerKind `protobuf:"varint,4,opt,name=container_kind,json=containerKind,proto3,enum=foundation.schema.ContainerKind" json:"container_kind,omitempty"`
	Output        []byte               `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_testing_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_internal_testing_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_internal_testing_storage_proto_rawDescGZIP(), []int{1}
}

func (x *Log) GetPackageName() string {
	if x != nil {
		return x.PackageName
	}
	return ""
}

func (x *Log) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

func (x *Log) GetContainerKind() schema.ContainerKind {
	if x != nil {
		return x.ContainerKind
	}
	return schema.ContainerKind(0)
}

func (x *Log) GetOutput() []byte {
	if x != nil {
		return x.Output
	}
	return nil
}

var File_internal_testing_storage_proto protoreflect.FileDescriptor

var file_internal_testing_storage_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x67, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x1b, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x14, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xca,
	0x01, 0x0a, 0x13, 0x50, 0x72, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74,
	0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x3b, 0x0a,
	0x08, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4c, 0x6f,
	0x67, 0x52, 0x07, 0x74, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x12, 0x3f, 0x0a, 0x0a, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x67, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4c, 0x6f, 0x67,
	0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x22, 0xb0, 0x01, 0x0a, 0x03,
	0x4c, 0x6f, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x47, 0x0a,
	0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x5f, 0x6b, 0x69, 0x6e, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x42, 0x2f,
	0x5a, 0x2d, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x6c, 0x61, 0x62, 0x73, 0x2e,
	0x64, 0x65, 0x76, 0x2f, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_testing_storage_proto_rawDescOnce sync.Once
	file_internal_testing_storage_proto_rawDescData = file_internal_testing_storage_proto_rawDesc
)

func file_internal_testing_storage_proto_rawDescGZIP() []byte {
	file_internal_testing_storage_proto_rawDescOnce.Do(func() {
		file_internal_testing_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_testing_storage_proto_rawDescData)
	})
	return file_internal_testing_storage_proto_rawDescData
}

var file_internal_testing_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_testing_storage_proto_goTypes = []interface{}{
	(*PreStoredTestBundle)(nil), // 0: foundation.internal.testing.PreStoredTestBundle
	(*Log)(nil),                 // 1: foundation.internal.testing.Log
	(*schema.TestResult)(nil),   // 2: foundation.schema.TestResult
	(schema.ContainerKind)(0),   // 3: foundation.schema.ContainerKind
}
var file_internal_testing_storage_proto_depIdxs = []int32{
	2, // 0: foundation.internal.testing.PreStoredTestBundle.result:type_name -> foundation.schema.TestResult
	1, // 1: foundation.internal.testing.PreStoredTestBundle.test_log:type_name -> foundation.internal.testing.Log
	1, // 2: foundation.internal.testing.PreStoredTestBundle.server_log:type_name -> foundation.internal.testing.Log
	3, // 3: foundation.internal.testing.Log.container_kind:type_name -> foundation.schema.ContainerKind
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_internal_testing_storage_proto_init() }
func file_internal_testing_storage_proto_init() {
	if File_internal_testing_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_testing_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreStoredTestBundle); i {
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
		file_internal_testing_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
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
			RawDescriptor: file_internal_testing_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_testing_storage_proto_goTypes,
		DependencyIndexes: file_internal_testing_storage_proto_depIdxs,
		MessageInfos:      file_internal_testing_storage_proto_msgTypes,
	}.Build()
	File_internal_testing_storage_proto = out.File
	file_internal_testing_storage_proto_rawDesc = nil
	file_internal_testing_storage_proto_goTypes = nil
	file_internal_testing_storage_proto_depIdxs = nil
}
