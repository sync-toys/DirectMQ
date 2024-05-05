package directmq

const PROTOCOL_VERSION = 1
const UNKNOWN_PROTOCOL_VERSION = 0

type networkEdgeStateConnecting struct {
	edge                 *networkEdge
	initializeConnection bool
}

var _ networkEdgeState = (*networkEdgeStateConnecting)(nil)

func (n *networkEdgeStateConnecting) GetStateName() edgeStateName {
	return stateConnecting
}

func (n *networkEdgeStateConnecting) OnSet() {
	if !n.initializeConnection {
		return
	}

	err := n.edge.protocol.SupportedProtocolVersions(SupportedProtocolVersionsMessage{
		DataFrame: DataFrame{
			TTL:       ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL,
			Traversed: []string{n.edge.network.config.HostID},
		},
		SupportedVersions: []uint32{PROTOCOL_VERSION},
	})

	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnected{n.edge, "Supported protocol version negotiation failed: " + err.Error(), nil})
	}
}

/* networkParticipant interface implementation */

func (n *networkEdgeStateConnecting) GetSubscribedTopics() []string {
	return []string{}
}

func (n *networkEdgeStateConnecting) HandlePublish(publication PublishMessage) (handled bool) {
	// we are connecting, we cannot handle any publications
	return false
}

func (n *networkEdgeStateConnecting) HandleSubscribe(subscription SubscribeMessage) {
	// we are connecting, we cannot handle any subscriptions
}

func (n *networkEdgeStateConnecting) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	// we are connecting, we cannot handle any unsubscriptions
}

func (n *networkEdgeStateConnecting) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Network terminated: " + terminate.Reason})
}

/* ProtocolDecoderHandler interface implementation */

func (n *networkEdgeStateConnecting) OnSupportedProtocolVersions(message SupportedProtocolVersionsMessage) {
	n.edge.info.BridgedNodeSupportedProtocolVersions = message.SupportedVersions
	n.edge.info.NegotiatedProtocolVersion = PROTOCOL_VERSION

	if message.TTL == ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL {
		n.respondWithSupportedProtocolVersions(message)
		return
	}

	n.initializeEdgeConnection()
}

func (n *networkEdgeStateConnecting) respondWithSupportedProtocolVersions(message SupportedProtocolVersionsMessage) {
	err := n.edge.protocol.SupportedProtocolVersions(SupportedProtocolVersionsMessage{
		DataFrame: DataFrame{
			TTL:       ONLY_DIRECT_CONNECTION_TTL,
			Traversed: append(message.Traversed, n.edge.network.config.HostID),
		},
		SupportedVersions: []uint32{PROTOCOL_VERSION},
	})

	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnected{n.edge, "Respond to supported protocol versions failed: " + err.Error(), nil})
	}
}

func (n *networkEdgeStateConnecting) initializeEdgeConnection() {
	err := n.edge.protocol.InitConnection(InitConnectionMessage{
		DataFrame: DataFrame{
			TTL:       ONLY_DIRECT_CONNECTION_TTL,
			Traversed: []string{n.edge.network.config.HostID},
		},
		MaxMessageSize: n.edge.network.config.HostMaxIncomingMessageSize,
	})

	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnected{n.edge, "Connection initialization failed: " + err.Error(), nil})
	}
}

func (n *networkEdgeStateConnecting) OnInitConnection(message InitConnectionMessage) {
	if n.edge.info.NegotiatedProtocolVersion == UNKNOWN_PROTOCOL_VERSION {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unknown protocol version, missing protocol negotiation"})
		return
	}

	if len(message.Traversed) != 1 {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unexpected number of traversed nodes in init connection message"})
		return
	}

	n.edge.info.BridgedNodeID = message.Traversed[0]
	n.edge.info.BridgedNodeMaxMessageSize = message.MaxMessageSize

	n.acceptEdgeConnection()
}

func (n *networkEdgeStateConnecting) acceptEdgeConnection() {
	err := n.edge.protocol.ConnectionAccepted(ConnectionAcceptedMessage{
		DataFrame: DataFrame{
			TTL:       ONLY_DIRECT_CONNECTION_TTL,
			Traversed: []string{n.edge.network.config.HostID},
		},
		MaxMessageSize: n.edge.network.config.HostMaxIncomingMessageSize,
	})

	if err != nil {
		n.edge.SetState(&networkEdgeStateDisconnected{n.edge, "Connection acceptance failed: " + err.Error(), nil})
		return
	}

	n.edge.SetState(&networkEdgeStateConnected{n.edge})
}

func (n *networkEdgeStateConnecting) OnConnectionAccepted(message ConnectionAcceptedMessage) {
	if n.edge.info.NegotiatedProtocolVersion == UNKNOWN_PROTOCOL_VERSION {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unknown protocol version, missing protocol negotiation"})
		return
	}

	if len(message.Traversed) != 1 {
		n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unexpected number of traversed nodes in init connection message"})
		return
	}

	n.edge.info.BridgedNodeID = message.Traversed[0]
	n.edge.info.BridgedNodeMaxMessageSize = message.MaxMessageSize

	n.edge.SetState(&networkEdgeStateConnected{n.edge})
}

func (n *networkEdgeStateConnecting) OnGracefullyClose(message GracefullyCloseMessage) {
	n.edge.SetState(&networkEdgeStateDisconnected{n.edge, message.Reason, nil})
}

func (n *networkEdgeStateConnecting) OnTerminateNetwork(message TerminateNetworkMessage) {
	n.edge.SetState(&networkEdgeStateDisconnected{n.edge, "Network terminated", nil})
	n.edge.network.Terminated(message, n.edge)
}

func (n *networkEdgeStateConnecting) OnPublish(message PublishMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unexpected publish message in connection process"})
}

func (n *networkEdgeStateConnecting) OnSubscribe(message SubscribeMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unexpected subscribe message in connection process"})
}

func (n *networkEdgeStateConnecting) OnUnsubscribe(message UnsubscribeMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Unexpected unsubscribe message in connection process"})
}

func (n *networkEdgeStateConnecting) OnMalformedMessage(message MalformedMessage) {
	n.edge.SetState(&networkEdgeStateDisconnecting{n.edge, "Malformed message in connection process"})
}
