package messages

import "fmt"

const MessageTypeWelcome = 2
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
	SessionID() int64
	Details() map[string]any
}

type welcomeFields struct {
	sessionID int64
	details   map[string]any
}

func NewWelcomeFields(sessionID int64, details map[string]any) WelcomeFields {
	return &welcomeFields{
		sessionID: sessionID,
		details:   details,
	}
}

func (w *welcomeFields) SessionID() int64 {
	return w.sessionID
}

func (w *welcomeFields) Details() map[string]any {
	return w.details
}

type Welcome struct {
	WelcomeFields
}

func NewWelcome(sessionID int64, details map[string]any) *Welcome {
	return &Welcome{WelcomeFields: NewWelcomeFields(sessionID, details)}
}

func NewWelcomeWithFields(fields WelcomeFields) *Welcome { return &Welcome{WelcomeFields: fields} }

func (w *Welcome) Type() int {
	return MessageTypeWelcome
}

func (w *Welcome) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, welcomeValidationSpec)
	if err != nil {
		return fmt.Errorf("welcome: failed to validate message %s: %w", MessageNameWelcome, err)
	}

	w.WelcomeFields = NewWelcomeFields(fields.SessionID, fields.Details)

	return nil
}

func (w *Welcome) Marshal() []any {
	return []any{MessageTypeWelcome, w.SessionID(), w.Details()}
}
