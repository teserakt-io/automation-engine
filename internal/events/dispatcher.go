package events

//go:generate mockgen -destination=../mocks/events_dispatcher.go -package=mocks gitlab.com/teserakt/c2se/internal/events Dispatcher

import (
	"log"
	"sync"
)

// Type defines a custom type for defining events
type Type int

// List of availabe event types
const (
	UnknowEventType Type = iota
	SchedulerTickType
	RulesModifiedType
	ClientSubscribedType
	ClientUnsubscribedType
)

// EventStrings list all events in a human readable form
var EventStrings = map[Type]string{
	UnknowEventType:        "unknow",
	SchedulerTickType:      "schedulerTick",
	RulesModifiedType:      "ruleModified",
	ClientSubscribedType:   "clientSubscribed",
	ClientUnsubscribedType: "clientUnsubscribed",
}

// Event is the data transmitted to each listeners once dispatched
type Event interface {
	Type() Type
	Source() interface{}
	Value() interface{}
}

// event is the internal implementaiton of the Event interface
type event struct {
	eventType Type
	source    interface{}
	value     interface{}
}

func (e *event) Type() Type {
	return e.eventType
}

func (e *event) Source() interface{} {
	return e.source
}

func (e *event) Value() interface{} {
	return e.value
}

// Listener interface is used to suscribe to event types and receive them
type Listener interface {
	OnEvent(Event) error
}

// Dispatcher is an interface used to defines system able to send events to its registered listeners
type Dispatcher interface {
	Start()
	Register(Type, Listener)
	Dispatch(Type, interface{}, interface{})
	ClearListeners()
}

type dispatcher struct {
	sync.Mutex
	eventChannel chan Event

	listeners        map[Type][]Listener
	listenerChannels map[Listener]chan Event
}

var _ Dispatcher = &dispatcher{}
var _ Event = &event{}

// NewDispatcher creates a new event dispatcher
func NewDispatcher() Dispatcher {
	return &dispatcher{
		listeners:        make(map[Type][]Listener),
		eventChannel:     make(chan Event, 100), // small buffer on event channel to avoid blocking
		listenerChannels: make(map[Listener]chan Event),
	}
}

func (d *dispatcher) Start() {
	for evt := range d.eventChannel {
		for _, listener := range d.listeners[evt.Type()] {
			select {
			case d.listenerChannels[listener] <- evt:
			default:
				log.Printf("Discarding event %#v - listener busy", evt)
			}
		}
	}
}

// Register allow to add listeners to a given event type
func (d *dispatcher) Register(eventtType Type, listener Listener) {
	d.Lock()
	defer d.Unlock()

	d.listeners[eventtType] = append(d.listeners[eventtType], listener)
	d.listenerChannels[listener] = make(chan Event)
	go func() {
		log.Printf("Started listener %p event loop", listener)
		for evt := range d.listenerChannels[listener] {
			if err := listener.OnEvent(evt); err != nil {
				log.Printf("Error while processing event: %s", err)
			}
		}

		log.Printf("Stopped listener %p event loop", listener)
	}()
}

// Dispatch will notify each listener associated with this event
func (d *dispatcher) Dispatch(eventType Type, source interface{}, value interface{}) {
	evt := &event{
		eventType: eventType,
		source:    source,
		value:     value,
	}

	d.eventChannel <- evt
}

// ClearListeners removes all registered listeners from the dispatcher
func (d *dispatcher) ClearListeners() {
	d.Lock()
	defer d.Unlock()

	for _, channel := range d.listenerChannels {
		close(channel)
	}

	d.listeners = make(map[Type][]Listener)
	d.listenerChannels = make(map[Listener]chan Event)
}
