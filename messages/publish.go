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

type PublishFields interface {
	RequestID() int64
	Options() map[string]any
	Topic() string
	Args() []any
	KwArgs() map[string]any
}

type publishFields struct {
	requestID int64
	options   map[string]any
	topic     string
	args      []any
	kwArgs    map[string]any
}

func (e *publishFields) RequestID() int64 {
	return e.requestID
}

func (e *publishFields) Options() map[string]any {
	return e.options
}

func (e *publishFields) Topic() string {
	return e.topic
}

func (e *publishFields) Args() []any {
	return e.args
}

func (e *publishFields) KwArgs() map[string]any {
	return e.kwArgs
}

type Publish struct {
	PublishFields
}

func NewPublish(requestID int64, uri string, args []any, kwArgs map[string]any) *Publish {
	return &Publish{PublishFields: &publishFields{
		requestID: requestID,
		topic:     uri,
		args:      args,
		kwArgs:    kwArgs,
	}}
}

func NewPublishWithFields(fields PublishFields) *Publish { return &Publish{PublishFields: fields} }

func (e *Publish) Type() int {
	return MessageTypePublish
}

func (e *Publish) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, publishValidationSpec)
	if err != nil {
		return fmt.Errorf("publish: failed to validate message %s: %w", MessageNamePublish, err)
	}

	e.PublishFields = &publishFields{
		requestID: fields.RequestID,
		topic:     fields.Topic,
		args:      fields.Args,
		kwArgs:    fields.KwArgs,
	}

	return nil
}

func (e *Publish) Marshal() []any {
	result := []any{MessageTypePublish, e.RequestID(), e.Options(), e.Topic()}

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
