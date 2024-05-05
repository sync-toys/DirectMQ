package directmq

type networkEdgeStateDisconnecting struct {
	edge   *networkEdge
	reason string
}

var _ networkEdgeState = (*networkEdgeStateDisconnecting)(nil)

func (n *networkEdgeStateDisconnecting) GetStateName() edgeStateName {
	return stateDisconnecting
}

func (n *networkEdgeStateDisconnecting) OnSet() {
	n.edge.protocol.GracefullyClose(GracefullyCloseMessage{
		DataFrame: DataFrame{
			TTL:       ONLY_DIRECT_CONNECTION_TTL,
			Traversed: []string{n.edge.network.config.HostID},
		},
		Reason: n.reason,
	})

	n.edge.SetState(&networkEdgeStateDisconnected{n.edge, n.reason, nil})
}

/* networkParticipant interface implementation */

func (n *networkEdgeStateDisconnecting) GetSubscribedTopics() []string {
	return []string{}
}

func (n *networkEdgeStateDisconnecting) HandlePublish(publication PublishMessage) (handled bool) {
	// we are disconnecting, we cannot handle any publications
	return false
}

func (n *networkEdgeStateDisconnecting) HandleSubscribe(subscription SubscribeMessage) {
	// we are disconnecting, we cannot handle any subscriptions
}

func (n *networkEdgeStateDisconnecting) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	// we are disconnecting, we cannot handle any unsubscriptions
}

func (n *networkEdgeStateDisconnecting) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	// we are disconnecting, we cannot handle any termination messages
}

/* ProtocolDecoderHandler interface implementation */

func (n *networkEdgeStateDisconnecting) OnSupportedProtocolVersions(message SupportedProtocolVersionsMessage) {
	// we are disconnecting, we cannot handle any supported protocol versions messages
}

func (n *networkEdgeStateDisconnecting) OnInitConnection(message InitConnectionMessage) {
	// we are disconnecting, we cannot handle any init connection messages
}

func (n *networkEdgeStateDisconnecting) OnConnectionAccepted(message ConnectionAcceptedMessage) {
	// we are disconnecting, we cannot handle any connection accepted messages
}

func (n *networkEdgeStateDisconnecting) OnGracefullyClose(message GracefullyCloseMessage) {
	// we are disconnecting, we cannot handle any gracefully close messages
}

func (n *networkEdgeStateDisconnecting) OnTerminateNetwork(message TerminateNetworkMessage) {
	// we are disconnecting, we cannot handle any terminate network messages
}

func (n *networkEdgeStateDisconnecting) OnPublish(message PublishMessage) {
	// we are disconnecting, we cannot handle any publish messages
}

func (n *networkEdgeStateDisconnecting) OnSubscribe(message SubscribeMessage) {
	// we are disconnecting, we cannot handle any subscribe messages
}

func (n *networkEdgeStateDisconnecting) OnUnsubscribe(message UnsubscribeMessage) {
	// we are disconnecting, we cannot handle any unsubscribe messages
}

func (n *networkEdgeStateDisconnecting) OnMalformedMessage(message MalformedMessage) {
	// we are disconnecting, we cannot handle any malformed messages
}
