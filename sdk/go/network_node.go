package directmq

type EdgeManager interface {
	AddListeningEdge(portal Portal) error
	AddConnectingEdge(portal Portal) error
	RemoveEdge(portal Portal, reason string)
}

type NetworkNode interface {
	NativeAPI
	DiagnosticsAPI
	EdgeManager

	GetBridgedNodeIDs() []string
	CloseNode(reason string)
}

type networkNode struct {
	network *globalNetwork
	edges   []*networkEdge

	api         *nativeAPI
	diagnostics *diagnosticsAPI

	createProtocolInstance ProtocolFactory

	onConnectionLostCallback func(bridgedNodeID, reason string, portal Portal)
}

var _ NetworkNode = (*networkNode)(nil)

func newNetworkNode(networkConfig NetworkNodeConfig, protocol ProtocolFactory) *networkNode {
	diagnosticsAPI := &diagnosticsAPI{}
	nativeAPI := newNativeAPI()
	globalNetwork := newGlobalNetwork(networkConfig, nativeAPI, diagnosticsAPI)

	nativeAPI.network = globalNetwork

	node := &networkNode{
		network: globalNetwork,
		edges:   make([]*networkEdge, 0),

		api:         nativeAPI,
		diagnostics: diagnosticsAPI,

		createProtocolInstance: protocol,
	}

	diagnosticsAPI.OnConnectionLost(node.handleEdgeConnectionLost)

	return node
}

func NewNetworkNode(networkConfig NetworkNodeConfig, protocol ProtocolFactory) NetworkNode {
	return newNetworkNode(networkConfig, protocol)
}

/* EdgeManager interface implementation */

func (n *networkNode) AddListeningEdge(portal Portal) error {
	edge := newNetworkEdge(portal, n.network)
	edge.protocol = n.createProtocolInstance(edge, portal)
	edge.SetState(&networkEdgeStateConnecting{edge, false})
	n.edges = append(n.edges, edge)
	n.network.participants = append(n.network.participants, edge)
	return edge.Run()
}

func (n *networkNode) AddConnectingEdge(portal Portal) error {
	edge := newNetworkEdge(portal, n.network)
	edge.protocol = n.createProtocolInstance(edge, portal)
	edge.SetState(&networkEdgeStateConnecting{edge, true})
	n.edges = append(n.edges, edge)
	n.network.participants = append(n.network.participants, edge)
	return edge.Run()
}

func (n *networkNode) RemoveEdge(portal Portal, reason string) {
	edge, index := n.findEdgeByPortal(portal)
	if edge != nil {
		n.disconnectEdgeFromNetwork(edge, reason)
		n.edges = append(n.edges[:index], n.edges[index+1:]...)
	}
}

func (n *networkNode) disconnectEdgeFromNetwork(edge *networkEdge, reason string) {
	n.removeNetworkParticipant(edge)

	if edge.GetStateName() == stateDisconnecting || edge.GetStateName() == stateDisconnected {
		return
	}

	edge.SetState(&networkEdgeStateDisconnecting{edge, reason})
}

func (n *networkNode) removeNetworkParticipant(participant networkParticipant) {
	for i, p := range n.network.participants {
		if p == participant {
			n.network.participants = append(n.network.participants[:i], n.network.participants[i+1:]...)
			break
		}
	}
}

func (n *networkNode) handleEdgeConnectionLost(bridgedNodeID, reason string, portal Portal) {
	// difference from the RemoveEdge method is that
	// that this method is not setting the edge state to Disconnecting,
	// as we can assume that the edge is already disconnected

	edge, index := n.findEdgeByPortal(portal)
	if edge != nil {
		n.removeNetworkParticipant(edge)
		n.edges = append(n.edges[:index], n.edges[index+1:]...)
	}

	if n.onConnectionLostCallback != nil {
		n.onConnectionLostCallback(bridgedNodeID, reason, portal)
	}
}

func (n *networkNode) findEdgeByPortal(portal Portal) (edge *networkEdge, index int) {
	for i, edge := range n.edges {
		if edge.portal == portal {
			return edge, i
		}
	}

	return nil, -1
}

func (n *networkNode) GetBridgedNodeIDs() []string {
	ids := make([]string, 0, len(n.edges))
	for _, edge := range n.edges {
		if edge.GetStateName() == stateConnected {
			ids = append(ids, edge.info.BridgedNodeID)
		}
	}

	return ids
}

func (n *networkNode) CloseNode(reason string) {
	for _, edge := range n.edges {
		n.disconnectEdgeFromNetwork(edge, reason)
	}

	n.edges = nil
}

/* NativeAPI interface implementation */

func (n *networkNode) Publish(topic string, payload []byte, deliveryStrategy DeliveryStrategy) {
	n.api.Publish(topic, payload, deliveryStrategy)
}

func (n *networkNode) Subscribe(topic string, handler func(payload []byte)) SubscriptionID {
	return n.api.Subscribe(topic, handler)
}

func (n *networkNode) Unsubscribe(id SubscriptionID) {
	n.api.Unsubscribe(id)
}

/* DiagnosticsAPI interface implementation */

func (n *networkNode) OnConnectionEstablished(callback func(bridgedNodeID string, portal Portal)) {
	n.diagnostics.OnConnectionEstablished(callback)
}

func (n *networkNode) OnConnectionLost(callback func(bridgedNodeID, reason string, portal Portal)) {
	n.onConnectionLostCallback = callback
}

func (n *networkNode) OnPublication(callback func(message PublishMessage)) {
	n.diagnostics.OnPublication(callback)
}

func (n *networkNode) OnSubscription(callback func(message SubscribeMessage)) {
	n.diagnostics.OnSubscription(callback)
}

func (n *networkNode) OnUnsubscribe(callback func(message UnsubscribeMessage)) {
	n.diagnostics.OnUnsubscribe(callback)
}

func (n *networkNode) OnTerminateNetwork(callback func(message TerminateNetworkMessage)) {
	n.diagnostics.OnTerminateNetwork(callback)
}
