package watchers

import (
	"fmt"
	"log"
	"time"

	"github.com/gorhill/cronexpr"
	"gitlab.com/teserakt/c2se/internal/events"
	"gitlab.com/teserakt/c2se/internal/models"

	"gitlab.com/teserakt/c2se/internal/pb"
)

//go:generate mockgen -destination=triggers_mocks.go -package watchers -self_package gitlab.com/teserakt/c2se/internal/engine/watchers gitlab.com/teserakt/c2se/internal/engine/watchers TriggerWatcherFactory,TriggerWatcher

// TriggerWatcher defines an interface for types watching on a trigger
type TriggerWatcher interface {
	Start()
	Stop() error
	UpdateLastExecuted(time.Time) error
}

// TriggerWatcherFactory allows to create trigger watchers from a given trigger
// independently of the trigger type
type TriggerWatcherFactory interface {
	Create(
		trigger models.Trigger,
		lastExecuted time.Time,
		triggeredChan chan<- events.TriggerEvent,
		errorChan chan<- error,
	) (TriggerWatcher, error)
}

type triggerWatcherFactory struct {
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error
	stopChan      chan bool
}

var _ TriggerWatcherFactory = &triggerWatcherFactory{}
var _ TriggerWatcher = &schedulerWatcher{}
var _ TriggerWatcher = &clientSubscribedWatcher{}
var _TriggerWatcher = &clientUnsubscribedWatcher{}

// NewTriggerWatcherFactory creates a new watcher factory for given trigger
func NewTriggerWatcherFactory() TriggerWatcherFactory {
	return &triggerWatcherFactory{}
}

func (l *triggerWatcherFactory) Create(
	trigger models.Trigger,
	lastExecuted time.Time,
	triggeredChan chan<- events.TriggerEvent,
	errorChan chan<- error,
) (TriggerWatcher, error) {

	var watcher TriggerWatcher

	switch trigger.TriggerType {
	case pb.TriggerType_TIME_INTERVAL:
		watcher = &schedulerWatcher{
			trigger:       trigger,
			triggeredChan: triggeredChan,
			errorChan:     errorChan,
			stopChan:      make(chan bool, 1),
			updateChan:    make(chan time.Time),
			lastExecuted:  lastExecuted,
		}

	case pb.TriggerType_CLIENT_SUBSCRIBED:
		watcher = &clientSubscribedWatcher{
			trigger:       trigger,
			triggeredChan: triggeredChan,
			errorChan:     errorChan,
			stopChan:      make(chan bool, 1),
			updateChan:    make(chan time.Time),
			lastExecuted:  lastExecuted,
		}

	case pb.TriggerType_CLIENT_UNSUBSCRIBED:
		watcher = &clientUnsubscribedWatcher{
			trigger:       trigger,
			triggeredChan: triggeredChan,
			errorChan:     errorChan,
			stopChan:      make(chan bool, 1),
			updateChan:    make(chan time.Time),
			lastExecuted:  lastExecuted,
		}

	default:
		return nil, fmt.Errorf("TriggerWatcherFactory don't know how to handle trigger type %s", trigger.TriggerType)
	}

	return watcher, nil
}

// InvalidTriggerSettings describe an error returned when the trigger settings are invalid
type InvalidTriggerSettings struct {
	Err error
}

func (e InvalidTriggerSettings) Error() string {
	return e.Err.Error()
}

// InvalidCronExpr describe an error returned when the trigger have an invalid cron expression set
type InvalidCronExpr struct {
	Err error
}

func (e InvalidCronExpr) Error() string {
	return e.Err.Error()
}

type schedulerWatcher struct {
	trigger       models.Trigger
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error

	stopChan     chan bool
	updateChan   chan time.Time
	lastExecuted time.Time
}

func (w *schedulerWatcher) Start() {
	log.Printf("Started trigger watcher for trigger %d (Rule %d)", w.trigger.ID, w.trigger.RuleID)

	settings := &pb.TriggerSettingsTimeInterval{}
	if err := settings.Decode(w.trigger.Settings); err != nil {
		w.errorChan <- InvalidTriggerSettings{fmt.Errorf("failed to decode trigger settings: %s", err)}

		return
	}

	for {
		var delay time.Duration
		expr, err := cronexpr.Parse(settings.Expr)
		if err != nil {
			w.errorChan <- InvalidCronExpr{fmt.Errorf("failed to parse cron expression: %s", err)}

			return
		}
		nextTime := expr.Next(w.lastExecuted)

		if now := time.Now(); nextTime.After(now) {
			delay = nextTime.Sub(now)
		}

		trigger := time.After(delay)
		select {
		case <-trigger:
			now := time.Now()

			w.triggeredChan <- events.TriggerEvent{
				Trigger: w.trigger,
				Time:    now,
			}
			w.lastExecuted = now
		case w.lastExecuted = <-w.updateChan:
		case <-w.stopChan:
			return
		}
	}
}

func (w *schedulerWatcher) Stop() error {
	select {
	case w.stopChan <- true:
		log.Printf("Stopped schedulerWatcher for trigger %d (rule %d)", w.trigger.ID, w.trigger.RuleID)
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("Couldn't stop schedulerWatcher for trigger %d (rule %d), maybe it's already stopped ?", w.trigger.ID, w.trigger.RuleID)
	}

	return nil
}

func (w *schedulerWatcher) UpdateLastExecuted(lastExecuted time.Time) error {
	select {
	case w.updateChan <- lastExecuted:
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("Couldn't update lastExecuted on schedulerWatcher for trigger %d (rule %d), maybe it's already stopped ?", w.trigger.ID, w.trigger.RuleID)
	}

	return nil
}

type clientSubscribedWatcher struct {
	trigger       models.Trigger
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error
	stopChan      chan bool
	updateChan    chan time.Time
	lastExecuted  time.Time
}

func (w *clientSubscribedWatcher) Start() {
	// TODO
}

func (w *clientSubscribedWatcher) Stop() error {
	// TODO
	return nil
}

func (w *clientSubscribedWatcher) UpdateLastExecuted(time.Time) error {
	// TODO
	return nil
}

type clientUnsubscribedWatcher struct {
	trigger       models.Trigger
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error
	stopChan      chan bool
	updateChan    chan time.Time
	lastExecuted  time.Time
}

func (w *clientUnsubscribedWatcher) Start() {
	// TODO
}

func (w *clientUnsubscribedWatcher) Stop() error {
	// TODO
	return nil
}

func (w *clientUnsubscribedWatcher) UpdateLastExecuted(time.Time) error {
	// TODO
	return nil
}
