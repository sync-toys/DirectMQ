package dmqtests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
	dmqspecagents "github.com/sync-toys/DirectMQ/spec/agents"
	testbench "github.com/sync-toys/DirectMQ/spec/test_bench"
)

var _ = Describe("Central point topology", func() {
	var bench *testbench.CentralPointTopoTestBench

	BeforeEach(func() {
		bench = testbench.NewGinkgoCentralPointTopoTestBench(testbench.CentralPointTopoTestBenchConfig{
			TopSpawn:     dmqspecagents.GolangAgent("top", dmqspecagents.NO_DEBUGGING),
			LeftSpawn:    dmqspecagents.GolangAgent("left", dmqspecagents.NO_DEBUGGING),
			RightSpawn:   dmqspecagents.GolangAgent("right", dmqspecagents.NO_DEBUGGING),
			CentralSpawn: dmqspecagents.GolangAgent("central", dmqspecagents.NO_DEBUGGING),

			TopTTL:     directmq.DEFAULT_TTL,
			LeftTTL:    directmq.DEFAULT_TTL,
			RightTTL:   directmq.DEFAULT_TTL,
			CentralTTL: directmq.DEFAULT_TTL,

			TopMaxMessageSize:     directmq.NO_MAX_MESSAGE_SIZE,
			LeftMaxMessageSize:    directmq.NO_MAX_MESSAGE_SIZE,
			RightMaxMessageSize:   directmq.NO_MAX_MESSAGE_SIZE,
			CentralMaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,

			LogTopToCentralCommunication:   true,
			LogCentralToTopCommunication:   true,
			LogLeftToCentralCommunication:  true,
			LogCentralToLeftCommunication:  true,
			LogCentralToRightCommunication: true,
			LogRightToCentralCommunication: true,

			LogTopLogs:     true,
			LogLeftLogs:    true,
			LogRightLogs:   true,
			LogCentralLogs: true,

			DisableAllLogs: false,
		})

		bench.Start()
	})

	AfterEach(func() {
		bench.Stop("test ended")
	})

	Context("subscriptions management", func() {
		When("simple subscription", func() {
			It("should propagate subscription to all agents", func() {
				subscriptionPropagatedToLeft := make(chan struct{})
				bench.Left.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagatedToLeft <- struct{}{}
				})

				subscriptionPropagatedToRight := make(chan struct{})
				bench.Right.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagatedToRight <- struct{}{}
				})

				subscriptionPropagatedToCentral := make(chan struct{})
				bench.Central.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagatedToCentral <- struct{}{}
				})

				log("subscribing to test from top")
				bench.Top.Subscribe(dmqspecagent.SubscribeTopicCommand{
					Topic: "test",
				})

				log("waiting for subscription to be propagated to central")
				<-subscriptionPropagatedToCentral

				log("waiting for subscription to be propagated to left")
				<-subscriptionPropagatedToLeft

				log("waiting for subscription to be propagated to right")
				<-subscriptionPropagatedToRight
			})

			It("should propagate unsubscription to all agents", func() {
				unsubscriptionPropagatedToLeft := make(chan struct{})
				bench.Left.OnUnsubscribe(func(notification dmqspecagent.OnUnsubscribeNotification) {
					unsubscriptionPropagatedToLeft <- struct{}{}
				})

				unsubscriptionPropagatedToRight := make(chan struct{})
				bench.Right.OnUnsubscribe(func(notification dmqspecagent.OnUnsubscribeNotification) {
					unsubscriptionPropagatedToRight <- struct{}{}
				})

				unsubscriptionPropagatedToCentral := make(chan struct{})
				bench.Central.OnUnsubscribe(func(notification dmqspecagent.OnUnsubscribeNotification) {
					unsubscriptionPropagatedToCentral <- struct{}{}
				})

				subscriptionPropagated := make(chan struct{}, 3)
				bench.Left.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagated <- struct{}{}
				})

				bench.Central.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagated <- struct{}{}
				})

				bench.Right.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagated <- struct{}{}
				})

				log("subscribing to test from top")
				subscriptionID := bench.Top.Subscribe(dmqspecagent.SubscribeTopicCommand{
					Topic: "test",
				})

				log("waiting for subscription to be propagated to all agents")
				for i := 0; i < 3; i++ {
					<-subscriptionPropagated
				}

				log("unsubscribing from test from top")
				bench.Top.Unsubscribe(dmqspecagent.UnsubscribeTopicCommand{
					SubscriptionID: subscriptionID,
				})

				log("waiting for unsubscription to be propagated to central")
				<-unsubscriptionPropagatedToCentral

				log("waiting for unsubscription to be propagated to left")
				<-unsubscriptionPropagatedToLeft

				log("waiting for unsubscription to be propagated to right")
				<-unsubscriptionPropagatedToRight
			})
		})

		When("complex subscriptions", func() {
			var subscriptionPropagatedToTop chan string
			var unsubscriptionPropagatedToTop chan string

			var subscriptionPropagatedToLeft chan string
			var unsubscriptionPropagatedToLeft chan string

			var subscriptionPropagatedToRight chan string
			var unsubscriptionPropagatedToRight chan string

			BeforeEach(func() {
				log("preparing left node")

				subscriptionPropagatedToTop = make(chan string, 1)
				unsubscriptionPropagatedToTop = make(chan string, 1)

				bench.Top.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagatedToTop <- notification.Topic
				})

				bench.Top.OnUnsubscribe(func(notification dmqspecagent.OnUnsubscribeNotification) {
					unsubscriptionPropagatedToTop <- notification.Topic
				})

				log("preparing right node")

				subscriptionPropagatedToLeft = make(chan string, 1)
				unsubscriptionPropagatedToLeft = make(chan string, 1)

				bench.Left.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagatedToLeft <- notification.Topic
				})

				bench.Left.OnUnsubscribe(func(notification dmqspecagent.OnUnsubscribeNotification) {
					unsubscriptionPropagatedToLeft <- notification.Topic
				})

				log("preparing right node")

				subscriptionPropagatedToRight = make(chan string, 1)
				unsubscriptionPropagatedToRight = make(chan string, 1)

				bench.Right.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
					subscriptionPropagatedToRight <- notification.Topic
				})

				bench.Right.OnUnsubscribe(func(notification dmqspecagent.OnUnsubscribeNotification) {
					unsubscriptionPropagatedToRight <- notification.Topic
				})
			})

			It("should manage complex subscriptions structure", func() {
				defer bench.SnapshotRecordings("cpt/complex/%s")

				log("subscribing to test/1 from left")
				bench.Left.Subscribe(dmqspecagent.SubscribeTopicCommand{
					Topic: "test/1",
				})

				log("waiting for subscription to be propagated to left")
				Expect(<-subscriptionPropagatedToLeft).To(Equal("test/1"))

				log("waiting for subscription to be propagated to top")
				Expect(<-subscriptionPropagatedToTop).To(Equal("test/1"))

				log("waiting for subscription to be propagated to right")
				Expect(<-subscriptionPropagatedToRight).To(Equal("test/1"))

				log("subscribing to test/* from top")
				bench.Top.Subscribe(dmqspecagent.SubscribeTopicCommand{
					Topic: "test/*",
				})

				log("waiting for subscription to be propagated to left")
				Expect(<-subscriptionPropagatedToLeft).To(Equal("test/*"))

				log("waiting for subscription to be propagated to right")
				Expect(<-subscriptionPropagatedToRight).To(Equal("test/*"))

				log("waiting for unsubscribe from test/1 sent to right")
				Expect(<-unsubscriptionPropagatedToRight).To(Equal("test/1"))
			})
		})
	})
})
