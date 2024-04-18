package directmq

import (
	"github.com/gobwas/glob"
	"github.com/google/uuid"
)

type deduplicatedSubscription struct {
	topic   string
	matcher glob.Glob
}

type deduplicatedSubscriptionOnce struct {
	subscriptionId string
	topic          string
	matcher        glob.Glob
}

type SubscriptionList struct {
	subscriptions     []deduplicatedSubscription
	subscriptionsOnce []deduplicatedSubscriptionOnce
}

func NewSubscriptionList() *SubscriptionList {
	return &SubscriptionList{
		subscriptions:     make([]deduplicatedSubscription, 0),
		subscriptionsOnce: make([]deduplicatedSubscriptionOnce, 0),
	}
}

func (c *SubscriptionList) Subscribe(topic string) error {
	matcher, err := glob.Compile(topic)
	if err != nil {
		return err
	}

	c.subscriptions = append(c.subscriptions, deduplicatedSubscription{
		topic,
		matcher,
	})

	return nil
}

func (c *SubscriptionList) SubscribeOnce(topic string) (subscriptionId string, err error) {
	matcher, err := glob.Compile(topic)
	if err != nil {
		return
	}

	subscriptionId = uuid.New().String()

	c.subscriptionsOnce = append(c.subscriptionsOnce, deduplicatedSubscriptionOnce{
		subscriptionId,
		topic,
		matcher,
	})

	return
}

func (c *SubscriptionList) Unsubscribe(topic string) (removedSubscriptions uint, err error) {
	for i, sub := range c.subscriptions {
		if sub.topic != topic {
			continue
		}

		c.subscriptions = append(c.subscriptions[:i], c.subscriptions[i+1:]...)
		removedSubscriptions++
	}

	return
}

func (c *SubscriptionList) UnsubscribeOnceById(subscriptionId string) (removedSubscriptions uint, err error) {
	for i, sub := range c.subscriptionsOnce {
		if sub.subscriptionId != subscriptionId {
			continue
		}

		c.subscriptionsOnce = append(c.subscriptionsOnce[:i], c.subscriptionsOnce[i+1:]...)
		removedSubscriptions++
	}

	return
}

func (c *SubscriptionList) GetSubscribedTopics() []string {
	topics := make([]string, len(c.subscriptions))
	for i, sub := range c.subscriptions {
		topics[i] = sub.topic
	}
	return topics
}

func (c *SubscriptionList) GetSubscribedOnceTopics() map[string]string {
	topics := make(map[string]string, len(c.subscriptionsOnce))
	for _, sub := range c.subscriptionsOnce {
		topics[sub.topic] = sub.subscriptionId
	}
	return topics
}

func (c *SubscriptionList) ShouldHandle(publication Publish) bool {
	for _, sub := range c.subscriptions {
		if !sub.matcher.Match(publication.Topic) {
			continue
		}

		return true
	}

	for _, sub := range c.subscriptionsOnce {
		if !sub.matcher.Match(publication.Topic) {
			continue
		}

		return true
	}

	return false
}
