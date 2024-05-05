package directmq

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionList", func() {
	noopHandler := func() {}

	Context("when creating a new subscription list", func() {
		It("should be empty", func() {
			subscriptions := newSubscriptionList[func()]()
			Expect(subscriptions.GetSubscriptions()).To(HaveLen(0))
		})
	})

	Context("when adding a subscription", func() {
		It("should return a subscription ID", func() {
			subscriptions := newSubscriptionList[func()]()
			id := subscriptions.AddSubscription("topic", &noopHandler)
			Expect(id).ToNot(BeNil())
		})

		It("should add a subscription to the list", func() {
			subscriptions := newSubscriptionList[func()]()
			subscriptions.AddSubscription("topic", &noopHandler)
			Expect(subscriptions.GetSubscriptions()).To(HaveLen(1))
		})

		It("should panic if the handler is nil", func() {
			subscriptions := newSubscriptionList[func()]()
			Expect(func() {
				subscriptions.AddSubscription("topic", nil)
			}).To(Panic())
		})

		It("should panic if the topic pattern is incorrect", func() {
			subscriptions := newSubscriptionList[func()]()
			Expect(func() {
				subscriptions.AddSubscription("topic!", &noopHandler)
			}).To(Panic())
		})

		It("should return a different ID for each subscription", func() {
			subscriptions := newSubscriptionList[func()]()
			id1 := subscriptions.AddSubscription("topic", &noopHandler)
			id2 := subscriptions.AddSubscription("topic", &noopHandler)
			Expect(id1).ToNot(Equal(id2))
		})
	})

	Context("when removing a subscription", func() {
		It("should remove the subscription from the list", func() {
			subscriptions := newSubscriptionList[func()]()
			id := subscriptions.AddSubscription("topic", &noopHandler)
			subscriptions.RemoveSubscription(id)
			Expect(subscriptions.GetSubscriptions()).To(HaveLen(0))
		})

		It("should not remove other subscriptions", func() {
			subscriptions := newSubscriptionList[func()]()
			id1 := subscriptions.AddSubscription("topic1", &noopHandler)
			id2 := subscriptions.AddSubscription("topic2", &noopHandler)
			subscriptions.RemoveSubscription(id1)
			Expect(subscriptions.GetSubscriptions()).To(HaveLen(1))
			Expect(subscriptions.GetSubscriptions()[0].ID).To(Equal(id2))
		})
	})

	Context("when getting triggered subscriptions", func() {
		It("should return a list of triggered subscriptions", func() {
			subscriptions := newSubscriptionList[func()]()
			id1 := subscriptions.AddSubscription("topic1", &noopHandler)
			subscriptions.AddSubscription("topic2", &noopHandler)
			triggeredSubscriptions := subscriptions.GetTriggeredSubscriptions("topic1")
			Expect(triggeredSubscriptions).To(HaveLen(1))
			Expect(triggeredSubscriptions[0].ID).To(Equal(id1))
		})

		It("should return an empty list if no subscriptions are triggered", func() {
			subscriptions := newSubscriptionList[func()]()
			subscriptions.AddSubscription("topic1", &noopHandler)
			subscriptions.AddSubscription("topic2", &noopHandler)
			triggeredSubscriptions := subscriptions.GetTriggeredSubscriptions("topic3")
			Expect(triggeredSubscriptions).To(HaveLen(0))
		})

		It("should return a list of multiple triggered subscriptions", func() {
			subscriptions := newSubscriptionList[func()]()
			id1 := subscriptions.AddSubscription("topic1", &noopHandler)
			id2 := subscriptions.AddSubscription("topic1", &noopHandler)
			triggeredSubscriptions := subscriptions.GetTriggeredSubscriptions("topic1")
			Expect(triggeredSubscriptions).To(HaveLen(2))
			Expect(triggeredSubscriptions[0].ID).To(Equal(id1))
			Expect(triggeredSubscriptions[1].ID).To(Equal(id2))
		})
	})

	Context("when getting all subscriptions", func() {
		It("should return all subscriptions", func() {
			subscriptions := newSubscriptionList[func()]()
			id1 := subscriptions.AddSubscription("topic1", &noopHandler)
			id2 := subscriptions.AddSubscription("topic2", &noopHandler)
			allSubscriptions := subscriptions.GetSubscriptions()
			Expect(allSubscriptions).To(HaveLen(2))
			Expect(allSubscriptions[0].ID).To(Equal(id1))
			Expect(allSubscriptions[1].ID).To(Equal(id2))
		})
	})

	Context("when getting list of unique subscribed topics", func() {
		It("should return a list of unique topics", func() {
			subscriptions := newSubscriptionList[func()]()
			subscriptions.AddSubscription("topic1", &noopHandler)
			subscriptions.AddSubscription("topic2", &noopHandler)
			subscriptions.AddSubscription("topic1", &noopHandler)
			topics := subscriptions.GetUniqueSubscribedTopics()
			Expect(topics).To(HaveLen(2))
			Expect(topics).To(ConsistOf("topic1", "topic2"))
		})
	})

	Context("when getting only top level subscribed topics", func() {
		It("should return deduplicated list of only top level topics", func() {
			subscriptions := newSubscriptionList[func()]()
			subscriptions.AddSubscription("topic1/*", &noopHandler)
			subscriptions.AddSubscription("topic1/level", &noopHandler)
			subscriptions.AddSubscription("topic2", &noopHandler)
			topics := subscriptions.GetOnlyTopLevelSubscribedTopics()
			Expect(topics).To(HaveLen(2))
			Expect(topics).To(ConsistOf("topic1/*", "topic2"))
		})

		It("should return correct result in case of wildcard and super wildcard mixing", func() {
			subscriptions := newSubscriptionList[func()]()
			subscriptions.AddSubscription("topic1/*/something/**", &noopHandler)
			subscriptions.AddSubscription("topic1/**", &noopHandler)
			subscriptions.AddSubscription("topic1/level", &noopHandler)
			topics := subscriptions.GetOnlyTopLevelSubscribedTopics()
			Expect(topics).To(HaveLen(1))
			Expect(topics).To(ConsistOf("topic1/**"))
		})
	})
})
