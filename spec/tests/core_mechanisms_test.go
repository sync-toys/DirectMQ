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

var _ = Describe("Core mechanisms", func() {
	Context("TTL", func() {
		It("should not forward message if TTL exceeded", func() {
			bench := testbench.NewGinkgoLineTopoTestBench(testbench.LineTopoTestBenchConfig{
				LeftSpawn:   dmqspecagents.GolangAgent("left", dmqspecagents.NO_DEBUGGING),
				MiddleSpawn: dmqspecagents.GolangAgent("middle", dmqspecagents.NO_DEBUGGING),
				RightSpawn:  dmqspecagents.GolangAgent("right", dmqspecagents.NO_DEBUGGING),

				LeftTTL:   3,
				MiddleTTL: 2,
				RightTTL:  2,

				LogLeftToMiddleCommunication:  true,
				LogMiddleToRightCommunication: true,
				LogRightToMiddleCommunication: true,
				LogMiddleToLeftCommunication:  true,

				LogLeftLogs:   true,
				LogMiddleLogs: true,
				LogRightLogs:  true,

				DisableAllLogs: false,
			})

			bench.Start()

			subscriptionPropagatedToRight := make(chan struct{})
			bench.Right.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagatedToRight <- struct{}{}
			})

			log("subscribing to test from left")
			bench.Left.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to right")
			<-subscriptionPropagatedToRight

			log("subscribing to test from middle")
			bench.Middle.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("not waiting for subscription to be propagated to right, already propagated")

			messageReceivedOn := make(chan string)
			bench.Left.OnMessageReceived(func(notification dmqspecagent.MessageReceivedNotification) {
				messageReceivedOn <- "left"
			})

			bench.Middle.OnMessageReceived(func(notification dmqspecagent.MessageReceivedNotification) {
				messageReceivedOn <- "middle"
			})

			log("publishing to test from right")
			bench.Right.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				Payload:          []byte("test"),
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
			})

			log("waiting for message to be received by middle")
			Expect(<-messageReceivedOn).To(Equal("middle"))

			log("waiting one second to ensure that message is not received by left")
			time.Sleep(time.Second)

			log("expecting that message is not received by left")
			Expect(messageReceivedOn).To(BeEmpty())
		})
	})

	Context("Max message size", func() {
		var bench *testbench.PairTopoTestBench
		var messageReceived chan struct{} = make(chan struct{})

		prepareBench := func(masterMaxMessageSize, salveMaxMessageSize uint64) {
			bench = testbench.NewGinkgoPairTopoTestBench(testbench.PairTopoTestBenchConfig{
				MasterSpawn: dmqspecagents.GolangAgent("master", dmqspecagents.NO_DEBUGGING),
				SalveSpawn:  dmqspecagents.GolangAgent("salve", dmqspecagents.NO_DEBUGGING),

				MasterTTL:            directmq.DEFAULT_TTL,
				MasterMaxMessageSize: masterMaxMessageSize,
				LogMasterLogs:        true,

				SalveTTL:            directmq.DEFAULT_TTL,
				SalveMaxMessageSize: salveMaxMessageSize,
				LogSalveLogs:        true,

				LogMasterToSalveCommunication: true,
				LogSalveToMasterCommunication: true,

				DisableAllLogs: false,
			})

			bench.Start()

			bench.Master.OnMessageReceived(func(notification dmqspecagent.MessageReceivedNotification) {
				log("MESSAGE RECEIVED")
				messageReceived <- struct{}{}
			})

			subscriptionPropagated := make(chan struct{})
			bench.Salve.OnSubscription(func(notification dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagated <- struct{}{}
			})

			log("subscribing to test topic from master")
			bench.Master.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for subscription to be propagated to salve")
			<-subscriptionPropagated
		}

		AfterEach(func() {
			bench.Stop("test ended")
		})

		It("should forward message if there is no message size limits", func() {
			prepareBench(
				directmq.NO_MAX_MESSAGE_SIZE,
				directmq.NO_MAX_MESSAGE_SIZE,
			)

			log("publishing to test topic from salve")
			bench.Salve.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				Payload:          []byte("test"),
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
			})

			log("waiting for message to be received by master")

			select {
			case <-messageReceived:
				log("message received")
			case <-time.After(5 * time.Second):
				Fail("message not received")
			}
		})

		It("should forward message if message size is within limits", func() {
			prepareBench(
				5, // <- master max message size
				directmq.NO_MAX_MESSAGE_SIZE,
			)

			log("publishing to test topic from salve")
			bench.Salve.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				Payload:          []byte("test"),
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
			})

			log("waiting for message to be received by master")

			select {
			case <-messageReceived:
				log("message received")
				return
			case <-time.After(5 * time.Second):
				Fail("message not received")
			}
		})

		It("should not forward message if message size exceeded", func() {
			prepareBench(
				3, // <- master max message size
				directmq.NO_MAX_MESSAGE_SIZE,
			)

			log("publishing to test topic from salve")
			bench.Salve.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				Payload:          []byte("test"),
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
			})

			log("waiting for message to be received by master")

			select {
			case <-messageReceived:
				Fail("message received")
			case <-time.After(1 * time.Second):
				log("message not received")
				return
			}
		})
	})

	Context("Loop detection", func() {
		It("should terminate network if loop has been detected", func() {
			bench := testbench.NewGinkgoLoopTopoTestBench(testbench.LoopTopoTestBenchConfig{
				LeftSpawn:  dmqspecagents.GolangAgent("left", dmqspecagents.NO_DEBUGGING),
				TopSpawn:   dmqspecagents.GolangAgent("top", dmqspecagents.NO_DEBUGGING),
				RightSpawn: dmqspecagents.GolangAgent("right", dmqspecagents.NO_DEBUGGING),

				LeftTTL:  directmq.DEFAULT_TTL,
				TopTTL:   directmq.DEFAULT_TTL,
				RightTTL: directmq.DEFAULT_TTL,

				LeftMaxMessageSize:  directmq.NO_MAX_MESSAGE_SIZE,
				TopMaxMessageSize:   directmq.NO_MAX_MESSAGE_SIZE,
				RightMaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,

				LogLeftToTopCommunication:  true,
				LogTopToRightCommunication: true,
				LogTopToLeftCommunication:  true,
				LogRightToTopCommunication: true,

				LogLeftLogs:  true,
				LogTopLogs:   true,
				LogRightLogs: true,

				DisableAllLogs: false,
			})

			defer bench.Stop("test ended")
			bench.Start()

			networkTerminated := make(chan struct{})
			bench.Left.OnNetworkTermination(func(_ dmqspecagent.OnNetworkTerminationNotification) {
				networkTerminated <- struct{}{}
			})

			log("subscribing to test from left")
			bench.Left.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			log("waiting for network to be terminated")
			select {
			case <-networkTerminated:
				log("network terminated")
			case <-time.After(5 * time.Second):
				Fail("network not terminated")
			}
		})
	})
})
