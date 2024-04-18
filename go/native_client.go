package directmq

import (
	"reflect"

	"github.com/gobwas/glob"
	"github.com/google/uuid"
)

type subscription struct {
	topic   string
	matcher glob.Glob
	handler SubscriptionHandler
}

type subscriptionOnce struct {
	subscriptionId string
	topic          string
	matcher        glob.Glob
	handler        SubscriptionOnceHandler
}

type nativeClient struct {
	global   *GlobalQueue
	clientId string
	ttl      int32

	subscriptions     []subscription
	subscriptionsOnce []subscriptionOnce
}

var _ NativeAPI = (*nativeClient)(nil)
var _ GlobalClient = (*nativeClient)(nil)

func newNativeClient(queue *GlobalQueue, clientId string, ttl int32) *nativeClient {
	return &nativeClient{
		global:   queue,
		clientId: clientId,
		ttl:      ttl,

		subscriptions:     make([]subscription, 0),
		subscriptionsOnce: make([]subscriptionOnce, 0),
	}
}

/* Implementing the Client interface */

func (c *nativeClient) Publish(topic string, mode DeliveryMode, payload []byte) error {
	publishCommand := Publish{
		DataFrame: DataFrame{TTL: c.ttl},
		ClientId:  c.clientId,
		Topic:     topic,
		Mode:      mode,
		Payload:   payload,
	}

	handled := c.HandlePublish(publishCommand)
	if mode == AT_MOST_ONCE && handled {
		return nil
	}

	c.global.Published(publishCommand)
	return nil
}

func (c *nativeClient) Subscribe(topic string, handler SubscriptionHandler) error {
	matcher, err := glob.Compile(topic)
	if err != nil {
		return err
	}

	c.subscriptions = append(c.subscriptions, subscription{
		topic,
		matcher,
		handler,
	})

	subscribe := Subscribe{
		DataFrame: DataFrame{TTL: c.ttl},
		ClientId:  c.clientId,
		Topic:     topic,
	}

	c.global.Subscribed(subscribe)

	return nil
}

func (c *nativeClient) SubscribeOnce(topic string, handler SubscriptionOnceHandler) (subscriptionId string, err error) {
	matcher, err := glob.Compile(topic)
	if err != nil {
		return
	}

	subscriptionId = uuid.New().String()

	c.subscriptionsOnce = append(c.subscriptionsOnce, subscriptionOnce{
		subscriptionId,
		topic,
		matcher,
		handler,
	})

	subscribeOnce := SubscribeOnce{
		DataFrame:      DataFrame{TTL: c.ttl},
		ClientId:       c.clientId,
		Topic:          topic,
		SubscriptionId: subscriptionId,
	}

	c.global.SubscribedOnce(subscribeOnce)

	return
}

func (c *nativeClient) Unsubscribe(topic string, handler SubscriptionHandler) (removedSubscriptions uint, err error) {
	unsubscribe := Unsubscribe{
		DataFrame: DataFrame{TTL: c.ttl},
		ClientId:  c.clientId,
		Topic:     topic,
	}

	c.global.Unsubscribed(unsubscribe)

	for i, sub := range c.subscriptions {
		if sub.topic != topic {
			continue
		}

		registeredHandlerAddress := reflect.ValueOf(sub.handler).Pointer()
		givenHandlerAddress := reflect.ValueOf(handler).Pointer()

		if registeredHandlerAddress != givenHandlerAddress {
			continue
		}

		c.subscriptions = append(c.subscriptions[:i], c.subscriptions[i+1:]...)
		removedSubscriptions++
	}

	return
}

func (c *nativeClient) UnsubscribeOnce(topic string, handler SubscriptionOnceHandler) (removedSubscriptions uint, err error) {
	for i, sub := range c.subscriptionsOnce {
		if sub.topic != topic {
			continue
		}

		registeredHandlerAddress := reflect.ValueOf(sub.handler).Pointer()
		givenHandlerAddress := reflect.ValueOf(handler).Pointer()

		if registeredHandlerAddress != givenHandlerAddress {
			continue
		}

		c.subscriptions = append(c.subscriptions[:i], c.subscriptions[i+1:]...)
		removedSubscriptions++

		unsubscribe := UnsubscribeOnce{
			DataFrame:      DataFrame{TTL: c.ttl},
			ClientId:       c.clientId,
			Topic:          topic,
			SubscriptionId: sub.subscriptionId,
		}

		c.global.UnsubscribedOnce(unsubscribe)
	}

	return
}

func (c *nativeClient) UnsubscribeOnceById(subscriptionId string) (removedSubscriptions uint, err error) {
	for _, sub := range c.subscriptionsOnce {
		if sub.subscriptionId != subscriptionId {
			continue
		}

		remSubs, err := c.UnsubscribeOnce(sub.topic, sub.handler)
		removedSubscriptions += remSubs

		if err != nil {
			return removedSubscriptions, err
		}
	}

	return
}

/* Implementing the GlobalClient interface */

func (c *nativeClient) IsConnected() bool {
	return true
}

func (c *nativeClient) GetClientId() string {
	return c.clientId
}

func (c *nativeClient) GetSubscribedTopics() []string {
	topics := make([]string, 0, len(c.subscriptions))
	for _, sub := range c.subscriptions {
		topics = append(topics, sub.topic)
	}
	return topics
}

func (c *nativeClient) GetSubscribedOnceTopics() map[string]string {
	topics := make(map[string]string, len(c.subscriptionsOnce))
	for _, sub := range c.subscriptionsOnce {
		topics[sub.topic] = sub.subscriptionId
	}
	return topics
}

func (c *nativeClient) HandlePublish(publication Publish) (handled bool) {
	handled = c.handleSubscriptionsOnce(publication)
	if publication.Mode == AT_MOST_ONCE && handled {
		return
	}

	handled = c.handleSubscriptions(publication)
	return
}

func (c *nativeClient) handleSubscriptions(publication Publish) (handled bool) {
	for _, sub := range c.subscriptions {
		if !sub.matcher.Match(publication.Topic) {
			continue
		}

		sub.handler(publication, Subscribe{
			DataFrame: publication.DataFrame,
			ClientId:  c.clientId,
			Topic:     publication.Topic,
		})

		handled = true

		if publication.Mode == AT_MOST_ONCE {
			return
		}
	}

	return
}

func (c *nativeClient) handleSubscriptionsOnce(publication Publish) (handled bool) {
	for _, sub := range c.subscriptionsOnce {
		if !sub.matcher.Match(publication.Topic) {
			continue
		}

		c.UnsubscribeOnce(sub.topic, sub.handler)

		sub.handler(publication, SubscribeOnce{
			DataFrame:      publication.DataFrame,
			ClientId:       c.clientId,
			Topic:          publication.Topic,
			SubscriptionId: sub.subscriptionId,
		})

		handled = true

		if publication.Mode == AT_MOST_ONCE {
			return
		}
	}

	return
}

func (c *nativeClient) HandleSubscribe(subscription Subscribe) {
	/* No-op */
}

func (c *nativeClient) HandleSubscribeOnce(subscription SubscribeOnce) {
	/* No-op */
}

func (c *nativeClient) HandleUnsubscribe(unsubscribe Unsubscribe) {
	/* No-op */
}

func (c *nativeClient) HandleUnsubscribeOnce(unsubscribe UnsubscribeOnce) {
	/* No-op */
}

func (c *nativeClient) HandleUnsubscribeOnceById(subscriptionId string) {
	/* No-op */
}

func (c *nativeClient) HandleLoopDetection(loopDetection LoopDetection) {
	/* No-op */
}

func (c *nativeClient) HandleLoopDetected(loopDetected LoopDetected) {
	/* No-op */
}
