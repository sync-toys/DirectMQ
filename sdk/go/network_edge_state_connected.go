package directmq

type networkEdgeStateConnected struct {
	edge *networkEdge
}

var _ networkEdgeState = (*networkEdgeStateConnected)(nil)

func (n *networkEdgeStateConnected) GetStateName() edgeStateName {
	return stateConnected
}

func (n *networkEdgeStateConnected) OnSet() {
	n.edge.network.diag.HandleConnectionEstablished(n.edge.info.BridgedNodeID, n.edge.portal)
	n.exchangeAllNodeSubscriptions()
}

func (n *networkEdgeStateConnected) exchangeAllNodeSubscriptions() {
	for _, topic := range n.edge.network.GetAllSubscribedTopics() {
		err := n.edge.protocol.Subscribe(SubscribeMessage{Topic: topic})
		if err != nil {
			n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Failed to exchange subscriptions: " + err.Error()})
			return
		}
	}
}

/* networkParticipant interface implementation */

func (n *networkEdgeStateConnected) GetSubscribedTopics() []string {
	return n.edge.bridgedNodeSubscriptions.GetOnlyTopLevelSubscribedTopics()
}

func (n *networkEdgeStateConnected) WillHandleTopic(topic string) bool {
	return n.edge.bridgedNodeSubscriptions.WillHandleTopic(topic)
}

func (n *networkEdgeStateConnected) IsOriginOfFrame(frame DataFrame) bool {
	panic("this method should not be used, use the networkEdge.IsOriginOfFrame method instead")
}

func (n *networkEdgeStateConnected) HandlePublish(publication PublishMessage) (handled bool) {
	if n.edge.IsOriginOfFrame(publication.DataFrame) {
		return false
	}

	if triggered := n.edge.bridgedNodeSubscriptions.GetTriggeredSubscriptions(publication.Topic); len(triggered) == 0 {
		return false
	}

	publicationToForward := PublishMessage{
		DataFrame:        n.edge.updateFrame(publication.DataFrame),
		Topic:            publication.Topic,
		Payload:          publication.Payload,
		DeliveryStrategy: publication.DeliveryStrategy,
	}

	if !n.edge.shouldForwardMessage(publicationToForward.DataFrame) {
		return false
	}

	if n.edge.info.BridgedNodeMaxMessageSize != NO_MAX_MESSAGE_SIZE && uint64(len(publicationToForward.Payload)) > n.edge.info.BridgedNodeMaxMessageSize {
		return false
	}

	err := n.edge.protocol.Publish(publicationToForward)
	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Failed to publish message: " + err.Error()})
		return false
	}

	return true
}

func (n *networkEdgeStateConnected) HandleSubscribe(subscription SubscribeMessage) {
	if n.edge.IsOriginOfFrame(subscription.DataFrame) {
		return
	}

	subscriptionToForward := SubscribeMessage{
		DataFrame: n.edge.updateFrame(subscription.DataFrame),
		Topic:     subscription.Topic,
	}

	if !n.edge.shouldForwardMessage(subscriptionToForward.DataFrame) {
		return
	}

	err := n.edge.protocol.Subscribe(subscriptionToForward)
	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Failed to subscribe: " + err.Error()})
	}
}

func (n *networkEdgeStateConnected) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	if n.edge.IsOriginOfFrame(unsubscribe.DataFrame) {
		return
	}

	unsubscriptionToForward := UnsubscribeMessage{
		DataFrame: n.edge.updateFrame(unsubscribe.DataFrame),
		Topic:     unsubscribe.Topic,
	}

	if !n.edge.shouldForwardMessage(unsubscriptionToForward.DataFrame) {
		return
	}

	err := n.edge.protocol.Unsubscribe(unsubscriptionToForward)
	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Failed to unsubscribe: " + err.Error()})
	}
}

func (n *networkEdgeStateConnected) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	err := n.edge.protocol.TerminateNetwork(TerminateNetworkMessage{
		DataFrame: DataFrame{
			TTL:       ONLY_DIRECT_CONNECTION_TTL,
			Traversed: append(terminate.Traversed, n.edge.network.config.HostID),
		},
		Reason: terminate.Reason,
	})

	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnected{n.edge, "Failed to terminate network edge: " + err.Error(), nil})
		return
	}

	n.edge.SetState(&networkEdgeStateDisconnected{n.edge, terminate.Reason, nil})
}

/* ProtocolDecoderHandler interface implementation */

func (n *networkEdgeStateConnected) OnSupportedProtocolVersions(message SupportedProtocolVersionsMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Received supported protocol versions message in connected state"})
}

func (n *networkEdgeStateConnected) OnInitConnection(message InitConnectionMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Received connection initialize message in connected state"})
}

func (n *networkEdgeStateConnected) OnConnectionAccepted(message ConnectionAcceptedMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Received connection accepted message in connected state"})
}

func (n *networkEdgeStateConnected) OnGracefullyClose(message GracefullyCloseMessage) {
	n.edge.SetState(&networkEdgeStateDisconnected{n.edge, message.Reason, nil})
}

func (n *networkEdgeStateConnected) OnTerminateNetwork(message TerminateNetworkMessage) {
	n.edge.SetState(&networkEdgeStateDisconnected{n.edge, message.Reason, nil})
	n.edge.network.Terminated(message)
}

func (n *networkEdgeStateConnected) OnPublish(message PublishMessage) {
	n.edge.network.Published(message)
}

func (n *networkEdgeStateConnected) OnSubscribe(message SubscribeMessage) {
	oldTopics := n.edge.bridgedNodeSubscriptions.GetOnlyTopLevelSubscribedTopics()
	n.edge.bridgedNodeSubscriptions.AddSubscription(message.Topic, &struct{}{})
	n.updateSubscriptions(oldTopics, message.DataFrame)
}

func (n *networkEdgeStateConnected) OnUnsubscribe(message UnsubscribeMessage) {
	subscriptionID, subscriptionFound := n.findSubscriptionUsingTopicPattern(message.Topic)
	if !subscriptionFound {
		return
	}

	oldTopics := n.edge.bridgedNodeSubscriptions.GetOnlyTopLevelSubscribedTopics()
	n.edge.bridgedNodeSubscriptions.RemoveSubscription(subscriptionID)
	n.updateSubscriptions(oldTopics, message.DataFrame)
}

func (n *networkEdgeStateConnected) findSubscriptionUsingTopicPattern(topic string) (SubscriptionID, bool) {
	subscriptions := n.edge.bridgedNodeSubscriptions.GetSubscriptions()

	for _, subscription := range subscriptions {
		if subscription.TopicPattern == topic {
			return subscription.ID, true
		}
	}

	return 0, false
}

func (n *networkEdgeStateConnected) updateSubscriptions(oldTopics []string, frame DataFrame) {
	topicsToUnsubscribe, topicsToSubscribe := GetDeduplicatedOverlappingTopicsDiff(oldTopics, n.edge.bridgedNodeSubscriptions.GetOnlyTopLevelSubscribedTopics())

	for _, topic := range topicsToSubscribe {
		n.edge.network.Subscribed(SubscribeMessage{
			DataFrame: frame,
			Topic:     topic,
		})
	}

	for _, topic := range topicsToUnsubscribe {
		n.edge.network.Unsubscribed(UnsubscribeMessage{
			DataFrame: frame,
			Topic:     topic,
		})
	}
}

func (n *networkEdgeStateConnected) OnMalformedMessage(message MalformedMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Malformed message received"})
}
