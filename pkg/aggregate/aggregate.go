package aggregate

import (
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/esja/pkg/event"
)

type Aggregate[T EventSourced] struct {
	id          ID
	es          T
	eventsQueue []event.Event
}

func (a *Aggregate[T]) RecordEvents(events ...event.Event) error {
	for _, event := range events {
		err := a.ApplyEvents(event)
		if err != nil {
			return err
		}
		a.eventsQueue = append(a.eventsQueue, event)
	}

	return nil
}

func (a *Aggregate[T]) ApplyEvents(events ...event.Event) error {
	for _, event := range events {
		err := a.es.Handle(event)
		if err != nil {
			return fmt.Errorf("error applying event '%s': %w", event.EventName(), err)
		}
	}
	return nil
}

type EventSourced interface {
	Handle(event event.Event) error
}

func NewAggregate[T EventSourced](id ID, es T) (*Aggregate[T], error) {
	var a *Aggregate[T]
	if id == "" {
		return a, errors.New("id must not be empty")
	}

	return &Aggregate[T]{
		id:          id,
		es:          es,
		eventsQueue: []event.Event{},
	}, nil
}

func (a Aggregate[T]) ID() ID {
	return a.id
}

func (a Aggregate[T]) Base() T {
	return a.es
}

func (a *Aggregate[T]) PopEvents() []event.Event {
	var tmp = make([]event.Event, len(a.eventsQueue))
	copy(tmp, a.eventsQueue)
	a.eventsQueue = []event.Event{}

	return tmp
}
