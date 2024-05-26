package dmqtests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
	dmqspecagents "github.com/sync-toys/DirectMQ/spec/agents"
	testbench "github.com/sync-toys/DirectMQ/spec/test_bench"
)

var _ = Describe("Subscriptions hierarchy", func() {
	type propagation struct {
		Topic            string
		IsUnsubscription bool
	}

	var bench *testbench.PairTopoTestBench

	BeforeEach(func() {
		bench = testbench.NewGinkgoPairTopoTestBench(
			testbench.PairTopoTestBenchConfig{
				MasterSpawn: dmqspecagents.GolangAgent("master", dmqspecagents.NO_DEBUGGING),
				SalveSpawn:  dmqspecagents.GolangAgent("salve", dmqspecagents.NO_DEBUGGING),

				MasterTTL:            directmq.DEFAULT_TTL,
				MasterMaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,

				SalveTTL:            directmq.DEFAULT_TTL,
				SalveMaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,

				LogMasterToSalveCommunication: true,
				LogSalveToMasterCommunication: true,

				LogMasterLogs: true,
				LogSalveLogs:  true,

				DisableAllLogs: false,
			},
		)

		bench.Start()
	})

	AfterEach(func() {
		bench.Stop("test ended")
	})

	When("subscribing to a higher level topic", func() {
		It("should override lower level subscriptions", func() {
			defer bench.SnapshotRecording("hierarchy-override-lls")

			/////////// phase 0 - setup

			propagated := make(chan propagation)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				if bench.Stopped() {
					return
				}

				propagated <- propagation{Topic: subscription.Topic, IsUnsubscription: false}
			})
			bench.Master.OnUnsubscribe(func(unsub dmqspecagent.OnUnsubscribeNotification) {
				if bench.Stopped() {
					return
				}

				propagated <- propagation{Topic: unsub.Topic, IsUnsubscription: true}
			})

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "a/b/c",
			})
			<-propagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "a/b/d",
			})
			<-propagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "x/y/z",
			})
			<-propagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "x/y/z/0",
			})
			<-propagated

			/////////// phase 1

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "x/y/**",
			})
			xySubscription := <-propagated
			xyzUnsubscription := <-propagated
			xyz0Unsubscription := <-propagated

			/////////// phase 2

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "a/b/*",
			})
			abSubscription := <-propagated
			abdUnsubscription := <-propagated
			abcUnsubscription := <-propagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "x/**",
			})
			xSubscription := <-propagated
			xyUnsubscription := <-propagated

			/////////// results checking phase

			Expect(xySubscription).To(Equal(propagation{Topic: "x/y/**", IsUnsubscription: false}))
			Expect(xyzUnsubscription).To(Equal(propagation{Topic: "x/y/z", IsUnsubscription: true}))
			Expect(xyz0Unsubscription).To(Equal(propagation{Topic: "x/y/z/0", IsUnsubscription: true}))

			Expect(abSubscription).To(Equal(propagation{Topic: "a/b/*", IsUnsubscription: false}))
			Expect(abcUnsubscription).To(Equal(propagation{Topic: "a/b/c", IsUnsubscription: true}))
			Expect(abdUnsubscription).To(Equal(propagation{Topic: "a/b/d", IsUnsubscription: true}))

			Expect(xSubscription).To(Equal(propagation{Topic: "x/**", IsUnsubscription: false}))
			Expect(xyUnsubscription).To(Equal(propagation{Topic: "x/y/**", IsUnsubscription: true}))
		})

		It("should not propagate lower level subscriptions on subscribe", func() {
			defer bench.SnapshotRecording("hierarchy-override-lls-x")

			controlReceived := make(chan struct{})
			bench.Master.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				if bench.Stopped() {
					return
				}

				if string(message.Payload) == "synchronization" {
					controlReceived <- struct{}{}
				}
			})

			subscriptionPropagated := make(chan struct{})
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				if bench.Stopped() {
					return
				}

				if subscription.Topic != "internal/control" {
					subscriptionPropagated <- struct{}{}
				}
			})

			controlSubscriptionPropagated := make(chan struct{})
			bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				if bench.Stopped() {
					return
				}

				if subscription.Topic == "internal/control" {
					controlSubscriptionPropagated <- struct{}{}
				}
			})

			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "internal/control",
			})
			<-controlSubscriptionPropagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test/**",
			})
			<-subscriptionPropagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test/a",
			})

			bench.Salve.Publish(dmqspecagent.PublishCommand{
				Topic:            "internal/control",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("synchronization"),
			})

			select {
			case <-controlReceived:
				break
			case <-time.After(1 * time.Second):
				Fail("control message not received")
			case <-subscriptionPropagated:
				Fail("received lower level subscription")
			}
		})
	})

	When("unsubscribing from a higher level topic", func() {
		It("should subscribe to lower level subscriptions", func() {
			defer bench.SnapshotRecording("hierarchy-override-uls")

			/////////// phase 0 - setup

			propagated := make(chan propagation)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				if bench.Stopped() {
					return
				}

				propagated <- propagation{Topic: subscription.Topic, IsUnsubscription: false}
			})
			bench.Master.OnUnsubscribe(func(unsub dmqspecagent.OnUnsubscribeNotification) {
				if bench.Stopped() {
					return
				}

				propagated <- propagation{Topic: unsub.Topic, IsUnsubscription: true}
			})

			topLevelSubscriptionID := bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "a/b/*",
			})
			<-propagated

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "a/b/c",
			})
			// no propagation, covered by a/b/*

			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "a/b/d",
			})
			// no propagation, covered by a/b/*

			/////////// phase 1

			bench.Salve.Unsubscribe(dmqspecagent.UnsubscribeTopicCommand{
				SubscriptionID: topLevelSubscriptionID,
			})
			abcSubscription := <-propagated
			abdSubscription := <-propagated
			abUnsubscription := <-propagated

			/////////// results checking phase

			Expect(abcSubscription).To(Equal(propagation{Topic: "a/b/c", IsUnsubscription: false}))
			Expect(abdSubscription).To(Equal(propagation{Topic: "a/b/d", IsUnsubscription: false}))
			Expect(abUnsubscription).To(Equal(propagation{Topic: "a/b/*", IsUnsubscription: true}))
		})
	})
})
