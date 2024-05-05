package directmq

import (
	"github.com/sync-toys/DirectMQ/sdk/go/protocol"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type DeliveryStrategy uint8

const (
	AT_LEAST_ONCE DeliveryStrategy = 0
	AT_MOST_ONCE  DeliveryStrategy = 1
)

type DataFrame struct {
	TTL       int32
	Traversed []string
}

type SupportedProtocolVersionsMessage struct {
	DataFrame
	SupportedVersions []uint32
}

type InitConnectionMessage struct {
	DataFrame
	MaxMessageSize uint64
}

type ConnectionAcceptedMessage struct {
	DataFrame
	MaxMessageSize uint64
}

type GracefullyCloseMessage struct {
	DataFrame
	Reason string
}

type TerminateNetworkMessage struct {
	DataFrame
	Reason string
}

type PublishMessage struct {
	DataFrame
	Topic            string
	DeliveryStrategy DeliveryStrategy
	Payload          []byte
}

type SubscribeMessage struct {
	DataFrame
	Topic string
}

type UnsubscribeMessage struct {
	DataFrame
	Topic string
}

type MalformedMessage struct {
	Message []byte
}

type ProtocolEncoder interface {
	SupportedProtocolVersions(message SupportedProtocolVersionsMessage) error
	InitConnection(message InitConnectionMessage) error
	ConnectionAccepted(message ConnectionAcceptedMessage) error
	GracefullyClose(message GracefullyCloseMessage) error
	TerminateNetwork(message TerminateNetworkMessage) error
	Publish(message PublishMessage) error
	Subscribe(message SubscribeMessage) error
	Unsubscribe(message UnsubscribeMessage) error
}

type ProtocolDecoder interface {
	ReadFrom(pr PacketReader) error
}

type ProtocolDecoderHandler interface {
	OnSupportedProtocolVersions(message SupportedProtocolVersionsMessage)
	OnInitConnection(message InitConnectionMessage)
	OnConnectionAccepted(message ConnectionAcceptedMessage)
	OnGracefullyClose(message GracefullyCloseMessage)
	OnTerminateNetwork(message TerminateNetworkMessage)
	OnPublish(message PublishMessage)
	OnSubscribe(message SubscribeMessage)
	OnUnsubscribe(message UnsubscribeMessage)
	OnMalformedMessage(message MalformedMessage)
}

type Protocol interface {
	ProtocolEncoder
	ProtocolDecoder
}

type ProtocolFactory func(ProtocolDecoderHandler, PacketWriter) Protocol

type ProtobufProtocolFormat uint

const (
	PROTOBUF_FORMAT_BINARY ProtobufProtocolFormat = 0
	PROTOBUF_FORMAT_JSON   ProtobufProtocolFormat = 1
)

type ProtobufProtocol struct {
	format  ProtobufProtocolFormat
	handler ProtocolDecoderHandler
	writer  PacketWriter
}

var _ ProtocolEncoder = (*ProtobufProtocol)(nil)
var _ ProtocolDecoder = (*ProtobufProtocol)(nil)

func NewProtobufBinaryProtocol() ProtocolFactory {
	return func(handler ProtocolDecoderHandler, writer PacketWriter) Protocol {
		return &ProtobufProtocol{
			format:  PROTOBUF_FORMAT_BINARY,
			handler: handler,
			writer:  writer,
		}
	}
}

func NewProtobufJSONProtocol() ProtocolFactory {
	return func(handler ProtocolDecoderHandler, writer PacketWriter) Protocol {
		return &ProtobufProtocol{
			format:  PROTOBUF_FORMAT_JSON,
			handler: handler,
			writer:  writer,
		}
	}
}

func (p *ProtobufProtocol) writeFrame(message proto.Message) error {
	encoded, err := p.marshal(message)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) SupportedProtocolVersions(message SupportedProtocolVersionsMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_SupportedProtocolVersions{
			SupportedProtocolVersions: &protocol.SupportedProtocolVersions{
				SupportedProtocolVersions: message.SupportedVersions,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) InitConnection(message InitConnectionMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_InitConnection{
			InitConnection: &protocol.InitConnection{
				MaxMessageSize: message.MaxMessageSize,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) ConnectionAccepted(message ConnectionAcceptedMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_ConnectionAccepted{
			ConnectionAccepted: &protocol.ConnectionAccepted{
				MaxMessageSize: message.MaxMessageSize,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) GracefullyClose(message GracefullyCloseMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_GracefullyClose{
			GracefullyClose: &protocol.GracefullyClose{
				Reason: message.Reason,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) TerminateNetwork(message TerminateNetworkMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_TerminateNetwork{
			TerminateNetwork: &protocol.TerminateNetwork{
				Reason: message.Reason,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) Publish(message PublishMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_Publish{
			Publish: &protocol.Publish{
				Topic:            message.Topic,
				DeliveryStrategy: protocol.DeliveryStrategy(message.DeliveryStrategy),
				Payload:          message.Payload,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) Subscribe(message SubscribeMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_Subscribe{
			Subscribe: &protocol.Subscribe{
				Topic: message.Topic,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) Unsubscribe(message UnsubscribeMessage) error {
	frame := protocol.DataFrame{
		Ttl:       message.TTL,
		Traversed: message.Traversed,
		Message: &protocol.DataFrame_Unsubscribe{
			Unsubscribe: &protocol.Unsubscribe{
				Topic: message.Topic,
			},
		},
	}

	return p.writeFrame(&frame)
}

func (p *ProtobufProtocol) ReadFrom(pr PacketReader) error {
	data, err := pr.ReadPacket()
	if err != nil {
		return err
	}

	frame := new(protocol.DataFrame)
	if err := p.unmarshal(data, frame); err != nil {
		return err
	}

	switch frame.Message.(type) {
	case *protocol.DataFrame_SupportedProtocolVersions:
		message := frame.Message.(*protocol.DataFrame_SupportedProtocolVersions).SupportedProtocolVersions
		p.handler.OnSupportedProtocolVersions(SupportedProtocolVersionsMessage{
			DataFrame:         frameToDataFrame(frame),
			SupportedVersions: message.SupportedProtocolVersions,
		})

	case *protocol.DataFrame_InitConnection:
		message := frame.Message.(*protocol.DataFrame_InitConnection).InitConnection
		p.handler.OnInitConnection(InitConnectionMessage{
			DataFrame:      frameToDataFrame(frame),
			MaxMessageSize: message.MaxMessageSize,
		})

	case *protocol.DataFrame_ConnectionAccepted:
		message := frame.Message.(*protocol.DataFrame_ConnectionAccepted).ConnectionAccepted
		p.handler.OnConnectionAccepted(ConnectionAcceptedMessage{
			DataFrame:      frameToDataFrame(frame),
			MaxMessageSize: message.MaxMessageSize,
		})

	case *protocol.DataFrame_GracefullyClose:
		message := frame.Message.(*protocol.DataFrame_GracefullyClose).GracefullyClose
		p.handler.OnGracefullyClose(GracefullyCloseMessage{
			DataFrame: frameToDataFrame(frame),
			Reason:    message.Reason,
		})

	case *protocol.DataFrame_TerminateNetwork:
		message := frame.Message.(*protocol.DataFrame_TerminateNetwork).TerminateNetwork
		p.handler.OnTerminateNetwork(TerminateNetworkMessage{
			DataFrame: frameToDataFrame(frame),
			Reason:    message.Reason,
		})

	case *protocol.DataFrame_Publish:
		message := frame.Message.(*protocol.DataFrame_Publish).Publish
		p.handler.OnPublish(PublishMessage{
			DataFrame:        frameToDataFrame(frame),
			Topic:            message.Topic,
			DeliveryStrategy: DeliveryStrategy(message.DeliveryStrategy),
			Payload:          message.Payload,
		})

	case *protocol.DataFrame_Subscribe:
		message := frame.Message.(*protocol.DataFrame_Subscribe).Subscribe
		p.handler.OnSubscribe(SubscribeMessage{
			DataFrame: frameToDataFrame(frame),
			Topic:     message.Topic,
		})

	case *protocol.DataFrame_Unsubscribe:
		message := frame.Message.(*protocol.DataFrame_Unsubscribe).Unsubscribe
		p.handler.OnUnsubscribe(UnsubscribeMessage{
			DataFrame: frameToDataFrame(frame),
			Topic:     message.Topic,
		})

	default:
		p.handler.OnMalformedMessage(MalformedMessage{
			Message: data,
		})
	}

	return nil
}

func frameToDataFrame(frame *protocol.DataFrame) DataFrame {
	return DataFrame{
		TTL:       frame.Ttl,
		Traversed: frame.Traversed,
	}
}

func (p *ProtobufProtocol) marshal(message proto.Message) (encoded []byte, err error) {
	switch p.format {
	case PROTOBUF_FORMAT_BINARY:
		return proto.Marshal(message)
	case PROTOBUF_FORMAT_JSON:
		return protojson.Marshal(message)
	default:
		panic("Unknown protobuf format")
	}
}

func (p *ProtobufProtocol) unmarshal(data []byte, message proto.Message) error {
	switch p.format {
	case PROTOBUF_FORMAT_BINARY:
		return proto.Unmarshal(data, message)
	case PROTOBUF_FORMAT_JSON:
		return protojson.Unmarshal(data, message)
	default:
		panic("Unknown protobuf format")
	}
}
