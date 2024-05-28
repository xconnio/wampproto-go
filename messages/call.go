package messages

import "fmt"

const MessageTypeCall = 48
const MessageNameCall = "CALL"

var callValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 6,
	Message:   MessageNameCall,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
		3: ValidateURI,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type CallFields interface {
	RequestID() int64
	Options() map[string]any
	Procedure() string
	Args() []any
	KwArgs() map[string]any
}

type callFields struct {
	requestID int64
	options   map[string]any
	procedure string
	args      []any
	kwArgs    map[string]any
}

func NewCallFields(requestID int64, options map[string]any, procedure string, args []any,
	kwArgs map[string]any) CallFields {
	return &callFields{
		requestID: requestID,
		options:   options,
		procedure: procedure,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *callFields) RequestID() int64 {
	return e.requestID
}

func (e *callFields) Options() map[string]any {
	return e.options
}

func (e *callFields) Procedure() string {
	return e.procedure
}

func (e *callFields) Args() []any {
	return e.args
}

func (e *callFields) KwArgs() map[string]any {
	return e.kwArgs
}

type Call struct {
	CallFields
}

func NewCall(fields CallFields) *Call {
	return &Call{CallFields: fields}
}

func (e *Call) Type() int {
	return MessageTypeCall
}

func (e *Call) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, callValidationSpec)
	if err != nil {
		return fmt.Errorf("call: failed to validate message %s: %w", MessageNameCall, err)
	}

	e.CallFields = NewCallFields(fields.RequestID, fields.Options, fields.URI, fields.Args, fields.KwArgs)

	return nil
}

func (e *Call) Marshal() []any {
	result := []any{MessageTypeCall, e.RequestID(), e.Options(), e.Procedure()}

	if e.Args() != nil {
		result = append(result, e.Args())
	}

	if e.KwArgs() != nil {
		if e.Args() == nil {
			result = append(result, []any{})
		}

		result = append(result, e.KwArgs())
	}

	return result
}
