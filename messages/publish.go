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
	RequestID() uint64
	Options() map[string]any
	Topic() string
	Args() []any
	KwArgs() map[string]any

	BinaryPayload
}

type publishFields struct {
	requestID uint64
	options   map[string]any
	topic     string
	args      []any
	kwArgs    map[string]any
}

func (e *publishFields) RequestID() uint64 {
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

func (e *publishFields) PayloadIsBinary() bool {
	return false
}

func (e *publishFields) Payload() []byte {
	return nil
}

func (e *publishFields) PayloadSerializer() uint64 {
	return 0
}

type Publish struct {
	PublishFields
}

func NewPublish(requestID uint64, option map[string]any, uri string, args []any, kwArgs map[string]any) *Publish {
	if option == nil {
		option = make(map[string]any)
	}

	return &Publish{PublishFields: &publishFields{
		requestID: requestID,
		options:   option,
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
		topic:     fields.URI,
		options:   fields.Options,
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
