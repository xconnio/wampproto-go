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

type Call interface {
	Message

	RequestID() int64
	Options() map[string]any
	Procedure() string
	Args() []any
	KwArgs() map[string]any
}

type call struct {
	requestID int64
	options   map[string]any
	procedure string
	args      []any
	kwArgs    map[string]any
}

func NewEmptyCall() Call {
	return &call{}
}

func NewCall(requestID int64, options map[string]any, procedure string, args []any, kwArgs map[string]any) Call {
	return &call{
		requestID: requestID,
		options:   options,
		procedure: procedure,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *call) RequestID() int64 {
	return e.requestID
}

func (e *call) Options() map[string]any {
	return e.options
}

func (e *call) Procedure() string {
	return e.procedure
}

func (e *call) Args() []any {
	return e.args
}

func (e *call) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *call) Type() int {
	return MessageTypeCall
}

func (e *call) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, callValidationSpec)
	if err != nil {
		return fmt.Errorf("call: failed to validate message %s: %w", MessageNameCall, err)
	}

	e.requestID = fields.RequestID
	e.options = fields.Options
	e.procedure = fields.URI
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *call) Marshal() []any {
	result := []any{MessageTypeCall, e.requestID, e.options, e.procedure}

	if e.args != nil {
		result = append(result, e.args)
	}

	if e.kwArgs != nil {
		if e.args == nil {
			result = append(result, []any{})
		}

		result = append(result, e.kwArgs)
	}

	return result
}
