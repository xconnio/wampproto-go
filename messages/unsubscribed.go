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

type UnSubscribed interface {
	Message

	RequestID() int64
}

type unSubscribed struct {
	requestID int64
}

func NewEmptyUnSubscribed() UnSubscribed {
	return &unSubscribed{}
}

func NewUnSubscribed(requestID int64) UnSubscribed {
	return &unSubscribed{
		requestID: requestID,
	}
}

func (r *unSubscribed) Type() int {
	return MessageTypeUnSubscribed
}

func (r *unSubscribed) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unSubscribedValidationSpec)
	if err != nil {
		return fmt.Errorf("unsubscribed: failed to validate message %s: %w", MessageNameUnSubscribed, err)
	}

	r.requestID = fields.RequestID

	return nil
}

func (r *unSubscribed) Marshal() []any {
	return []any{MessageTypeUnSubscribed, r.requestID}
}

func (r *unSubscribed) RequestID() int64 {
	return r.requestID
}
