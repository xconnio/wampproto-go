package messages

import "fmt"

const MessageTypeYield uint64 = 70
const MessageNameYield = "YIELD"

var yieldValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 5,
	Message:   MessageNameYield,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
		3: ValidateArgs,
		4: ValidateKwArgs,
	},
}

type YieldFields interface {
	RequestID() uint64
	Options() map[string]any
	Args() []any
	KwArgs() map[string]any

	BinaryPayload
}

type yieldFields struct {
	requestID uint64
	options   map[string]any
	args      []any
	kwArgs    map[string]any

	binary     bool
	serializer uint64
	payload    []byte
}

func (e *yieldFields) RequestID() uint64 {
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

func (e *yieldFields) PayloadIsBinary() bool {
	return e.binary
}

func (e *yieldFields) Payload() []byte {
	return e.payload
}

func (e *yieldFields) PayloadSerializer() uint64 {
	return e.serializer
}

type Yield struct {
	YieldFields
}

func NewYield(requestID uint64, options map[string]any, args []any, kwArgs map[string]any) *Yield {
	if options == nil {
		options = make(map[string]any)
	}

	return &Yield{YieldFields: &yieldFields{requestID: requestID, options: options, args: args, kwArgs: kwArgs}}
}

func NewYieldWithFields(fields YieldFields) *Yield {
	return &Yield{YieldFields: fields}
}

func NewYieldBinary(requestID uint64, options map[string]any, payload []byte, serializer uint64) *Yield {
	if options == nil {
		options = make(map[string]any)
	}

	return &Yield{YieldFields: &yieldFields{
		requestID:  requestID,
		options:    options,
		binary:     true,
		serializer: serializer,
		payload:    payload,
	}}
}

func (e *Yield) Type() uint64 {
	return MessageTypeYield
}

func (e *Yield) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, yieldValidationSpec)
	if err != nil {
		return fmt.Errorf("yield: failed to validate message %s: %w", MessageNameYield, err)
	}

	e.YieldFields = &yieldFields{
		requestID: fields.RequestID,
		options:   fields.Options,
		args:      fields.Args,
		kwArgs:    fields.KwArgs,
	}

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
