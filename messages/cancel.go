package messages

import "fmt"

const MessageTypeCancel uint64 = 49
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
	RequestID() uint64
	Options() map[string]any
}

type cancelFields struct {
	requestID uint64
	options   map[string]any
}

func (c *cancelFields) RequestID() uint64 {
	return c.requestID
}

func (c *cancelFields) Options() map[string]any {
	return c.options
}

type Cancel struct {
	CancelFields
}

func NewCancel(requestID uint64, options map[string]any) *Cancel {
	if options == nil {
		options = make(map[string]any)
	}

	return &Cancel{CancelFields: &cancelFields{requestID: requestID, options: options}}
}

func NewCancelWithFields(fields CancelFields) *Cancel {
	return &Cancel{CancelFields: fields}
}

func (c *Cancel) Type() uint64 {
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
