package directmq

type DirectBridge struct {
	bridge *Bridge
}

func newDirectBridge(
	hostClientId string,
	hostMaxMessageSize uint64,
	ttl int32,

	queue *GlobalQueue,
	portal Portal,
	protocol ProtocolFactory,
) *DirectBridge {
	return &DirectBridge{
		bridge: NewBridge(
			queue,
			portal,
			protocol,
			hostClientId,
			hostMaxMessageSize,
			ttl,
		),
	}
}

func (c *DirectBridge) Connect() error {
	return c.bridge.Connect()
}

func (c *DirectBridge) Listen() error {
	return c.bridge.Listen()
}

func (c *DirectBridge) Disconnect(reason string) []error {
	return c.bridge.Disconnect(reason)
}

func (c *DirectBridge) OnReady(handler ReadyHandler) {
	c.bridge.onReady = handler
}

func (c *DirectBridge) OnConnectionEstablished(handler ConnectionEstablishedHandler) {
	c.bridge.onConnectionEstablished = handler
}

func (c *DirectBridge) OnConnectionLost(handler ConnectionLostHandler) {
	c.bridge.onConnectionLost = handler
}
