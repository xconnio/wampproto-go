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

type Yield interface {
	Message

	RequestID() int64
	Options() map[string]any
	Args() []any
	KwArgs() map[string]any
}

type yield struct {
	requestID int64
	options   map[string]any
	args      []any
	kwArgs    map[string]any
}

func NewEmptyYield() Yield {
	return &yield{}
}

func NewYield(requestID int64, options map[string]any, procedure string, args []any, kwArgs map[string]any) Yield {
	return &yield{
		requestID: requestID,
		options:   options,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *yield) RequestID() int64 {
	return e.requestID
}

func (e *yield) Options() map[string]any {
	return e.options
}

func (e *yield) Args() []any {
	return e.args
}

func (e *yield) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *yield) Type() int {
	return MessageTypeYield
}

func (e *yield) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, yieldValidationSpec)
	if err != nil {
		return fmt.Errorf("yield: failed to validate message %s: %w", MessageNameYield, err)
	}

	e.requestID = fields.RequestID
	e.options = fields.Options
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *yield) Marshal() []any {
	result := []any{MessageTypeYield, e.requestID, e.options}

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
