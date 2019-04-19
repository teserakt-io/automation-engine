package engine

//go:generate mockgen -destination=../mocks/engine_scheduler.go -package=mocks gitlab.com/teserakt/c2se/internal/engine Scheduler

import (
	"errors"
	"log"
	"sync"
	"time"

	"gitlab.com/teserakt/c2se/internal/events"
)

// Scheduler describe available opterations on an engine scheduler
type Scheduler interface {
	Start()
	Stop() error
}

type scheduler struct {
	sync.Mutex

	tickInterval time.Duration
	dispatcher   events.Dispatcher

	ticker *time.Ticker

	started bool
}

var _ Scheduler = &scheduler{}

var (
	// ErrNotStarted is returned when trying to stop a scheduler that is not started
	ErrNotStarted = errors.New("scheduler is not started")
)

// NewScheduler creates a new scheduler which will tick at given interval
func NewScheduler(tickInterval time.Duration, dispatcher events.Dispatcher) Scheduler {
	return &scheduler{
		tickInterval: tickInterval,
		dispatcher:   dispatcher,
	}
}

// Start will make the scheduler call its Tick method for every configured time interval
func (s *scheduler) Start() {
	s.Lock()

	s.ticker = time.NewTicker(s.tickInterval)
	s.started = true
	log.Printf("Scheduler started at %s\n", time.Now())

	s.Unlock()

	for t := range s.ticker.C {
		s.dispatcher.Dispatch(events.SchedulerTickType, s, events.SchedulerEventValue{
			Time: t,
		})
	}
}

func (s *scheduler) Stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.started {
		return ErrNotStarted
	}

	s.ticker.Stop()
	s.started = false

	log.Printf("Scheduler stopped at %s\n", time.Now())

	return nil
}
