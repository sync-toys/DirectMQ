package directmq

type diagnosticsClient struct {
	onPublication      func(publication Publish)
	onSubscription     func(subscription Subscribe)
	onSubscriptionOnce func(subscription SubscribeOnce)
	onUnsubscribe      func(unsubscribe Unsubscribe)
	onUnsubscribeOnce  func(unsubscribe UnsubscribeOnce)
	onLoopDetection    func(loopDetection LoopDetection)
	onLoopDetected     func(loopDetected LoopDetected)
}

var _ GlobalClient = (*diagnosticsClient)(nil)

func newDiagnosticsClient() *diagnosticsClient {
	return &diagnosticsClient{}
}

func (d *diagnosticsClient) IsConnected() bool {
	return true
}

func (d *diagnosticsClient) GetClientId() string {
	return "diagnostics"
}

func (d *diagnosticsClient) GetSubscribedTopics() []string {
	return []string{}
}

func (d *diagnosticsClient) GetSubscribedOnceTopics() map[string]string {
	return map[string]string{}
}

func (d *diagnosticsClient) HandlePublish(publication Publish) (handled bool) {
	if d.onPublication != nil {
		d.onPublication(publication)
	}

	return false
}

func (d *diagnosticsClient) HandleSubscribe(subscription Subscribe) {
	if d.onSubscription != nil {
		d.onSubscription(subscription)
	}
}

func (d *diagnosticsClient) HandleSubscribeOnce(subscription SubscribeOnce) {
	if d.onSubscriptionOnce != nil {
		d.onSubscriptionOnce(subscription)
	}
}

func (d *diagnosticsClient) HandleUnsubscribe(unsubscribe Unsubscribe) {
	if d.onUnsubscribe != nil {
		d.onUnsubscribe(unsubscribe)
	}
}

func (d *diagnosticsClient) HandleUnsubscribeOnce(unsubscribe UnsubscribeOnce) {
	if d.onUnsubscribeOnce != nil {
		d.onUnsubscribeOnce(unsubscribe)
	}
}

func (d *diagnosticsClient) HandleLoopDetection(loopDetection LoopDetection) {
	if d.onLoopDetection != nil {
		d.onLoopDetection(loopDetection)
	}
}

func (d *diagnosticsClient) HandleLoopDetected(loopDetected LoopDetected) {
	if d.onLoopDetected != nil {
		d.onLoopDetected(loopDetected)
	}
}
