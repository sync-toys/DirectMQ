package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	directmq "github.com/sync-toys/DirectMQ/sdk/go"
	dmqportals "github.com/sync-toys/DirectMQ/sdk/go/portals"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

var node directmq.NetworkNode
var ctx context.Context

func startCommandLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		log("Waiting for command")
		command, err := reader.ReadString('\n')
		if err != nil {
			fatal("Error during command reading: " + err.Error())
			return
		}

		var cmd dmqspecagent.UniversalCommand
		if err := json.Unmarshal([]byte(command), &cmd); err != nil {
			fatal("Error during command unmarshalling: " + err.Error())
		}

		handleIncomingCommand(cmd)
	}
}

func handleIncomingCommand(cmd dmqspecagent.UniversalCommand) {
	if cmd.Setup != nil {
		handleSetupCommand(*cmd.Setup)
	}

	if cmd.Listen != nil {
		handleListenCommand(*cmd.Listen)
	}

	if cmd.Connect != nil {
		handleConnectCommand(*cmd.Connect)
	}

	if cmd.Stop != nil {
		handleStopCommand(*cmd.Stop)
	}

	if cmd.Publish != nil {
		handlePublishCommand(*cmd.Publish)
	}

	if cmd.SubscribeTopic != nil {
		handleSubscribeCommand(*cmd.SubscribeTopic)
	}

	if cmd.UnsubscribeTopic != nil {
		handleUnsubscribeCommand(*cmd.UnsubscribeTopic)
	}
}

func handleSetupCommand(cmd dmqspecagent.SetupCommand) {
	log("Setting up DirectMQ node")

	node = directmq.NewNetworkNode(
		directmq.NetworkNodeConfig{
			HostID:                     cmd.NodeID,
			HostTTL:                    directmq.TTL(cmd.TTL),
			HostMaxIncomingMessageSize: cmd.MaxMessageSize,
		},
		directmq.NewProtobufJSONProtocol(),
	)

	registerDiagnosticsHandlers()

	log("Setup complete")
}

func handleListenCommand(cmd dmqspecagent.ListenCommand) {
	log("Listening as server at " + cmd.Address)

	u, err := url.Parse(cmd.Address)
	if err != nil {
		fatal("Failed to parse URL: " + err.Error())
	}

	ctx = context.Background()

	go func() {
		err = dmqportals.WebsocketListen(u, dmqportals.TextMessages, node, ctx)
		if err != nil {
			fatal("Failed to listen on websocket: " + err.Error())
		}
	}()
}

func handleConnectCommand(cmd dmqspecagent.ConnectCommand) {
	log("Connecting as client to " + cmd.Address)

	u, err := url.Parse(cmd.Address)
	if err != nil {
		fatal("Failed to parse URL: " + err.Error())
	}

	log("Connecting websocket")
	wsPortal, err := dmqportals.WebsocketConnect(u, dmqportals.TextMessages)
	if err != nil {
		fatal("Failed to connect to websocket: " + err.Error())
	}

	log("Starting connection protocol")
	go func() {
		log("Starting connection protocol")
		if err := node.AddConnectingEdge(wsPortal); !dmqportals.IsNetworkConnectionClosedError(err) {
			fatal("Websocket connection failure: " + err.Error())
		}
	}()
}

func handleStopCommand(cmd dmqspecagent.StopCommand) {
	log("Stopping DirectMQ node")
	node.CloseNode(cmd.Reason)

	log("Sending stopped notification")
	sendNotification(dmqspecagent.UniversalNotification{
		Stopped: &dmqspecagent.StoppedNotification{
			Reason: cmd.Reason,
		},
	})

	log("Clean exit 0")
	os.Exit(0)
}

func handlePublishCommand(cmd dmqspecagent.PublishCommand) {
	log("Publishing message to topic: " + cmd.Topic)
	node.Publish(cmd.Topic, cmd.Payload, cmd.DeliveryStrategy)
}

func handleSubscribeCommand(cmd dmqspecagent.SubscribeTopicCommand) {
	log("Subscribing to topic: " + cmd.Topic)
	handler := createSubscriptionHandler(cmd.Topic)
	subscriptionID := node.Subscribe(cmd.Topic, handler)
	log("Subscription ID: " + strconv.Itoa(int(subscriptionID)))
	sendNotification(dmqspecagent.UniversalNotification{
		Subscribed: &dmqspecagent.SubscribedNotification{
			SubscriptionID: subscriptionID,
		},
	})
}

func handleUnsubscribeCommand(cmd dmqspecagent.UnsubscribeTopicCommand) {
	node.Unsubscribe(cmd.SubscriptionID)
}

func createSubscriptionHandler(topic string) func(payload []byte) {
	return func(payload []byte) {
		sendNotification(dmqspecagent.UniversalNotification{
			MessageReceived: &dmqspecagent.MessageReceivedNotification{
				Topic:   topic,
				Payload: payload,
			},
		})
	}
}

func sendNotification(notification dmqspecagent.UniversalNotification) {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		panic("Failed to marshal notification: " + err.Error())
	}

	notificationBytes = append(notificationBytes, '\n')

	if _, err := os.Stdout.Write(notificationBytes); err != nil {
		panic("Failed to write notification to stdout: " + err.Error())
	}
}

func registerDiagnosticsHandlers() {
	log("Registering diagnostics handlers")

	node.OnConnectionEstablished(func(bridgedNodeID string, portal directmq.Portal) {
		sendNotification(dmqspecagent.UniversalNotification{
			ConnectionEstablished: &dmqspecagent.ConnectionEstablishedNotification{
				BridgedNodeId: bridgedNodeID,
			},
		})
	})

	node.OnConnectionLost(func(bridgedNodeID, reason string, portal directmq.Portal) {
		sendNotification(dmqspecagent.UniversalNotification{
			ConnectionLost: &dmqspecagent.ConnectionLostNotification{
				BridgedNodeId: bridgedNodeID,
				Reason:        reason,
			},
		})
	})

	node.OnPublication(func(publication directmq.PublishMessage) {
		sendNotification(dmqspecagent.UniversalNotification{
			OnPublication: &dmqspecagent.OnPublicationNotification{
				TTL:              publication.TTL,
				Traversed:        publication.Traversed,
				Topic:            publication.Topic,
				DeliveryStrategy: publication.DeliveryStrategy,
				Payload:          publication.Payload,
			},
		})
	})

	node.OnSubscription(func(subscription directmq.SubscribeMessage) {
		sendNotification(dmqspecagent.UniversalNotification{
			OnSubscription: &dmqspecagent.OnSubscriptionNotification{
				TTL:       subscription.TTL,
				Traversed: subscription.Traversed,
				Topic:     subscription.Topic,
			},
		})
	})

	node.OnUnsubscribe(func(unsubscribe directmq.UnsubscribeMessage) {
		sendNotification(dmqspecagent.UniversalNotification{
			OnUnsubscribe: &dmqspecagent.OnUnsubscribeNotification{
				TTL:       unsubscribe.TTL,
				Traversed: unsubscribe.Traversed,
				Topic:     unsubscribe.Topic,
			},
		})
	})

	node.OnTerminateNetwork(func(termination directmq.TerminateNetworkMessage) {
		sendNotification(dmqspecagent.UniversalNotification{
			OnNetworkTermination: &dmqspecagent.OnNetworkTerminationNotification{
				TTL:       termination.TTL,
				Traversed: termination.Traversed,
				Reason:    termination.Reason,
			},
		})
	})
}

func fatal(reason string) {
	defer os.Exit(1)

	sendNotification(dmqspecagent.UniversalNotification{
		Fatal: &dmqspecagent.FatalNotification{
			Err: reason,
		},
	})

	panic(reason)
}

func ready() {
	sendNotification(dmqspecagent.UniversalNotification{
		Ready: &dmqspecagent.ReadyNotification{
			Time: time.Now().String(),
		},
	})
}

func log(entry string) {
	fmt.Println(entry)
}

func main() {
	log("Starting DirectMQ testing agent")
	go ready()
	log("Starting command loop")
	startCommandLoop()
}
