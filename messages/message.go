package messages

type Message interface {
	Parse([]any) error
	Marshal() []any
}
