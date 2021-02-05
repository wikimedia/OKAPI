// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: protos/projects.proto

package protos

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Index io description
type IndexRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IndexRequest) Reset() {
	*x = IndexRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_projects_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexRequest) ProtoMessage() {}

func (x *IndexRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_projects_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexRequest.ProtoReflect.Descriptor instead.
func (*IndexRequest) Descriptor() ([]byte, []int) {
	return file_protos_projects_proto_rawDescGZIP(), []int{0}
}

type IndexResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total  int32 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Errors int32 `protobuf:"varint,2,opt,name=errors,proto3" json:"errors,omitempty"`
}

func (x *IndexResponse) Reset() {
	*x = IndexResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_projects_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexResponse) ProtoMessage() {}

func (x *IndexResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_projects_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexResponse.ProtoReflect.Descriptor instead.
func (*IndexResponse) Descriptor() ([]byte, []int) {
	return file_protos_projects_proto_rawDescGZIP(), []int{1}
}

func (x *IndexResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *IndexResponse) GetErrors() int32 {
	if x != nil {
		return x.Errors
	}
	return 0
}

// Projects io description
type FetchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FetchRequest) Reset() {
	*x = FetchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_projects_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchRequest) ProtoMessage() {}

func (x *FetchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_projects_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchRequest.ProtoReflect.Descriptor instead.
func (*FetchRequest) Descriptor() ([]byte, []int) {
	return file_protos_projects_proto_rawDescGZIP(), []int{2}
}

type FetchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FetchResponse) Reset() {
	*x = FetchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_projects_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchResponse) ProtoMessage() {}

func (x *FetchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_projects_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchResponse.ProtoReflect.Descriptor instead.
func (*FetchResponse) Descriptor() ([]byte, []int) {
	return file_protos_projects_proto_rawDescGZIP(), []int{3}
}

var File_protos_projects_proto protoreflect.FileDescriptor

var file_protos_projects_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x73, 0x22, 0x0e, 0x0a, 0x0c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x3d, 0x0a, 0x0d, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x22, 0x0e, 0x0a, 0x0c, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x0f, 0x0a, 0x0d, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x7e, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x38, 0x0a,
	0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x73, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x05, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x1b, 0x5a, 0x19, 0x6f, 0x6b, 0x61, 0x70, 0x69, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_projects_proto_rawDescOnce sync.Once
	file_protos_projects_proto_rawDescData = file_protos_projects_proto_rawDesc
)

func file_protos_projects_proto_rawDescGZIP() []byte {
	file_protos_projects_proto_rawDescOnce.Do(func() {
		file_protos_projects_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_projects_proto_rawDescData)
	})
	return file_protos_projects_proto_rawDescData
}

var file_protos_projects_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protos_projects_proto_goTypes = []interface{}{
	(*IndexRequest)(nil),  // 0: projects.IndexRequest
	(*IndexResponse)(nil), // 1: projects.IndexResponse
	(*FetchRequest)(nil),  // 2: projects.FetchRequest
	(*FetchResponse)(nil), // 3: projects.FetchResponse
}
var file_protos_projects_proto_depIdxs = []int32{
	0, // 0: projects.Projects.Index:input_type -> projects.IndexRequest
	2, // 1: projects.Projects.Fetch:input_type -> projects.FetchRequest
	1, // 2: projects.Projects.Index:output_type -> projects.IndexResponse
	3, // 3: projects.Projects.Fetch:output_type -> projects.FetchResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_projects_proto_init() }
func file_protos_projects_proto_init() {
	if File_protos_projects_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_projects_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexRequest); i {
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
		file_protos_projects_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexResponse); i {
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
		file_protos_projects_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchRequest); i {
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
		file_protos_projects_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchResponse); i {
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
			RawDescriptor: file_protos_projects_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_projects_proto_goTypes,
		DependencyIndexes: file_protos_projects_proto_depIdxs,
		MessageInfos:      file_protos_projects_proto_msgTypes,
	}.Build()
	File_protos_projects_proto = out.File
	file_protos_projects_proto_rawDesc = nil
	file_protos_projects_proto_goTypes = nil
	file_protos_projects_proto_depIdxs = nil
}
