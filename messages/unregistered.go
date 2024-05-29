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

type UnRegisteredFields interface {
	RequestID() int64
}

type unRegisteredFields struct {
	requestID int64
}

func (ur *unRegisteredFields) RequestID() int64 {
	return ur.requestID
}

type UnRegistered struct {
	UnRegisteredFields
}

func NewUnRegistered(requestID int64) *UnRegistered {
	return &UnRegistered{UnRegisteredFields: &unRegisteredFields{requestID: requestID}}
}

func NewUnRegisteredWithFields(fields UnRegisteredFields) *UnRegistered {
	return &UnRegistered{UnRegisteredFields: fields}
}

func (ur *UnRegistered) Type() int {
	return MessageTypeUnRegistered
}

func (ur *UnRegistered) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unRegisteredValidationSpec)
	if err != nil {
		return fmt.Errorf("unregistered: failed to validate message %s: %w", MessageNameUnRegistered, err)
	}

	ur.UnRegisteredFields = &unRegisteredFields{requestID: fields.RequestID}

	return nil
}

func (ur *UnRegistered) Marshal() []any {
	return []any{MessageTypeUnRegistered, ur.RequestID()}
}
