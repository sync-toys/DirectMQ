package dmqportals

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"net/url"

	directmq "github.com/sync-toys/DirectMQ/sdk/go"
)

type TcpPortal struct {
	conn   net.Conn
	closed bool
}

var _ directmq.Portal = (*TcpPortal)(nil)

func newTcpPortal(conn net.Conn) TcpPortal {
	return TcpPortal{
		conn:   conn,
		closed: false,
	}
}

func (p *TcpPortal) Close() error {
	if p.closed {
		return nil
	}

	p.closed = true
	return p.conn.Close()
}

func ReadFullPacket(conn net.Conn) ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	messageSize := binary.BigEndian.Uint32(header)

	message := make([]byte, messageSize)
	_, err = io.ReadFull(conn, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func WriteFullPacket(conn net.Conn, packet []byte) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(packet)))

	_, err := conn.Write(header)
	if err != nil {
		return err
	}

	_, err = conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func (p *TcpPortal) ReadPacket() ([]byte, error) {
	return ReadFullPacket(p.conn)
}

func (p *TcpPortal) WritePacket(packet []byte) error {
	if err := WriteFullPacket(p.conn, packet); err != nil {
		p.closed = true
		return err
	}

	return nil
}

func TcpConnect(u *url.URL, ctx context.Context) (*TcpPortal, error) {
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}

	portal := newTcpPortal(conn)

	go func() {
		<-ctx.Done()
		portal.Close()
	}()

	return &portal, nil
}

func TcpListen(u *url.URL, dmq directmq.NetworkNode, ctx context.Context) error {
	l, err := net.Listen("tcp", u.Host)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		portal := newTcpPortal(conn)
		go dmq.AddListeningEdge(&portal)
	}
}
