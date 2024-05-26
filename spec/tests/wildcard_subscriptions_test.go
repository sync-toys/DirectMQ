package dmqtests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
	dmqspecagents "github.com/sync-toys/DirectMQ/spec/agents"
	testbench "github.com/sync-toys/DirectMQ/spec/test_bench"
)

// wildcard subscription
// multiple wildcard subscription
// many wildcard subscriptions
// many multiple wildcard subscriptions

// super wildcard subscription
// super wildcard subscription with multiple wildcards
// super wildcard subscription with many wildcards
// super wildcard subscription with many multiple wildcards

// mixed wildcard and super wildcard subscriptions

// those tests are made in certain way
// so if salve received its message, it means that
// the host-local subscription already had time to process it

var _ = Describe("Wildcard subscriptions", func() {
	runTestCases := func(subscription string, tests ...interface{}) {
		When("subscribing to "+subscription, func() {
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

			controlTopic := "internal/control"
			controlMessage := "control message"

			for i := 0; i < len(tests); i += 2 {
				topic := tests[i].(string)
				expected := tests[i+1].(bool)

				description := "receive message on topic " + topic
				if !expected {
					description = "should not " + description
				} else {
					description = "should " + description
				}

				It(description, func() {
					subscriptionPropagatedToMaster := make(chan struct{}, 2)
					bench.Master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
						subscriptionPropagatedToMaster <- struct{}{}
					})

					bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
						Topic: subscription,
					})
					bench.Salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
						Topic: controlTopic,
					})

					log("waiting for subscription ** to be propagated to master")
					receiveWithTimeout(1, subscriptionPropagatedToMaster)

					log("not waiting for internal/control subscription")
					log("subscription already covered by **")

					messageReceivedBySalve := make(chan []byte, 2)
					bench.Salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
						messageReceivedBySalve <- message.Payload
					})

					log("publishing message")
					bench.Master.Publish(dmqspecagent.PublishCommand{
						Topic:            topic,
						DeliveryStrategy: directmq.AT_LEAST_ONCE,
						Payload:          []byte("Hello, World!"),
					})

					log("publishing control message")
					bench.Master.Publish(dmqspecagent.PublishCommand{
						Topic:            controlTopic,
						DeliveryStrategy: directmq.AT_LEAST_ONCE,
						Payload:          []byte(controlMessage),
					})

					if expected {
						log("waiting for message to be received by salve")
						Expect(receiveWithTimeout(1, messageReceivedBySalve)).To(Equal([]byte("Hello, World!")))
					}

					log("waiting for control message to be received by salve")
					Expect(receiveWithTimeout(1, messageReceivedBySalve)).To(Equal([]byte(controlMessage)))
				})
			}
		})

	}

	Context("*", func() {
		runTestCases("*",
			"a", true,
			"a/b", false,
		)

		runTestCases("test/*",
			"test/a", true,
			"test", false,
			"test/a/b", false,
		)

		runTestCases("test/*/a",
			"test/x/a", true,
			"test/test/a", true,
			"test/test/a/b", false,
			"test/a/x", false,
		)

		runTestCases("test/*/*",
			"test/a/b", true,
			"test/a/b/c", false,
			"test/a", false,
		)

		runTestCases("test/*/*/a",
			"test/c/b/a", true,
			"test/x/a", false,
			"test/x/x/x", false,
			"test/c/b/a/x", false,
		)

		runTestCases("test/a/b/*/d",
			"test/a/b/c/d", true,
			"test/a/b/c/x", false,
		)
	})

	Context("**", func() {
		runTestCases("**",
			"a", true,
			"a/b", true,
		)

		runTestCases("test/**",
			"test/a", true,
			"test/a/b", true,
			"test", false,
		)

		runTestCases("**/a",
			"x/a", true,
			"x/y/z/a", true,
			"x/y/z/b", false,
			"x/a/b", false,
		)

		runTestCases("**/a/**",
			"x/a/b", true,
			"x/y/z/a/b", true,
			"x/y/z/a", false,
			"a/x/y/z", false,
		)

		runTestCases("test/**/a",
			"test/x/a", true,
			"test/x/y/z/a", true,
			"test/x/y/z", false,
		)

		runTestCases("test/**/a/**/x",
			"test/x/a/y/x", true,
			"test/x/y/z/a/b/x", true,
			"test/x/y/z/a/b", false,
		)
	})

	Context("mixed", func() {
		runTestCases("test/**/a/*",
			"test/x/a/b", true,
			"test/x/y/a/b", true,
			"test/x/y/a", false,
			"test/x/y/a/b/c", false,
		)
	})
})
