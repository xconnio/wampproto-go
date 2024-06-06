package serializers

import (
	"reflect"

	"github.com/fxamacker/cbor/v2"

	"github.com/xconnio/wampproto-go/messages"
)

var cborEncoder cbor.DecMode //nolint:gochecknoglobals

type CBORSerializer struct{}

func init() {
	opt := cbor.DecOptions{DefaultMapType: reflect.TypeOf(map[string]any(nil))}
	cborEncoder, _ = opt.DecMode()
}

func (c *CBORSerializer) Serialize(message messages.Message) ([]byte, error) {
	return cbor.Marshal(message.Marshal())
}

func (c *CBORSerializer) Deserialize(payload []byte) (messages.Message, error) {
	var msgRaw []any
	if err := cborEncoder.Unmarshal(payload, &msgRaw); err != nil {
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

func (c *CBORSerializer) Static() bool {
	return false
}
