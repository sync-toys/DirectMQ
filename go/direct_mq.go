package directmq

type DirectMQ struct {
	queue       *GlobalQueue
	nativeAPI   NativeAPI
	Diagnostics *DirectDiagnostics
	bridges     []*DirectBridge
	protocol    ProtocolFactory

	hostClientId       string
	hostMaxMessageSize uint64
	ttl                int32
}

func NewDirectMQ(hostClientId string, hostMaxMessageSize uint64, ttl int32, protocolFactory ProtocolFactory) *DirectMQ {
	queue := NewGlobalQueue()

	nativeClient := newNativeClient(queue, hostClientId, ttl)
	queue.clients = append(queue.clients, nativeClient)

	directDiagnostics := newDirectDiagnostics()
	queue.clients = append(queue.clients, directDiagnostics.client)

	return &DirectMQ{
		queue:       queue,
		nativeAPI:   nativeClient,
		Diagnostics: directDiagnostics,
		bridges:     make([]*DirectBridge, 0),
		protocol:    protocolFactory,

		hostClientId:       hostClientId,
		hostMaxMessageSize: hostMaxMessageSize,
		ttl:                ttl,
	}
}

func (d *DirectMQ) AddBridge(portal Portal) *DirectBridge {
	directBridge := newDirectBridge(
		d.hostClientId,
		d.hostMaxMessageSize,
		d.ttl,

		d.queue,
		portal,
		d.protocol,
	)

	d.bridges = append(d.bridges, directBridge)
	d.queue.clients = append(d.queue.clients, directBridge.bridge)

	return directBridge
}

func (d *DirectMQ) Publish(topic string, mode DeliveryMode, payload []byte) error {
	return d.nativeAPI.Publish(topic, mode, payload)
}

func (d *DirectMQ) Subscribe(topic string, handler SubscriptionHandler) error {
	return d.nativeAPI.Subscribe(topic, handler)
}

func (d *DirectMQ) SubscribeOnce(topic string, handler SubscriptionOnceHandler) (string, error) {
	return d.nativeAPI.SubscribeOnce(topic, handler)
}

func (d *DirectMQ) Unsubscribe(topic string, handler SubscriptionHandler) (uint, error) {
	return d.nativeAPI.Unsubscribe(topic, handler)
}

func (d *DirectMQ) UnsubscribeOnce(topic string, handler SubscriptionOnceHandler) (uint, error) {
	return d.nativeAPI.UnsubscribeOnce(topic, handler)
}

func (d *DirectMQ) UnsubscribeOnceById(subscriptionId string) (uint, error) {
	return d.nativeAPI.UnsubscribeOnceById(subscriptionId)
}
