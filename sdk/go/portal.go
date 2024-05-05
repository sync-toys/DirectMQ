package directmq

import "io"

type PacketReader interface {
	ReadPacket() ([]byte, error)
}

type PacketWriter interface {
	WritePacket([]byte) error
}

type Portal interface {
	PacketReader
	PacketWriter
	io.Closer
}
