// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.15.2
// source: space/api/pb/v1/video.proto

package v1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GetVideoRequest 获取视频请求
type GetVideoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id 视频id
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetVideoRequest) Reset() {
	*x = GetVideoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_space_api_pb_v1_video_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVideoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVideoRequest) ProtoMessage() {}

func (x *GetVideoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_space_api_pb_v1_video_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVideoRequest.ProtoReflect.Descriptor instead.
func (*GetVideoRequest) Descriptor() ([]byte, []int) {
	return file_space_api_pb_v1_video_proto_rawDescGZIP(), []int{0}
}

func (x *GetVideoRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// Video 视频信息
type Video struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id 视频id
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// CreatedAt 创建时间
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt 修改时间
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// DeletedAt 删除时间
	DeletedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}

func (x *Video) Reset() {
	*x = Video{}
	if protoimpl.UnsafeEnabled {
		mi := &file_space_api_pb_v1_video_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Video) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Video) ProtoMessage() {}

func (x *Video) ProtoReflect() protoreflect.Message {
	mi := &file_space_api_pb_v1_video_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Video.ProtoReflect.Descriptor instead.
func (*Video) Descriptor() ([]byte, []int) {
	return file_space_api_pb_v1_video_proto_rawDescGZIP(), []int{1}
}

func (x *Video) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Video) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Video) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Video) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

// TestResponse 测试响应
type TestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id 视频id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Video 视频信息
	Video *Video `protobuf:"bytes,2,opt,name=video,proto3" json:"video,omitempty"`
}

func (x *TestResponse) Reset() {
	*x = TestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_space_api_pb_v1_video_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponse) ProtoMessage() {}

func (x *TestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_space_api_pb_v1_video_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponse.ProtoReflect.Descriptor instead.
func (*TestResponse) Descriptor() ([]byte, []int) {
	return file_space_api_pb_v1_video_proto_rawDescGZIP(), []int{2}
}

func (x *TestResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TestResponse) GetVideo() *Video {
	if x != nil {
		return x.Video
	}
	return nil
}

var File_space_api_pb_v1_video_proto protoreflect.FileDescriptor

var file_space_api_pb_v1_video_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x2f, 0x76,
	0x31, 0x2f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x73, 0x6d, 0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2e, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x21, 0x0a, 0x0f, 0x47, 0x65, 0x74,
	0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0xc8, 0x01, 0x0a,
	0x05, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x5b, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3b, 0x0a, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x73, 0x6d, 0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2e, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x05, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x32, 0xfa, 0x01, 0x0a, 0x0c, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7d, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x2f, 0x2e, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x73, 0x6d, 0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2e, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e,
	0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x73, 0x6d, 0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2e, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x56,
	0x69, 0x64, 0x65, 0x6f, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2f, 0x67, 0x65, 0x74, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x12, 0x6b, 0x0a, 0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x73, 0x6d,
	0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2e, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x22, 0x12, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x3a, 0x01,
	0x2a, 0x42, 0xb3, 0x01, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x73, 0x6d, 0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2e, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x2e, 0x73, 0x61, 0x62, 0x65,
	0x72, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x73, 0x6d, 0x63, 0x6d, 0x61, 0x6c, 0x6c, 0x2f, 0x61, 0x70,
	0x70, 0x2f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x2f, 0x76,
	0x31, 0x92, 0x41, 0x43, 0x12, 0x1a, 0x0a, 0x0c, 0xe8, 0xa7, 0x86, 0xe9, 0xa2, 0x91, 0xe6, 0x9c,
	0x8d, 0xe5, 0x8a, 0xa1, 0x2a, 0x05, 0x0a, 0x03, 0x4d, 0x49, 0x54, 0x32, 0x03, 0x31, 0x2e, 0x30,
	0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_space_api_pb_v1_video_proto_rawDescOnce sync.Once
	file_space_api_pb_v1_video_proto_rawDescData = file_space_api_pb_v1_video_proto_rawDesc
)

func file_space_api_pb_v1_video_proto_rawDescGZIP() []byte {
	file_space_api_pb_v1_video_proto_rawDescOnce.Do(func() {
		file_space_api_pb_v1_video_proto_rawDescData = protoimpl.X.CompressGZIP(file_space_api_pb_v1_video_proto_rawDescData)
	})
	return file_space_api_pb_v1_video_proto_rawDescData
}

var file_space_api_pb_v1_video_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_space_api_pb_v1_video_proto_goTypes = []interface{}{
	(*GetVideoRequest)(nil),       // 0: vector.smcmall.space.api.pb.v1.GetVideoRequest
	(*Video)(nil),                 // 1: vector.smcmall.space.api.pb.v1.Video
	(*TestResponse)(nil),          // 2: vector.smcmall.space.api.pb.v1.TestResponse
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 4: google.protobuf.Empty
}
var file_space_api_pb_v1_video_proto_depIdxs = []int32{
	3, // 0: vector.smcmall.space.api.pb.v1.Video.created_at:type_name -> google.protobuf.Timestamp
	3, // 1: vector.smcmall.space.api.pb.v1.Video.updated_at:type_name -> google.protobuf.Timestamp
	3, // 2: vector.smcmall.space.api.pb.v1.Video.deleted_at:type_name -> google.protobuf.Timestamp
	1, // 3: vector.smcmall.space.api.pb.v1.TestResponse.video:type_name -> vector.smcmall.space.api.pb.v1.Video
	0, // 4: vector.smcmall.space.api.pb.v1.VideoService.Get:input_type -> vector.smcmall.space.api.pb.v1.GetVideoRequest
	4, // 5: vector.smcmall.space.api.pb.v1.VideoService.Test:input_type -> google.protobuf.Empty
	1, // 6: vector.smcmall.space.api.pb.v1.VideoService.Get:output_type -> vector.smcmall.space.api.pb.v1.Video
	2, // 7: vector.smcmall.space.api.pb.v1.VideoService.Test:output_type -> vector.smcmall.space.api.pb.v1.TestResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_space_api_pb_v1_video_proto_init() }
func file_space_api_pb_v1_video_proto_init() {
	if File_space_api_pb_v1_video_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_space_api_pb_v1_video_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVideoRequest); i {
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
		file_space_api_pb_v1_video_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Video); i {
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
		file_space_api_pb_v1_video_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResponse); i {
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
			RawDescriptor: file_space_api_pb_v1_video_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_space_api_pb_v1_video_proto_goTypes,
		DependencyIndexes: file_space_api_pb_v1_video_proto_depIdxs,
		MessageInfos:      file_space_api_pb_v1_video_proto_msgTypes,
	}.Build()
	File_space_api_pb_v1_video_proto = out.File
	file_space_api_pb_v1_video_proto_rawDesc = nil
	file_space_api_pb_v1_video_proto_goTypes = nil
	file_space_api_pb_v1_video_proto_depIdxs = nil
}
