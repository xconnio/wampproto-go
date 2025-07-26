package wampproto

import (
	"fmt"

	"github.com/xconnio/wampproto-go/internal"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

type Session struct {
	serializer serializers.Serializer

	// data structures for RPC
	callRequests       internal.Map[uint64, uint64]
	registerRequests   internal.Map[uint64, uint64]
	registrations      internal.Map[uint64, uint64]
	invocationRequests internal.Map[uint64, uint64]
	unregisterRequests internal.Map[uint64, uint64]

	// data structures for PubSub
	publishRequests     internal.Map[uint64, uint64]
	subscribeRequests   internal.Map[uint64, uint64]
	subscriptions       internal.Map[uint64, uint64]
	unsubscribeRequests internal.Map[uint64, uint64]
}

func NewSession(serializer serializers.Serializer) *Session {
	if serializer == nil {
		serializer = &serializers.JSONSerializer{}
	}

	return &Session{
		serializer: serializer,

		callRequests:       internal.Map[uint64, uint64]{},
		registerRequests:   internal.Map[uint64, uint64]{},
		registrations:      internal.Map[uint64, uint64]{},
		invocationRequests: internal.Map[uint64, uint64]{},
		unregisterRequests: internal.Map[uint64, uint64]{},

		publishRequests:     internal.Map[uint64, uint64]{},
		subscribeRequests:   internal.Map[uint64, uint64]{},
		subscriptions:       internal.Map[uint64, uint64]{},
		unsubscribeRequests: internal.Map[uint64, uint64]{},
	}
}

func (w *Session) SendMessage(msg messages.Message) ([]byte, error) {
	data, err := w.serializer.Serialize(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	switch msg.Type() {
	case messages.MessageTypeCall:
		call := msg.(*messages.Call)
		w.callRequests.Store(call.RequestID(), call.RequestID())

		return data, nil
	case messages.MessageTypeYield:
		yield := msg.(*messages.Yield)
		progress, _ := yield.Options()[OptionProgress].(bool)
		if !progress {
			w.invocationRequests.Delete(yield.RequestID())
		}

		return data, nil
	case messages.MessageTypeRegister:
		register := msg.(*messages.Register)
		w.registerRequests.Store(register.RequestID(), register.RequestID())

		return data, nil
	case messages.MessageTypeUnregister:
		unregister := msg.(*messages.Unregister)
		w.unregisterRequests.Store(unregister.RequestID(), unregister.RegistrationID())

		return data, nil
	case messages.MessageTypePublish:
		publish := msg.(*messages.Publish)
		acknowledge, ok := publish.Options()["acknowledge"].(bool)
		if ok && acknowledge {
			w.publishRequests.Store(publish.RequestID(), publish.RequestID())
		}

		return data, nil
	case messages.MessageTypeSubscribe:
		subscribe := msg.(*messages.Subscribe)
		w.subscribeRequests.Store(subscribe.RequestID(), subscribe.RequestID())

		return data, nil
	case messages.MessageTypeUnsubscribe:
		unsubscribe := msg.(*messages.Unsubscribe)
		_, exists := w.subscriptions.Load(unsubscribe.SubscriptionID())
		if !exists {
			return nil, fmt.Errorf("unsubscribe request for non existent subscription %d",
				unsubscribe.SubscriptionID())
		}

		w.unsubscribeRequests.Store(unsubscribe.RequestID(), unsubscribe.SubscriptionID())

		return data, nil
	case messages.MessageTypeError:
		errorMsg := msg.(*messages.Error)
		if errorMsg.MessageType() != messages.MessageTypeInvocation {
			return nil, fmt.Errorf("send only supported for invocation error")
		}

		w.invocationRequests.Delete(errorMsg.RequestID())
		return data, nil
	case messages.MessageTypeGoodbye:
		return data, nil
	default:
		return nil, fmt.Errorf("send not supported for message of type %T", msg)
	}
}

func (w *Session) Receive(data []byte) (messages.Message, error) {
	msg, err := w.serializer.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %w", err)
	}

	return w.ReceiveMessage(msg)
}

func (w *Session) ReceiveMessage(msg messages.Message) (messages.Message, error) {
	switch msg.Type() {
	case messages.MessageTypeResult:
		result := msg.(*messages.Result)
		_, exists := w.callRequests.Load(result.RequestID())
		if !exists {
			return nil, fmt.Errorf("received RESULT for invalid requestID")
		}

		progress, _ := result.Details()[OptionProgress].(bool)
		if !progress {
			w.callRequests.Delete(result.RequestID())
		}

		return result, nil
	case messages.MessageTypeRegistered:
		registered := msg.(*messages.Registered)
		_, exists := w.registerRequests.LoadAndDelete(registered.RequestID())
		if !exists {
			return nil, fmt.Errorf("received REGISTERED for invalid requestID")
		}

		w.registrations.Store(registered.RegistrationID(), registered.RegistrationID())
		return registered, nil
	case messages.MessageTypeUnregistered:
		unregistered := msg.(*messages.Unregistered)
		registrationID, exists := w.unregisterRequests.Load(unregistered.RequestID())
		if !exists {
			return nil, fmt.Errorf("received UNREGISTERED for invalid requestID")
		}

		_, exists = w.registrations.LoadAndDelete(registrationID)
		if !exists {
			return nil, fmt.Errorf("received UNREGISTERED for invalid registrationID")
		}

		return unregistered, nil
	case messages.MessageTypeInvocation:
		invocation := msg.(*messages.Invocation)
		_, exists := w.registrations.Load(invocation.RegistrationID())
		if !exists {
			return nil, fmt.Errorf("received INVOCATION for invalid registrationID")
		}

		w.invocationRequests.Store(invocation.RequestID(), invocation.RequestID())

		return invocation, nil
	case messages.MessageTypePublished:
		published := msg.(*messages.Published)
		_, exists := w.publishRequests.LoadAndDelete(published.RequestID())
		if !exists {
			return nil, fmt.Errorf("received PUBLISHED for invalid requestID")
		}

		return published, nil
	case messages.MessageTypeSubscribed:
		subscribed := msg.(*messages.Subscribed)
		_, exists := w.subscribeRequests.Load(subscribed.RequestID())
		if !exists {
			return nil, fmt.Errorf("received SUBSCRIBED for invalid requestID")
		}

		w.subscriptions.Store(subscribed.SubscriptionID(), subscribed.SubscriptionID())

		return subscribed, nil
	case messages.MessageTypeUnsubscribed:
		unsubscribed := msg.(*messages.Unsubscribed)
		subscriptionID, exists := w.unsubscribeRequests.Load(unsubscribed.RequestID())
		if !exists {
			return nil, fmt.Errorf("received UNSUBSCRIBED for invalid requestID")
		}

		_, exists = w.subscriptions.LoadAndDelete(subscriptionID)
		if !exists {
			return nil, fmt.Errorf("received UNSUBSCRIBED for invalid subscriptionID %d", subscriptionID)
		}

		return unsubscribed, nil
	case messages.MessageTypeEvent:
		event := msg.(*messages.Event)
		_, exists := w.subscriptions.Load(event.SubscriptionID())
		if !exists {
			return nil, fmt.Errorf("received EVENT for invalid subscriptionID")
		}

		return event, nil
	case messages.MessageTypeError:
		errorMsg := msg.(*messages.Error)
		switch errorMsg.MessageType() {
		case messages.MessageTypeCall:
			_, exists := w.callRequests.LoadAndDelete(errorMsg.RequestID())
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid call request")
			}

			return errorMsg, nil
		case messages.MessageTypeRegister:
			_, exists := w.registerRequests.LoadAndDelete(errorMsg.RequestID())
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid register request")
			}

			return errorMsg, nil
		case messages.MessageTypeUnregister:
			_, exists := w.unregisterRequests.LoadAndDelete(errorMsg.RequestID())
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid unregister request")
			}

			return errorMsg, nil
		case messages.MessageTypeSubscribe:
			_, exists := w.subscribeRequests.LoadAndDelete(errorMsg.RequestID())
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid subscribe request")
			}

			return errorMsg, nil
		case messages.MessageTypeUnsubscribe:
			_, exists := w.unsubscribeRequests.LoadAndDelete(errorMsg.RequestID())
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid unsubscribe request")
			}

			return errorMsg, nil
		case messages.MessageTypePublish:
			_, exists := w.publishRequests.LoadAndDelete(errorMsg.RequestID())
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid publish request")
			}

			return errorMsg, nil
		default:
			return nil, fmt.Errorf("unknown error message type %T", msg)
		}
	case messages.MessageTypeGoodbye:
		return msg, nil
	default:
		return nil, fmt.Errorf("unknown message type %T", msg)
	}
}
