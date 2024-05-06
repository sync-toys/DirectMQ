package dmqspecagent

import directmq "github.com/sync-toys/DirectMQ/sdk/go"

type AgentConnectionAPI interface {
	Listen(ListenCommand, UniversalSpawn) error
	Connect(ConnectCommand, UniversalSpawn) error
	Stop(StopCommand)
	Kill()
}

type AgentAPI interface {
	OnConnectionEstablished(func(ConnectionEstablishedNotification))
	OnConnectionLost(func(ConnectionLostNotification))
	OnReady(func(ReadyNotification))
	OnFatal(func(FatalNotification))
	OnLogEntry(func(string))
	OnMessageReceived(func(MessageReceivedNotification))
}

type NativeAPI interface {
	Publish(PublishCommand)
	Subscribe(SubscribeTopicCommand) directmq.SubscriptionID
	Unsubscribe(UnsubscribeTopicCommand)
}

type DiagnosticsAPI interface {
	OnPublication(func(OnPublicationNotification))
	OnSubscription(func(OnSubscriptionNotification))
	OnUnsubscribe(func(OnUnsubscribeNotification))
	OnNetworkTermination(func(OnNetworkTerminationNotification))
}

type Agent interface {
	AgentConnectionAPI
	AgentAPI
	NativeAPI
	DiagnosticsAPI

	GetNodeID() string
}
