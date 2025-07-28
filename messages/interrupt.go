package messages

import "fmt"

const MessageTypeInterrupt uint64 = 69
const MessageNameInterrupt = "INTERRUPT"

var interruptValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameInterrupt,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
	},
}

type InterruptFields interface {
	RequestID() uint64
	Options() map[string]any
}

type interruptFields struct {
	requestID uint64
	options   map[string]any
}

func (c *interruptFields) RequestID() uint64 {
	return c.requestID
}

func (c *interruptFields) Options() map[string]any {
	return c.options
}

type Interrupt struct {
	InterruptFields
}

func NewInterruptWithFields(fields InterruptFields) *Interrupt {
	return &Interrupt{InterruptFields: fields}
}

func NewInterrupt(requestID uint64, options map[string]any) *Interrupt {
	return &Interrupt{InterruptFields: &interruptFields{requestID: requestID, options: options}}
}

func (c *Interrupt) Type() uint64 {
	return MessageTypeInterrupt
}

func (c *Interrupt) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, interruptValidationSpec)
	if err != nil {
		return fmt.Errorf("interrupt: failed to validate message %s: %w", MessageNameInterrupt, err)
	}

	c.InterruptFields = &interruptFields{requestID: fields.RequestID, options: fields.Options}

	return nil
}

func (c *Interrupt) Marshal() []any {
	return []any{MessageTypeInterrupt, c.RequestID(), c.Options()}
}
