package messages

import "fmt"

const MessageTypeSubscribe = 32
const MessageNameSubscribe = "SUBSCRIBE"

var subscribeValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 4,
	Message:   MessageNameSubscribe,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateOptions,
		3: ValidateURI,
	},
}

type Subscribe interface {
	Message

	RequestID() int64
	Options() map[string]any
	Topic() string
}

type subscribe struct {
	requestID int64
	options   map[string]any
	topic     string
}

func NewEmptySubscribe() Subscribe {
	return &subscribe{}
}

func NewSubscribe(requestID int64, options map[string]any, uri string) Subscribe {
	return &subscribe{
		requestID: requestID,
		options:   options,
		topic:     uri,
	}
}

func (r *subscribe) Type() int {
	return MessageTypeSubscribe
}

func (r *subscribe) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, subscribeValidationSpec)
	if err != nil {
		return fmt.Errorf("subscribe: failed to validate message %s: %w", MessageNameSubscribe, err)
	}

	r.requestID = fields.SessionID
	r.options = fields.Options
	r.topic = fields.URI

	return nil
}

func (r *subscribe) Marshal() []any {
	return []any{MessageTypeSubscribe, r.requestID, r.options, r.topic}
}

func (r *subscribe) RequestID() int64 {
	return r.requestID
}

func (r *subscribe) Options() map[string]any {
	return r.options
}

func (r *subscribe) Topic() string {
	return r.topic
}
