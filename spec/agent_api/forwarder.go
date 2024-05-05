package dmqspecagent

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Forwarder interface {
	Abort()
	OnMessage(handler func(message string))
}

type websocketForwarder struct {
	from *url.URL
	to   *url.URL

	incomingConn *websocket.Conn
	outgoingConn *websocket.Conn

	messageHandler func(message string, sender, receiver *url.URL)
	server         *http.Server
}

func NewWebsocketForwarder(from, to *url.URL) *websocketForwarder {
	forwarder := &websocketForwarder{
		from: from,
		to:   to,

		incomingConn: nil,
		outgoingConn: nil,
	}

	forwarder.initWebsocketServer()

	return forwarder
}

func (f *websocketForwarder) Abort() {
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

func (f *websocketForwarder) OnMessage(handler func(message string, sender, receiver *url.URL)) {
	f.messageHandler = handler
}

func (f *websocketForwarder) StartWebsocketForwarder() error {
	err := f.server.ListenAndServe()

	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (f *websocketForwarder) initWebsocketServer() {
	mux := http.NewServeMux()
	mux.HandleFunc(f.from.Path, f.createWebsocketHandler())

	server := &http.Server{
		Addr:    f.from.Host,
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

		f.outgoingConn, _, err = websocket.DefaultDialer.Dial(f.to.String(), nil)
		if err != nil {
			panic("unable to dial to the target websocket: " + err.Error())
		}

		done := make(chan struct{}, 2)

		go f.runForwardingRoutine(f.incomingConn, f.outgoingConn, f.from, f.to, done)
		go f.runForwardingRoutine(f.outgoingConn, f.incomingConn, f.to, f.from, done)

		<-done
		<-done

		f.Abort()
	}
}

func (f *websocketForwarder) runForwardingRoutine(from, to *websocket.Conn, fromURL, toURL *url.URL, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	for {
		_, data, err := from.ReadMessage()
		if err != nil {
			return
		}

		if f.messageHandler != nil {
			f.messageHandler(string(data), fromURL, toURL)
		}

		if err := to.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}
