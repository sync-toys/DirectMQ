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

	When("agent A publishes a message", func() {
		It("should be received by agent B", func() {
			bench := testbench.NewGinkgoPairTopoTestBench(
				testbench.PairTopoTestBenchConfig{
					MasterSpawn: dmqspecagents.GolangAgent("master", dmqspecagents.NO_DEBUGGING),
					SalveSpawn:  dmqspecagents.GolangAgent("salve", dmqspecagents.NO_DEBUGGING),

					MasterTTL:            32,
					MasterMaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,

					SalveTTL:            32,
					SalveMaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,

					LogMasterToSalveCommunication: true,
					LogSalveToMasterCommunication: true,

					LogMasterLogs: true,
					LogSalveLogs:  true,

					DisableAllLogs: false,
				},
			)

			bench.Start()
			defer bench.Stop("test ended")
			defer bench.SnapshotRecording("pubsub")

			subscriptionPropagated := make(chan struct{}, 1)
			bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagated <- struct{}{}
			})

			log("subscribing to test topic from salve")
			bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to master")
			<-subscriptionPropagated

			result := make(chan []byte)
			bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				result <- message.Payload
			})

			log("publishing message from master on test topic")
			bench.Master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			log("waiting for message to be received by salve")
			Expect(<-result).To(Equal([]byte("Hello, World!")))
		})
	})
})
