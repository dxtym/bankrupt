// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.12.4
// source: service_bankrupt.proto

package pb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_bankrupt_proto protoreflect.FileDescriptor

var file_service_bankrupt_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x62, 0x61, 0x6e, 0x6b, 0x72, 0x75,
	0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32,
	0xc1, 0x02, 0x0a, 0x08, 0x42, 0x61, 0x6e, 0x6b, 0x72, 0x75, 0x70, 0x74, 0x12, 0x88, 0x01, 0x0a,
	0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x70, 0x62,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4b, 0x92, 0x41, 0x2e, 0x12,
	0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x75, 0x73, 0x65, 0x72,
	0x1a, 0x1b, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x20, 0x74, 0x6f, 0x20, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x75, 0x73, 0x65, 0x72, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x14, 0x3a, 0x01, 0x2a, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x12, 0xa9, 0x01, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x62,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x6f, 0x92, 0x41, 0x53, 0x12, 0x19, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x20, 0x75,
	0x73, 0x65, 0x72, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x67, 0x65, 0x74, 0x20, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x73, 0x1a, 0x36, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x20, 0x74, 0x6f, 0x20,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x72, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x67,
	0x65, 0x74, 0x20, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20, 0x26, 0x20, 0x72, 0x65, 0x66, 0x65,
	0x72, 0x73, 0x68, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13,
	0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x42, 0x95, 0x01, 0x92, 0x41, 0x74, 0x12, 0x72, 0x0a, 0x0c, 0x42, 0x61, 0x6e,
	0x6b, 0x72, 0x75, 0x70, 0x74, 0x20, 0x41, 0x50, 0x49, 0x22, 0x5d, 0x0a, 0x14, 0x44, 0x69, 0x6c,
	0x6d, 0x75, 0x72, 0x6f, 0x64, 0x20, 0x41, 0x62, 0x64, 0x75, 0x73, 0x61, 0x6d, 0x61, 0x64, 0x6f,
	0x76, 0x12, 0x21, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x78, 0x74, 0x79, 0x6d, 0x2f, 0x62, 0x61, 0x6e, 0x6b,
	0x72, 0x75, 0x70, 0x74, 0x1a, 0x22, 0x64, 0x69, 0x6c, 0x6d, 0x75, 0x72, 0x6f, 0x64, 0x2e, 0x61,
	0x62, 0x64, 0x75, 0x73, 0x61, 0x6d, 0x61, 0x64, 0x6f, 0x76, 0x32, 0x30, 0x30, 0x34, 0x40, 0x67,
	0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6d, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x5a, 0x1c, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x78, 0x74, 0x79, 0x6d, 0x2f,
	0x62, 0x61, 0x6e, 0x6b, 0x72, 0x75, 0x70, 0x74, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var file_service_bankrupt_proto_goTypes = []interface{}{
	(*CreateUserRequest)(nil),  // 0: pb.CreateUserRequest
	(*LoginUserRequest)(nil),   // 1: pb.LoginUserRequest
	(*CreateUserResponse)(nil), // 2: pb.CreateUserResponse
	(*LoginUserResponse)(nil),  // 3: pb.LoginUserResponse
}
var file_service_bankrupt_proto_depIdxs = []int32{
	0, // 0: pb.Bankrupt.CreateUser:input_type -> pb.CreateUserRequest
	1, // 1: pb.Bankrupt.LoginUser:input_type -> pb.LoginUserRequest
	2, // 2: pb.Bankrupt.CreateUser:output_type -> pb.CreateUserResponse
	3, // 3: pb.Bankrupt.LoginUser:output_type -> pb.LoginUserResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_bankrupt_proto_init() }
func file_service_bankrupt_proto_init() {
	if File_service_bankrupt_proto != nil {
		return
	}
	file_create_user_proto_init()
	file_login_user_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_bankrupt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_bankrupt_proto_goTypes,
		DependencyIndexes: file_service_bankrupt_proto_depIdxs,
	}.Build()
	File_service_bankrupt_proto = out.File
	file_service_bankrupt_proto_rawDesc = nil
	file_service_bankrupt_proto_goTypes = nil
	file_service_bankrupt_proto_depIdxs = nil
}