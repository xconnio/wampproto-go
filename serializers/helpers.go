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
		msg = &messages.Abort{}
	case messages.MessageTypeAuthenticate:
		msg = &messages.Authenticate{}
	case messages.MessageTypeCall:
		msg = &messages.Call{}
	case messages.MessageTypeCancel:
		msg = &messages.Cancel{}
	case messages.MessageTypeChallenge:
		msg = &messages.Challenge{}
	case messages.MessageTypeError:
		msg = &messages.Error{}
	case messages.MessageTypeEvent:
		msg = &messages.Event{}
	case messages.MessageTypeGoodbye:
		msg = &messages.GoodBye{}
	case messages.MessageTypeHello:
		msg = &messages.Hello{}
	case messages.MessageTypeInterrupt:
		msg = &messages.Interrupt{}
	case messages.MessageTypeInvocation:
		msg = &messages.Invocation{}
	case messages.MessageTypePublish:
		msg = &messages.Publish{}
	case messages.MessageTypePublished:
		msg = &messages.Published{}
	case messages.MessageTypeRegister:
		msg = &messages.Register{}
	case messages.MessageTypeRegistered:
		msg = &messages.Registered{}
	case messages.MessageTypeResult:
		msg = &messages.Result{}
	case messages.MessageTypeSubscribe:
		msg = &messages.Subscribe{}
	case messages.MessageTypeSubscribed:
		msg = &messages.Subscribed{}
	case messages.MessageTypeUnSubscribe:
		msg = &messages.UnSubscribe{}
	case messages.MessageTypeUnSubscribed:
		msg = &messages.UnSubscribed{}
	case messages.MessageTypeUnRegister:
		msg = &messages.UnRegister{}
	case messages.MessageTypeUnRegistered:
		msg = &messages.UnRegistered{}
	case messages.MessageTypeWelcome:
		msg = &messages.Welcome{}
	case messages.MessageTypeYield:
		msg = &messages.Yield{}
	default:
		return nil, fmt.Errorf("unknown message %T", wampMsg[0])
	}

	if err := msg.Parse(wampMsg); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	return msg, nil
}
