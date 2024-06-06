package serializers

import (
	"encoding/json"

	"github.com/xconnio/wampproto-go/messages"
)

type JSONSerializer struct{}

func (j *JSONSerializer) Serialize(message messages.Message) ([]byte, error) {
	return json.Marshal(message.Marshal())
}

func (j *JSONSerializer) Deserialize(payload []byte) (messages.Message, error) {
	var msgRaw []any
	if err := json.Unmarshal(payload, &msgRaw); err != nil {
		return nil, err
	}

	msg, err := ToMessage(msgRaw)
	if err != nil {
		return nil, err
	}

	if err = msg.Parse(msgRaw); err != nil {
		return nil, err
	}

	return msg, nil
}

func (j *JSONSerializer) Static() bool {
	return false
}
