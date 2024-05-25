package messages

import "fmt"

const MessageTypeUnRegister = 66
const MessageNameUnRegister = "UNREGISTER"

var unRegisterValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameUnRegister,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateRegistrationID,
	},
}

type UnRegister interface {
	Message

	RequestID() int64
	RegistrationID() int64
}

type unRegister struct {
	requestID      int64
	registrationID int64
}

func NewEmptyUnRegister() UnRegister {
	return &unRegister{}
}

func NewUnRegister(requestID, registrationID int64) UnRegister {
	return &unRegister{
		requestID:      requestID,
		registrationID: registrationID,
	}
}

func (r *unRegister) Type() int {
	return MessageTypeUnRegister
}

func (r *unRegister) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unRegisterValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnRegister, err)
	}

	r.requestID = fields.RequestID
	r.registrationID = fields.RegistrationID

	return nil
}

func (r *unRegister) Marshal() []any {
	return []any{MessageTypeUnRegister, r.requestID, r.registrationID}
}

func (r *unRegister) RequestID() int64 {
	return r.requestID
}

func (r *unRegister) RegistrationID() int64 {
	return r.registrationID
}
