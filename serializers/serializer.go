package serializers

import "github.com/xconnio/wampproto-go/messages"

type Serializer interface {
	Serialize(message messages.Message) ([]byte, error)
	Deserialize([]byte) (messages.Message, error)
}
