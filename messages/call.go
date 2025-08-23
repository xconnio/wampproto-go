package messages

import "fmt"

const MessageTypeCall uint64 = 48
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
	RequestID() uint64
	Options() map[string]any
	Procedure() string
	Args() []any
	KwArgs() map[string]any

	BinaryPayload
}

type callFields struct {
	requestID uint64
	options   map[string]any
	procedure string
	args      []any
	kwArgs    map[string]any

	binary     bool
	serializer uint64
	payload    []byte
}

func (e *callFields) RequestID() uint64 {
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

func (e *callFields) PayloadIsBinary() bool {
	return e.binary
}

func (e *callFields) Payload() []byte {
	return e.payload
}

func (e *callFields) PayloadSerializer() uint64 {
	return e.serializer
}

type Call struct {
	CallFields
}

func NewCall(requestID uint64, options map[string]any, procedure string, args []any, kwArgs map[string]any) *Call {
	if options == nil {
		options = make(map[string]any)
	}

	return &Call{CallFields: &callFields{
		requestID: requestID,
		options:   options,
		procedure: procedure,
		args:      args,
		kwArgs:    kwArgs,
	}}
}

func NewCallWithFields(fields CallFields) *Call { return &Call{CallFields: fields} }

func NewCallBinary(requestID uint64, options map[string]any, procedure string, payload []byte,
	serializer uint64) *Call {
	if options == nil {
		options = make(map[string]any)
	}

	options["x_payload_serializer"] = serializer

	return &Call{CallFields: &callFields{
		requestID:  requestID,
		options:    options,
		procedure:  procedure,
		binary:     true,
		payload:    payload,
		serializer: serializer,
	}}
}

func (e *Call) Type() uint64 {
	return MessageTypeCall
}

func (e *Call) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, callValidationSpec)
	if err != nil {
		return fmt.Errorf("call: failed to validate message %s: %w", MessageNameCall, err)
	}

	e.CallFields = &callFields{
		requestID: fields.RequestID,
		options:   fields.Options,
		procedure: fields.URI,
		args:      fields.Args,
		kwArgs:    fields.KwArgs,
	}

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
