package events

//go:generate mockgen -destination=listener_mocks.go -package events -self_package github.com/teserakt-io/automation-engine/internal/events github.com/teserakt-io/automation-engine/internal/events StreamListenerFactory,StreamListener

import (
	c2pb "github.com/teserakt-io/c2/pkg/pb"

	pb "github.com/teserakt-io/automation-engine/internal/pb"
)

var (
	// DefaultListenerBufSize defines the default buffer size for the stream listeners internal channel
	DefaultListenerBufSize = 1000
)

// StreamListenerFactory defines a factory creating StreamListeners
type StreamListenerFactory interface {
	Create(eventChanBufSize int, eventTypeWhitelist ...pb.EventType) StreamListener
}

type streamListenerFactory struct {
	streamer Streamer
}

var _ StreamListenerFactory = (*streamListenerFactory)(nil)

// NewStreamListenerFactory creates a new StreamListener factory
func NewStreamListenerFactory(streamer Streamer) StreamListenerFactory {
	return &streamListenerFactory{
		streamer: streamer,
	}
}

func (f *streamListenerFactory) Create(eventChanBufSize int, eventTypeWhitelist ...pb.EventType) StreamListener {
	lis := &streamListener{
		eventChan:          make(chan c2pb.Event, eventChanBufSize),
		eventTypeWhitelist: eventTypeWhitelist,
		streamer:           f.streamer,
	}

	f.streamer.AddListener(lis)

	return lis
}

// StreamListener defines a type able to listen for stream events
type StreamListener interface {
	onEvent(c2pb.Event)
	C() <-chan c2pb.Event
	Close() error
}

type streamListener struct {
	eventChan          chan c2pb.Event
	eventTypeWhitelist []pb.EventType
	streamer           Streamer
}

var _ StreamListener = (*streamListener)(nil)

func (l *streamListener) onEvent(evt c2pb.Event) {
	var whitelistedType bool
	for _, t := range l.eventTypeWhitelist {
		if t == pb.EventType(evt.Type.String()) {
			whitelistedType = true
			continue
		}
	}

	if !whitelistedType {
		return
	}

	select {
	case l.eventChan <- evt:
	default:
		<-l.eventChan
		l.eventChan <- evt
	}
}

// C returns an event channel, containing only listener's whitelisted types event
// From the Streamer the listener has been registered on.
func (l *streamListener) C() <-chan c2pb.Event {
	return l.eventChan
}

func (l *streamListener) Close() error {
	return l.streamer.RemoveListener(l)
}
