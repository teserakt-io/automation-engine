package events

//go:generate mockgen -destination=streamer_mocks.go -package events -self_package github.com/teserakt-io/automation-engine/internal/events github.com/teserakt-io/automation-engine/internal/events Streamer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/teserakt-io/automation-engine/internal/services"
)

// events errors
var (
	ErrListenerNotFound = errors.New("listener not found")
)

// Streamer defines an interface to stream C2 events
type Streamer interface {
	StartStream(context.Context) error
	AddListener(listener StreamListener)
	RemoveListener(listener StreamListener) error
	Listeners() []StreamListener
}

type streamer struct {
	c2Client services.C2
	logger   log.FieldLogger

	listeners []StreamListener
	lock      sync.RWMutex
}

var _ Streamer = (*streamer)(nil)

// NewStreamer creates a new streamer factory
func NewStreamer(c2Client services.C2, logger log.FieldLogger) Streamer {
	return &streamer{
		c2Client:  c2Client,
		logger:    logger,
		listeners: []StreamListener{},
	}
}

func (s *streamer) AddListener(listener StreamListener) {
	s.lock.Lock()
	s.listeners = append(s.listeners, listener)
	s.lock.Unlock()
	s.logger.Info("added listener to event streamer")
}

func (s *streamer) RemoveListener(listener StreamListener) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, lis := range s.listeners {
		if lis == listener {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			s.logger.Info("removed listener to event streamer")

			return nil
		}
	}

	return ErrListenerNotFound
}

func (s *streamer) Listeners() []StreamListener {
	return s.listeners
}

// StartStream will open a stream from the C2 clients, and
// fan out every events it receive to all registered listeners.
func (s *streamer) StartStream(ctx context.Context) error {
	stream, err := s.c2Client.SubscribeToEventStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to start event stream: %v", err)
	}

	s.logger.Info("started event streamer")

	for {
		select {
		case <-ctx.Done():
			s.logger.WithError(ctx.Err()).Warn("stopped event stream")
			return ctx.Err()
		default:
		}

		evt, err := stream.Recv()
		if err != nil {
			return err
		}

		s.lock.Lock()
		for _, lis := range s.listeners {
			go lis.onEvent(*evt)
		}
		s.lock.Unlock()
	}
}
