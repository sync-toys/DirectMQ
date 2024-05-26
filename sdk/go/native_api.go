package directmq

type NativeAPI interface {
	Publish(topic string, payload []byte, deliveryStrategy DeliveryStrategy)
	Subscribe(topic string, handler func(payload []byte)) SubscriptionID
	Unsubscribe(id SubscriptionID)
}

type nativeAPI struct {
	network       *globalNetwork
	subscriptions *subscriptionList[func(payload []byte)]
}

var _ networkParticipant = (*nativeAPI)(nil)
var _ NativeAPI = (*nativeAPI)(nil)

func newNativeAPI() *nativeAPI {
	return &nativeAPI{
		subscriptions: newSubscriptionList[func(payload []byte)](),
	}
}

/* Native API interface */

func (n *nativeAPI) Publish(topic string, payload []byte, deliveryStrategy DeliveryStrategy) {
	if !IsCorrectTopicPattern(topic) {
		panic("incorrect topic pattern")
	}

	if len(payload) == 0 {
		panic("forbidden payload value")
	}

	message := PublishMessage{
		DataFrame:        n.getInitialDataFrame(),
		Topic:            topic,
		Payload:          payload,
		DeliveryStrategy: deliveryStrategy,
	}

	n.network.Published(message)
}

func (n *nativeAPI) Subscribe(topic string, handler func(payload []byte)) SubscriptionID {
	if handler == nil {
		panic("handler cannot be nil")
	}

	oldTopics := n.subscriptions.GetOnlyTopLevelSubscribedTopics()
	subscriptionID := n.subscriptions.AddSubscription(topic, &handler)
	n.updateSubscriptions(oldTopics)
	return subscriptionID
}

func (n *nativeAPI) Unsubscribe(id SubscriptionID) {
	oldTopics := n.subscriptions.GetOnlyTopLevelSubscribedTopics()
	n.subscriptions.RemoveSubscription(id)
	n.updateSubscriptions(oldTopics)
}

func (n *nativeAPI) updateSubscriptions(oldTopics []string) {
	topicsToUnsubscribe, topicsToSubscribe := GetDeduplicatedOverlappingTopicsDiff(oldTopics, n.subscriptions.GetOnlyTopLevelSubscribedTopics())

	for _, topic := range topicsToSubscribe {
		n.network.Subscribed(SubscribeMessage{
			DataFrame: n.getInitialDataFrame(),
			Topic:     topic,
		})
	}

	for _, topic := range topicsToUnsubscribe {
		n.network.Unsubscribed(UnsubscribeMessage{
			DataFrame: n.getInitialDataFrame(),
			Topic:     topic,
		})
	}
}

func (n *nativeAPI) getInitialDataFrame() DataFrame {
	return DataFrame{
		TTL:       int32(n.network.config.HostTTL),
		Traversed: []string{},
	}
}

/* Network participant interface */

func (n *nativeAPI) GetSubscribedTopics() []string {
	return n.subscriptions.GetOnlyTopLevelSubscribedTopics()
}

func (n *nativeAPI) WillHandleTopic(topic string) bool {
	return n.subscriptions.WillHandleTopic(topic)
}

func (n *nativeAPI) IsOriginOfFrame(frame DataFrame) bool {
	return len(frame.Traversed) == 0
}

func (n *nativeAPI) HandlePublish(publication PublishMessage) (handled bool) {
	subscribers := randomOrder(n.subscriptions.GetTriggeredSubscriptions(publication.Topic))
	if len(subscribers) == 0 {
		return false
	}

	// Handle AT_MOST_ONCE delivery strategy
	if publication.DeliveryStrategy == AT_MOST_ONCE {
		subscribers[0].Handler(publication.Payload)
		return true
	}

	// Handle AT_LEAST_ONCE delivery strategy
	for _, subscriber := range subscribers {
		subscriber.Handler(publication.Payload)
	}

	return true
}

func (n *nativeAPI) HandleSubscribe(subscription SubscribeMessage) {
	/* No-op */
}

func (n *nativeAPI) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	/* No-op */
}

func (n *nativeAPI) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	/* No-op */
}
