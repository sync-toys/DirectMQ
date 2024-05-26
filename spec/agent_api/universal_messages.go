package dmqspecagent

import directmq "github.com/sync-toys/DirectMQ/sdk/go"

type UniversalCommand struct {
	// Connection API
	Setup   *SetupCommand
	Listen  *ListenCommand
	Connect *ConnectCommand
	Stop    *StopCommand

	// Native API
	Publish          *PublishCommand
	SubscribeTopic   *SubscribeTopicCommand
	UnsubscribeTopic *UnsubscribeTopicCommand
}

type UniversalNotification struct {
	// Agent API
	Ready           *ReadyNotification
	Fatal           *FatalNotification
	MessageReceived *MessageReceivedNotification
	Subscribed      *SubscribedNotification
	Stopped         *StoppedNotification

	// Diagnostics API
	ConnectionEstablished *ConnectionEstablishedNotification
	ConnectionLost        *ConnectionLostNotification
	OnPublication         *OnPublicationNotification
	OnSubscription        *OnSubscriptionNotification
	OnUnsubscribe         *OnUnsubscribeNotification
	OnNetworkTermination  *OnNetworkTerminationNotification
}

// Connection API

type SetupCommand struct {
	TTL            int32
	NodeID         string
	MaxMessageSize uint64
}

type ListenCommand struct {
	Address string
}

type ConnectCommand struct {
	Address string
}

type StopCommand struct {
	Reason string
}

// Agent API

type ConnectionEstablishedNotification struct {
	BridgedNodeId string
}

type ConnectionLostNotification struct {
	BridgedNodeId string
	Reason        string
}

type ReadyNotification struct {
	Time string
}

type FatalNotification struct {
	Err string
}

type MessageReceivedNotification struct {
	Topic   string
	Payload []byte
}

type SubscribedNotification struct {
	SubscriptionID directmq.SubscriptionID
}

type StoppedNotification struct {
	Reason string
}

// Native API

type PublishCommand struct {
	Topic            string
	DeliveryStrategy directmq.DeliveryStrategy
	Payload          []byte
}
type SubscribeTopicCommand struct {
	Topic string
}

type UnsubscribeTopicCommand struct {
	SubscriptionID directmq.SubscriptionID
}

// Diagnostics API

type OnPublicationNotification struct {
	TTL              int32
	Traversed        []string
	Topic            string
	DeliveryStrategy directmq.DeliveryStrategy
	Payload          []byte
}

type OnSubscriptionNotification struct {
	TTL       int32
	Traversed []string
	Topic     string
}

type OnUnsubscribeNotification struct {
	TTL       int32
	Traversed []string
	Topic     string
}

type OnNetworkTerminationNotification struct {
	TTL       int32
	Traversed []string
	Reason    string
}
