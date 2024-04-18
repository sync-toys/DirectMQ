package directmq

type GlobalClient interface {
	IsConnected() bool
	GetClientId() string
	GetSubscribedTopics() []string
	GetSubscribedOnceTopics() map[string]string
	HandlePublish(publication Publish) (handled bool)
	HandleSubscribe(subscription Subscribe)
	HandleSubscribeOnce(subscription SubscribeOnce)
	HandleUnsubscribe(unsubscribe Unsubscribe)
	HandleUnsubscribeOnce(unsubscribe UnsubscribeOnce)
	HandleLoopDetection(loopDetection LoopDetection)
	HandleLoopDetected(loopDetected LoopDetected)
}

type GlobalQueue struct {
	clients []GlobalClient
}

func NewGlobalQueue() *GlobalQueue {
	return &GlobalQueue{
		clients: make([]GlobalClient, 0),
	}
}

func (d *GlobalQueue) GetAllSubscribedTopics() []string {
	topics := make([]string, 0)
	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		for _, topic := range client.GetSubscribedTopics() {
			if !contains(topics, topic) {
				topics = append(topics, topic)
			}
		}
	}

	return topics
}

func (d *GlobalQueue) GetAllSubscribedOnceTopics() map[string][]string {
	topics := make(map[string][]string)

	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		for topic, subscriptionId := range client.GetSubscribedOnceTopics() {
			if _, ok := topics[topic]; !ok {
				topics[topic] = make([]string, 0)
			}

			topics[topic] = append(topics[topic], subscriptionId)
		}
	}

	return topics
}

func (d *GlobalQueue) Published(message Publish) {
	for _, client := range randomOrder(d.clients) {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == message.ClientId {
			continue
		}

		handled := client.HandlePublish(message)
		if message.Mode == AT_MOST_ONCE && handled {
			return
		}
	}
}

func (d *GlobalQueue) Subscribed(message Subscribe) {
	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == message.ClientId {
			continue
		}

		client.HandleSubscribe(message)
	}
}

func (d *GlobalQueue) SubscribedOnce(message SubscribeOnce) {
	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == message.ClientId {
			continue
		}

		client.HandleSubscribeOnce(message)
	}
}

func (d *GlobalQueue) Unsubscribed(message Unsubscribe) {
	allClientsNotSubscribe := true

	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == message.ClientId {
			continue
		}

		topics := client.GetSubscribedTopics()
		for _, topic := range topics {
			if topic == message.Topic {
				allClientsNotSubscribe = false
				break
			}
		}
	}

	if allClientsNotSubscribe {
		return
	}

	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == message.ClientId {
			continue
		}

		client.HandleUnsubscribe(message)
	}
}

func (d *GlobalQueue) UnsubscribedOnce(message UnsubscribeOnce) {
	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == message.ClientId {
			continue
		}

		client.HandleUnsubscribeOnce(message)
	}
}

func (d *GlobalQueue) LoopDetection(message LoopDetection, callingClientId string) {
	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		if client.GetClientId() == callingClientId {
			continue
		}

		client.HandleLoopDetection(message)
	}
}

func (d *GlobalQueue) LoopDetected(message LoopDetected) {
	for _, client := range d.clients {
		if !client.IsConnected() {
			continue
		}

		client.HandleLoopDetected(message)
	}
}
