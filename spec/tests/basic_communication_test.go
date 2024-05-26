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

var _ = Describe("Basic communication", func() {
	When("master publishes a message on simple topic", func() {
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

		It("should be received by salve subscription", func() {
			defer bench.SnapshotRecording("bc-ps-mtos")

			subscriptionPropagatedToMaster := make(chan struct{}, 1)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToMaster <- struct{}{}
			})

			log("subscribing to test topic from salve")
			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to master")
			<-subscriptionPropagatedToMaster

			messageReceivedBySalve := make(chan []byte)
			bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedBySalve <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			log("waiting for message to be received by salve")
			Expect(<-messageReceivedBySalve).To(Equal([]byte("Hello, World!")))
		})

		It("should be received by master subscription", func() {
			defer bench.SnapshotRecording("bc-ps-mtom")

			subscriptionPropagatedToSalve := make(chan struct{}, 1)
			bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToSalve <- struct{}{}
			})

			log("subscribing to test topic from master")
			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to salve")
			<-subscriptionPropagatedToSalve

			messageReceivedByMaster := make(chan []byte)
			bench.Master.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedByMaster <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			log("waiting for message to be received by master")
			Expect(<-messageReceivedByMaster).To(Equal([]byte("Hello, World!")))
		})

		It("should be received all subscribers", func() {
			defer bench.SnapshotRecording("bc-ps-all")

			subscriptionPropagatedToSalve := make(chan struct{}, 1)
			bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToSalve <- struct{}{}
			})

			log("subscribing to test topic from master")
			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to salve")
			<-subscriptionPropagatedToSalve

			subscriptionPropagatedToMaster := make(chan struct{}, 1)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToMaster <- struct{}{}
			})

			log("subscribing to test topic from salve")
			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to master")
			<-subscriptionPropagatedToMaster

			messageReceivedByMaster := make(chan []byte)
			bench.Master.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedByMaster <- message.Payload
			})

			messageReceivedBySalve := make(chan []byte)
			bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedBySalve <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World! (master)"),
			})

			log("waiting for message to be received by master")
			Expect(<-messageReceivedByMaster).To(Equal([]byte("Hello, World! (master)")))

			log("waiting for message to be received by salve")
			Expect(<-messageReceivedBySalve).To(Equal([]byte("Hello, World! (master)")))

			log("publishing message from salve on test topic")
			bench.Salve.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World! (salve)"),
			})

			log("waiting for message to be received by salve")
			Expect(<-messageReceivedBySalve).To(Equal([]byte("Hello, World! (salve)")))

			log("waiting for message to be received by master")
			Expect(<-messageReceivedByMaster).To(Equal([]byte("Hello, World! (salve)")))
		})
	})

	When("master publishes a message with at most once delivery strategy", func() {
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

		It("should contain unchanged payload", func() {
			defer bench.SnapshotRecording("bc-mq-payload")

			subscriptionPropagatedToMaster := make(chan struct{}, 1)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToMaster <- struct{}{}
			})

			log("subscribing to test topic from salve")
			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to master")
			<-subscriptionPropagatedToMaster

			messageReceivedBySalve := make(chan []byte)
			bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedBySalve <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_MOST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			log("waiting for message to be received by salve")
			Expect(<-messageReceivedBySalve).To(Equal([]byte("Hello, World!")))
		})

		It("should be received only by a single subscriber", func() {
			// we are not making snapshots for this test
			// as it is not deterministic in terms of
			// which subscriber will receive the message

			subscriptionPropagatedToSalve := make(chan struct{}, 1)
			bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToSalve <- struct{}{}
			})

			log("subscribing to test topic from master")
			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to salve")
			<-subscriptionPropagatedToSalve

			subscriptionPropagatedToMaster := make(chan struct{}, 1)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToMaster <- struct{}{}
			})

			log("subscribing to test topic from salve")
			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to master")
			<-subscriptionPropagatedToMaster

			messageReceivedByMaster := make(chan []byte)
			bench.Master.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedByMaster <- message.Payload
			})

			messageReceivedBySalve := make(chan []byte)
			bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedBySalve <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_MOST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			log("waiting for message to be received by any subscriber")
			select {
			case <-messageReceivedByMaster:
			case <-messageReceivedBySalve:
				break

			case <-time.After(5 * time.Second):
				Fail("no subscriber received the message after 5 seconds")
			}

			log("waiting 1 second to ensure no other subscriber receives the message")
			select {
			case <-messageReceivedByMaster:
			case <-messageReceivedBySalve:
				Fail("both subscribers received the message")

			case <-time.After(1 * time.Second):
				break
			}
		})
	})

	When("master publishes a message on topic already unsubscribed by salve", func() {
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

		It("should not be received by salve", func() {
			defer bench.SnapshotRecording("bc-unsub")

			subscriptionPropagatedToMaster := make(chan struct{}, 1)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToMaster <- struct{}{}
			})

			log("subscribing to test topic from salve")
			subscriptionID := bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to master")
			<-subscriptionPropagatedToMaster

			unsubscriptionPropagatedToMaster := make(chan struct{}, 1)
			bench.Master.OnUnsubscribe(func(unsubscribe dmqspecagent.OnUnsubscribeNotification) {
				unsubscriptionPropagatedToMaster <- struct{}{}
			})

			log("unsubscribing from test topic from salve")
			bench.Salve.Unsubscribe(dmqspecagent.UnsubscribeTopicCommand{
				SubscriptionID: subscriptionID,
			})

			log("waiting for unsubscription to be propagated to master")
			<-unsubscriptionPropagatedToMaster

			messageReceivedBySalve := make(chan []byte)
			bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedBySalve <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			log("waiting 1 second to ensure no message is received by salve")
			select {
			case <-messageReceivedBySalve:
				Fail("salve received the message")

			case <-time.After(1 * time.Second):
				break
			}
		})
	})

	When("salve publishes a binary message", func() {
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

		It("should be received correctly by master", func() {
			defer bench.SnapshotRecording("bc-binary")

			subscriptionPropagatedToSalve := make(chan struct{}, 1)
			bench.Salve.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToSalve <- struct{}{}
			})

			log("subscribing to test topic from master")
			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to salve")
			<-subscriptionPropagatedToSalve

			messageReceivedByMaster := make(chan []byte)
			bench.Master.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				messageReceivedByMaster <- message.Payload
			})

			log("publishing binary message from salve")
			bench.Salve.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte{0x00, 0x01, 0x02, 0x03},
			})

			log("waiting for message to be received by master")
			Expect(<-messageReceivedByMaster).To(Equal([]byte{0x00, 0x01, 0x02, 0x03}))
		})
	})
})

// TODO: should not unsubscribe if there is any listener to given topic in network node
