// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: shop/post/v1/post.proto

package postv1

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PostStatus int32

const (
	PostStatus_POST_STATUS_UNSPECIFIED PostStatus = 0
	PostStatus_POST_STATUS_PUBLISHED   PostStatus = 1
)

// Enum value maps for PostStatus.
var (
	PostStatus_name = map[int32]string{
		0: "POST_STATUS_UNSPECIFIED",
		1: "POST_STATUS_PUBLISHED",
	}
	PostStatus_value = map[string]int32{
		"POST_STATUS_UNSPECIFIED": 0,
		"POST_STATUS_PUBLISHED":   1,
	}
)

func (x PostStatus) Enum() *PostStatus {
	p := new(PostStatus)
	*p = x
	return p
}

func (x PostStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PostStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_shop_post_v1_post_proto_enumTypes[0].Descriptor()
}

func (PostStatus) Type() protoreflect.EnumType {
	return &file_shop_post_v1_post_proto_enumTypes[0]
}

func (x PostStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PostStatus.Descriptor instead.
func (PostStatus) EnumDescriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{0}
}

type Post struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Content       string                 `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	Status        PostStatus             `protobuf:"varint,4,opt,name=status,proto3,enum=post.v1.PostStatus" json:"status,omitempty"`
	CreateTime    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	UpdateTime    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	DeleteTime    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=delete_time,json=deleteTime,proto3" json:"delete_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Post) Reset() {
	*x = Post{}
	mi := &file_shop_post_v1_post_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Post) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Post) ProtoMessage() {}

func (x *Post) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Post.ProtoReflect.Descriptor instead.
func (*Post) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{0}
}

func (x *Post) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Post) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Post) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Post) GetStatus() PostStatus {
	if x != nil {
		return x.Status
	}
	return PostStatus_POST_STATUS_UNSPECIFIED
}

func (x *Post) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Post) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *Post) GetDeleteTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DeleteTime
	}
	return nil
}

type GetPostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPostRequest) Reset() {
	*x = GetPostRequest{}
	mi := &file_shop_post_v1_post_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPostRequest) ProtoMessage() {}

func (x *GetPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPostRequest.ProtoReflect.Descriptor instead.
func (*GetPostRequest) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{1}
}

func (x *GetPostRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetPostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Post          *Post                  `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPostResponse) Reset() {
	*x = GetPostResponse{}
	mi := &file_shop_post_v1_post_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPostResponse) ProtoMessage() {}

func (x *GetPostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPostResponse.ProtoReflect.Descriptor instead.
func (*GetPostResponse) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{2}
}

func (x *GetPostResponse) GetPost() *Post {
	if x != nil {
		return x.Post
	}
	return nil
}

type CreatePostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Post          *Post                  `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreatePostRequest) Reset() {
	*x = CreatePostRequest{}
	mi := &file_shop_post_v1_post_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostRequest) ProtoMessage() {}

func (x *CreatePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostRequest.ProtoReflect.Descriptor instead.
func (*CreatePostRequest) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{3}
}

func (x *CreatePostRequest) GetPost() *Post {
	if x != nil {
		return x.Post
	}
	return nil
}

type CreatePostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreatePostResponse) Reset() {
	*x = CreatePostResponse{}
	mi := &file_shop_post_v1_post_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostResponse) ProtoMessage() {}

func (x *CreatePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostResponse.ProtoReflect.Descriptor instead.
func (*CreatePostResponse) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{4}
}

type GetPostErrorRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPostErrorRequest) Reset() {
	*x = GetPostErrorRequest{}
	mi := &file_shop_post_v1_post_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPostErrorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPostErrorRequest) ProtoMessage() {}

func (x *GetPostErrorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPostErrorRequest.ProtoReflect.Descriptor instead.
func (*GetPostErrorRequest) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{5}
}

type GetPostErrorResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPostErrorResponse) Reset() {
	*x = GetPostErrorResponse{}
	mi := &file_shop_post_v1_post_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPostErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPostErrorResponse) ProtoMessage() {}

func (x *GetPostErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shop_post_v1_post_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPostErrorResponse.ProtoReflect.Descriptor instead.
func (*GetPostErrorResponse) Descriptor() ([]byte, []int) {
	return file_shop_post_v1_post_proto_rawDescGZIP(), []int{6}
}

var File_shop_post_v1_post_proto protoreflect.FileDescriptor

const file_shop_post_v1_post_proto_rawDesc = "" +
	"\n" +
	"\x17shop/post/v1/post.proto\x12\apost.v1\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\xaa\x02\n" +
	"\x04Post\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12\x18\n" +
	"\acontent\x18\x03 \x01(\tR\acontent\x12+\n" +
	"\x06status\x18\x04 \x01(\x0e2\x13.post.v1.PostStatusR\x06status\x12;\n" +
	"\vcreate_time\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"createTime\x12;\n" +
	"\vupdate_time\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"updateTime\x12;\n" +
	"\vdelete_time\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"deleteTime\"9\n" +
	"\x0eGetPostRequest\x12'\n" +
	"\x02id\x18\x01 \x01(\x03B\x17\x92A\x142\x12The post id field.R\x02id\"4\n" +
	"\x0fGetPostResponse\x12!\n" +
	"\x04post\x18\x01 \x01(\v2\r.post.v1.PostR\x04post\"6\n" +
	"\x11CreatePostRequest\x12!\n" +
	"\x04post\x18\x01 \x01(\v2\r.post.v1.PostR\x04post\"\x14\n" +
	"\x12CreatePostResponse\"\x15\n" +
	"\x13GetPostErrorRequest\"\x16\n" +
	"\x14GetPostErrorResponse*D\n" +
	"\n" +
	"PostStatus\x12\x1b\n" +
	"\x17POST_STATUS_UNSPECIFIED\x10\x00\x12\x19\n" +
	"\x15POST_STATUS_PUBLISHED\x10\x012\xaf\x02\n" +
	"\vPostService\x12W\n" +
	"\aGetPost\x12\x17.post.v1.GetPostRequest\x1a\x18.post.v1.GetPostResponse\"\x19\x82\xd3\xe4\x93\x02\x13\x12\x11/api/v1/post/{id}\x12^\n" +
	"\n" +
	"CreatePost\x12\x1a.post.v1.CreatePostRequest\x1a\x1b.post.v1.CreatePostResponse\"\x17\x82\xd3\xe4\x93\x02\x11:\x01*\"\f/api/v1/post\x12g\n" +
	"\fGetPostError\x12\x1c.post.v1.GetPostErrorRequest\x1a\x1d.post.v1.GetPostErrorResponse\"\x1a\x82\xd3\xe4\x93\x02\x14\x12\x12/api/v1/post-errorB\x9e\x03\x92A\xff\x01\x12\xb5\x01\n" +
	"\bShop API\"K\n" +
	"\fshop project\x12)https://github.com/miiy/goc/examples/shop\x1a\x10none@example.com*W\n" +
	"\x14BSD 3-Clause License\x12?https://github.com/miiy/goc/blob/main/examples/shop/LICENSE.txt2\x031.0*\x03\x01\x02\x04r@\n" +
	"\x0eMore about goc\x12.https://github.com/grpc-ecosystem/grpc-gateway\n" +
	"\vcom.post.v1B\tPostProtoP\x01ZDgithub.com/miiy/goc-quickstart/apis/gen/proto/go/shop/post/v1;postv1\xa2\x02\x03PXX\xaa\x02\aPost.V1\xca\x02\aPost\\V1\xe2\x02\x13Post\\V1\\GPBMetadata\xea\x02\bPost::V1b\x06proto3"

var (
	file_shop_post_v1_post_proto_rawDescOnce sync.Once
	file_shop_post_v1_post_proto_rawDescData []byte
)

func file_shop_post_v1_post_proto_rawDescGZIP() []byte {
	file_shop_post_v1_post_proto_rawDescOnce.Do(func() {
		file_shop_post_v1_post_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_shop_post_v1_post_proto_rawDesc), len(file_shop_post_v1_post_proto_rawDesc)))
	})
	return file_shop_post_v1_post_proto_rawDescData
}

var file_shop_post_v1_post_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_shop_post_v1_post_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_shop_post_v1_post_proto_goTypes = []any{
	(PostStatus)(0),               // 0: post.v1.PostStatus
	(*Post)(nil),                  // 1: post.v1.Post
	(*GetPostRequest)(nil),        // 2: post.v1.GetPostRequest
	(*GetPostResponse)(nil),       // 3: post.v1.GetPostResponse
	(*CreatePostRequest)(nil),     // 4: post.v1.CreatePostRequest
	(*CreatePostResponse)(nil),    // 5: post.v1.CreatePostResponse
	(*GetPostErrorRequest)(nil),   // 6: post.v1.GetPostErrorRequest
	(*GetPostErrorResponse)(nil),  // 7: post.v1.GetPostErrorResponse
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_shop_post_v1_post_proto_depIdxs = []int32{
	0, // 0: post.v1.Post.status:type_name -> post.v1.PostStatus
	8, // 1: post.v1.Post.create_time:type_name -> google.protobuf.Timestamp
	8, // 2: post.v1.Post.update_time:type_name -> google.protobuf.Timestamp
	8, // 3: post.v1.Post.delete_time:type_name -> google.protobuf.Timestamp
	1, // 4: post.v1.GetPostResponse.post:type_name -> post.v1.Post
	1, // 5: post.v1.CreatePostRequest.post:type_name -> post.v1.Post
	2, // 6: post.v1.PostService.GetPost:input_type -> post.v1.GetPostRequest
	4, // 7: post.v1.PostService.CreatePost:input_type -> post.v1.CreatePostRequest
	6, // 8: post.v1.PostService.GetPostError:input_type -> post.v1.GetPostErrorRequest
	3, // 9: post.v1.PostService.GetPost:output_type -> post.v1.GetPostResponse
	5, // 10: post.v1.PostService.CreatePost:output_type -> post.v1.CreatePostResponse
	7, // 11: post.v1.PostService.GetPostError:output_type -> post.v1.GetPostErrorResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_shop_post_v1_post_proto_init() }
func file_shop_post_v1_post_proto_init() {
	if File_shop_post_v1_post_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_shop_post_v1_post_proto_rawDesc), len(file_shop_post_v1_post_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shop_post_v1_post_proto_goTypes,
		DependencyIndexes: file_shop_post_v1_post_proto_depIdxs,
		EnumInfos:         file_shop_post_v1_post_proto_enumTypes,
		MessageInfos:      file_shop_post_v1_post_proto_msgTypes,
	}.Build()
	File_shop_post_v1_post_proto = out.File
	file_shop_post_v1_post_proto_goTypes = nil
	file_shop_post_v1_post_proto_depIdxs = nil
}
