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

type Welcome interface {
	Message

	SessionID() int64
	Details() map[string]any
}

type welcome struct {
	sessionID int64
	details   map[string]any
}

func NewEmptyWelcome() Welcome {
	return &welcome{}
}

func NewWelcome(sessionID int64, details map[string]any) Welcome {
	return &welcome{
		sessionID: sessionID,
		details:   details,
	}
}

func (w *welcome) Type() int {
	return MessageTypeWelcome
}

func (w *welcome) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, welcomeValidationSpec)
	if err != nil {
		return fmt.Errorf("welcome: failed to validate message %s: %w", MessageNameWelcome, err)
	}

	w.sessionID = fields.SessionID
	w.details = fields.Details

	return nil
}

func (w *welcome) Marshal() []any {
	return []any{MessageTypeWelcome, w.sessionID, w.details}
}

func (w *welcome) SessionID() int64 {
	return w.sessionID
}

func (w *welcome) Details() map[string]any {
	return w.details
}
