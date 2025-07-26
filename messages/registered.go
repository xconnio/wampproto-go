package messages

import "fmt"

const MessageTypeRegistered = 65
const MessageNameRegistered = "REGISTERED"

var registeredValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameRegistered,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateRegistrationID,
	},
}

type RegisteredFields interface {
	RequestID() uint64
	RegistrationID() uint64
}

type registeredFields struct {
	requestID      uint64
	registrationID uint64
}

func (r *registeredFields) RequestID() uint64 {
	return r.requestID
}

func (r *registeredFields) RegistrationID() uint64 {
	return r.registrationID
}

type Registered struct {
	RegisteredFields
}

func NewRegistered(requestID, registrationID uint64) *Registered {
	return &Registered{RegisteredFields: &registeredFields{requestID: requestID, registrationID: registrationID}}
}

func NewRegisteredWithFields(fields RegisteredFields) *Registered {
	return &Registered{RegisteredFields: fields}
}

func (r *Registered) Type() int {
	return MessageTypeRegistered
}

func (r *Registered) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, registeredValidationSpec)
	if err != nil {
		return fmt.Errorf("registered: failed to validate message %s: %w", MessageNameRegistered, err)
	}

	r.RegisteredFields = &registeredFields{requestID: fields.RequestID, registrationID: fields.RegistrationID}

	return nil
}

func (r *Registered) Marshal() []any {
	return []any{MessageTypeRegistered, r.RequestID(), r.RegistrationID()}
}
