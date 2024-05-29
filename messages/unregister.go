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

type UnRegisterFields interface {
	RequestID() int64
	RegistrationID() int64
}

type unRegisterFields struct {
	requestID      int64
	registrationID int64
}

func NewUnRegisterFields(requestID, registrationID int64) UnRegisterFields {
	return &unRegisterFields{
		requestID:      requestID,
		registrationID: registrationID,
	}
}

func (ur *unRegisterFields) RequestID() int64 {
	return ur.requestID
}

func (ur *unRegisterFields) RegistrationID() int64 {
	return ur.registrationID
}

type UnRegister struct {
	UnRegisterFields
}

func NewUnRegister(requestID, registrationID int64) *UnRegister {
	return &UnRegister{UnRegisterFields: NewUnRegisterFields(requestID, registrationID)}
}

func NewUnRegisterWithFields(fields UnRegisterFields) *UnRegister {
	return &UnRegister{UnRegisterFields: fields}
}

func (ur *UnRegister) Type() int {
	return MessageTypeUnRegister
}

func (ur *UnRegister) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, unRegisterValidationSpec)
	if err != nil {
		return fmt.Errorf("unregister: failed to validate message %s: %w", MessageNameUnRegister, err)
	}

	ur.UnRegisterFields = NewUnRegisterFields(fields.RequestID, fields.RegistrationID)

	return nil
}

func (ur *UnRegister) Marshal() []any {
	return []any{MessageTypeUnRegister, ur.RequestID(), ur.RegistrationID()}
}
