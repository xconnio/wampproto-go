package messages

import "fmt"

const MessageTypeAbort = 3
const MessageNameAbort = "ABORT"

var abortValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameAbort,
	Spec: Spec{
		1: ValidateDetails,
		2: ValidateReason,
	},
}

type Abort interface {
	Message

	Details() map[string]any
	Reason() string
	Args() []any
	KwArgs() map[string]any
}

type abort struct {
	details map[string]any
	reason  string
	args    []any
	kwArgs  map[string]any
}

func (a *abort) Type() int {
	return MessageTypeAbort
}

func (a *abort) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, abortValidationSpec)
	if err != nil {
		return fmt.Errorf("abort: failed to validate message %s: %w", MessageNameAbort, err)
	}

	a.details = fields.Details
	a.reason = fields.Reason
	a.args = fields.Args
	a.kwArgs = fields.KwArgs

	return nil
}

func (a *abort) Marshal() []any {
	return []any{MessageTypeAbort, a.details, a.reason}
}

func (a *abort) Details() map[string]any {
	return a.details
}

func (a *abort) Reason() string {
	return a.reason
}

func (a *abort) Args() []any {
	return a.args
}

func (a *abort) KwArgs() map[string]any {
	return a.kwArgs
}

func NewEmptyAbort() Abort {
	return &abort{}
}

func NewAbort(details map[string]any, reason string, args []any, KwArgs map[string]any) Abort {
	return &abort{
		details: details,
		reason:  reason,
		args:    args,
		kwArgs:  KwArgs,
	}
}
