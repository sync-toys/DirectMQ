package dmqportals

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	directmq "github.com/sync-toys/DirectMQ/sdk/go"

	"github.com/gorilla/websocket"
)

type MessagesType int

const (
	TextMessages   MessagesType = websocket.TextMessage
	BinaryMessages              = websocket.BinaryMessage
)

type WebsocketPortal struct {
	c      *websocket.Conn
	closed bool

	messagesType MessagesType
}

var _ directmq.Portal = (*WebsocketPortal)(nil)

func (p *WebsocketPortal) Close() error {
	if p.closed {
		return nil
	}

	p.closed = true

	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := p.c.WriteMessage(websocket.CloseMessage, msg); err != nil {
		p.c.Close()
		return err
	}

	return p.c.Close()
}

func (p *WebsocketPortal) ReadPacket() ([]byte, error) {
	messageType, data, err := p.c.ReadMessage()

	if websocket.IsCloseError(err, websocket.CloseNormalClosure) || websocket.IsUnexpectedCloseError(err) {
		p.closed = true
	}

	if err != nil {
		return nil, err
	}

	if messageType != int(p.messagesType) {
		panic("Unexpected message type")
	}

	return data, nil
}

func (p *WebsocketPortal) WritePacket(packet []byte) error {
	err := p.c.WriteMessage(int(p.messagesType), packet)

	if websocket.IsCloseError(err, websocket.CloseNormalClosure) || websocket.IsUnexpectedCloseError(err) {
		p.closed = true
	}

	return err
}

func WebsocketConnect(u *url.URL, messagesType MessagesType) (*WebsocketPortal, error) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &WebsocketPortal{c: c, messagesType: messagesType}, nil
}

func WebsocketListen(u *url.URL, messagesType MessagesType, dmq directmq.NetworkNode, ctx context.Context) error {
	var upgrader = websocket.Upgrader{}

	mux := http.NewServeMux()
	mux.HandleFunc(u.Path, func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		portal := &WebsocketPortal{c: c, messagesType: messagesType}
		dmq.AddListeningEdge(portal)
	})

	server := &http.Server{
		Addr:    u.Host,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func IsNetworkConnectionClosedError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
