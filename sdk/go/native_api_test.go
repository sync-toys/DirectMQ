package directmq

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("nativeAPI", func() {
	var networkConfig NetworkNodeConfig = NetworkNodeConfig{
		HostTTL:                    DEFAULT_TTL,
		HostMaxIncomingMessageSize: NO_MAX_MESSAGE_SIZE,
		HostID:                     "host",
	}

	var node *networkNode
	noopHandler := func([]byte) {}

	getInitialMessageDataFrame := func() DataFrame {
		return DataFrame{
			TTL:       int32(networkConfig.HostTTL),
			Traversed: []string{networkConfig.HostID},
		}
	}

	BeforeEach(func() {
		node = newNetworkNode(networkConfig, NewProtobufJSONProtocol())
	})

	Context("when creating a new native API", func() {
		It("should be empty", func() {
			Expect(node.api.GetSubscribedTopics()).To(HaveLen(0))
		})
	})

	Context("when adding a subscription", func() {
		It("should return a subscription ID", func() {
			id := node.api.Subscribe("topic", noopHandler)
			Expect(id).ToNot(BeNil())
		})

		It("should add a subscription to the list", func() {
			node.api.Subscribe("topic", noopHandler)
			Expect(node.api.GetSubscribedTopics()).To(HaveLen(1))
		})

		It("should panic if the handler is nil", func() {
			Expect(func() {
				node.api.Subscribe("topic", nil)
			}).To(Panic())
		})

		It("should panic if the topic pattern is incorrect", func() {
			Expect(func() {
				node.api.Subscribe("topic!", noopHandler)
			}).To(Panic())
		})

		It("should return a different ID for each subscription", func() {
			id1 := node.api.Subscribe("topic", noopHandler)
			id2 := node.api.Subscribe("topic", noopHandler)
			Expect(id1).ToNot(Equal(id2))
		})

		It("should update the global network", func() {
			result := make(chan SubscribeMessage)
			node.OnSubscription(func(message SubscribeMessage) {
				result <- message
			})

			go node.api.Subscribe("topic", noopHandler)

			Eventually(result).Should(Receive(Equal(SubscribeMessage{
				DataFrame: getInitialMessageDataFrame(),
				Topic:     "topic",
			})))
		})

		Context("in case of complex subscription list", func() {
			It("should replace the existing subscription with new higher level subscription", func() {
				node.api.Subscribe("topic1", noopHandler)
				node.api.Subscribe("topic1/subtopic1", noopHandler)
				node.api.Subscribe("topic1/subtopic2", noopHandler)
				node.api.Subscribe("topic2/subtopic", noopHandler)

				subscriptionMessage := make(chan SubscribeMessage)
				node.OnSubscription(func(message SubscribeMessage) {
					subscriptionMessage <- message
				})

				unsubscribeMessage := make(chan UnsubscribeMessage)
				node.OnUnsubscribe(func(message UnsubscribeMessage) {
					unsubscribeMessage <- message
				})

				go node.api.Subscribe("topic1/*", noopHandler)

				Eventually(subscriptionMessage).Should(Receive(Equal(SubscribeMessage{
					DataFrame: getInitialMessageDataFrame(),
					Topic:     "topic1/*",
				})))

				Eventually(unsubscribeMessage).Should(Receive(Equal(UnsubscribeMessage{
					DataFrame: getInitialMessageDataFrame(),
					Topic:     "topic1/subtopic1",
				})))

				Eventually(unsubscribeMessage).Should(Receive(Equal(UnsubscribeMessage{
					DataFrame: getInitialMessageDataFrame(),
					Topic:     "topic1/subtopic2",
				})))
			})
		})
	})

	Context("when removing a subscription", func() {
		It("should remove the subscription from the list", func() {
			id := node.api.Subscribe("topic", noopHandler)
			node.api.Unsubscribe(id)
			Expect(node.api.GetSubscribedTopics()).To(HaveLen(0))
		})

		It("should not remove other subscriptions", func() {
			id1 := node.api.Subscribe("topic1", noopHandler)
			node.api.Subscribe("topic2", noopHandler)
			node.api.Unsubscribe(id1)
			Expect(node.api.GetSubscribedTopics()).To(HaveLen(1))
			Expect(node.api.GetSubscribedTopics()[0]).To(Equal("topic2"))
		})

		It("should update the global network", func() {
			result := make(chan UnsubscribeMessage)
			node.OnUnsubscribe(func(message UnsubscribeMessage) {
				result <- message
			})

			id := node.api.Subscribe("topic", noopHandler)
			go node.api.Unsubscribe(id)

			Eventually(result).Should(Receive(Equal(UnsubscribeMessage{
				DataFrame: getInitialMessageDataFrame(),
				Topic:     "topic",
			})))
		})

		It("should not panic if the subscription ID is incorrect", func() {
			node.api.Unsubscribe(0)
		})

		Context("in case of complex subscription list", func() {
			It("should replace the higher level subscription with existing lower level subscriptions", func() {
				topLevelSubID := node.api.Subscribe("topic1/*", noopHandler)
				node.api.Subscribe("topic1/subtopic1", noopHandler)
				node.api.Subscribe("topic1/subtopic2", noopHandler)
				node.api.Subscribe("topic2/subtopic", noopHandler)

				subscriptionMessage := make(chan SubscribeMessage)
				node.OnSubscription(func(message SubscribeMessage) {
					subscriptionMessage <- message
				})

				unsubscribeMessage := make(chan UnsubscribeMessage)
				node.OnUnsubscribe(func(message UnsubscribeMessage) {
					unsubscribeMessage <- message
				})

				go node.api.Unsubscribe(topLevelSubID)

				Eventually(subscriptionMessage).Should(Receive(Equal(SubscribeMessage{
					DataFrame: getInitialMessageDataFrame(),
					Topic:     "topic1/subtopic1",
				})))

				Eventually(subscriptionMessage).Should(Receive(Equal(SubscribeMessage{
					DataFrame: getInitialMessageDataFrame(),
					Topic:     "topic1/subtopic2",
				})))

				Eventually(unsubscribeMessage).Should(Receive(Equal(UnsubscribeMessage{
					DataFrame: getInitialMessageDataFrame(),
					Topic:     "topic1/*",
				})))
			})
		})
	})

	Context("when publishing a message", func() {
		It("should panic if the topic pattern is incorrect", func() {
			Expect(func() {
				node.api.Publish("topic!", []byte{0}, AT_MOST_ONCE)
			}).To(Panic())
		})

		It("should panic if the payload is empty", func() {
			Expect(func() {
				node.api.Publish("topic", []byte{}, AT_MOST_ONCE)
			}).To(Panic())
		})

		It("should update the global network", func() {
			result := make(chan PublishMessage, 2)
			node.OnPublication(func(message PublishMessage) {
				result <- message
			})

			go node.api.Publish("topic", []byte{0}, AT_MOST_ONCE)

			Eventually(result).Should(Receive(Equal(PublishMessage{
				DataFrame:        getInitialMessageDataFrame(),
				Topic:            "topic",
				Payload:          []byte{0},
				DeliveryStrategy: AT_MOST_ONCE,
			})))
		})
	})
})
