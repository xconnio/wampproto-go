package messages

import "fmt"

const MessageTypePublish = 16
const MessageNamePublish = "PUBLISH"

var publishValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 6,
	Message:   MessageNamePublish,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
		3: ValidateURI,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type Publish interface {
	Message

	RequestID() int64
	Options() map[string]any
	Topic() string
	Args() []any
	KwArgs() map[string]any
}

type publish struct {
	requestID int64
	options   map[string]any
	topic     string
	args      []any
	kwArgs    map[string]any
}

func NewEmptyPublish() Publish {
	return &publish{}
}

func NewPublish(requestID int64, uri string, args []any, kwArgs map[string]any) Publish {
	return &publish{
		requestID: requestID,
		topic:     uri,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *publish) RequestID() int64 {
	return e.requestID
}

func (e *publish) Options() map[string]any {
	return e.options
}

func (e *publish) Topic() string {
	return e.topic
}

func (e *publish) Args() []any {
	return e.args
}

func (e *publish) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *publish) Type() int {
	return MessageTypePublish
}

func (e *publish) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, publishValidationSpec)
	if err != nil {
		return fmt.Errorf("publish: failed to validate message %s: %w", MessageNamePublish, err)
	}

	e.requestID = fields.RequestID
	e.options = fields.Options
	e.topic = fields.URI
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *publish) Marshal() []any {
	result := []any{MessageTypePublish, e.requestID, e.options, e.topic}

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
