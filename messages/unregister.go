package messages

import "fmt"

const MessageTypeUnregister = 66
const MessageNameUnregister = "UNREGISTER"

var unregisterValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameUnregister,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateRegistrationID,
	},
}

type UnregisterFields interface {
	RequestID() uint64
	RegistrationID() uint64
}

type unregisterFields struct {
	requestID      uint64
	registrationID uint64
}

func (ur *unregisterFields) RequestID() uint64 {
	return ur.requestID
}

func (ur *unregisterFields) RegistrationID() uint64 {
	return ur.registrationID
}

type Unregister struct {
	UnregisterFields
}

func NewUnregister(requestID, registrationID uint64) *Unregister {
	return &Unregister{UnregisterFields: &unregisterFields{requestID: requestID, registrationID: registrationID}}
}

func NewUnregisterWithFields(fields UnregisterFields) *Unregister {
	return &Unregister{UnregisterFields: fields}
}

func (ur *Unregister) Type() int {
	return MessageTypeUnregister
}

func (ur *Unregister) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unregisterValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnregister, err)
	}

	ur.UnregisterFields = &unregisterFields{requestID: fields.RequestID, registrationID: fields.RegistrationID}

	return nil
}

func (ur *Unregister) Marshal() []any {
	return []any{MessageTypeUnregister, ur.RequestID(), ur.RegistrationID()}
}
