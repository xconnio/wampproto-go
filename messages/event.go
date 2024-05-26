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

type Event interface {
	Message

	SubscriptionID() int64
	PublicationID() int64
	Details() map[string]any
	Args() []any
	KwArgs() map[string]any
}

type event struct {
	subscriptionID int64
	publicationID  int64
	details        map[string]any
	args           []any
	kwArgs         map[string]any
}

func NewEmptyEvent() Event {
	return &event{}
}

func NewEvent(subscriptionID, publicationID int64, details map[string]any, args []any, kwArgs map[string]any) Event {
	return &event{
		subscriptionID: subscriptionID,
		publicationID:  publicationID,
		details:        details,
		args:           args,
		kwArgs:         kwArgs,
	}
}

func (e *event) SubscriptionID() int64 {
	return e.subscriptionID
}

func (e *event) PublicationID() int64 {
	return e.publicationID
}

func (e *event) Details() map[string]any {
	return e.details
}

func (e *event) Args() []any {
	return e.args
}

func (e *event) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *event) Type() int {
	return MessageTypeEvent
}

func (e *event) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, eventValidationSpec)
	if err != nil {
		return fmt.Errorf("event: failed to validate message %s: %w", MessageNameEvent, err)
	}

	e.subscriptionID = fields.SubscriptionID
	e.publicationID = fields.PublicationID
	e.details = fields.Details
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *event) Marshal() []any {
	result := []any{MessageTypeEvent, e.subscriptionID, e.publicationID, e.details}

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
