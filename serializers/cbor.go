package serializers

import (
	"fmt"
	"reflect"

	"github.com/fxamacker/cbor/v2"

	"github.com/xconnio/wampproto-go/messages"
)

const CBORSerializerID = 3

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

func EncodeCBOR(args []any, kwargs map[string]any) ([]byte, error) {
	var payload []any
	payload = append(payload, args)
	payload = append(payload, kwargs)
	payloadData, err := cbor.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadData, nil
}

func DecodeCBOR(data []byte) ([]any, map[string]any, error) {
	var payload []any
	if err := cborEncoder.Unmarshal(data, &payload); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal cbor payload: %w", err)
	}

	args, ok := payload[0].([]any)
	if !ok {
		return nil, nil, fmt.Errorf("invalid args type: %T", payload[0])
	}

	kwargs, ok := payload[1].(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("invalid kwargs type: %T", payload[1])
	}

	return args, kwargs, nil
}
