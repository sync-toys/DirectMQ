package directmq

import "math/rand"

type SubscriptionID int

type subscription[THandler any] struct {
	ID           SubscriptionID
	TopicPattern string
	Handler      THandler
}

type subscriptionList[THandler any] struct {
	subscriptions []subscription[THandler]
}

func newSubscriptionList[THandler any]() *subscriptionList[THandler] {
	return &subscriptionList[THandler]{
		subscriptions: make([]subscription[THandler], 0),
	}
}

func (l *subscriptionList[THandler]) getRandomID() SubscriptionID {
	id := SubscriptionID(rand.Int())

	for _, subscription := range l.subscriptions {
		if subscription.ID == id {
			return l.getRandomID()
		}
	}

	return id
}

func (l *subscriptionList[THandler]) AddSubscription(topic string, handler *THandler) SubscriptionID {
	if handler == nil {
		panic("handler must not be nil")
	}

	if !IsCorrectTopicPattern(topic) {
		panic("incorrect topic pattern")
	}

	subscription := subscription[THandler]{
		ID:           l.getRandomID(),
		TopicPattern: topic,
		Handler:      *handler,
	}

	l.subscriptions = append(l.subscriptions, subscription)

	return subscription.ID
}

func (l *subscriptionList[THandler]) RemoveSubscription(id SubscriptionID) {
	for i, subscription := range l.subscriptions {
		if subscription.ID == id {
			l.subscriptions = append(l.subscriptions[:i], l.subscriptions[i+1:]...)
			return
		}
	}
}

func (l *subscriptionList[THandler]) GetTriggeredSubscriptions(topic string) []subscription[THandler] {
	subscriptions := make([]subscription[THandler], 0)
	for _, subscription := range l.subscriptions {
		if MatchTopicPattern(subscription.TopicPattern, topic) {
			subscriptions = append(subscriptions, subscription)
		}
	}

	return subscriptions
}

func (l *subscriptionList[THandler]) GetSubscriptions() []subscription[THandler] {
	subscriptions := make([]subscription[THandler], len(l.subscriptions))
	copy(subscriptions, l.subscriptions)

	return subscriptions
}

func (l *subscriptionList[THandler]) GetUniqueSubscribedTopics() []string {
	topics := make([]string, 0)
	for _, subscription := range l.subscriptions {
		topics = append(topics, subscription.TopicPattern)
	}

	return unique(topics)
}

func (l *subscriptionList[THandler]) GetOnlyTopLevelSubscribedTopics() []string {
	uniqueTopics := l.GetUniqueSubscribedTopics()
	topLevelTopics := DeduplicateOverlappingTopics(uniqueTopics)

	return topLevelTopics
}
