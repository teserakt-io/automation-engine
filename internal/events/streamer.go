package events

//go:generate mockgen -destination=streamer_mocks.go -package events -self_package gitlab.com/teserakt/c2ae/internal/events gitlab.com/teserakt/c2ae/internal/events Streamer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-kit/kit/log"

	c2pb "gitlab.com/teserakt/c2/pkg/pb"
	"gitlab.com/teserakt/c2ae/internal/services"
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
	logger   log.Logger

	listeners []StreamListener
	stream    c2pb.C2_SubscribeToEventStreamClient
	lock      sync.RWMutex
}

var _ Streamer = (*streamer)(nil)

// NewStreamer creates a new streamer factory
func NewStreamer(c2Client services.C2, logger log.Logger) Streamer {
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
	s.logger.Log("msg", "added listener to event streamer")
}

func (s *streamer) RemoveListener(listener StreamListener) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, lis := range s.listeners {
		if lis == listener {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			s.logger.Log("msg", "removed listener to event streamer")

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

	s.logger.Log("msg", "started event streamer")

	for {
		select {
		case <-ctx.Done():
			s.logger.Log("msg", "stopped event stream", "reason", ctx.Err())
			return ctx.Err()
		default:
		}

		evt, err := stream.Recv()
		if err != nil {
			return err
		}

		for _, lis := range s.listeners {
			go lis.onEvent(*evt)
		}
	}
}
