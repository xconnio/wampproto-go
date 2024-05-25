package messages

type Message interface {
	Type() int
	Parse([]any) error
	Marshal() []any
}
