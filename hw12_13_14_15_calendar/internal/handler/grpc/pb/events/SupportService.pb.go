// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: SupportService.proto

package events

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title     string                 `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	Date      *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=Date,proto3" json:"Date,omitempty"`
	Duration  *durationpb.Duration   `protobuf:"bytes,4,opt,name=Duration,proto3" json:"Duration,omitempty"`
	UserName  string                 `protobuf:"bytes,5,opt,name=UserName,proto3" json:"UserName,omitempty"`
	UserEmail string                 `protobuf:"bytes,6,opt,name=UserEmail,proto3" json:"UserEmail,omitempty"`
}

func (x *Notification) Reset() {
	*x = Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SupportService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_SupportService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_SupportService_proto_rawDescGZIP(), []int{0}
}

func (x *Notification) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Notification) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Notification) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *Notification) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *Notification) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *Notification) GetUserEmail() string {
	if x != nil {
		return x.UserEmail
	}
	return ""
}

type Notifies struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*Notification `protobuf:"bytes,1,rep,name=List,proto3" json:"List,omitempty"`
}

func (x *Notifies) Reset() {
	*x = Notifies{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SupportService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notifies) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notifies) ProtoMessage() {}

func (x *Notifies) ProtoReflect() protoreflect.Message {
	mi := &file_SupportService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notifies.ProtoReflect.Descriptor instead.
func (*Notifies) Descriptor() ([]byte, []int) {
	return file_SupportService_proto_rawDescGZIP(), []int{1}
}

func (x *Notifies) GetList() []*Notification {
	if x != nil {
		return x.List
	}
	return nil
}

type NotifyIDList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IDList []string `protobuf:"bytes,1,rep,name=IDList,proto3" json:"IDList,omitempty"`
}

func (x *NotifyIDList) Reset() {
	*x = NotifyIDList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SupportService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyIDList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyIDList) ProtoMessage() {}

func (x *NotifyIDList) ProtoReflect() protoreflect.Message {
	mi := &file_SupportService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyIDList.ProtoReflect.Descriptor instead.
func (*NotifyIDList) Descriptor() ([]byte, []int) {
	return file_SupportService_proto_rawDescGZIP(), []int{2}
}

func (x *NotifyIDList) GetIDList() []string {
	if x != nil {
		return x.IDList
	}
	return nil
}

var File_SupportService_proto protoreflect.FileDescriptor

var file_SupportService_proto_rawDesc = []byte{
	0x0a, 0x14, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5, 0x01, 0x0a, 0x0c, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x2e, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x44, 0x61, 0x74, 0x65,
	0x12, 0x35, 0x0a, 0x08, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x55, 0x73, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x22, 0x31, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x65, 0x73, 0x12, 0x25, 0x0a,
	0x04, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x26, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x49, 0x44,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x32, 0xc8, 0x01, 0x0a,
	0x07, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x3b, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x73, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x64, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x79, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x44, 0x0a, 0x10, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x75, 0x70, 0x4f, 0x6c, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x21, 0x5a, 0x1f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x62, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_SupportService_proto_rawDescOnce sync.Once
	file_SupportService_proto_rawDescData = file_SupportService_proto_rawDesc
)

func file_SupportService_proto_rawDescGZIP() []byte {
	file_SupportService_proto_rawDescOnce.Do(func() {
		file_SupportService_proto_rawDescData = protoimpl.X.CompressGZIP(file_SupportService_proto_rawDescData)
	})
	return file_SupportService_proto_rawDescData
}

var file_SupportService_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_SupportService_proto_goTypes = []interface{}{
	(*Notification)(nil),          // 0: api.Notification
	(*Notifies)(nil),              // 1: api.Notifies
	(*NotifyIDList)(nil),          // 2: api.NotifyIDList
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 4: google.protobuf.Duration
	(*emptypb.Empty)(nil),         // 5: google.protobuf.Empty
}
var file_SupportService_proto_depIdxs = []int32{
	3, // 0: api.Notification.Date:type_name -> google.protobuf.Timestamp
	4, // 1: api.Notification.Duration:type_name -> google.protobuf.Duration
	0, // 2: api.Notifies.List:type_name -> api.Notification
	5, // 3: api.support.GetNotifications:input_type -> google.protobuf.Empty
	2, // 4: api.support.SetNotified:input_type -> api.NotifyIDList
	5, // 5: api.support.CleanupOldEvents:input_type -> google.protobuf.Empty
	1, // 6: api.support.GetNotifications:output_type -> api.Notifies
	5, // 7: api.support.SetNotified:output_type -> google.protobuf.Empty
	5, // 8: api.support.CleanupOldEvents:output_type -> google.protobuf.Empty
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_SupportService_proto_init() }
func file_SupportService_proto_init() {
	if File_SupportService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_SupportService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notification); i {
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
		file_SupportService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notifies); i {
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
		file_SupportService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyIDList); i {
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
			RawDescriptor: file_SupportService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_SupportService_proto_goTypes,
		DependencyIndexes: file_SupportService_proto_depIdxs,
		MessageInfos:      file_SupportService_proto_msgTypes,
	}.Build()
	File_SupportService_proto = out.File
	file_SupportService_proto_rawDesc = nil
	file_SupportService_proto_goTypes = nil
	file_SupportService_proto_depIdxs = nil
}
