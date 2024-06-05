package messages

import "fmt"

const MessageTypeResult = 50
const MessageNameResult = "RESULT"

var resultValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 5,
	Message:   MessageNameResult,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateDetails,
		3: ValidateArgs,
		4: ValidateKwArgs,
	},
}

type ResultFields interface {
	RequestID() int64
	Details() map[string]any
	Args() []any
	KwArgs() map[string]any

	BinaryPayload
}

type resultFields struct {
	requestID int64
	details   map[string]any
	args      []any
	kwArgs    map[string]any
}

func (e *resultFields) RequestID() int64 {
	return e.requestID
}

func (e *resultFields) Details() map[string]any {
	return e.details
}

func (e *resultFields) Args() []any {
	return e.args
}

func (e *resultFields) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *resultFields) PayloadIsBinary() bool {
	return false
}

func (e *resultFields) Payload() []byte {
	return nil
}

func (e *resultFields) PayloadSerializer() int {
	return 0
}

type Result struct {
	ResultFields
}

func NewResult(requestID int64, details map[string]any, args []any, kwArgs map[string]any) *Result {
	if details == nil {
		details = make(map[string]any)
	}

	return &Result{ResultFields: &resultFields{requestID: requestID, details: details, args: args, kwArgs: kwArgs}}
}

func NewResultWithFields(field ResultFields) *Result {
	return &Result{ResultFields: field}
}

func (e *Result) Type() int {
	return MessageTypeResult
}

func (e *Result) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, resultValidationSpec)
	if err != nil {
		return fmt.Errorf("result: failed to validate message %s: %w", MessageNameResult, err)
	}

	e.ResultFields = &resultFields{
		requestID: fields.RequestID,
		details:   fields.Details,
		args:      fields.Args,
		kwArgs:    fields.KwArgs,
	}

	return nil
}

func (e *Result) Marshal() []any {
	result := []any{MessageTypeResult, e.RequestID(), e.Details()}

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
