package dmqtests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
	dmqspecagents "github.com/sync-toys/DirectMQ/spec/agents"
	testbench "github.com/sync-toys/DirectMQ/spec/test_bench"
)

var _ = Describe("Synchronization", func() {
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
	})

	When("connection is established", func() {
		It("should synchronize all subscriptions", func() {
			defer bench.SnapshotRecording("synchro-on-estab")
			defer bench.Stop("test finished")

			onSubscriptionReceived := make(chan string)

			bench.OnBenchUp(func() {
				log("registering salve subscription handler")
				bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
					onSubscriptionReceived <- subscription.Topic
				})

				log("subscribing to topics without active nodes connection")
				bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
					Topic: "topic1",
				})

				bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
					Topic: "topic2",
				})
			})

			bench.Start()

			log("waiting for first subscription to be received by salve")
			Expect(receiveWithTimeout(1, onSubscriptionReceived)).To(Equal("topic1"))
			log("waiting for second subscription to be received by salve")
			Expect(receiveWithTimeout(1, onSubscriptionReceived)).To(Equal("topic2"))
		})
	})

	When("connection is lost", func() {
		It("should unsubscribe all subscriptions", func() {
			defer bench.SnapshotRecording("synchro-on-lost")

			bench.Start()

			onSubscribeReceived := make(chan struct{}, 10)
			onUnsubscribeReceived := make(chan string, 10)

			log("registering salve subscription and unsubscription handlers")
			bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				onSubscribeReceived <- struct{}{}
			})

			bench.Salve.OnUnsubscribe(func(unsubscribe dmqspecagent.OnUnsubscribeNotification) {
				onUnsubscribeReceived <- unsubscribe.Topic
			})

			log("subscribing to topics from master")
			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "topic1",
			})

			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "topic2",
			})

			log("waiting for first subscription to be received by salve")
			receiveWithTimeout(1, onSubscribeReceived)

			log("waiting for second subscription to be received by salve")
			receiveWithTimeout(1, onSubscribeReceived)

			bench.Stop("testing subscription synchronization")

			log("waiting for first unsubscription to be received by salve")
			Expect(receiveWithTimeout(1, onUnsubscribeReceived)).To((Equal("topic1")))
			log("waiting for second unsubscription to be received by salve")
			Expect(receiveWithTimeout(1, onUnsubscribeReceived)).To(Equal("topic2"))
		})
	})
})
