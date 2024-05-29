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
	RequestID() int64
	RegistrationID() int64
}

type registeredFields struct {
	requestID      int64
	registrationID int64
}

func NewRegisteredFields(requestID, registrationID int64) RegisteredFields {
	return &registeredFields{
		requestID:      requestID,
		registrationID: registrationID,
	}
}

func (r *registeredFields) RequestID() int64 {
	return r.requestID
}

func (r *registeredFields) RegistrationID() int64 {
	return r.registrationID
}

type Registered struct {
	RegisteredFields
}

func NewRegistered(requestID, registrationID int64) *Registered {
	return &Registered{RegisteredFields: NewRegisteredFields(requestID, registrationID)}
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

	r.RegisteredFields = NewRegisteredFields(fields.RequestID, fields.RegistrationID)

	return nil
}

func (r *Registered) Marshal() []any {
	return []any{MessageTypeRegistered, r.RequestID(), r.RegistrationID()}
}
