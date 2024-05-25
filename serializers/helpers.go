package serializers

import (
	"fmt"

	"github.com/xconnio/wampproto-go/messages"
)

func ToMessage(wampMsg []any) (messages.Message, error) {
	messageType, _ := messages.AsInt64(wampMsg[0])
	var msg messages.Message
	switch messageType {
	case messages.MessageTypeAbort:
		msg = messages.NewEmptyAbort()
	default:
		return nil, fmt.Errorf("unknown message %T", wampMsg[0])
	}

	if err := msg.Parse(wampMsg); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	return msg, nil
}
