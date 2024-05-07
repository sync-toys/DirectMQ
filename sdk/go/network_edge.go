package directmq

import (
	"fmt"
	"strings"
)

type edgeStateName int

const (
	stateConnecting edgeStateName = iota
	stateConnected
	stateDisconnecting
	stateDisconnected
)

type networkEdgeState interface {
	ProtocolDecoderHandler
	networkParticipant

	GetStateName() edgeStateName
	OnSet()
}

type edgeInfo struct {
	BridgedNodeID                        string
	BridgedNodeMaxMessageSize            uint64
	BridgedNodeSupportedProtocolVersions []uint32
	NegotiatedProtocolVersion            uint32
}

type networkEdge struct {
	network *globalNetwork

	portal   Portal
	protocol Protocol

	state networkEdgeState
	info  edgeInfo

	bridgedNodeSubscriptions *subscriptionList[struct{}]
}

var _ networkParticipant = (*networkEdge)(nil)
var _ ProtocolDecoderHandler = (*networkEdge)(nil)

func newNetworkEdge(portal Portal, network *globalNetwork) *networkEdge {
	edge := &networkEdge{
		network: network,

		portal:   portal,
		protocol: nil,

		state: nil,
		info: edgeInfo{
			BridgedNodeID:                        "",
			BridgedNodeMaxMessageSize:            NO_MAX_MESSAGE_SIZE,
			BridgedNodeSupportedProtocolVersions: []uint32{},
			NegotiatedProtocolVersion:            UNKNOWN_PROTOCOL_VERSION,
		},

		bridgedNodeSubscriptions: newSubscriptionList[struct{}](),
	}

	return edge
}

/* networkEdge API */

func (n *networkEdge) GetStateName() edgeStateName {
	return n.state.GetStateName()
}

func (n *networkEdge) SetState(state networkEdgeState) {
	fmt.Println("Setting edge state to: ", state.GetStateName())
	n.state = state
	n.state.OnSet()
}

func (n *networkEdge) Run() error {
	for {
		if err := n.protocol.ReadFrom(n.portal); err != nil {
			return err
		}
	}
}

/* networkParticipant interface implementation */

func (n *networkEdge) GetSubscribedTopics() []string {
	return n.state.GetSubscribedTopics()
}

func (n *networkEdge) HandlePublish(publication PublishMessage) (handled bool) {
	return n.state.HandlePublish(publication)
}

func (n *networkEdge) HandleSubscribe(subscription SubscribeMessage) {
	n.state.HandleSubscribe(subscription)
}

func (n *networkEdge) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	n.state.HandleUnsubscribe(unsubscribe)
}

func (n *networkEdge) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	n.state.HandleTerminateNetwork(terminate)
}

/* ProtocolDecoderHandler interface implementation */

func (n *networkEdge) OnSupportedProtocolVersions(message SupportedProtocolVersionsMessage) {
	n.state.OnSupportedProtocolVersions(message)
}

func (n *networkEdge) OnInitConnection(message InitConnectionMessage) {
	n.state.OnInitConnection(message)
}

func (n *networkEdge) OnConnectionAccepted(message ConnectionAcceptedMessage) {
	n.state.OnConnectionAccepted(message)
}

func (n *networkEdge) OnGracefullyClose(message GracefullyCloseMessage) {
	n.state.OnGracefullyClose(message)
}

func (n *networkEdge) OnTerminateNetwork(message TerminateNetworkMessage) {
	n.state.OnTerminateNetwork(message)
}

func (n *networkEdge) OnPublish(message PublishMessage) {
	n.state.OnPublish(message)
}

func (n *networkEdge) OnSubscribe(message SubscribeMessage) {
	n.state.OnSubscribe(message)
}

func (n *networkEdge) OnUnsubscribe(message UnsubscribeMessage) {
	n.state.OnUnsubscribe(message)
}

func (n *networkEdge) OnMalformedMessage(message MalformedMessage) {
	n.state.OnMalformedMessage(message)
}

/* networkEdge utility methods */

func (n *networkEdge) shouldForwardMessage(frame DataFrame) bool {
	loopDetected := n.checkForNetworkLoops(frame)
	if loopDetected {
		// report loop, terminate whole network
		n.network.Terminated(TerminateNetworkMessage{
			DataFrame: DataFrame{
				TTL:       ONLY_DIRECT_CONNECTION_TTL,
				Traversed: []string{},
			},
			Reason: "Loop detected: " + strings.Join(frame.Traversed, " -> "),
		})
	}

	return !loopDetected && frame.TTL > 0
}

func (n *networkEdge) checkForNetworkLoops(frame DataFrame) bool {
	existingHosts := make(map[string]struct{})

	for _, hostID := range frame.Traversed {
		if _, exists := existingHosts[hostID]; exists {
			return true
		}
		existingHosts[hostID] = struct{}{}
	}

	return false
}

func (n *networkEdge) updateFrame(frame DataFrame) DataFrame {
	return DataFrame{
		TTL:       frame.TTL - 1,
		Traversed: append(frame.Traversed, n.network.config.HostID),
	}
}

func (n *networkEdge) isCurrentEdgeOriginOfFrame(frame DataFrame) bool {
	if len(frame.Traversed) == 0 {
		return false
	}

	return frame.Traversed[len(frame.Traversed)-1] == n.info.BridgedNodeID
}
