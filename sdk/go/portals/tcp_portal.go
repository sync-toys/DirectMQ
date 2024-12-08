package dmqportals

import (
	"context"
	"net"
	"net/url"

	"github.com/Lobaro/slip"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
)

type TcpPortal struct {
	conn   net.Conn
	closed bool

	reader *slip.Reader
	writer *slip.Writer
}

var _ directmq.Portal = (*TcpPortal)(nil)

func newTcpPortal(conn net.Conn) TcpPortal {
	return TcpPortal{
		conn:   conn,
		closed: false,
		reader: slip.NewReader(conn),
		writer: slip.NewWriter(conn),
	}
}

func (p *TcpPortal) Close() error {
	if p.closed {
		return nil
	}

	p.closed = true
	return p.conn.Close()
}

func ReadFullSlipPacket(reader *slip.Reader) ([]byte, error) {
	isLastFrame := false
	fullPacket := []byte{}

	for !isLastFrame {
		bytes, isPrefix, err := reader.ReadPacket()
		if err != nil {
			return nil, err
		}

		isLastFrame = !isPrefix
		fullPacket = append(fullPacket, bytes...)
	}

	return fullPacket, nil
}

func (p *TcpPortal) ReadPacket() ([]byte, error) {
	return ReadFullSlipPacket(p.reader)
}

func (p *TcpPortal) WritePacket(packet []byte) error {
	if err := p.writer.WritePacket(packet); err != nil {
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
		dmq.AddListeningEdge(&portal)
	}
}
