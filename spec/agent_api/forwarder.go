package dmqspecagent

import (
	"net"
	"net/url"
	"sync"

	"github.com/Lobaro/slip"
	dmqportals "github.com/sync-toys/DirectMQ/sdk/go/portals"
)

type ForwardedMessage struct {
	FromURL   *url.URL
	FromAlias string

	ToURL   *url.URL
	ToAlias string

	Message []byte
}

type Forwarder interface {
	StartForwarder() error
	OnMessage(handler func(message ForwardedMessage))
	Abort()
}

type TcpForwarder interface {
	StartTcpForwarder() error
	OnMessage(handler func(message ForwardedMessage))
	Abort()
}

type tcpForwarder struct {
	fromURL   *url.URL
	fromAlias string

	toURL   *url.URL
	toAlias string

	incomingConn net.Conn
	outgoingConn net.Conn

	messageHandler func(message ForwardedMessage)

	abortMutex sync.Mutex
}

var _ Forwarder = (*tcpForwarder)(nil)

func NewTcpForwarder(from *url.URL, fromAlias string, to *url.URL, toAlias string) *tcpForwarder {
	forwarder := &tcpForwarder{
		fromURL:   from,
		fromAlias: fromAlias,

		toURL:   to,
		toAlias: toAlias,

		incomingConn: nil,
		outgoingConn: nil,

		abortMutex: sync.Mutex{},
	}

	return forwarder
}

func (f *tcpForwarder) StartForwarder() error {
	listener, err := net.Listen("tcp", f.fromURL.Host)
	if err != nil {
		return err
	}

	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	f.incomingConn = conn
	return f.startTcpBridge()
}

func (f *tcpForwarder) startTcpBridge() error {
	outgoing, err := net.Dial("tcp", f.toURL.Host)
	if err != nil {
		f.incomingConn.Close()
		return err
	}

	f.outgoingConn = outgoing

	done := make(chan struct{}, 2)

	go f.runForwardingRoutine(f.incomingConn, f.outgoingConn, f.fromURL, f.fromAlias, f.toURL, f.toAlias, done)
	go f.runForwardingRoutine(f.outgoingConn, f.incomingConn, f.toURL, f.toAlias, f.fromURL, f.fromAlias, done)

	<-done
	<-done

	return nil
}

func (f *tcpForwarder) OnMessage(handler func(message ForwardedMessage)) {
	f.messageHandler = handler
}

func (f *tcpForwarder) Abort() {
}

func (f *tcpForwarder) runForwardingRoutine(
	from, to net.Conn,
	fromURL *url.URL, fromAlias string,
	toURL *url.URL, toAlias string,
	done chan struct{},
) {
	defer func() { done <- struct{}{} }()

	reader := slip.NewReader(from)
	writer := slip.NewWriter(to)

	for {
		data, err := dmqportals.ReadFullSlipPacket(reader)
		if err != nil {
			return
		}

		message := ForwardedMessage{
			FromURL:   fromURL,
			FromAlias: fromAlias,

			ToURL:   toURL,
			ToAlias: toAlias,

			Message: data,
		}

		if f.messageHandler != nil {
			f.messageHandler(message)
		}

		if err := writer.WritePacket(data); err != nil {
			return
		}
	}
}
