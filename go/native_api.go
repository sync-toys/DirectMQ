package directmq

type NativeAPI interface {
	Publish(topic string, mode DeliveryMode, payload []byte) error
	Subscribe(topic string, handler SubscriptionHandler) error
	SubscribeOnce(topic string, handler SubscriptionOnceHandler) (subscriptionId string, err error)
	Unsubscribe(topic string, handler SubscriptionHandler) (removedSubscriptions uint, err error)
	UnsubscribeOnce(topic string, handler SubscriptionOnceHandler) (removedSubscriptions uint, err error)
	UnsubscribeOnceById(subscriptionId string) (removedSubscriptions uint, err error)
}
