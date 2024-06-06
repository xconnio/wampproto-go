package serializers

import (
	"github.com/vmihailenco/msgpack/v5"

	"github.com/xconnio/wampproto-go/messages"
)

type MsgPackSerializer struct{}

func (m *MsgPackSerializer) Serialize(message messages.Message) ([]byte, error) {
	return msgpack.Marshal(message.Marshal())
}

func (m *MsgPackSerializer) Deserialize(payload []byte) (messages.Message, error) {
	var msgRaw []any
	if err := msgpack.Unmarshal(payload, &msgRaw); err != nil {
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

func (m *MsgPackSerializer) Static() bool {
	return false
}
