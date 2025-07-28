package messages

import "fmt"

const MessageTypeAbort uint64 = 3
const MessageNameAbort = "ABORT"

var abortValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 5,
	Message:   MessageNameAbort,
	Spec: Spec{
		1: ValidateDetails,
		2: ValidateReason,
		3: ValidateArgs,
		4: ValidateKwArgs,
	},
}

type AbortFields interface {
	Details() map[string]any
	Reason() string
	Args() []any
	KwArgs() map[string]any
}

type abortFields struct {
	details map[string]any
	reason  string
	args    []any
	kwArgs  map[string]any
}

func (a *abortFields) Details() map[string]any {
	return a.details
}

func (a *abortFields) Reason() string {
	return a.reason
}

func (a *abortFields) Args() []any {
	return a.args
}

func (a *abortFields) KwArgs() map[string]any {
	return a.kwArgs
}

type Abort struct {
	AbortFields
}

func NewAbort(details map[string]any, reason string, args []any, KwArgs map[string]any) *Abort {
	return &Abort{AbortFields: &abortFields{
		details: details,
		reason:  reason,
		args:    args,
		kwArgs:  KwArgs,
	}}
}

func NewAbortWithFields(fields AbortFields) *Abort { return &Abort{AbortFields: fields} }

func (a *Abort) Type() uint64 {
	return MessageTypeAbort
}

func (a *Abort) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, abortValidationSpec)
	if err != nil {
		return fmt.Errorf("abort: failed to validate message %s: %w", MessageNameAbort, err)
	}

	a.AbortFields = &abortFields{
		details: fields.Details,
		reason:  fields.Reason,
		args:    fields.Args,
		kwArgs:  fields.KwArgs,
	}

	return nil
}

func (a *Abort) Marshal() []any {
	payload := []any{MessageTypeAbort, a.Details(), a.Reason()}

	if a.Args() != nil {
		payload = append(payload, a.Args())
	}

	if a.KwArgs() != nil {
		if a.Args() == nil {
			payload = append(payload, a.Args())
		}

		payload = append(payload, a.KwArgs())
	}

	return payload
}
