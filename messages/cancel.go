package messages

import "fmt"

const MessageTypeCancel = 49
const MessageNameCancel = "CANCEL"

var cancelValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameCancel,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
	},
}

type Cancel interface {
	Message

	RequestID() int64
	Options() map[string]any
}

type cancel struct {
	requestID int64
	options   map[string]any
}

func NewEmptyCancel() Cancel {
	return &cancel{}
}

func NewCancel(requestID int64, options map[string]any) Cancel {
	return &cancel{
		requestID: requestID,
		options:   options,
	}
}

func (c *cancel) RequestID() int64 {
	return c.requestID
}

func (c *cancel) Options() map[string]any {
	return c.options
}

func (c *cancel) Type() int {
	return MessageTypeCancel
}

func (c *cancel) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, cancelValidationSpec)
	if err != nil {
		return fmt.Errorf("cancel: failed to validate message %s: %w", MessageNameCancel, err)
	}

	c.requestID = fields.RequestID
	c.options = fields.Options

	return nil
}

func (c *cancel) Marshal() []any {
	return []any{MessageTypeCancel, c.requestID, c.options}
}
