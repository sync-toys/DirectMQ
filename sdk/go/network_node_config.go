package directmq

const (
	DEFAULT_TTL                              = 32
	ONLY_DIRECT_CONNECTION_TTL               = 1
	ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL = 2
	NO_ANSWER_TTL                            = 0
	NO_MAX_MESSAGE_SIZE                      = 0
)

type TTL int32

type NetworkNodeConfig struct {
	HostTTL                    TTL
	HostMaxIncomingMessageSize uint64
	HostID                     string
}
