package directmq

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type MessagesType int

const (
	TextMessages MessagesType = iota
	BinaryMessages
)

type WebsocketPortal struct {
	c *websocket.Conn

	messagesType MessagesType
}

var _ Portal = (*WebsocketPortal)(nil)

func (p *WebsocketPortal) Close() error {
	return p.c.Close()
}

func (p *WebsocketPortal) ReadPacket() ([]byte, error) {
	_, data, err := p.c.ReadMessage()
	return data, err
}

func (p *WebsocketPortal) WritePacket(packet []byte) error {
	return p.c.WriteMessage(websocket.TextMessage, packet)
}

func WebsocketConnect(u *url.URL, messagesType MessagesType) (*WebsocketPortal, error) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &WebsocketPortal{c: c, messagesType: messagesType}, nil
}

func WebsocketListen(u *url.URL, messagesType MessagesType, dmq *DirectMQ, ctx context.Context) error {
	var upgrader = websocket.Upgrader{}

	mux := http.NewServeMux()
	mux.HandleFunc(u.Path, func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		portal := &WebsocketPortal{c: c, messagesType: messagesType}
		bridge := dmq.AddBridge(portal)

		defer bridge.Disconnect("websocket connection ended")

		bridge.Listen()
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
