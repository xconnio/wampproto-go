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

type Registered interface {
	Message

	RequestID() int64
	RegistrationID() int64
}

type registered struct {
	requestID      int64
	registrationID int64
}

func NewEmptyRegistered() Registered {
	return &registered{}
}

func NewRegistered(requestID, registrationID int64) Registered {
	return &registered{
		requestID:      requestID,
		registrationID: registrationID,
	}
}

func (r *registered) Type() int {
	return MessageTypeRegistered
}

func (r *registered) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, registeredValidationSpec)
	if err != nil {
		return fmt.Errorf("registered: failed to validate message %s: %w", MessageNameRegistered, err)
	}

	r.requestID = fields.RequestID
	r.registrationID = fields.RegistrationID

	return nil
}

func (r *registered) Marshal() []any {
	return []any{MessageTypeRegistered, r.requestID, r.registrationID}
}

func (r *registered) RequestID() int64 {
	return r.requestID
}

func (r *registered) RegistrationID() int64 {
	return r.registrationID
}
