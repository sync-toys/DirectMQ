// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: directmq/v1/publish.proto

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

type DeliveryStrategy int32

const (
	DeliveryStrategy_DELIVERY_STRATEGY_AT_LEAST_ONCE_UNSPECIFIED DeliveryStrategy = 0
	DeliveryStrategy_DELIVERY_STRATEGY_AT_MOST_ONCE              DeliveryStrategy = 1
)

// Enum value maps for DeliveryStrategy.
var (
	DeliveryStrategy_name = map[int32]string{
		0: "DELIVERY_STRATEGY_AT_LEAST_ONCE_UNSPECIFIED",
		1: "DELIVERY_STRATEGY_AT_MOST_ONCE",
	}
	DeliveryStrategy_value = map[string]int32{
		"DELIVERY_STRATEGY_AT_LEAST_ONCE_UNSPECIFIED": 0,
		"DELIVERY_STRATEGY_AT_MOST_ONCE":              1,
	}
)

func (x DeliveryStrategy) Enum() *DeliveryStrategy {
	p := new(DeliveryStrategy)
	*p = x
	return p
}

func (x DeliveryStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeliveryStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_directmq_v1_publish_proto_enumTypes[0].Descriptor()
}

func (DeliveryStrategy) Type() protoreflect.EnumType {
	return &file_directmq_v1_publish_proto_enumTypes[0]
}

func (x DeliveryStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeliveryStrategy.Descriptor instead.
func (DeliveryStrategy) EnumDescriptor() ([]byte, []int) {
	return file_directmq_v1_publish_proto_rawDescGZIP(), []int{0}
}

type Publish struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic            string           `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	DeliveryStrategy DeliveryStrategy `protobuf:"varint,2,opt,name=delivery_strategy,json=deliveryStrategy,proto3,enum=directmq.v1.DeliveryStrategy" json:"delivery_strategy,omitempty"`
	Size             uint64           `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Payload          []byte           `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Publish) Reset() {
	*x = Publish{}
	if protoimpl.UnsafeEnabled {
		mi := &file_directmq_v1_publish_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publish) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publish) ProtoMessage() {}

func (x *Publish) ProtoReflect() protoreflect.Message {
	mi := &file_directmq_v1_publish_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publish.ProtoReflect.Descriptor instead.
func (*Publish) Descriptor() ([]byte, []int) {
	return file_directmq_v1_publish_proto_rawDescGZIP(), []int{0}
}

func (x *Publish) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Publish) GetDeliveryStrategy() DeliveryStrategy {
	if x != nil {
		return x.DeliveryStrategy
	}
	return DeliveryStrategy_DELIVERY_STRATEGY_AT_LEAST_ONCE_UNSPECIFIED
}

func (x *Publish) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Publish) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_directmq_v1_publish_proto protoreflect.FileDescriptor

var file_directmq_v1_publish_proto_rawDesc = []byte{
	0x0a, 0x19, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6d, 0x71, 0x2e, 0x76, 0x31, 0x22, 0x99, 0x01, 0x0a, 0x07, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x4a, 0x0a, 0x11, 0x64, 0x65,
	0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6d, 0x71,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x52, 0x10, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x53, 0x74,
	0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x2a, 0x67, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79,
	0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x2f, 0x0a, 0x2b, 0x44, 0x45, 0x4c, 0x49,
	0x56, 0x45, 0x52, 0x59, 0x5f, 0x53, 0x54, 0x52, 0x41, 0x54, 0x45, 0x47, 0x59, 0x5f, 0x41, 0x54,
	0x5f, 0x4c, 0x45, 0x41, 0x53, 0x54, 0x5f, 0x4f, 0x4e, 0x43, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x22, 0x0a, 0x1e, 0x44, 0x45, 0x4c,
	0x49, 0x56, 0x45, 0x52, 0x59, 0x5f, 0x53, 0x54, 0x52, 0x41, 0x54, 0x45, 0x47, 0x59, 0x5f, 0x41,
	0x54, 0x5f, 0x4d, 0x4f, 0x53, 0x54, 0x5f, 0x4f, 0x4e, 0x43, 0x45, 0x10, 0x01, 0x42, 0x0c, 0x5a,
	0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_directmq_v1_publish_proto_rawDescOnce sync.Once
	file_directmq_v1_publish_proto_rawDescData = file_directmq_v1_publish_proto_rawDesc
)

func file_directmq_v1_publish_proto_rawDescGZIP() []byte {
	file_directmq_v1_publish_proto_rawDescOnce.Do(func() {
		file_directmq_v1_publish_proto_rawDescData = protoimpl.X.CompressGZIP(file_directmq_v1_publish_proto_rawDescData)
	})
	return file_directmq_v1_publish_proto_rawDescData
}

var file_directmq_v1_publish_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_directmq_v1_publish_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_directmq_v1_publish_proto_goTypes = []interface{}{
	(DeliveryStrategy)(0), // 0: directmq.v1.DeliveryStrategy
	(*Publish)(nil),       // 1: directmq.v1.Publish
}
var file_directmq_v1_publish_proto_depIdxs = []int32{
	0, // 0: directmq.v1.Publish.delivery_strategy:type_name -> directmq.v1.DeliveryStrategy
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_directmq_v1_publish_proto_init() }
func file_directmq_v1_publish_proto_init() {
	if File_directmq_v1_publish_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_directmq_v1_publish_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Publish); i {
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
			RawDescriptor: file_directmq_v1_publish_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_directmq_v1_publish_proto_goTypes,
		DependencyIndexes: file_directmq_v1_publish_proto_depIdxs,
		EnumInfos:         file_directmq_v1_publish_proto_enumTypes,
		MessageInfos:      file_directmq_v1_publish_proto_msgTypes,
	}.Build()
	File_directmq_v1_publish_proto = out.File
	file_directmq_v1_publish_proto_rawDesc = nil
	file_directmq_v1_publish_proto_goTypes = nil
	file_directmq_v1_publish_proto_depIdxs = nil
}
