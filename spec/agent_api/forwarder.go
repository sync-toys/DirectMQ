package dmqspecagent

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ForwardedMessage struct {
	FromURL   *url.URL
	FromAlias string

	ToURL   *url.URL
	ToAlias string

	Message []byte
}
type Forwarder interface {
	StartWebsocketForwarder() error
	OnMessage(handler func(message ForwardedMessage))
	Abort()
}

type websocketForwarder struct {
	fromURL   *url.URL
	fromAlias string

	toURL   *url.URL
	toAlias string

	incomingConn *websocket.Conn
	outgoingConn *websocket.Conn

	messageHandler func(message ForwardedMessage)
	server         *http.Server

	abortMutex sync.Mutex
}

var _ Forwarder = (*websocketForwarder)(nil)

func NewWebsocketForwarder(from *url.URL, fromAlias string, to *url.URL, toAlias string) *websocketForwarder {
	forwarder := &websocketForwarder{
		fromURL:   from,
		fromAlias: fromAlias,

		toURL:   to,
		toAlias: toAlias,

		incomingConn: nil,
		outgoingConn: nil,

		abortMutex: sync.Mutex{},
	}

	forwarder.initWebsocketServer()

	return forwarder
}

func (f *websocketForwarder) StartWebsocketForwarder() error {
	err := f.server.ListenAndServe()

	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (f *websocketForwarder) OnMessage(handler func(message ForwardedMessage)) {
	f.messageHandler = handler
}

func (f *websocketForwarder) Abort() {
	f.abortMutex.Lock()
	defer f.abortMutex.Unlock()

	f.server.Close()
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")

	if f.incomingConn != nil {
		f.incomingConn.WriteMessage(websocket.CloseMessage, msg)
		f.incomingConn.Close()
		f.incomingConn = nil
	}

	if f.outgoingConn != nil {
		f.outgoingConn.WriteMessage(websocket.CloseMessage, msg)
		f.outgoingConn.Close()
		f.outgoingConn = nil
	}
}

func (f *websocketForwarder) initWebsocketServer() {
	mux := http.NewServeMux()
	mux.HandleFunc(f.fromURL.Path, f.createWebsocketHandler())

	server := &http.Server{
		Addr:    f.fromURL.Host,
		Handler: mux,
	}

	f.server = server
}

func (f *websocketForwarder) createWebsocketHandler() http.HandlerFunc {
	upgrader := websocket.Upgrader{}

	return func(w http.ResponseWriter, r *http.Request) {
		if f.incomingConn != nil {
			panic("websocket forwarder already connected")
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic("unable to upgrade connection: " + err.Error())
		}

		f.incomingConn = c

		retry(5, 100*time.Millisecond, func() struct{} {
			f.outgoingConn, _, err = websocket.DefaultDialer.Dial(f.toURL.String(), nil)
			if err != nil {
				panic("unable to dial to the target websocket: " + err.Error())
			}

			return struct{}{}
		})

		done := make(chan struct{}, 2)

		go f.runForwardingRoutine(f.incomingConn, f.outgoingConn, f.fromURL, f.fromAlias, f.toURL, f.toAlias, done)
		go f.runForwardingRoutine(f.outgoingConn, f.incomingConn, f.toURL, f.toAlias, f.fromURL, f.fromAlias, done)

		<-done
		<-done

		f.Abort()
	}
}

func (f *websocketForwarder) runForwardingRoutine(
	from, to *websocket.Conn,
	fromURL *url.URL, fromAlias string,
	toURL *url.URL, toAlias string,
	done chan struct{},
) {
	defer func() { done <- struct{}{} }()

	for {
		_, data, err := from.ReadMessage()
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

		if err := to.WriteMessage(websocket.BinaryMessage, data); err != nil {
			return
		}
	}
}
