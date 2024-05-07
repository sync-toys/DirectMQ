package dmqtests

import (
	"io"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqspec "github.com/sync-toys/DirectMQ/spec"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
	dmqspecagents "github.com/sync-toys/DirectMQ/spec/agents"
	testbench "github.com/sync-toys/DirectMQ/spec/test_bench"
)

var _ = Describe("Pubsub", func() {
	log := dmqspec.CreateLogger("test", func(message string, writer io.Writer) {
		color.New(color.FgGreen).Fprintln(writer, message)
	})

	Context("basic communication", func() {
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
				defer bench.SnapshotRecording("pubsub-mtos")

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
				defer bench.SnapshotRecording("pubsub-mtom")

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
				defer bench.SnapshotRecording("pubsub-all")

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
	})
})
