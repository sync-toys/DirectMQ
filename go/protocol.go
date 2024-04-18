package directmq

import (
	"github.com/sync-toys/DirectMQ/go/protocol/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type DeliveryMode uint8

const (
	AT_LEAST_ONCE DeliveryMode = 0
	AT_MOST_ONCE  DeliveryMode = 1
)

type DataFrame struct {
	TTL int32
}

type SupportedProtocolVersions struct {
	DataFrame
	SupportedVersions []uint32
}

type InitConnection struct {
	DataFrame
	MaxMessageSize uint64
	ClientId       string
}

type ConnectionAccepted struct {
	DataFrame
	MaxMessageSize uint64
	ClientId       string
}

type GracefullyClose struct {
	DataFrame
	Reason   string
	ClientId string
}

type LoopDetection struct {
	DataFrame
	Path []string
}

type LoopDetected struct {
	DataFrame
	Path []string
}

type Publish struct {
	DataFrame
	ClientId string
	Topic    string
	Mode     DeliveryMode
	Payload  []byte
}

type Subscribe struct {
	DataFrame
	ClientId string
	Topic    string
}

type SubscribeOnce struct {
	DataFrame
	SubscriptionId string
	ClientId       string
	Topic          string
}

type Unsubscribe struct {
	DataFrame
	ClientId string
	Topic    string
}

type UnsubscribeOnce struct {
	DataFrame
	SubscriptionId string
	ClientId       string
	Topic          string
}

type MalformedMessage struct {
	Message []byte
}

type ProtocolEncoder interface {
	SupportedProtocolVersions(message SupportedProtocolVersions) error
	InitConnection(message InitConnection) error
	ConnectionAccepted(message ConnectionAccepted) error
	GracefullyClose(message GracefullyClose) error
	LoopDetection(message LoopDetection) error
	LoopDetected(message LoopDetected) error
	Publish(message Publish) error
	Subscribe(message Subscribe) error
	SubscribeOnce(message SubscribeOnce) error
	Unsubscribe(message Unsubscribe) error
	UnsubscribeOnce(message UnsubscribeOnce) error
}

type ProtocolDecoder interface {
	ReadFrom(pr PacketReader) error
}

type ProtocolDecoderHandler interface {
	OnSupportedProtocolVersions(message SupportedProtocolVersions)
	OnInitConnection(message InitConnection)
	OnConnectionAccepted(message ConnectionAccepted)
	OnGracefullyClose(message GracefullyClose)
	OnLoopDetection(message LoopDetection)
	OnLoopDetected(message LoopDetected)
	OnPublish(message Publish)
	OnSubscribe(message Subscribe)
	OnSubscribeOnce(message SubscribeOnce)
	OnUnsubscribe(message Unsubscribe)
	OnUnsubscribeOnce(message UnsubscribeOnce)
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
	return func(handler ProtocolDecoderHandler, output PacketWriter) Protocol {
		return newProtobufProtocol(PROTOBUF_FORMAT_BINARY, handler, output)
	}
}

func NewProtobufJSONProtocol() ProtocolFactory {
	return func(handler ProtocolDecoderHandler, output PacketWriter) Protocol {
		return newProtobufProtocol(PROTOBUF_FORMAT_JSON, handler, output)
	}
}

func newProtobufProtocol(format ProtobufProtocolFormat, handler ProtocolDecoderHandler, output PacketWriter) *ProtobufProtocol {
	return &ProtobufProtocol{
		format:  format,
		handler: handler,
		writer:  output,
	}
}

func (p *ProtobufProtocol) SupportedProtocolVersions(message SupportedProtocolVersions) error {
	// TODO: extract SPV from DataFrame as SPV is meant to be universal
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_SupportedProtocolVersions{
			SupportedProtocolVersions: &protocol.SupportedProtocolVersions{
				SupportedProtocolVersions: message.SupportedVersions,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) InitConnection(message InitConnection) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_InitConnection{
			InitConnection: &protocol.InitConnection{
				MaxMessageSize: message.MaxMessageSize,
				ClientId:       message.ClientId,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) ConnectionAccepted(message ConnectionAccepted) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_ConnectionAccepted{
			ConnectionAccepted: &protocol.ConnectionAccepted{
				MaxMessageSize: message.MaxMessageSize,
				ClientId:       message.ClientId,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) GracefullyClose(message GracefullyClose) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_GracefullyClose{
			GracefullyClose: &protocol.GracefullyClose{
				Reason:   message.Reason,
				ClientId: message.ClientId,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) LoopDetection(message LoopDetection) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_LoopDetection{
			LoopDetection: &protocol.LoopDetection{
				Path: message.Path,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) LoopDetected(message LoopDetected) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_LoopDetected{
			LoopDetected: &protocol.LoopDetected{
				Path: message.Path,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) Publish(message Publish) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_Publish{
			Publish: &protocol.Publish{
				ClientId: message.ClientId,
				Topic:    message.Topic,
				Mode:     protocol.DeliveryMode(message.Mode),
				Payload:  message.Payload,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) Subscribe(message Subscribe) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_Subscribe{
			Subscribe: &protocol.Subscribe{
				ClientId: message.ClientId,
				Topic:    message.Topic,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) SubscribeOnce(message SubscribeOnce) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_SubscribeOnce{
			SubscribeOnce: &protocol.SubscribeOnce{
				SubscriptionId: message.SubscriptionId,
				ClientId:       message.ClientId,
				Topic:          message.Topic,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) Unsubscribe(message Unsubscribe) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_Unsubscribe{
			Unsubscribe: &protocol.Unsubscribe{
				ClientId: message.ClientId,
				Topic:    message.Topic,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
}

func (p *ProtobufProtocol) UnsubscribeOnce(message UnsubscribeOnce) error {
	frame := protocol.DataFrame{
		Ttl: message.TTL,
		Command: &protocol.DataFrame_UnsubscribeOnce{
			UnsubscribeOnce: &protocol.UnsubscribeOnce{
				SubscriptionId: message.SubscriptionId,
				ClientId:       message.ClientId,
				Topic:          message.Topic,
			},
		},
	}

	encoded, err := p.marshal(&frame)
	if err != nil {
		return err
	}

	return p.writer.WritePacket(encoded)
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

	switch frame.Command.(type) {
	case *protocol.DataFrame_SupportedProtocolVersions:
		message := frame.Command.(*protocol.DataFrame_SupportedProtocolVersions).SupportedProtocolVersions
		p.handler.OnSupportedProtocolVersions(SupportedProtocolVersions{
			DataFrame:         frameToDataFrame(frame),
			SupportedVersions: message.SupportedProtocolVersions,
		})

	case *protocol.DataFrame_InitConnection:
		message := frame.Command.(*protocol.DataFrame_InitConnection).InitConnection
		p.handler.OnInitConnection(InitConnection{
			DataFrame:      frameToDataFrame(frame),
			MaxMessageSize: message.MaxMessageSize,
			ClientId:       message.ClientId,
		})

	case *protocol.DataFrame_ConnectionAccepted:
		message := frame.Command.(*protocol.DataFrame_ConnectionAccepted).ConnectionAccepted
		p.handler.OnConnectionAccepted(ConnectionAccepted{
			DataFrame:      frameToDataFrame(frame),
			MaxMessageSize: message.MaxMessageSize,
			ClientId:       message.ClientId,
		})

	case *protocol.DataFrame_GracefullyClose:
		message := frame.Command.(*protocol.DataFrame_GracefullyClose).GracefullyClose
		p.handler.OnGracefullyClose(GracefullyClose{
			DataFrame: frameToDataFrame(frame),
			Reason:    message.Reason,
			ClientId:  message.ClientId,
		})

	case *protocol.DataFrame_LoopDetection:
		message := frame.Command.(*protocol.DataFrame_LoopDetection).LoopDetection
		p.handler.OnLoopDetection(LoopDetection{
			DataFrame: frameToDataFrame(frame),
			Path:      message.Path,
		})

	case *protocol.DataFrame_LoopDetected:
		message := frame.Command.(*protocol.DataFrame_LoopDetected).LoopDetected
		p.handler.OnLoopDetected(LoopDetected{
			DataFrame: frameToDataFrame(frame),
			Path:      message.Path,
		})

	case *protocol.DataFrame_Publish:
		message := frame.Command.(*protocol.DataFrame_Publish).Publish
		p.handler.OnPublish(Publish{
			DataFrame: frameToDataFrame(frame),
			ClientId:  message.ClientId,
			Topic:     message.Topic,
			Mode:      DeliveryMode(message.Mode),
			Payload:   message.Payload,
		})

	case *protocol.DataFrame_Subscribe:
		message := frame.Command.(*protocol.DataFrame_Subscribe).Subscribe
		p.handler.OnSubscribe(Subscribe{
			DataFrame: frameToDataFrame(frame),
			ClientId:  message.ClientId,
			Topic:     message.Topic,
		})

	case *protocol.DataFrame_SubscribeOnce:
		message := frame.Command.(*protocol.DataFrame_SubscribeOnce).SubscribeOnce
		p.handler.OnSubscribeOnce(SubscribeOnce{
			DataFrame:      frameToDataFrame(frame),
			SubscriptionId: message.SubscriptionId,
			ClientId:       message.ClientId,
			Topic:          message.Topic,
		})

	case *protocol.DataFrame_Unsubscribe:
		message := frame.Command.(*protocol.DataFrame_Unsubscribe).Unsubscribe
		p.handler.OnUnsubscribe(Unsubscribe{
			DataFrame: frameToDataFrame(frame),
			ClientId:  message.ClientId,
			Topic:     message.Topic,
		})

	case *protocol.DataFrame_UnsubscribeOnce:
		message := frame.Command.(*protocol.DataFrame_UnsubscribeOnce).UnsubscribeOnce
		p.handler.OnUnsubscribeOnce(UnsubscribeOnce{
			DataFrame:      frameToDataFrame(frame),
			SubscriptionId: message.SubscriptionId,
			ClientId:       message.ClientId,
			Topic:          message.Topic,
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
		TTL: frame.Ttl,
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
