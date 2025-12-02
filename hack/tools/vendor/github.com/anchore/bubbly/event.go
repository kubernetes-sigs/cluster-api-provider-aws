package bubbly

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wagoodman/go-partybus"
)

var _ EventHandler = (*EventDispatcher)(nil)

type EventHandlerFn func(partybus.Event) []tea.Model

type EventHandler interface {
	partybus.Responder
	Handle(partybus.Event) []tea.Model
}

type EventDispatcher struct {
	dispatch map[partybus.EventType]EventHandlerFn
	types    []partybus.EventType
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		dispatch: map[partybus.EventType]EventHandlerFn{},
	}
}

func (d *EventDispatcher) AddHandlers(handlers map[partybus.EventType]EventHandlerFn) {
	for t, h := range handlers {
		d.AddHandler(t, h)
	}
}

func (d *EventDispatcher) AddHandler(t partybus.EventType, fn EventHandlerFn) {
	d.dispatch[t] = fn
	d.types = append(d.types, t)
}

func (d EventDispatcher) RespondsTo() []partybus.EventType {
	return d.types
}

func (d EventDispatcher) Handle(e partybus.Event) []tea.Model {
	if fn, ok := d.dispatch[e.Type]; ok {
		return fn(e)
	}
	return nil
}
