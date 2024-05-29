package messages

import "fmt"

const MessageTypeUnSubscribed = 35
const MessageNameUnSubscribed = "UNSUBSCRIBED"

var unSubscribedValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 2,
	MaxLength: 2,
	Message:   MessageNameUnSubscribed,
	Spec: Spec{
		1: ValidateRequestID,
	},
}

type UnSubscribedFields interface {
	RequestID() int64
}

type unSubscribedFields struct {
	requestID int64
}

func NewUnSubscribedFields(requestID int64) UnSubscribedFields {
	return &unSubscribedFields{
		requestID: requestID,
	}
}

func (us *unSubscribedFields) RequestID() int64 {
	return us.requestID
}

type UnSubscribed struct {
	UnSubscribedFields
}

func NewUnSubscribed(requestID int64) *UnSubscribed {
	return &UnSubscribed{UnSubscribedFields: NewUnSubscribedFields(requestID)}
}

func NewUnSubscribedWithFields(fields UnSubscribedFields) *UnSubscribed {
	return &UnSubscribed{UnSubscribedFields: fields}
}

func (us *UnSubscribed) Type() int {
	return MessageTypeUnSubscribed
}

func (us *UnSubscribed) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unSubscribedValidationSpec)
	if err != nil {
		return fmt.Errorf("unsubscribed: failed to validate message %s: %w", MessageNameUnSubscribed, err)
	}

	us.UnSubscribedFields = NewUnSubscribedFields(fields.RequestID)

	return nil
}

func (us *UnSubscribed) Marshal() []any {
	return []any{MessageTypeUnSubscribed, us.RequestID()}
}
