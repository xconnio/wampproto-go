package messages

import "fmt"

const MessageTypeWelcome uint64 = 2
const MessageNameWelcome = "WELCOME"

var welcomeValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameWelcome,
	Spec: Spec{
		1: ValidateSessionID,
		2: ValidateDetails,
	},
}

type WelcomeFields interface {
	SessionID() uint64
	Details() map[string]any
}

type welcomeFields struct {
	sessionID uint64
	details   map[string]any
}

func (w *welcomeFields) SessionID() uint64 {
	return w.sessionID
}

func (w *welcomeFields) Details() map[string]any {
	return w.details
}

type Welcome struct {
	WelcomeFields
}

func NewWelcome(sessionID uint64, details map[string]any) *Welcome {
	return &Welcome{WelcomeFields: &welcomeFields{sessionID: sessionID, details: details}}
}

func NewWelcomeWithFields(fields WelcomeFields) *Welcome { return &Welcome{WelcomeFields: fields} }

func (w *Welcome) Type() uint64 {
	return MessageTypeWelcome
}

func (w *Welcome) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, welcomeValidationSpec)
	if err != nil {
		return fmt.Errorf("welcome: failed to validate message %s: %w", MessageNameWelcome, err)
	}

	w.WelcomeFields = &welcomeFields{sessionID: fields.SessionID, details: fields.Details}

	return nil
}

func (w *Welcome) Marshal() []any {
	return []any{MessageTypeWelcome, w.SessionID(), w.Details()}
}
