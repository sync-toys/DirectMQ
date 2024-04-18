package directmq

type SubscriptionHandler func(publish Publish, subscription Subscribe)

type SubscriptionOnceHandler func(publish Publish, subscription SubscribeOnce)

type PublicationHandler interface {
	HandlePublish(publication Publish)
}
