package messages

import "fmt"

const MessageTypePublished uint64 = 17
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
	RequestID() uint64
	PublicationID() uint64
}

type publishedFields struct {
	requestID     uint64
	publicationID uint64
}

func (p *publishedFields) RequestID() uint64 {
	return p.requestID
}

func (p *publishedFields) PublicationID() uint64 {
	return p.publicationID
}

type Published struct {
	PublishedFields
}

func NewPublished(requestID, publicationID uint64) *Published {
	return &Published{PublishedFields: &publishedFields{requestID: requestID, publicationID: publicationID}}
}

func NewPublishedWithFields(fields PublishedFields) *Published {
	return &Published{PublishedFields: fields}
}

func (p *Published) Type() uint64 {
	return MessageTypePublished
}

func (p *Published) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, publishedValidationSpec)
	if err != nil {
		return fmt.Errorf("published: failed to validate message %s: %w", MessageNamePublished, err)
	}

	p.PublishedFields = &publishedFields{
		requestID:     fields.RequestID,
		publicationID: fields.PublicationID,
	}

	return nil
}

func (p *Published) Marshal() []any {
	return []any{MessageTypePublished, p.RequestID(), p.PublicationID()}
}
