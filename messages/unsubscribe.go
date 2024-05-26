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

type UnSubscribe interface {
	Message

	RequestID() int64
	SubscriptionID() int64
}

type unSubscribe struct {
	requestID      int64
	subscriptionID int64
}

func NewEmptyUnSubscribe() UnSubscribe {
	return &unSubscribe{}
}

func NewUnSubscribe(requestID, subscriptionID int64) UnSubscribe {
	return &unSubscribe{
		requestID:      requestID,
		subscriptionID: subscriptionID,
	}
}

func (r *unSubscribe) Type() int {
	return MessageTypeUnSubscribe
}

func (r *unSubscribe) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unSubscribeValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnSubscribe, err)
	}

	r.requestID = fields.RequestID
	r.subscriptionID = fields.SubscriptionID

	return nil
}

func (r *unSubscribe) Marshal() []any {
	return []any{MessageTypeUnSubscribe, r.requestID, r.subscriptionID}
}

func (r *unSubscribe) RequestID() int64 {
	return r.requestID
}

func (r *unSubscribe) SubscriptionID() int64 {
	return r.subscriptionID
}
