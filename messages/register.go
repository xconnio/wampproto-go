package messages

import "fmt"

const MessageTypeRegister = 64
const MessageNameRegister = "REGISTER"

var registerValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 4,
	Message:   MessageNameRegister,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
		3: ValidateURI,
	},
}

type RegisterFields interface {
	RequestID() int64
	Options() map[string]any
	Procedure() string
}

type registerFields struct {
	requestID int64
	options   map[string]any
	procedure string
}

func NewRegisterFields(requestID int64, options map[string]any, uri string) RegisterFields {
	return &registerFields{
		requestID: requestID,
		options:   options,
		procedure: uri,
	}
}

func (r *registerFields) RequestID() int64 {
	return r.requestID
}

func (r *registerFields) Options() map[string]any {
	return r.options
}

func (r *registerFields) Procedure() string {
	return r.procedure
}

type Register struct {
	RegisterFields
}

func NewRegister(fields RegisterFields) *Register {
	return &Register{RegisterFields: fields}
}

func (r *Register) Type() int {
	return MessageTypeRegister
}

func (r *Register) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, registerValidationSpec)
	if err != nil {
		return fmt.Errorf("registerFields: failed to validate message %s: %w", MessageNameRegister, err)
	}

	r.RegisterFields = NewRegisterFields(fields.SessionID, fields.Options, fields.URI)

	return nil
}

func (r *Register) Marshal() []any {
	return []any{MessageTypeRegister, r.RequestID(), r.Options(), r.Procedure()}
}
