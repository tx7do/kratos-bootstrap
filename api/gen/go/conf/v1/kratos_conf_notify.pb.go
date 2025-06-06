// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: conf/v1/kratos_conf_notify.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 通知消息
type Notification struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sms           *Notification_SMS      `protobuf:"bytes,1,opt,name=sms,proto3,oneof" json:"sms,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Notification) Reset() {
	*x = Notification{}
	mi := &file_conf_v1_kratos_conf_notify_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_conf_v1_kratos_conf_notify_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_conf_v1_kratos_conf_notify_proto_rawDescGZIP(), []int{0}
}

func (x *Notification) GetSms() *Notification_SMS {
	if x != nil {
		return x.Sms
	}
	return nil
}

// 短信
type Notification_SMS struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Endpoint        string                 `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`                                        // 公网接入地址
	RegionId        string                 `protobuf:"bytes,2,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`                        // 地域ID
	AccessKeyId     string                 `protobuf:"bytes,3,opt,name=access_key_id,json=accessKeyId,proto3" json:"access_key_id,omitempty"`             // 访问密钥ID
	AccessKeySecret string                 `protobuf:"bytes,4,opt,name=access_key_secret,json=accessKeySecret,proto3" json:"access_key_secret,omitempty"` // 访问密钥
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Notification_SMS) Reset() {
	*x = Notification_SMS{}
	mi := &file_conf_v1_kratos_conf_notify_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notification_SMS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification_SMS) ProtoMessage() {}

func (x *Notification_SMS) ProtoReflect() protoreflect.Message {
	mi := &file_conf_v1_kratos_conf_notify_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification_SMS.ProtoReflect.Descriptor instead.
func (*Notification_SMS) Descriptor() ([]byte, []int) {
	return file_conf_v1_kratos_conf_notify_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Notification_SMS) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Notification_SMS) GetRegionId() string {
	if x != nil {
		return x.RegionId
	}
	return ""
}

func (x *Notification_SMS) GetAccessKeyId() string {
	if x != nil {
		return x.AccessKeyId
	}
	return ""
}

func (x *Notification_SMS) GetAccessKeySecret() string {
	if x != nil {
		return x.AccessKeySecret
	}
	return ""
}

var File_conf_v1_kratos_conf_notify_proto protoreflect.FileDescriptor

var file_conf_v1_kratos_conf_notify_proto_rawDesc = string([]byte{
	0x0a, 0x20, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x63, 0x6f, 0x6e, 0x66, 0x22, 0xd6, 0x01, 0x0a, 0x0c, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x03, 0x73, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x4d, 0x53, 0x48, 0x00,
	0x52, 0x03, 0x73, 0x6d, 0x73, 0x88, 0x01, 0x01, 0x1a, 0x8e, 0x01, 0x0a, 0x03, 0x53, 0x4d, 0x53,
	0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79, 0x49, 0x64, 0x12, 0x2a, 0x0a,
	0x11, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x4b, 0x65, 0x79, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x73, 0x6d,
	0x73, 0x42, 0x87, 0x01, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x42, 0x15,
	0x4b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x78, 0x37, 0x64, 0x6f, 0x2f, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73,
	0x2d, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x43, 0x58, 0x58, 0xaa, 0x02, 0x04, 0x43, 0x6f, 0x6e, 0x66, 0xca, 0x02, 0x04, 0x43, 0x6f, 0x6e,
	0x66, 0xe2, 0x02, 0x10, 0x43, 0x6f, 0x6e, 0x66, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x04, 0x43, 0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_conf_v1_kratos_conf_notify_proto_rawDescOnce sync.Once
	file_conf_v1_kratos_conf_notify_proto_rawDescData []byte
)

func file_conf_v1_kratos_conf_notify_proto_rawDescGZIP() []byte {
	file_conf_v1_kratos_conf_notify_proto_rawDescOnce.Do(func() {
		file_conf_v1_kratos_conf_notify_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_conf_v1_kratos_conf_notify_proto_rawDesc), len(file_conf_v1_kratos_conf_notify_proto_rawDesc)))
	})
	return file_conf_v1_kratos_conf_notify_proto_rawDescData
}

var file_conf_v1_kratos_conf_notify_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_conf_v1_kratos_conf_notify_proto_goTypes = []any{
	(*Notification)(nil),     // 0: conf.Notification
	(*Notification_SMS)(nil), // 1: conf.Notification.SMS
}
var file_conf_v1_kratos_conf_notify_proto_depIdxs = []int32{
	1, // 0: conf.Notification.sms:type_name -> conf.Notification.SMS
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_conf_v1_kratos_conf_notify_proto_init() }
func file_conf_v1_kratos_conf_notify_proto_init() {
	if File_conf_v1_kratos_conf_notify_proto != nil {
		return
	}
	file_conf_v1_kratos_conf_notify_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_conf_v1_kratos_conf_notify_proto_rawDesc), len(file_conf_v1_kratos_conf_notify_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_conf_v1_kratos_conf_notify_proto_goTypes,
		DependencyIndexes: file_conf_v1_kratos_conf_notify_proto_depIdxs,
		MessageInfos:      file_conf_v1_kratos_conf_notify_proto_msgTypes,
	}.Build()
	File_conf_v1_kratos_conf_notify_proto = out.File
	file_conf_v1_kratos_conf_notify_proto_goTypes = nil
	file_conf_v1_kratos_conf_notify_proto_depIdxs = nil
}
