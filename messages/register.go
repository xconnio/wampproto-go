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

type Register interface {
	Message

	RequestID() int64
	Options() map[string]any
	Procedure() string
}

type register struct {
	requestID int64
	options   map[string]any
	procedure string
}

func NewEmptyRegister() Register {
	return &register{}
}

func NewRegister(requestID int64, options map[string]any, uri string) Register {
	return &register{
		requestID: requestID,
		options:   options,
		procedure: uri,
	}
}

func (r *register) Type() int {
	return MessageTypeRegister
}

func (r *register) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, registerValidationSpec)
	if err != nil {
		return fmt.Errorf("register: failed to validate message %s: %w", MessageNameRegister, err)
	}

	r.requestID = fields.SessionID
	r.options = fields.Options
	r.procedure = fields.URI

	return nil
}

func (r *register) Marshal() []any {
	return []any{MessageTypeRegister, r.requestID, r.options, r.procedure}
}

func (r *register) RequestID() int64 {
	return r.requestID
}

func (r *register) Options() map[string]any {
	return r.options
}

func (r *register) Procedure() string {
	return r.procedure
}
