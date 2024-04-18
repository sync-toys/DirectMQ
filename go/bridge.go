package directmq

import "strings"

const (
	DEFAULT_TTL                              = 32
	ONLY_DIRECT_CONNECTION_TTL               = 1
	ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL = 2
	NO_MAX_MESSAGE_SIZE                      = 0
	UNKNOWN_PROTOCOL_VERSION                 = 0
	CURRENT_PROTOCOL_VERSION                 = 1
)

type ReadyHandler func()

type ConnectionEstablishedHandler func(
	bridgedClientId string,
	bridgedProtocolVersion uint32,
	bridgedMaxMessageSize uint64,
	ttl int32,
)

type ConnectionLostHandler func(bridgedClientId, reason string)

type Bridge struct {
	hostClientId       string
	hostMaxMessageSize uint64

	bridgedClientId        *string
	bridgedProtocolVersion uint32
	bridgedMaxMessageSize  uint64
	ttl                    int32

	queue         *GlobalQueue
	subscriptions *SubscriptionList

	encoder ProtocolEncoder
	decoder ProtocolDecoder

	portal Portal

	opened             chan bool
	gotProtocolVersion chan struct{}

	onReady                 ReadyHandler
	onConnectionEstablished ConnectionEstablishedHandler
	onConnectionLost        ConnectionLostHandler
}

var _ GlobalClient = (*Bridge)(nil)
var _ ProtocolDecoderHandler = (*Bridge)(nil)

func NewBridge(
	queue *GlobalQueue,
	portal Portal,
	protocolFactory ProtocolFactory,
	hostClientId string,
	hostMaxMessageSize uint64,
	ttl int32,
) *Bridge {
	bridge := &Bridge{
		hostClientId:       hostClientId,
		hostMaxMessageSize: hostMaxMessageSize,

		bridgedClientId:        nil,
		bridgedProtocolVersion: UNKNOWN_PROTOCOL_VERSION,
		bridgedMaxMessageSize:  NO_MAX_MESSAGE_SIZE,
		ttl:                    ttl,

		queue:         queue,
		subscriptions: NewSubscriptionList(),

		portal: portal,

		opened:             make(chan bool, 1),
		gotProtocolVersion: make(chan struct{}, 1),
	}

	protocol := protocolFactory(bridge, portal)

	bridge.encoder = protocol
	bridge.decoder = protocol

	return bridge
}

func (b *Bridge) Connect() error {
	errChan := make(chan error)

	go func() {
		opened := <-b.opened
		if !opened {
			return
		}

		err := b.initializeConnection()
		errChan <- err
	}()

	go func() {
		err := b.Listen()
		errChan <- err
	}()

	return <-errChan
}

func (b *Bridge) Listen() error {
	b.opened <- true
	if b.onReady != nil {
		b.onReady()
	}

	for {
		if err := b.decoder.ReadFrom(b.portal); err != nil {
			b.Disconnect(ReadError)
			return err
		}
	}
}

func (b *Bridge) Disconnect(reason string) (errors []error) {
	if !b.IsConnected() {
		return nil
	}

	if b.onConnectionLost != nil {
		defer b.onConnectionLost(*b.bridgedClientId, reason)
	}

	for _, topic := range b.subscriptions.GetSubscribedTopics() {
		b.queue.Unsubscribed(Unsubscribe{
			ClientId: b.GetClientId(),
			Topic:    topic,
		})

		b.subscriptions.Unsubscribe(topic)
	}

	for topic, subscriptionId := range b.subscriptions.GetSubscribedOnceTopics() {
		b.queue.UnsubscribedOnce(UnsubscribeOnce{
			ClientId:       b.GetClientId(),
			Topic:          topic,
			SubscriptionId: subscriptionId,
		})

		b.subscriptions.UnsubscribeOnceById(subscriptionId)
	}

	gcEncodingErr := b.encoder.GracefullyClose(GracefullyClose{
		DataFrame: DataFrame{
			TTL: ONLY_DIRECT_CONNECTION_TTL,
		},
		Reason:   reason,
		ClientId: b.GetClientId(),
	})
	if gcEncodingErr != nil {
		errors = append(errors, gcEncodingErr)
	}

	closeErr := b.portal.Close()
	if closeErr != nil {
		errors = append(errors, closeErr)
	}

	b.bridgedClientId = nil
	b.bridgedProtocolVersion = UNKNOWN_PROTOCOL_VERSION
	b.bridgedMaxMessageSize = NO_MAX_MESSAGE_SIZE
	b.opened = make(chan bool, 1)
	b.gotProtocolVersion = make(chan struct{}, 1)

	return
}

/* GlobalClient */

func (b *Bridge) IsConnected() bool {
	return b.bridgedClientId != nil && b.bridgedProtocolVersion != UNKNOWN_PROTOCOL_VERSION
}

func (b *Bridge) GetClientId() string {
	if !b.IsConnected() {
		return ""
	}

	return *b.bridgedClientId
}

func (b *Bridge) GetSubscribedTopics() []string {
	return b.subscriptions.GetSubscribedTopics()
}

func (b *Bridge) GetSubscribedOnceTopics() map[string]string {
	return b.subscriptions.GetSubscribedOnceTopics()
}

func (b *Bridge) HandlePublish(publication Publish) (handled bool) {
	if !b.IsConnected() {
		return false
	}

	if publication.TTL < 1 {
		return false
	}

	if !b.subscriptions.ShouldHandle(publication) {
		return false
	}

	if b.bridgedMaxMessageSize != NO_MAX_MESSAGE_SIZE && uint64(len(publication.Payload)) > b.bridgedMaxMessageSize {
		return false
	}

	err := b.encoder.Publish(publication)
	if err != nil {
		b.Disconnect(WriteError)
	}
	return err == nil
}

func (b *Bridge) HandleSubscribe(subscription Subscribe) {
	if !b.IsConnected() {
		return
	}

	if subscription.TTL < 1 {
		return
	}

	err := b.encoder.Subscribe(subscription)
	if err != nil {
		b.Disconnect(WriteError)
	}
}

func (b *Bridge) HandleSubscribeOnce(subscription SubscribeOnce) {
	if !b.IsConnected() {
		return
	}

	if subscription.TTL < 1 {
		return
	}

	err := b.encoder.SubscribeOnce(subscription)
	if err != nil {
		b.Disconnect(WriteError)
	}
}

func (b *Bridge) HandleUnsubscribe(unsubscribe Unsubscribe) {
	if !b.IsConnected() {
		return
	}

	if unsubscribe.TTL < 1 {
		return
	}

	err := b.encoder.Unsubscribe(unsubscribe)
	if err != nil {
		b.Disconnect(WriteError)
	}
}

func (b *Bridge) HandleUnsubscribeOnce(unsubscribe UnsubscribeOnce) {
	if !b.IsConnected() {
		return
	}

	if unsubscribe.TTL < 1 {
		return
	}

	err := b.encoder.UnsubscribeOnce(unsubscribe)
	if err != nil {
		b.Disconnect(WriteError)
	}
}

func (b *Bridge) HandleLoopDetection(loopDetection LoopDetection) {
	if !b.IsConnected() {
		return
	}

	if loopDetection.TTL < 1 {
		return
	}

	err := b.encoder.LoopDetection(loopDetection)
	if err != nil {
		b.Disconnect(WriteError)
	}
}

func (b *Bridge) HandleLoopDetected(loopDetected LoopDetected) {
	if !b.IsConnected() {
		return
	}

	if loopDetected.TTL < 1 {
		b.Disconnect(ConnectionLoopDetected)
		return
	}

	b.encoder.LoopDetected(loopDetected)
	b.Disconnect(ConnectionLoopDetected)
}

/* ProtocolDecoderHandler */

func (b *Bridge) OnSupportedProtocolVersions(message SupportedProtocolVersions) {
	if b.bridgedProtocolVersion != UNKNOWN_PROTOCOL_VERSION {
		// already negotiated
		return
	}

	if len(message.SupportedVersions) == 0 {
		b.encoder.GracefullyClose(GracefullyClose{
			DataFrame: DataFrame{
				TTL: ONLY_DIRECT_CONNECTION_TTL,
			},
			Reason:   NoSupportedProtocolVersions,
			ClientId: b.GetClientId(),
		})
		return
	}

	if message.SupportedVersions[0] != CURRENT_PROTOCOL_VERSION {
		// TODO: Implement proper version negotiation
		b.encoder.GracefullyClose(GracefullyClose{
			DataFrame: DataFrame{
				TTL: ONLY_DIRECT_CONNECTION_TTL,
			},
			Reason:   UnsupportedProtocolVersion,
			ClientId: b.GetClientId(),
		})
		return
	}

	b.bridgedProtocolVersion = message.SupportedVersions[0]
	b.gotProtocolVersion <- struct{}{}

	message.TTL--
	if message.TTL < 1 {
		return
	}

	err := b.encoder.SupportedProtocolVersions(SupportedProtocolVersions{
		DataFrame: DataFrame{
			TTL: ONLY_DIRECT_CONNECTION_TTL,
		},
		SupportedVersions: []uint32{CURRENT_PROTOCOL_VERSION},
	})

	if err != nil {
		b.Disconnect(WriteError)
	}
}

func (b *Bridge) OnInitConnection(message InitConnection) {
	if b.bridgedProtocolVersion == UNKNOWN_PROTOCOL_VERSION {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	// TODO: trim all invisible characters, allow only printable ASCII characters
	if strings.ReplaceAll(message.ClientId, " ", "") == "" {
		b.Disconnect(InvalidClientId)
		return
	}

	b.bridgedClientId = &message.ClientId
	b.bridgedMaxMessageSize = message.MaxMessageSize

	err := b.encoder.ConnectionAccepted(ConnectionAccepted{
		DataFrame: DataFrame{
			TTL: ONLY_DIRECT_CONNECTION_TTL,
		},
		MaxMessageSize: b.hostMaxMessageSize,
		ClientId:       b.hostClientId,
	})

	if err != nil {
		b.Disconnect(WriteError)
		return
	}

	if err := b.exchangeSubscriptionsWithNetwork(); err != nil {
		b.Disconnect(WriteError)
		return
	}

	if b.onConnectionEstablished != nil {
		b.onConnectionEstablished(*b.bridgedClientId, b.bridgedProtocolVersion, b.bridgedMaxMessageSize, b.ttl)
	}
}

func (b *Bridge) OnConnectionAccepted(message ConnectionAccepted) {
	if b.bridgedProtocolVersion == UNKNOWN_PROTOCOL_VERSION {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	b.bridgedClientId = &message.ClientId
	b.bridgedMaxMessageSize = message.MaxMessageSize

	err := b.encoder.LoopDetection(LoopDetection{
		DataFrame: DataFrame{
			TTL: b.ttl,
		},
		Path: []string{b.GetClientId()},
	})

	if err != nil {
		b.Disconnect(WriteError)
		return
	}

	err = b.exchangeSubscriptionsWithNetwork()
	if err != nil {
		b.Disconnect(WriteError)
		return
	}

	if b.onConnectionEstablished != nil {
		b.onConnectionEstablished(*b.bridgedClientId, b.bridgedProtocolVersion, b.bridgedMaxMessageSize, b.ttl)
	}
}

func (b *Bridge) OnGracefullyClose(message GracefullyClose) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	b.Disconnect(NormalClose)
}

func (b *Bridge) OnLoopDetection(message LoopDetection) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	message.TTL--
	if message.TTL < 0 {
		b.Disconnect(TTLExceeded)
		return
	}

	message.Path = append(message.Path, b.GetClientId())

	knownClients := make([]string, 0, len(message.Path))
	foundDuplicatedClient := false
	for _, client := range message.Path {
		for _, knownClient := range knownClients {
			if client == knownClient {
				foundDuplicatedClient = true
			}

			knownClients = append(knownClients, client)
		}
	}

	if foundDuplicatedClient {
		b.OnLoopDetected(LoopDetected{
			DataFrame: DataFrame{
				TTL: b.ttl,
			},
			Path: message.Path,
		})
		return
	}

	if message.TTL < 1 {
		return
	}

	b.queue.LoopDetection(message, *b.bridgedClientId)
}

func (b *Bridge) OnLoopDetected(message LoopDetected) {
	if !b.IsConnected() {
		b.Disconnect(ConnectionLoopDetected)
		return
	}

	b.queue.LoopDetected(LoopDetected{
		DataFrame: DataFrame{
			TTL: ONLY_DIRECT_CONNECTION_TTL,
		},
		Path: message.Path,
	})

	b.Disconnect(ConnectionLoopDetected)
}

func (b *Bridge) OnPublish(message Publish) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	message.TTL--
	if message.TTL < 0 {
		b.Disconnect(TTLExceeded)
		return
	}

	b.queue.Published(message)
}

func (b *Bridge) OnSubscribe(message Subscribe) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	message.TTL--
	if message.TTL < 0 {
		b.Disconnect(TTLExceeded)
		return
	}

	b.subscriptions.Subscribe(message.Topic)
	b.queue.Subscribed(message)
}

func (b *Bridge) OnSubscribeOnce(message SubscribeOnce) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	message.TTL--
	if message.TTL < 0 {
		b.Disconnect(TTLExceeded)
		return
	}

	b.subscriptions.SubscribeOnce(message.Topic)
	b.queue.SubscribedOnce(message)
}

func (b *Bridge) OnUnsubscribe(message Unsubscribe) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	message.TTL--
	if message.TTL < 0 {
		b.Disconnect(TTLExceeded)
		return
	}

	b.subscriptions.Unsubscribe(message.Topic)
	b.queue.Unsubscribed(message)
}

func (b *Bridge) OnUnsubscribeOnce(message UnsubscribeOnce) {
	if !b.IsConnected() {
		b.Disconnect(UnexpectedMessageReceived)
		return
	}

	message.TTL--
	if message.TTL < 0 {
		b.Disconnect(TTLExceeded)
		return
	}

	b.subscriptions.UnsubscribeOnceById(message.SubscriptionId)
	b.queue.UnsubscribedOnce(message)
}

func (b *Bridge) OnMalformedMessage(message MalformedMessage) {
	b.Disconnect(MalformedMessageReceived)
}

/* Private helper methods */

func (b *Bridge) initializeConnection() error {
	err := b.encoder.SupportedProtocolVersions(SupportedProtocolVersions{
		DataFrame: DataFrame{
			TTL: ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL,
		},
		SupportedVersions: []uint32{1},
	})

	if err != nil {
		b.Disconnect(WriteError)
		return err
	}

	<-b.gotProtocolVersion

	err = b.encoder.InitConnection(InitConnection{
		DataFrame: DataFrame{
			TTL: ONLY_DIRECT_CONNECTION_TTL,
		},
		ClientId:       b.hostClientId,
		MaxMessageSize: b.hostMaxMessageSize,
	})

	if err != nil {
		b.Disconnect(WriteError)
		return err
	}

	return nil
}

// TODO: exchange subscriptions as HOST or client which subscribed to the topic?
// TODO: topics deduplication?
func (b *Bridge) exchangeSubscriptionsWithNetwork() error {
	for _, topic := range b.queue.GetAllSubscribedTopics() {
		err := b.encoder.Subscribe(Subscribe{
			DataFrame: DataFrame{
				TTL: b.ttl,
			},
			ClientId: b.GetClientId(),
			Topic:    topic,
		})

		if err != nil {
			b.Disconnect(WriteError)
			return err
		}
	}

	for topic, subscriptionIds := range b.queue.GetAllSubscribedOnceTopics() {
		for _, subscriptionId := range subscriptionIds {
			err := b.encoder.SubscribeOnce(SubscribeOnce{
				DataFrame: DataFrame{
					TTL: b.ttl,
				},
				ClientId:       b.GetClientId(),
				SubscriptionId: subscriptionId,
				Topic:          topic,
			})

			if err != nil {
				b.Disconnect(WriteError)
				return err
			}
		}
	}

	return nil
}
