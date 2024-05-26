package directmq

type networkEdgeStateDisconnected struct {
	edge *networkEdge

	reason     string
	closeError error
}

var _ networkEdgeState = (*networkEdgeStateDisconnected)(nil)

func (n *networkEdgeStateDisconnected) GetStateName() edgeStateName {
	return stateDisconnected
}

func (n *networkEdgeStateDisconnected) OnSet() {
	n.closeError = n.edge.portal.Close()
	n.edge.network.diag.HandleConnectionLost(n.edge.info.BridgedNodeID, n.reason, n.edge.portal)

	n.revokeAllBridgedNodeSubscriptionsFromNetwork()

	n.edge.info.BridgedNodeID = ""
	n.edge.info.BridgedNodeMaxMessageSize = NO_MAX_MESSAGE_SIZE
	n.edge.info.BridgedNodeSupportedProtocolVersions = []uint32{}
	n.edge.info.NegotiatedProtocolVersion = UNKNOWN_PROTOCOL_VERSION
}

func (n *networkEdgeStateDisconnected) revokeAllBridgedNodeSubscriptionsFromNetwork() {
	for _, topic := range n.edge.bridgedNodeSubscriptions.GetOnlyTopLevelSubscribedTopics() {
		unsubscribeMessage := UnsubscribeMessage{
			DataFrame: DataFrame{
				TTL:       int32(n.edge.network.config.HostTTL),
				Traversed: []string{n.edge.info.BridgedNodeID, n.edge.network.config.HostID},
			},
			Topic: topic,
		}

		n.edge.network.Unsubscribed(unsubscribeMessage)
	}

	n.edge.bridgedNodeSubscriptions = newSubscriptionList[struct{}]()
}

/* networkParticipant interface implementation */

func (n *networkEdgeStateDisconnected) GetSubscribedTopics() []string {
	return []string{}
}

func (n *networkEdgeStateDisconnected) WillHandleTopic(topic string) bool {
	return false
}

func (n *networkEdgeStateDisconnected) IsOriginOfFrame(frame DataFrame) bool {
	panic("this method should not be used, use the networkEdge.IsOriginOfFrame method instead")
}

func (n *networkEdgeStateDisconnected) HandlePublish(publication PublishMessage) (handled bool) {
	// we are disconnected, we cannot handle any publications
	return false
}

func (n *networkEdgeStateDisconnected) HandleSubscribe(subscription SubscribeMessage) {
	// we are disconnected, we cannot handle any subscriptions
}

func (n *networkEdgeStateDisconnected) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	// we are disconnected, we cannot handle any unsubscriptions
}

func (n *networkEdgeStateDisconnected) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	// we are disconnected, we cannot handle any termination messages
}

/* ProtocolDecoderHandler interface implementation */

func (n *networkEdgeStateDisconnected) OnSupportedProtocolVersions(message SupportedProtocolVersionsMessage) {
	// we are disconnected, we cannot handle any supported protocol versions messages
}

func (n *networkEdgeStateDisconnected) OnInitConnection(message InitConnectionMessage) {
	// we are disconnected, we cannot handle any init connection messages
}

func (n *networkEdgeStateDisconnected) OnConnectionAccepted(message ConnectionAcceptedMessage) {
	// we are disconnected, we cannot handle any connection accepted messages
}

func (n *networkEdgeStateDisconnected) OnGracefullyClose(message GracefullyCloseMessage) {
	// we are disconnected, we cannot handle any gracefully close messages
}

func (n *networkEdgeStateDisconnected) OnTerminateNetwork(message TerminateNetworkMessage) {
	// we are disconnected, we cannot handle any terminate network messages
}

func (n *networkEdgeStateDisconnected) OnPublish(message PublishMessage) {
	// we are disconnected, we cannot handle any publish messages
}

func (n *networkEdgeStateDisconnected) OnSubscribe(message SubscribeMessage) {
	// we are disconnected, we cannot handle any subscribe messages
}

func (n *networkEdgeStateDisconnected) OnUnsubscribe(message UnsubscribeMessage) {
	// we are disconnected, we cannot handle any unsubscribe messages
}

func (n *networkEdgeStateDisconnected) OnMalformedMessage(message MalformedMessage) {
	// we are disconnected, we cannot handle any malformed messages
}
