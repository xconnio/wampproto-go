package messages

import "fmt"

const MessageTypeUnregistered = 67
const MessageNameUnregistered = "UNREGISTERED"

var unregisteredValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 2,
	MaxLength: 2,
	Message:   MessageNameUnregistered,
	Spec: Spec{
		1: ValidateRequestID,
	},
}

type UnregisteredFields interface {
	RequestID() uint64
}

type unregisteredFields struct {
	requestID uint64
}

func (ur *unregisteredFields) RequestID() uint64 {
	return ur.requestID
}

type Unregistered struct {
	UnregisteredFields
}

func NewUnregistered(requestID uint64) *Unregistered {
	return &Unregistered{UnregisteredFields: &unregisteredFields{requestID: requestID}}
}

func NewUnregisteredWithFields(fields UnregisteredFields) *Unregistered {
	return &Unregistered{UnregisteredFields: fields}
}

func (ur *Unregistered) Type() int {
	return MessageTypeUnregistered
}

func (ur *Unregistered) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unregisteredValidationSpec)
	if err != nil {
		return fmt.Errorf("unregistered: failed to validate message %s: %w", MessageNameUnregistered, err)
	}

	ur.UnregisteredFields = &unregisteredFields{requestID: fields.RequestID}

	return nil
}

func (ur *Unregistered) Marshal() []any {
	return []any{MessageTypeUnregistered, ur.RequestID()}
}
