package messages

import "fmt"

const MessageTypeUnSubscribe = 33
const MessageNameUnSubscribe = "UNSUBSCRIBE"

var unSubscribeValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameUnSubscribe,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateSubscriptionID,
	},
}

type UnSubscribeFields interface {
	RequestID() int64
	SubscriptionID() int64
}

type unSubscribeFields struct {
	requestID      int64
	subscriptionID int64
}

func NewUnSubscribeFields(requestID, subscriptionID int64) UnSubscribeFields {
	return &unSubscribeFields{
		requestID:      requestID,
		subscriptionID: subscriptionID,
	}
}

func (us *unSubscribeFields) RequestID() int64 {
	return us.requestID
}

func (us *unSubscribeFields) SubscriptionID() int64 {
	return us.subscriptionID
}

type UnSubscribe struct {
	UnSubscribeFields
}

func NewUnSubscribe(fields UnSubscribeFields) *UnSubscribe {
	return &UnSubscribe{
		UnSubscribeFields: fields,
	}
}

func (us *UnSubscribe) Type() int {
	return MessageTypeUnSubscribe
}

func (us *UnSubscribe) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unSubscribeValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnSubscribe, err)
	}

	us.UnSubscribeFields = NewUnSubscribeFields(fields.RequestID, fields.SubscriptionID)

	return nil
}

func (us *UnSubscribe) Marshal() []any {
	return []any{MessageTypeUnSubscribe, us.RequestID(), us.SubscriptionID()}
}
