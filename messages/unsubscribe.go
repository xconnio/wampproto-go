package messages

import "fmt"

const MessageTypeUnSubscribe = 34
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

func (us *unSubscribeFields) RequestID() int64 {
	return us.requestID
}

func (us *unSubscribeFields) SubscriptionID() int64 {
	return us.subscriptionID
}

type UnSubscribe struct {
	UnSubscribeFields
}

func NewUnSubscribe(requestID, subscriptionID int64) *UnSubscribe {
	return &UnSubscribe{UnSubscribeFields: &unSubscribeFields{requestID: requestID, subscriptionID: subscriptionID}}
}

func NewUnSubscribeWithFields(fields UnSubscribeFields) *UnSubscribe {
	return &UnSubscribe{UnSubscribeFields: fields}
}

func (us *UnSubscribe) Type() int {
	return MessageTypeUnSubscribe
}

func (us *UnSubscribe) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unSubscribeValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnSubscribe, err)
	}

	us.UnSubscribeFields = &unSubscribeFields{requestID: fields.RequestID, subscriptionID: fields.SubscriptionID}

	return nil
}

func (us *UnSubscribe) Marshal() []any {
	return []any{MessageTypeUnSubscribe, us.RequestID(), us.SubscriptionID()}
}
