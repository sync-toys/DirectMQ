package directmq

type networkParticipant interface {
	GetSubscribedTopics() []string
	WillHandleTopic(topic string) bool
	IsOriginOfFrame(message DataFrame) bool

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

func (d *globalNetwork) getSubscribedTopicsExcludingOriginOfMessage(frame DataFrame) []string {
	topics := make([]string, 0)
	for _, participant := range d.participants {
		if participant.IsOriginOfFrame(frame) {
			continue
		}

		topics = append(topics, participant.GetSubscribedTopics()...)
	}

	return unique(topics)
}

func (d *globalNetwork) Published(message PublishMessage) {
	d.diag.HandlePublish(message)

	for _, participant := range randomOrder(d.participants) {
		handled := participant.HandlePublish(message)
		if message.DeliveryStrategy == AT_MOST_ONCE && handled {
			return
		}
	}
}

func (d *globalNetwork) Subscribed(message SubscribeMessage) {
	d.diag.HandleSubscribe(message)

	// todo: we need to check if every edge has given topic subscribed
	topLevelSubscriptions := d.getSubscribedTopicsExcludingOriginOfMessage(message.DataFrame)

	for _, participant := range d.participants {
		participant.HandleSubscribe(message)
	}

	if len(message.DataFrame.Traversed) == 0 {
		// we are skipping the unsubscribe synchronization part
		// because current node is origin of the message
		// and synchronization will be done by native api itself
		return
	}

	updatedTopLevelSubscriptions := d.GetAllSubscribedTopics()
	topicsToUnsubscribe, _ := GetDeduplicatedOverlappingTopicsDiff(topLevelSubscriptions, updatedTopLevelSubscriptions)

	for _, topic := range topicsToUnsubscribe {
		message := UnsubscribeMessage{
			DataFrame: message.DataFrame,
			Topic:     topic,
		}

		for _, participant := range d.participants {
			if participant.WillHandleTopic(message.Topic) {
				continue
			}

			participant.HandleUnsubscribe(message)
		}
	}
}

func (d *globalNetwork) Unsubscribed(message UnsubscribeMessage) {
	d.diag.HandleUnsubscribe(message)

	// todo: we need to check if every edge has given topic unsubscribed
	for _, participant := range d.participants {
		if participant.WillHandleTopic(message.Topic) {
			return
		}
	}

	for _, participant := range d.participants {
		participant.HandleUnsubscribe(message)
	}
}

func (d *globalNetwork) Terminated(message TerminateNetworkMessage) {
	d.diag.HandleTerminateNetwork(message)

	for _, participant := range d.participants {
		participant.HandleTerminateNetwork(message)
	}
}
