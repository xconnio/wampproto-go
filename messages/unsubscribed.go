package messages

import "fmt"

const MessageTypeUnsubscribed = 35
const MessageNameUnsubscribed = "UNSUBSCRIBED"

var unsubscribedValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 2,
	MaxLength: 2,
	Message:   MessageNameUnsubscribed,
	Spec: Spec{
		1: ValidateRequestID,
	},
}

type UnsubscribedFields interface {
	RequestID() uint64
}

type unsubscribedFields struct {
	requestID uint64
}

func (us *unsubscribedFields) RequestID() uint64 {
	return us.requestID
}

type Unsubscribed struct {
	UnsubscribedFields
}

func NewUnsubscribed(requestID uint64) *Unsubscribed {
	return &Unsubscribed{UnsubscribedFields: &unsubscribedFields{requestID: requestID}}
}

func NewUnsubscribedWithFields(fields UnsubscribedFields) *Unsubscribed {
	return &Unsubscribed{UnsubscribedFields: fields}
}

func (us *Unsubscribed) Type() int {
	return MessageTypeUnsubscribed
}

func (us *Unsubscribed) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unsubscribedValidationSpec)
	if err != nil {
		return fmt.Errorf("unsubscribed: failed to validate message %s: %w", MessageNameUnsubscribed, err)
	}

	us.UnsubscribedFields = &unsubscribedFields{requestID: fields.RequestID}

	return nil
}

func (us *Unsubscribed) Marshal() []any {
	return []any{MessageTypeUnsubscribed, us.RequestID()}
}
