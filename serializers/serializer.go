package serializers

import "github.com/xconnio/wampproto-go/messages"

const NoneSerializerID = 0

type Serializer interface {
	Serialize(message messages.Message) ([]byte, error)
	Deserialize([]byte) (messages.Message, error)
	Static() bool
}
