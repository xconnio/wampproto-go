package messages

import "fmt"

const MessageTypeSubscribed = 33
const MessageNameSubscribed = "SUBSCRIBED"

var subscribedValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameSubscribed,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateSubscriptionID,
	},
}

type Subscribed interface {
	Message

	RequestID() int64
	SubscriptionID() int64
}

type subscribed struct {
	requestID      int64
	subscriptionID int64
}

func NewEmptySubscribed() Subscribed {
	return &subscribed{}
}

func NewSubscribed(requestID, subscriptionID int64) Subscribed {
	return &subscribed{
		requestID:      requestID,
		subscriptionID: subscriptionID,
	}
}

func (r *subscribed) Type() int {
	return MessageTypeSubscribed
}

func (r *subscribed) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, subscribedValidationSpec)
	if err != nil {
		return fmt.Errorf("subscribed: failed to validate message %s: %w", MessageNameSubscribed, err)
	}

	r.requestID = fields.RequestID
	r.subscriptionID = fields.SubscriptionID

	return nil
}

func (r *subscribed) Marshal() []any {
	return []any{MessageTypeSubscribed, r.requestID, r.subscriptionID}
}

func (r *subscribed) RequestID() int64 {
	return r.requestID
}

func (r *subscribed) SubscriptionID() int64 {
	return r.subscriptionID
}
