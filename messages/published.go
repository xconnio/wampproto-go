package messages

import "fmt"

const MessageTypePublished = 17
const MessageNamePublished = "PUBLISHED"

var publishedValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNamePublished,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidatePublicationID,
	},
}

type Published interface {
	Message

	RequestID() int64
	PublicationID() int64
}

type published struct {
	requestID     int64
	publicationID int64
}

func NewEmptyPublished() Published {
	return &published{}
}

func NewPublished(requestID, publicationID int64) Published {
	return &published{
		requestID:     requestID,
		publicationID: publicationID,
	}
}

func (r *published) RequestID() int64 {
	return r.requestID
}

func (r *published) PublicationID() int64 {
	return r.publicationID
}

func (r *published) Type() int {
	return MessageTypePublished
}

func (r *published) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, publishedValidationSpec)
	if err != nil {
		return fmt.Errorf("published: failed to validate message %s: %w", MessageNamePublished, err)
	}

	r.requestID = fields.RequestID
	r.publicationID = fields.PublicationID

	return nil
}

func (r *published) Marshal() []any {
	return []any{MessageTypePublished, r.requestID, r.publicationID}
}
