package serializers

import (
	"fmt"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/util"
)

func ToMessage(wampMsg []any) (messages.Message, error) {
	messageType, _ := util.AsInt64(wampMsg[0])
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
	case messages.MessageTypeUnsubscribe:
		msg = &messages.Unsubscribe{}
	case messages.MessageTypeUnsubscribed:
		msg = &messages.Unsubscribed{}
	case messages.MessageTypeUnregister:
		msg = &messages.Unregister{}
	case messages.MessageTypeUnregistered:
		msg = &messages.Unregistered{}
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
