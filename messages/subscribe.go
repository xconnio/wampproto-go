package messages

import "fmt"

const MessageTypeSubscribe uint64 = 32
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

type SubscribeFields interface {
	RequestID() uint64
	Options() map[string]any
	Topic() string
}

type subscribeFields struct {
	requestID uint64
	options   map[string]any
	topic     string
}

func (s *subscribeFields) Marshal() []any {
	return []any{MessageTypeSubscribe, s.requestID, s.options, s.topic}
}

func (s *subscribeFields) RequestID() uint64 {
	return s.requestID
}

func (s *subscribeFields) Options() map[string]any {
	return s.options
}

func (s *subscribeFields) Topic() string {
	return s.topic
}

type Subscribe struct {
	SubscribeFields
}

func NewSubscribe(requestID uint64, options map[string]any, uri string) *Subscribe {
	if options == nil {
		options = make(map[string]any)
	}

	return &Subscribe{SubscribeFields: &subscribeFields{requestID: requestID, options: options, topic: uri}}
}

func NewSubscribeWithFields(fields SubscribeFields) *Subscribe {
	return &Subscribe{SubscribeFields: fields}
}

func (s *Subscribe) Type() uint64 {
	return MessageTypeSubscribe
}

func (s *Subscribe) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, subscribeValidationSpec)
	if err != nil {
		return fmt.Errorf("subscribe: failed to validate message %s: %w", MessageNameSubscribe, err)
	}

	s.SubscribeFields = &subscribeFields{requestID: fields.RequestID, options: fields.Options, topic: fields.URI}

	return nil
}

func (s *Subscribe) Marshal() []any {
	return []any{MessageTypeSubscribe, s.RequestID(), s.Options(), s.Topic()}
}
