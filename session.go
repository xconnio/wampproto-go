package wampproto

import (
	"fmt"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

type Session struct {
	serializer serializers.Serializer

	// data structures for RPC
	callRequests       map[int64]int64
	registerRequests   map[int64]int64
	registrations      map[int64]int64
	invocationRequests map[int64]int64
	unregisterRequests map[int64]int64

	// data structures for PubSub
	publishRequests     map[int64]int64
	subscribeRequests   map[int64]int64
	subscriptions       map[int64]int64
	unsubscribeRequests map[int64]int64
}

func NewSession(serializer serializers.Serializer) *Session {
	if serializer == nil {
		serializer = &serializers.JSONSerializer{}
	}

	return &Session{
		serializer: serializer,

		callRequests:       make(map[int64]int64),
		registerRequests:   make(map[int64]int64),
		registrations:      make(map[int64]int64),
		invocationRequests: make(map[int64]int64),
		unregisterRequests: make(map[int64]int64),

		publishRequests:     make(map[int64]int64),
		subscribeRequests:   make(map[int64]int64),
		subscriptions:       make(map[int64]int64),
		unsubscribeRequests: make(map[int64]int64),
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
		w.callRequests[call.RequestID()] = call.RequestID()

		return data, nil
	case messages.MessageTypeYield:
		yield := msg.(*messages.Yield)
		delete(w.invocationRequests, yield.RequestID())

		return data, nil
	case messages.MessageTypeRegister:
		register := msg.(*messages.Register)
		w.registerRequests[register.RequestID()] = register.RequestID()

		return data, nil
	case messages.MessageTypeUnRegister:
		unregister := msg.(*messages.UnRegister)
		w.unregisterRequests[unregister.RequestID()] = unregister.RequestID()

		return data, nil
	case messages.MessageTypePublish:
		publish := msg.(*messages.Publish)
		acknowledge, ok := publish.Options()["acknowledge"].(bool)
		if ok && acknowledge {
			w.publishRequests[publish.RequestID()] = publish.RequestID()
		}

		return data, nil
	case messages.MessageTypeSubscribe:
		subscribe := msg.(*messages.Subscribe)
		w.subscribeRequests[subscribe.RequestID()] = subscribe.RequestID()

		return data, nil
	case messages.MessageTypeUnSubscribe:
		unsubscribe := msg.(*messages.UnSubscribe)
		_, exists := w.subscriptions[unsubscribe.SubscriptionID()]
		if !exists {
			return nil, fmt.Errorf("unsubscribe request for non existent subscription %d",
				unsubscribe.SubscriptionID())
		}

		w.unsubscribeRequests[unsubscribe.RequestID()] = unsubscribe.RequestID()

		return data, nil
	case messages.MessageTypeError:
		errorMsg := msg.(*messages.Error)
		if errorMsg.MessageType() != messages.MessageTypeInvocation {
			return nil, fmt.Errorf("send only supported for invocation error")
		}

		delete(w.invocationRequests, errorMsg.RequestID())
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
		_, exists := w.callRequests[result.RequestID()]
		if !exists {
			return nil, fmt.Errorf("received RESULT for invalid requestID")
		}

		delete(w.callRequests, result.RequestID())
		return result, nil
	case messages.MessageTypeRegistered:
		registered := msg.(*messages.Registered)
		_, exists := w.registerRequests[registered.RequestID()]
		if !exists {
			return nil, fmt.Errorf("received REGISTERED for invalid requestID")
		}

		delete(w.registerRequests, registered.RequestID())
		w.registrations[registered.RegistrationID()] = registered.RegistrationID()
		return registered, nil
	case messages.MessageTypeUnRegistered:
		unregistered := msg.(*messages.UnRegistered)
		registrationID, exists := w.unregisterRequests[unregistered.RequestID()]
		if !exists {
			return nil, fmt.Errorf("received UNREGISTERED for invalid requestID")
		}

		_, exists = w.registrations[registrationID]
		if !exists {
			return nil, fmt.Errorf("received UNREGISTERED for invalid registrationID")
		}

		delete(w.registrations, registrationID)
		return unregistered, nil
	case messages.MessageTypeInvocation:
		invocation := msg.(*messages.Invocation)
		_, exists := w.registrations[invocation.RegistrationID()]
		if !exists {
			return nil, fmt.Errorf("received INVOCATION for invalid registrationID")
		}

		w.invocationRequests[invocation.RequestID()] = invocation.RequestID()

		return invocation, nil
	case messages.MessageTypePublished:
		published := msg.(*messages.Published)
		_, exists := w.publishRequests[published.RequestID()]
		if !exists {
			return nil, fmt.Errorf("received PUBLISHED for invalid requestID")
		}

		delete(w.publishRequests, published.RequestID())

		return published, nil
	case messages.MessageTypeSubscribed:
		subscribed := msg.(*messages.Subscribed)
		_, exists := w.subscribeRequests[subscribed.RequestID()]
		if !exists {
			return nil, fmt.Errorf("received SUBSCRIBED for invalid requestID")
		}

		w.subscriptions[subscribed.SubscriptionID()] = subscribed.SubscriptionID()

		return subscribed, nil
	case messages.MessageTypeUnSubscribed:
		unsubscribed := msg.(*messages.UnSubscribed)
		subscriptionID, exists := w.unsubscribeRequests[unsubscribed.RequestID()]
		if !exists {
			return nil, fmt.Errorf("received UNSUBSCRIBED for invalid requestID")
		}

		_, exists = w.subscriptions[subscriptionID]
		if !exists {
			return nil, fmt.Errorf("received UNSUBSCRIBED for invalid subscriptionID")
		}

		delete(w.subscriptions, subscriptionID)

		return unsubscribed, nil
	case messages.MessageTypeEvent:
		event := msg.(*messages.Event)
		_, exists := w.subscriptions[event.SubscriptionID()]
		if !exists {
			return nil, fmt.Errorf("received EVENT for invalid subscriptionID")
		}

		return event, nil
	case messages.MessageTypeError:
		errorMsg := msg.(*messages.Error)
		switch errorMsg.MessageType() {
		case messages.MessageTypeCall:
			_, exists := w.callRequests[errorMsg.RequestID()]
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid call request")
			}

			delete(w.callRequests, errorMsg.RequestID())
			return errorMsg, nil
		case messages.MessageTypeRegister:
			_, exists := w.registerRequests[errorMsg.RequestID()]
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid register request")
			}

			delete(w.registerRequests, errorMsg.RequestID())
			return errorMsg, nil
		case messages.MessageTypeUnRegister:
			_, exists := w.unregisterRequests[errorMsg.RequestID()]
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid unregister request")
			}

			delete(w.unregisterRequests, errorMsg.RequestID())
			return errorMsg, nil
		case messages.MessageTypeSubscribe:
			_, exists := w.subscribeRequests[errorMsg.RequestID()]
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid subscribe request")
			}

			delete(w.subscribeRequests, errorMsg.RequestID())
			return errorMsg, nil
		case messages.MessageTypeUnSubscribe:
			_, exists := w.unsubscribeRequests[errorMsg.RequestID()]
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid unsubscribe request")
			}

			delete(w.unsubscribeRequests, errorMsg.RequestID())
			return errorMsg, nil
		case messages.MessageTypePublish:
			_, exists := w.publishRequests[errorMsg.RequestID()]
			if !exists {
				return nil, fmt.Errorf("received ERROR for invalid publish request")
			}

			delete(w.publishRequests, errorMsg.RequestID())
			return errorMsg, nil
		default:
			return nil, fmt.Errorf("unknown error message type %T", msg)
		}
	default:
		return nil, fmt.Errorf("unknown message type %T", msg)
	}
}
