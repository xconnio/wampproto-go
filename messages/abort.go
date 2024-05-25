package messages

const MessageTypeAbort = 3
const MessageNameAbort = "ABORT"

type Abort interface {
	Message

	Details() map[string]any
	Reason() string
	Arguments() []any
	KwArguments() map[string]any
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
	panic("implement me")
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

func (a *abort) Arguments() []any {
	return a.args
}

func (a *abort) KwArguments() map[string]any {
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
