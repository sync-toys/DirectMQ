package directmq

type DiagnosticsAPI interface {
	OnConnectionEstablished(callback func(bridgedNodeID string, portal Portal))
	OnConnectionLost(callback func(bridgedNodeID, reason string, portal Portal))
	OnPublication(callback func(publication PublishMessage))
	OnSubscription(callback func(subscription SubscribeMessage))
	OnUnsubscribe(callback func(unsubscribe UnsubscribeMessage))
	OnTerminateNetwork(callback func(terminate TerminateNetworkMessage))
}

// TODO: handle protocol writing errors
type diagnosticsAPI struct {
	onConnectionEstablished func(bridgedNodeID string, portal Portal)
	onConnectionLost        func(bridgedNodeID, reason string, portal Portal)

	onPublication      func(publication PublishMessage)
	onSubscription     func(subscription SubscribeMessage)
	onUnsubscribe      func(unsubscribe UnsubscribeMessage)
	onTerminateNetwork func(terminate TerminateNetworkMessage)
}

var _ networkParticipant = (*diagnosticsAPI)(nil)
var _ DiagnosticsAPI = (*diagnosticsAPI)(nil)

func (d *diagnosticsAPI) GetSubscribedTopics() []string {
	return []string{}
}

func (d *diagnosticsAPI) WillHandleTopic(topic string) bool {
	return false
}

func (d *diagnosticsAPI) IsOriginOfFrame(message DataFrame) bool {
	return false
}

func (d *diagnosticsAPI) HandleConnectionEstablished(bridgedNodeID string, portal Portal) {
	if d.onConnectionEstablished != nil {
		d.onConnectionEstablished(bridgedNodeID, portal)
	}
}

func (d *diagnosticsAPI) HandleConnectionLost(bridgedNodeID, reason string, portal Portal) {
	if d.onConnectionLost != nil {
		d.onConnectionLost(bridgedNodeID, reason, portal)
	}
}

func (d *diagnosticsAPI) HandlePublish(publication PublishMessage) (handled bool) {
	if d.onPublication != nil {
		d.onPublication(publication)
	}

	return false
}

func (d *diagnosticsAPI) HandleSubscribe(subscription SubscribeMessage) {
	if d.onSubscription != nil {
		d.onSubscription(subscription)
	}
}

func (d *diagnosticsAPI) HandleUnsubscribe(unsubscribe UnsubscribeMessage) {
	if d.onUnsubscribe != nil {
		d.onUnsubscribe(unsubscribe)
	}
}

func (d *diagnosticsAPI) HandleTerminateNetwork(terminate TerminateNetworkMessage) {
	if d.onTerminateNetwork != nil {
		d.onTerminateNetwork(terminate)
	}
}

func (d *diagnosticsAPI) OnConnectionEstablished(callback func(bridgedNodeID string, portal Portal)) {
	d.onConnectionEstablished = callback
}

func (d *diagnosticsAPI) OnConnectionLost(callback func(bridgedNodeID, reason string, portal Portal)) {
	d.onConnectionLost = callback
}

func (d *diagnosticsAPI) OnPublication(callback func(publication PublishMessage)) {
	d.onPublication = callback
}

func (d *diagnosticsAPI) OnSubscription(callback func(subscription SubscribeMessage)) {
	d.onSubscription = callback
}

func (d *diagnosticsAPI) OnUnsubscribe(callback func(unsubscribe UnsubscribeMessage)) {
	d.onUnsubscribe = callback
}

func (d *diagnosticsAPI) OnTerminateNetwork(callback func(terminate TerminateNetworkMessage)) {
	d.onTerminateNetwork = callback
}
