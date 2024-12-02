// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: grpc/games.proto

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

type Game struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ConstructorName string `protobuf:"bytes,3,opt,name=constructor_name,json=constructorName,proto3" json:"constructor_name,omitempty"`
}

func (x *Game) Reset() {
	*x = Game{}
	mi := &file_grpc_games_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Game) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Game) ProtoMessage() {}

func (x *Game) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_games_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Game.ProtoReflect.Descriptor instead.
func (*Game) Descriptor() ([]byte, []int) {
	return file_grpc_games_proto_rawDescGZIP(), []int{0}
}

func (x *Game) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Game) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Game) GetConstructorName() string {
	if x != nil {
		return x.ConstructorName
	}
	return ""
}

type UploadGameFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameName        string `protobuf:"bytes,1,opt,name=game_name,json=gameName,proto3" json:"game_name,omitempty"`
	FileContent     []byte `protobuf:"bytes,2,opt,name=file_content,json=fileContent,proto3" json:"file_content,omitempty"`
	ConstructorName string `protobuf:"bytes,3,opt,name=constructor_name,json=constructorName,proto3" json:"constructor_name,omitempty"`
}

func (x *UploadGameFileRequest) Reset() {
	*x = UploadGameFileRequest{}
	mi := &file_grpc_games_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadGameFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadGameFileRequest) ProtoMessage() {}

func (x *UploadGameFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_games_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadGameFileRequest.ProtoReflect.Descriptor instead.
func (*UploadGameFileRequest) Descriptor() ([]byte, []int) {
	return file_grpc_games_proto_rawDescGZIP(), []int{1}
}

func (x *UploadGameFileRequest) GetGameName() string {
	if x != nil {
		return x.GameName
	}
	return ""
}

func (x *UploadGameFileRequest) GetFileContent() []byte {
	if x != nil {
		return x.FileContent
	}
	return nil
}

func (x *UploadGameFileRequest) GetConstructorName() string {
	if x != nil {
		return x.ConstructorName
	}
	return ""
}

type UploadGameFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UploadGameFileResponse) Reset() {
	*x = UploadGameFileResponse{}
	mi := &file_grpc_games_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadGameFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadGameFileResponse) ProtoMessage() {}

func (x *UploadGameFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_games_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadGameFileResponse.ProtoReflect.Descriptor instead.
func (*UploadGameFileResponse) Descriptor() ([]byte, []int) {
	return file_grpc_games_proto_rawDescGZIP(), []int{2}
}

func (x *UploadGameFileResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_grpc_games_proto protoreflect.FileDescriptor

var file_grpc_games_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x55, 0x0a, 0x04, 0x47, 0x61, 0x6d,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x6f, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x82, 0x01, 0x0a, 0x15, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x47, 0x61, 0x6d, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61,
	0x6d, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67,
	0x61, 0x6d, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x66,
	0x69, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6f,
	0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x6f,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x32, 0x0a, 0x16, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x47,
	0x61, 0x6d, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x5f, 0x0a, 0x0c, 0x47, 0x61, 0x6d,
	0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x47, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x67, 0x61,
	0x6d, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x47, 0x61, 0x6d, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x73, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x47, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x12, 0x5a, 0x10, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x3b, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_games_proto_rawDescOnce sync.Once
	file_grpc_games_proto_rawDescData = file_grpc_games_proto_rawDesc
)

func file_grpc_games_proto_rawDescGZIP() []byte {
	file_grpc_games_proto_rawDescOnce.Do(func() {
		file_grpc_games_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_games_proto_rawDescData)
	})
	return file_grpc_games_proto_rawDescData
}

var file_grpc_games_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_grpc_games_proto_goTypes = []any{
	(*Game)(nil),                   // 0: games.Game
	(*UploadGameFileRequest)(nil),  // 1: games.UploadGameFileRequest
	(*UploadGameFileResponse)(nil), // 2: games.UploadGameFileResponse
}
var file_grpc_games_proto_depIdxs = []int32{
	1, // 0: games.GamesService.UploadGameFile:input_type -> games.UploadGameFileRequest
	2, // 1: games.GamesService.UploadGameFile:output_type -> games.UploadGameFileResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_games_proto_init() }
func file_grpc_games_proto_init() {
	if File_grpc_games_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_games_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_games_proto_goTypes,
		DependencyIndexes: file_grpc_games_proto_depIdxs,
		MessageInfos:      file_grpc_games_proto_msgTypes,
	}.Build()
	File_grpc_games_proto = out.File
	file_grpc_games_proto_rawDesc = nil
	file_grpc_games_proto_goTypes = nil
	file_grpc_games_proto_depIdxs = nil
}
