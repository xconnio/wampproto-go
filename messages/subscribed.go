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

type SubscribedFields interface {
	RequestID() int64
	SubscriptionID() int64
}

type subscribedFields struct {
	requestID      int64
	subscriptionID int64
}

func NewSubscribedFields(requestID, subscriptionID int64) SubscribedFields {
	return &subscribedFields{
		requestID:      requestID,
		subscriptionID: subscriptionID,
	}
}

func (s *subscribedFields) RequestID() int64 {
	return s.requestID
}

func (s *subscribedFields) SubscriptionID() int64 {
	return s.subscriptionID
}

type Subscribed struct {
	SubscribedFields
}

func NewSubscribed(fields SubscribedFields) *Subscribed {
	return &Subscribed{
		SubscribedFields: fields,
	}
}

func (s *Subscribed) Type() int {
	return MessageTypeSubscribed
}

func (s *Subscribed) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, subscribedValidationSpec)
	if err != nil {
		return fmt.Errorf("subscribed: failed to validate message %s: %w", MessageNameSubscribed, err)
	}

	s.SubscribedFields = NewSubscribedFields(fields.RequestID, fields.SubscriptionID)

	return nil
}

func (s *Subscribed) Marshal() []any {
	return []any{MessageTypeSubscribed, s.RequestID(), s.SubscriptionID()}
}
