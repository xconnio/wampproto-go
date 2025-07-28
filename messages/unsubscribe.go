package messages

import "fmt"

const MessageTypeUnsubscribe uint64 = 34
const MessageNameUnsubscribe = "UNSUBSCRIBE"

var unsubscribeValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameUnsubscribe,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateSubscriptionID,
	},
}

type UnsubscribeFields interface {
	RequestID() uint64
	SubscriptionID() uint64
}

type unsubscribeFields struct {
	requestID      uint64
	subscriptionID uint64
}

func (us *unsubscribeFields) RequestID() uint64 {
	return us.requestID
}

func (us *unsubscribeFields) SubscriptionID() uint64 {
	return us.subscriptionID
}

type Unsubscribe struct {
	UnsubscribeFields
}

func NewUnsubscribe(requestID, subscriptionID uint64) *Unsubscribe {
	return &Unsubscribe{UnsubscribeFields: &unsubscribeFields{requestID: requestID, subscriptionID: subscriptionID}}
}

func NewUnsubscribeWithFields(fields UnsubscribeFields) *Unsubscribe {
	return &Unsubscribe{UnsubscribeFields: fields}
}

func (us *Unsubscribe) Type() uint64 {
	return MessageTypeUnsubscribe
}

func (us *Unsubscribe) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unsubscribeValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnsubscribe, err)
	}

	us.UnsubscribeFields = &unsubscribeFields{requestID: fields.RequestID, subscriptionID: fields.SubscriptionID}

	return nil
}

func (us *Unsubscribe) Marshal() []any {
	return []any{MessageTypeUnsubscribe, us.RequestID(), us.SubscriptionID()}
}
