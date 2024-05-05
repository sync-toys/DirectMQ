package dmqtests

import (
	"io"
	"net/url"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqspec "github.com/sync-toys/DirectMQ/spec"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

var _ = Describe("Pubsub", func() {
	logTest := dmqspec.CreateLogger("test", func(message string, writer io.Writer) {
		color.New(color.FgGreen).Fprintln(writer, message)
	})

	logComm := dmqspec.CreateLogger("comm", func(message string, writer io.Writer) {
		color.New(color.FgHiYellow).Fprintln(writer, message)
	})

	logMaster := dmqspec.CreateLogger("master", func(message string, writer io.Writer) {
		color.New(color.FgHiBlue).Fprintln(writer, message)
	})
	logSalve := dmqspec.CreateLogger("salve", func(message string, writer io.Writer) {
		color.New(color.FgHiBlue).Fprintln(writer, message)
	})

	forwardFrom, _ := url.Parse("http://localhost:8081/")
	forwardTo, _ := url.Parse("ws://localhost:8080/")
	forwarder := dmqspecagent.NewWebsocketForwarder(forwardFrom, forwardTo)
	forwarder.OnMessage(func(message string, sender, receiver *url.URL) {
		senderAgent := "master"
		receiverAgent := "salve"
		if sender.Port() == "8081" {
			senderAgent = "salve"
			receiverAgent = "master"
		}

		logComm("%s->%s message: %s", senderAgent, receiverAgent, message)
	})

	master := dmqspecagent.NewUniversalAgent()
	master.OnLogEntry(func(entry string) {
		logMaster(entry)
	})
	master.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		logMaster("FATAL: %s\n", fatal.Err)
		Fail(fatal.Err)
	})

	salve := dmqspecagent.NewUniversalAgent()
	salve.OnLogEntry(func(entry string) {
		logSalve(entry)
	})
	salve.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		logSalve("FATAL: %s\n", fatal.Err)
		Fail(fatal.Err)
	})

	var masterExec dmqspecagent.UniversalSpawn = dmqspecagent.UniversalSpawn{
		// ExecutablePath: "dlv",
		// Arguments:      []string{"exec", "./agents/go/go-agent.exe", "--headless", "--listen=:2345", "--api-version=2"},

		ExecutablePath: "../agents/go/bin/go-agent",
		Arguments:      []string{},
	}

	var salveExec dmqspecagent.UniversalSpawn = dmqspecagent.UniversalSpawn{
		// ExecutablePath: "dlv",
		// Arguments:      []string{"exec", "./agents/go/go-agent.exe", "--headless", "--listen=:2345", "--api-version=2"},

		ExecutablePath: "../agents/go/bin/go-agent",
		Arguments:      []string{},
	}

	BeforeEach(func() {
		go forwarder.StartWebsocketForwarder()
	})

	AfterEach(func() {
		logTest("stopping agents and forwarder")

		logTest("stopping master agent")
		master.Stop(dmqspecagent.StopCommand{
			Reason: "test ended",
		})

		logTest("stopping salve agent")
		salve.Stop(dmqspecagent.StopCommand{
			Reason: "test ended",
		})

		logTest("aborting forwarder")
		forwarder.Abort()

		logTest("agents and forwarder stopped gracefully")
	})

	When("agent A publishes a message", func() {
		It("should be received by agent B", func() {
			masterReady := make(chan struct{}, 1)
			masterConnected := make(chan struct{}, 1)

			salveReady := make(chan struct{}, 1)
			salveConnected := make(chan struct{}, 1)

			master.OnReady(func(_ dmqspecagent.ReadyNotification) {
				logTest("master ready")
				masterReady <- struct{}{}
			})

			master.OnConnectionEstablished(func(notification dmqspecagent.ConnectionEstablishedNotification) {
				logTest("master connected")
				masterConnected <- struct{}{}
			})

			salve.OnReady(func(_ dmqspecagent.ReadyNotification) {
				logTest("salve ready")
				salveReady <- struct{}{}
			})

			salve.OnConnectionEstablished(func(notification dmqspecagent.ConnectionEstablishedNotification) {
				logTest("salve connected")
				salveConnected <- struct{}{}
			})

			logTest("listening on master agent")
			master.Listen(dmqspecagent.ListenCommand{
				TTL:            32,
				Address:        "http://localhost:8080/",
				AsClientId:     "master",
				MaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,
			}, masterExec)

			<-masterReady

			logTest("connecting from salve agent")
			salve.Connect(dmqspecagent.ConnectCommand{
				TTL:            32,
				Address:        "ws://localhost:8081/",
				AsClientId:     "salve",
				MaxMessageSize: directmq.NO_MAX_MESSAGE_SIZE,
			}, salveExec)

			<-salveReady

			logTest("waiting for master connection to be established")
			<-masterConnected

			logTest("waiting for salve connection to be established")
			<-salveConnected

			subscriptionPropagated := make(chan struct{}, 1)
			master.OnSubscription(func(subscription dmqspecagent.OnSubscriptionNotification) {
				subscriptionPropagated <- struct{}{}
			})

			logTest("subscribing to test topic from salve")
			salve.Subscribe(dmqspecagent.SubscribeTopicCommand{
				Topic: "test",
			})

			logTest("waiting for subscription to be propagated to master")
			<-subscriptionPropagated

			result := make(chan []byte)
			salve.OnMessageReceived(func(message dmqspecagent.MessageReceivedNotification) {
				result <- message.Payload
			})

			logTest("publishing message from master on test topic")
			master.Publish(dmqspecagent.PublishCommand{
				Topic:            "test",
				DeliveryStrategy: directmq.AT_LEAST_ONCE,
				Payload:          []byte("Hello, World!"),
			})

			logTest("waiting for message to be received by salve")
			Expect(<-result).To(Equal([]byte("Hello, World!")))
		})
	})
})
