// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: directmq/v1/data_frame.proto

package protocol

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

type DataFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ttl       int32    `protobuf:"varint,1,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Traversed []string `protobuf:"bytes,2,rep,name=traversed,proto3" json:"traversed,omitempty"`
	// Types that are assignable to Message:
	//
	//	*DataFrame_SupportedProtocolVersions
	//	*DataFrame_InitConnection
	//	*DataFrame_ConnectionAccepted
	//	*DataFrame_Publish
	//	*DataFrame_Subscribe
	//	*DataFrame_Unsubscribe
	//	*DataFrame_GracefullyClose
	//	*DataFrame_TerminateNetwork
	Message isDataFrame_Message `protobuf_oneof:"message"`
}

func (x *DataFrame) Reset() {
	*x = DataFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_directmq_v1_data_frame_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataFrame) ProtoMessage() {}

func (x *DataFrame) ProtoReflect() protoreflect.Message {
	mi := &file_directmq_v1_data_frame_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataFrame.ProtoReflect.Descriptor instead.
func (*DataFrame) Descriptor() ([]byte, []int) {
	return file_directmq_v1_data_frame_proto_rawDescGZIP(), []int{0}
}

func (x *DataFrame) GetTtl() int32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *DataFrame) GetTraversed() []string {
	if x != nil {
		return x.Traversed
	}
	return nil
}

func (m *DataFrame) GetMessage() isDataFrame_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *DataFrame) GetSupportedProtocolVersions() *SupportedProtocolVersions {
	if x, ok := x.GetMessage().(*DataFrame_SupportedProtocolVersions); ok {
		return x.SupportedProtocolVersions
	}
	return nil
}

func (x *DataFrame) GetInitConnection() *InitConnection {
	if x, ok := x.GetMessage().(*DataFrame_InitConnection); ok {
		return x.InitConnection
	}
	return nil
}

func (x *DataFrame) GetConnectionAccepted() *ConnectionAccepted {
	if x, ok := x.GetMessage().(*DataFrame_ConnectionAccepted); ok {
		return x.ConnectionAccepted
	}
	return nil
}

func (x *DataFrame) GetPublish() *Publish {
	if x, ok := x.GetMessage().(*DataFrame_Publish); ok {
		return x.Publish
	}
	return nil
}

func (x *DataFrame) GetSubscribe() *Subscribe {
	if x, ok := x.GetMessage().(*DataFrame_Subscribe); ok {
		return x.Subscribe
	}
	return nil
}

func (x *DataFrame) GetUnsubscribe() *Unsubscribe {
	if x, ok := x.GetMessage().(*DataFrame_Unsubscribe); ok {
		return x.Unsubscribe
	}
	return nil
}

func (x *DataFrame) GetGracefullyClose() *GracefullyClose {
	if x, ok := x.GetMessage().(*DataFrame_GracefullyClose); ok {
		return x.GracefullyClose
	}
	return nil
}

func (x *DataFrame) GetTerminateNetwork() *TerminateNetwork {
	if x, ok := x.GetMessage().(*DataFrame_TerminateNetwork); ok {
		return x.TerminateNetwork
	}
	return nil
}

type isDataFrame_Message interface {
	isDataFrame_Message()
}

type DataFrame_SupportedProtocolVersions struct {
	SupportedProtocolVersions *SupportedProtocolVersions `protobuf:"bytes,3,opt,name=supported_protocol_versions,json=supportedProtocolVersions,proto3,oneof"`
}

type DataFrame_InitConnection struct {
	InitConnection *InitConnection `protobuf:"bytes,4,opt,name=init_connection,json=initConnection,proto3,oneof"`
}

type DataFrame_ConnectionAccepted struct {
	ConnectionAccepted *ConnectionAccepted `protobuf:"bytes,5,opt,name=connection_accepted,json=connectionAccepted,proto3,oneof"`
}

type DataFrame_Publish struct {
	Publish *Publish `protobuf:"bytes,6,opt,name=publish,proto3,oneof"`
}

type DataFrame_Subscribe struct {
	Subscribe *Subscribe `protobuf:"bytes,7,opt,name=subscribe,proto3,oneof"`
}

type DataFrame_Unsubscribe struct {
	Unsubscribe *Unsubscribe `protobuf:"bytes,8,opt,name=unsubscribe,proto3,oneof"`
}

type DataFrame_GracefullyClose struct {
	GracefullyClose *GracefullyClose `protobuf:"bytes,9,opt,name=gracefully_close,json=gracefullyClose,proto3,oneof"`
}

type DataFrame_TerminateNetwork struct {
	TerminateNetwork *TerminateNetwork `protobuf:"bytes,10,opt,name=terminate_network,json=terminateNetwork,proto3,oneof"`
}

func (*DataFrame_SupportedProtocolVersions) isDataFrame_Message() {}

func (*DataFrame_InitConnection) isDataFrame_Message() {}

func (*DataFrame_ConnectionAccepted) isDataFrame_Message() {}

func (*DataFrame_Publish) isDataFrame_Message() {}

func (*DataFrame_Subscribe) isDataFrame_Message() {}

func (*DataFrame_Unsubscribe) isDataFrame_Message() {}

func (*DataFrame_GracefullyClose) isDataFrame_Message() {}

func (*DataFrame_TerminateNetwork) isDataFrame_Message() {}

var File_directmq_v1_data_frame_proto protoreflect.FileDescriptor

var file_directmq_v1_data_frame_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6d, 0x71, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6d, 0x71, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1d, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2f, 0x76, 0x31, 0x2f, 0x75,
	0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x8d, 0x05, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x74, 0x74, 0x6c,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x76, 0x65, 0x72, 0x73, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x09, 0x74, 0x72, 0x61, 0x76, 0x65, 0x72, 0x73, 0x65, 0x64, 0x12, 0x68,
	0x0a, 0x1b, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x48, 0x00, 0x52, 0x19, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x46, 0x0a, 0x0f, 0x69, 0x6e, 0x69, 0x74,
	0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00,
	0x52, 0x0e, 0x69, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x52, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x48, 0x00,
	0x52, 0x12, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63, 0x65,
	0x70, 0x74, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x07, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x48, 0x00, 0x52, 0x07, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12, 0x36, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x48, 0x00, 0x52, 0x09, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x3c,
	0x0a, 0x0b, 0x75, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x48, 0x00, 0x52,
	0x0b, 0x75, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x49, 0x0a, 0x10,
	0x67, 0x72, 0x61, 0x63, 0x65, 0x66, 0x75, 0x6c, 0x6c, 0x79, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d,
	0x71, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x72, 0x61, 0x63, 0x65, 0x66, 0x75, 0x6c, 0x6c, 0x79, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0f, 0x67, 0x72, 0x61, 0x63, 0x65, 0x66, 0x75, 0x6c,
	0x6c, 0x79, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x11, 0x74, 0x65, 0x72, 0x6d, 0x69,
	0x6e, 0x61, 0x74, 0x65, 0x5f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x48, 0x00, 0x52, 0x10, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_directmq_v1_data_frame_proto_rawDescOnce sync.Once
	file_directmq_v1_data_frame_proto_rawDescData = file_directmq_v1_data_frame_proto_rawDesc
)

func file_directmq_v1_data_frame_proto_rawDescGZIP() []byte {
	file_directmq_v1_data_frame_proto_rawDescOnce.Do(func() {
		file_directmq_v1_data_frame_proto_rawDescData = protoimpl.X.CompressGZIP(file_directmq_v1_data_frame_proto_rawDescData)
	})
	return file_directmq_v1_data_frame_proto_rawDescData
}

var file_directmq_v1_data_frame_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_directmq_v1_data_frame_proto_goTypes = []interface{}{
	(*DataFrame)(nil),                 // 0: directmq.v1.DataFrame
	(*SupportedProtocolVersions)(nil), // 1: directmq.v1.SupportedProtocolVersions
	(*InitConnection)(nil),            // 2: directmq.v1.InitConnection
	(*ConnectionAccepted)(nil),        // 3: directmq.v1.ConnectionAccepted
	(*Publish)(nil),                   // 4: directmq.v1.Publish
	(*Subscribe)(nil),                 // 5: directmq.v1.Subscribe
	(*Unsubscribe)(nil),               // 6: directmq.v1.Unsubscribe
	(*GracefullyClose)(nil),           // 7: directmq.v1.GracefullyClose
	(*TerminateNetwork)(nil),          // 8: directmq.v1.TerminateNetwork
}
var file_directmq_v1_data_frame_proto_depIdxs = []int32{
	1, // 0: directmq.v1.DataFrame.supported_protocol_versions:type_name -> directmq.v1.SupportedProtocolVersions
	2, // 1: directmq.v1.DataFrame.init_connection:type_name -> directmq.v1.InitConnection
	3, // 2: directmq.v1.DataFrame.connection_accepted:type_name -> directmq.v1.ConnectionAccepted
	4, // 3: directmq.v1.DataFrame.publish:type_name -> directmq.v1.Publish
	5, // 4: directmq.v1.DataFrame.subscribe:type_name -> directmq.v1.Subscribe
	6, // 5: directmq.v1.DataFrame.unsubscribe:type_name -> directmq.v1.Unsubscribe
	7, // 6: directmq.v1.DataFrame.gracefully_close:type_name -> directmq.v1.GracefullyClose
	8, // 7: directmq.v1.DataFrame.terminate_network:type_name -> directmq.v1.TerminateNetwork
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_directmq_v1_data_frame_proto_init() }
func file_directmq_v1_data_frame_proto_init() {
	if File_directmq_v1_data_frame_proto != nil {
		return
	}
	file_directmq_v1_connection_proto_init()
	file_directmq_v1_publish_proto_init()
	file_directmq_v1_subscribe_proto_init()
	file_directmq_v1_unsubscribe_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_directmq_v1_data_frame_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataFrame); i {
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
	file_directmq_v1_data_frame_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*DataFrame_SupportedProtocolVersions)(nil),
		(*DataFrame_InitConnection)(nil),
		(*DataFrame_ConnectionAccepted)(nil),
		(*DataFrame_Publish)(nil),
		(*DataFrame_Subscribe)(nil),
		(*DataFrame_Unsubscribe)(nil),
		(*DataFrame_GracefullyClose)(nil),
		(*DataFrame_TerminateNetwork)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_directmq_v1_data_frame_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_directmq_v1_data_frame_proto_goTypes,
		DependencyIndexes: file_directmq_v1_data_frame_proto_depIdxs,
		MessageInfos:      file_directmq_v1_data_frame_proto_msgTypes,
	}.Build()
	File_directmq_v1_data_frame_proto = out.File
	file_directmq_v1_data_frame_proto_rawDesc = nil
	file_directmq_v1_data_frame_proto_goTypes = nil
	file_directmq_v1_data_frame_proto_depIdxs = nil
}
