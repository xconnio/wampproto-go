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
	RequestID() uint64
	SubscriptionID() uint64
}

type subscribedFields struct {
	requestID      uint64
	subscriptionID uint64
}

func (s *subscribedFields) RequestID() uint64 {
	return s.requestID
}

func (s *subscribedFields) SubscriptionID() uint64 {
	return s.subscriptionID
}

type Subscribed struct {
	SubscribedFields
}

func NewSubscribed(requestID, subscriptionID uint64) *Subscribed {
	return &Subscribed{SubscribedFields: &subscribedFields{requestID: requestID, subscriptionID: subscriptionID}}
}

func NewSubscribedWithFields(fields SubscribedFields) *Subscribed {
	return &Subscribed{SubscribedFields: fields}
}

func (s *Subscribed) Type() int {
	return MessageTypeSubscribed
}

func (s *Subscribed) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, subscribedValidationSpec)
	if err != nil {
		return fmt.Errorf("subscribed: failed to validate message %s: %w", MessageNameSubscribed, err)
	}

	s.SubscribedFields = &subscribedFields{requestID: fields.RequestID, subscriptionID: fields.SubscriptionID}

	return nil
}

func (s *Subscribed) Marshal() []any {
	return []any{MessageTypeSubscribed, s.RequestID(), s.SubscriptionID()}
}
