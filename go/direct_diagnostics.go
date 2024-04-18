package directmq

type DirectDiagnostics struct {
	client *diagnosticsClient
}

func newDirectDiagnostics() *DirectDiagnostics {
	return &DirectDiagnostics{
		client: newDiagnosticsClient(),
	}
}

func (c *DirectDiagnostics) OnPublication(handler func(publication Publish)) {
	c.client.onPublication = handler
}

func (c *DirectDiagnostics) OnSubscription(handler func(subscription Subscribe)) {
	c.client.onSubscription = handler
}

func (c *DirectDiagnostics) OnSubscriptionOnce(handler func(subscription SubscribeOnce)) {
	c.client.onSubscriptionOnce = handler
}

func (c *DirectDiagnostics) OnUnsubscribe(handler func(unsubscribe Unsubscribe)) {
	c.client.onUnsubscribe = handler
}

func (c *DirectDiagnostics) OnUnsubscribeOnce(handler func(unsubscribe UnsubscribeOnce)) {
	c.client.onUnsubscribeOnce = handler
}

func (c *DirectDiagnostics) OnLoopDetection(handler func(loopDetection LoopDetection)) {
	c.client.onLoopDetection = handler
}

func (c *DirectDiagnostics) OnLoopDetected(handler func(loopDetected LoopDetected)) {
	c.client.onLoopDetected = handler
}
