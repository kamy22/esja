package transport

import (
	"context"

	"github.com/ThreeDotsLabs/esja/stream"
)

// StructAnonymizer is an interface of the anonymizer component.
type StructAnonymizer interface {
	// Anonymize encrypts struct properties using secrets
	// correlated with a provided stream.ID.
	Anonymize(ctx context.Context, key stream.ID, data any) (any, error)

	// Deanonymize decrypts struct properties using secrets
	// correlated with a provided stream.ID.
	Deanonymize(ctx context.Context, key stream.ID, data any) (any, error)
}

// Anonymizer is a wrapper to any transport.Mapper instance.
// Anonymizer will anonymize transport model properties
// using provided StructAnonymizer implementation.
type Anonymizer[T any] struct {
	mapper     Mapper[T]
	anonymizer StructAnonymizer
}

// NewAnonymizer returns a new instance of Anonymizer.
func NewAnonymizer[T any](
	mapper Mapper[T],
	anonymizer StructAnonymizer,
) *Anonymizer[T] {
	return &Anonymizer[T]{
		mapper:     mapper,
		anonymizer: anonymizer,
	}
}

func (a *Anonymizer[T]) New(name stream.EventName) (any, error) {
	return a.mapper.New(name)
}

func (a *Anonymizer[T]) FromTransport(
	ctx context.Context,
	streamID stream.ID,
	payload any,
) (stream.Event[T], error) {
	payload, err := a.anonymizer.Deanonymize(ctx, streamID, payload)
	if err != nil {
		return nil, err
	}

	event, err := a.mapper.FromTransport(ctx, streamID, payload)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (a *Anonymizer[T]) ToTransport(
	ctx context.Context,
	streamID stream.ID,
	event stream.Event[T],
) (any, error) {
	e, err := a.mapper.ToTransport(ctx, streamID, event)
	if err != nil {
		return nil, err
	}

	payload, err := a.anonymizer.Anonymize(ctx, streamID, e)
	if err != nil {
		return nil, err
	}

	return payload, nil
}