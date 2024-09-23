package dmqspecagent

import directmq "github.com/sync-toys/DirectMQ/sdk/go"

type UniversalCommand struct {
	// Connection API
	Setup   *SetupCommand   `json:"setup,omitempty"`
	Listen  *ListenCommand  `json:"listen,omitempty"`
	Connect *ConnectCommand `json:"connect,omitempty"`
	Stop    *StopCommand    `json:"stop,omitempty"`

	// Native API
	Publish          *PublishCommand          `json:"publish,omitempty"`
	SubscribeTopic   *SubscribeTopicCommand   `json:"subscribeTopic,omitempty"`
	UnsubscribeTopic *UnsubscribeTopicCommand `json:"unsubscribeTopic,omitempty"`
}

type UniversalNotification struct {
	// Agent API
	Ready           *ReadyNotification           `json:"ready,omitempty"`
	Fatal           *FatalNotification           `json:"fatal,omitempty"`
	MessageReceived *MessageReceivedNotification `json:"messageReceived,omitempty"`
	Subscribed      *SubscribedNotification      `json:"subscribed,omitempty"`
	Stopped         *StoppedNotification         `json:"stopped,omitempty"`

	// Diagnostics API
	ConnectionEstablished *ConnectionEstablishedNotification `json:"connectionEstablished,omitempty"`
	ConnectionLost        *ConnectionLostNotification        `json:"connectionLost,omitempty"`
	OnPublication         *OnPublicationNotification         `json:"onPublication,omitempty"`
	OnSubscription        *OnSubscriptionNotification        `json:"onSubscription,omitempty"`
	OnUnsubscribe         *OnUnsubscribeNotification         `json:"onUnsubscribe,omitempty"`
	OnNetworkTermination  *OnNetworkTerminationNotification  `json:"onNetworkTermination,omitempty"`
}

// Connection API

type SetupCommand struct {
	TTL            int32  `json:"ttl"`
	NodeID         string `json:"nodeId,omitempty"`
	MaxMessageSize uint64 `json:"maxMessageSize"`
}

type ListenCommand struct {
	Address string `json:"address,omitempty"`
}

type ConnectCommand struct {
	Address string `json:"address,omitempty"`
}

type StopCommand struct {
	Reason string `json:"reason,omitempty"`
}

// Agent API

type ConnectionEstablishedNotification struct {
	BridgedNodeId string `json:"bridgedNodeId,omitempty"`
}

type ConnectionLostNotification struct {
	BridgedNodeId string `json:"bridgedNodeId,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

type ReadyNotification struct {
	Time string `json:"time,omitempty"`
}

type FatalNotification struct {
	Err string `json:"err,omitempty"`
}

type MessageReceivedNotification struct {
	Topic   string `json:"topic,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}

type SubscribedNotification struct {
	SubscriptionID directmq.SubscriptionID `json:"subscriptionId,omitempty"`
}

type StoppedNotification struct {
	Reason string `json:"reason,omitempty"`
}

// Native API

type PublishCommand struct {
	Topic            string                    `json:"topic,omitempty"`
	DeliveryStrategy directmq.DeliveryStrategy `json:"deliveryStrategy,omitempty"`
	Payload          []byte                    `json:"payload,omitempty"`
}
type SubscribeTopicCommand struct {
	Topic string `json:"topic,omitempty"`
}

type UnsubscribeTopicCommand struct {
	SubscriptionID directmq.SubscriptionID `json:"subscriptionId,omitempty"`
}

// Diagnostics API

type OnPublicationNotification struct {
	TTL              int32                     `json:"ttl"`
	Traversed        []string                  `json:"traversed,omitempty"`
	Topic            string                    `json:"topic,omitempty"`
	DeliveryStrategy directmq.DeliveryStrategy `json:"deliveryStrategy,omitempty"`
	Payload          []byte                    `json:"payload,omitempty"`
}

type OnSubscriptionNotification struct {
	TTL       int32    `json:"ttl"`
	Traversed []string `json:"traversed,omitempty"`
	Topic     string   `json:"topic,omitempty"`
}

type OnUnsubscribeNotification struct {
	TTL       int32    `json:"ttl"`
	Traversed []string `json:"traversed,omitempty"`
	Topic     string   `json:"topic,omitempty"`
}

type OnNetworkTerminationNotification struct {
	TTL       int32    `json:"ttl"`
	Traversed []string `json:"traversed,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}
