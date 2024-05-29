package messages

import "fmt"

const MessageTypeEvent = 36
const MessageNameEvent = "EVENT"

var eventValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 6,
	Message:   MessageNameEvent,
	Spec: Spec{
		1: ValidateSubscriptionID,
		2: ValidatePublicationID,
		3: ValidateDetails,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type EventFields interface {
	SubscriptionID() int64
	PublicationID() int64
	Details() map[string]any
	Args() []any
	KwArgs() map[string]any
}

type eventFields struct {
	subscriptionID int64
	publicationID  int64
	details        map[string]any
	args           []any
	kwArgs         map[string]any
}

func (e *eventFields) SubscriptionID() int64 {
	return e.subscriptionID
}

func (e *eventFields) PublicationID() int64 {
	return e.publicationID
}

func (e *eventFields) Details() map[string]any {
	return e.details
}

func (e *eventFields) Args() []any {
	return e.args
}

func (e *eventFields) KwArgs() map[string]any {
	return e.kwArgs
}

type Event struct {
	EventFields
}

func NewEventWithFields(fields EventFields) *Event { return &Event{EventFields: fields} }

func NewEvent(subscriptionID, publicationID int64, details map[string]any, args []any, kwArgs map[string]any) *Event {
	if details == nil {
		details = make(map[string]any)
	}

	return &Event{EventFields: &eventFields{
		subscriptionID: subscriptionID,
		publicationID:  publicationID,
		details:        details,
		args:           args,
		kwArgs:         kwArgs,
	}}
}

func (e *Event) Type() int {
	return MessageTypeEvent
}

func (e *Event) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, eventValidationSpec)
	if err != nil {
		return fmt.Errorf("event: failed to validate message %s: %w", MessageNameEvent, err)
	}

	e.EventFields = &eventFields{
		subscriptionID: fields.SubscriptionID,
		publicationID:  fields.PublicationID,
		details:        fields.Details,
		args:           fields.Args,
		kwArgs:         fields.KwArgs,
	}

	return nil
}

func (e *Event) Marshal() []any {
	result := []any{MessageTypeEvent, e.SubscriptionID(), e.PublicationID(), e.Details()}

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
