package messages

import "fmt"

const MessageTypeInterrupt = 69
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

type Interrupt interface {
	Message

	RequestID() int64
	Options() map[string]any
}

type interrupt struct {
	requestID int64
	options   map[string]any
}

func NewEmptyInterrupt() Interrupt {
	return &interrupt{}
}

func NewInterrupt(requestID int64, options map[string]any) Interrupt {
	return &interrupt{
		requestID: requestID,
		options:   options,
	}
}

func (c *interrupt) RequestID() int64 {
	return c.requestID
}

func (c *interrupt) Options() map[string]any {
	return c.options
}

func (c *interrupt) Type() int {
	return MessageTypeInterrupt
}

func (c *interrupt) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, interruptValidationSpec)
	if err != nil {
		return fmt.Errorf("interrupt: failed to validate message %s: %w", MessageNameInterrupt, err)
	}

	c.requestID = fields.RequestID
	c.options = fields.Options

	return nil
}

func (c *interrupt) Marshal() []any {
	return []any{MessageTypeInterrupt, c.requestID, c.options}
}
