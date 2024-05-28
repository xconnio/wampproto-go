package messages

import "fmt"

const MessageTypeYield = 70
const MessageNameYield = "YIELD"

var yieldValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 5,
	Message:   MessageNameYield,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type YieldFields interface {
	RequestID() int64
	Options() map[string]any
	Args() []any
	KwArgs() map[string]any
}

type yieldFields struct {
	requestID int64
	options   map[string]any
	args      []any
	kwArgs    map[string]any
}

func NewYieldFields(requestID int64, options map[string]any, args []any, kwArgs map[string]any) YieldFields {
	return &yieldFields{
		requestID: requestID,
		options:   options,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *yieldFields) RequestID() int64 {
	return e.requestID
}

func (e *yieldFields) Options() map[string]any {
	return e.options
}

func (e *yieldFields) Args() []any {
	return e.args
}

func (e *yieldFields) KwArgs() map[string]any {
	return e.kwArgs
}

type Yield struct {
	YieldFields
}

func NewYield(fields YieldFields) *Yield {
	return &Yield{YieldFields: fields}
}

func (e *Yield) Type() int {
	return MessageTypeYield
}

func (e *Yield) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, yieldValidationSpec)
	if err != nil {
		return fmt.Errorf("yield: failed to validate message %s: %w", MessageNameYield, err)
	}

	e.YieldFields = NewYieldFields(fields.RequestID, fields.Options, fields.Args, fields.KwArgs)

	return nil
}

func (e *Yield) Marshal() []any {
	result := []any{MessageTypeYield, e.RequestID(), e.Options()}

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
