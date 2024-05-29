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

type PublishedFields interface {
	RequestID() int64
	PublicationID() int64
}

type publishedFields struct {
	requestID     int64
	publicationID int64
}

func NewPublishedFields(requestID, publicationID int64) PublishedFields {
	return &publishedFields{
		requestID:     requestID,
		publicationID: publicationID,
	}
}

func (p *publishedFields) RequestID() int64 {
	return p.requestID
}

func (p *publishedFields) PublicationID() int64 {
	return p.publicationID
}

type Published struct {
	PublishedFields
}

func NewPublished(fields PublishedFields) *Published {
	return &Published{
		PublishedFields: fields,
	}
}

func (p *Published) Type() int {
	return MessageTypePublished
}

func (p *Published) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, publishedValidationSpec)
	if err != nil {
		return fmt.Errorf("published: failed to validate message %s: %w", MessageNamePublished, err)
	}

	p.PublishedFields = NewPublishedFields(fields.RequestID, fields.PublicationID)

	return nil
}

func (p *Published) Marshal() []any {
	return []any{MessageTypePublished, p.RequestID(), p.PublicationID()}
}
