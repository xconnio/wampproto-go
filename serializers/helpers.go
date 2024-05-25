package serializers

import (
	"fmt"

	"github.com/xconnio/wampproto-go/messages"
)

func ToMessage(wampMsg []any) (messages.Message, error) {
	messageType, _ := AsInt64(wampMsg[0])
	var msg messages.Message
	switch messageType {
	case messages.MessageTypeAbort:
		msg = messages.NewEmptyAbort()
	default:
		return nil, fmt.Errorf("unknown message %T", wampMsg[0])
	}

	if err := msg.Parse(wampMsg); err != nil {
		return nil, err
	}

	return msg, nil
}

func AsInt64(i interface{}) (int64, bool) {
	switch v := i.(type) {
	case int64:
		return v, true
	case uint64:
		return int64(v), true
	case uint8:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint:
		return int64(v), true
	case uint32:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	}
	return 0, false
}
