package directmq

type networkParticipant interface {
	GetSubscribedTopics() []string

	HandlePublish(publication PublishMessage) (handled bool)
	HandleSubscribe(subscription SubscribeMessage)
	HandleUnsubscribe(unsubscribe UnsubscribeMessage)
	HandleTerminateNetwork(terminate TerminateNetworkMessage)
}

type globalNetwork struct {
	config       NetworkNodeConfig
	participants []networkParticipant
	diag         *diagnosticsAPI
}

func newGlobalNetwork(config NetworkNodeConfig, nativeAPI *nativeAPI, diag *diagnosticsAPI) *globalNetwork {
	return &globalNetwork{
		config:       config,
		participants: []networkParticipant{nativeAPI},
		diag:         diag,
	}
}

func (d *globalNetwork) GetAllSubscribedTopics() []string {
	topics := make([]string, 0)
	for _, participant := range d.participants {
		topics = append(topics, participant.GetSubscribedTopics()...)
	}

	return unique(topics)
}

func (d *globalNetwork) Published(message PublishMessage, callingParticipant networkParticipant) {
	d.diag.HandlePublish(message)

	for _, participant := range randomOrder(d.participants) {
		if participant == callingParticipant {
			continue
		}

		handled := participant.HandlePublish(message)
		if message.DeliveryStrategy == AT_MOST_ONCE && handled {
			return
		}
	}
}

func (d *globalNetwork) Subscribed(message SubscribeMessage, callingParticipant networkParticipant) {
	d.diag.HandleSubscribe(message)

	for _, participant := range d.participants {
		if participant == callingParticipant {
			continue
		}

		participant.HandleSubscribe(message)
	}
}

func (d *globalNetwork) Unsubscribed(message UnsubscribeMessage, callingParticipant networkParticipant) {
	d.diag.HandleUnsubscribe(message)

	for _, participant := range d.participants {
		if participant == callingParticipant {
			continue
		}

		participant.HandleUnsubscribe(message)
	}
}

func (d *globalNetwork) Terminated(message TerminateNetworkMessage, callingParticipant networkParticipant) {
	d.diag.HandleTerminateNetwork(message)

	for _, participant := range d.participants {
		if participant == callingParticipant {
			continue
		}

		participant.HandleTerminateNetwork(message)
	}
}
