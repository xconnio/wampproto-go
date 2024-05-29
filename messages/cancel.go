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

type CancelFields interface {
	RequestID() int64
	Options() map[string]any
}

type cancelFields struct {
	requestID int64
	options   map[string]any
}

func (c *cancelFields) RequestID() int64 {
	return c.requestID
}

func (c *cancelFields) Options() map[string]any {
	return c.options
}

type Cancel struct {
	CancelFields
}

func NewCancel(requestID int64, options map[string]any) *Cancel {
	return &Cancel{CancelFields: &cancelFields{requestID: requestID, options: options}}
}

func (c *Cancel) Type() int {
	return MessageTypeCancel
}

func (c *Cancel) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, cancelValidationSpec)
	if err != nil {
		return fmt.Errorf("cancel: failed to validate message %s: %w", MessageNameCancel, err)
	}

	c.CancelFields = &cancelFields{
		requestID: fields.RequestID,
		options:   fields.Options,
	}

	return nil
}

func (c *Cancel) Marshal() []any {
	return []any{MessageTypeCancel, c.RequestID(), c.Options()}
}
