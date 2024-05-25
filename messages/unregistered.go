package messages

import "fmt"

const MessageTypeUnRegistered = 67
const MessageNameUnRegistered = "UNREGISTERED"

var unRegisteredValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 2,
	MaxLength: 2,
	Message:   MessageNameUnRegistered,
	Spec: Spec{
		1: ValidateRequestID,
	},
}

type UnRegistered interface {
	Message

	RequestID() int64
}

type unRegistered struct {
	requestID int64
}

func NewEmptyUnRegistered() UnRegistered {
	return &unRegistered{}
}

func NewUnRegistered(requestID int64) UnRegistered {
	return &unRegistered{
		requestID: requestID,
	}
}

func (r *unRegistered) Type() int {
	return MessageTypeUnRegistered
}

func (r *unRegistered) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unRegisteredValidationSpec)
	if err != nil {
		return fmt.Errorf("unregistered: failed to validate message %s: %w", MessageNameUnRegistered, err)
	}

	r.requestID = fields.RequestID

	return nil
}

func (r *unRegistered) Marshal() []any {
	return []any{MessageTypeUnRegistered, r.requestID}
}

func (r *unRegistered) RequestID() int64 {
	return r.requestID
}
