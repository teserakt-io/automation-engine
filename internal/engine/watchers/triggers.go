package watchers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorhill/cronexpr"

	"gitlab.com/teserakt/c2ae/internal/events"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
)

//go:generate mockgen -destination=triggers_mocks.go -package watchers -self_package gitlab.com/teserakt/c2ae/internal/engine/watchers gitlab.com/teserakt/c2ae/internal/engine/watchers TriggerWatcherFactory,TriggerWatcher

// TriggerWatcher defines an interface for types watching on a trigger
type TriggerWatcher interface {
	Start(context.Context)
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
	logger        log.Logger
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error
}

var _ TriggerWatcherFactory = &triggerWatcherFactory{}
var _ TriggerWatcher = &schedulerWatcher{}
var _ TriggerWatcher = &clientSubscribedWatcher{}
var _TriggerWatcher = &clientUnsubscribedWatcher{}

// NewTriggerWatcherFactory creates a new watcher factory for given trigger
func NewTriggerWatcherFactory(logger log.Logger) TriggerWatcherFactory {
	return &triggerWatcherFactory{
		logger: logger,
	}
}

func (f *triggerWatcherFactory) Create(
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
			updateChan:    make(chan time.Time, 1),
			lastExecuted:  lastExecuted,
			logger:        f.logger,
		}

	case pb.TriggerType_CLIENT_SUBSCRIBED:
		watcher = &clientSubscribedWatcher{
			trigger:       trigger,
			triggeredChan: triggeredChan,
			errorChan:     errorChan,
			updateChan:    make(chan time.Time),
			lastExecuted:  lastExecuted,
			logger:        f.logger,
		}

	case pb.TriggerType_CLIENT_UNSUBSCRIBED:
		watcher = &clientUnsubscribedWatcher{
			trigger:       trigger,
			triggeredChan: triggeredChan,
			errorChan:     errorChan,
			updateChan:    make(chan time.Time),
			lastExecuted:  lastExecuted,
			logger:        f.logger,
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
	logger        log.Logger

	updateChan   chan time.Time
	lastExecuted time.Time
}

func (w *schedulerWatcher) Start(ctx context.Context) {
	w.logger.Log("msg", "started trigger schedulerWatcher", "trigger", w.trigger.ID, "rule", w.trigger.RuleID)

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
		case <-ctx.Done():
			w.logger.Log("msg", "stopping trigger schedulerWatcher", "trigger", w.trigger.ID, "rule", w.trigger.RuleID, "reason", ctx.Err())
			return
		}
	}
}

func (w *schedulerWatcher) UpdateLastExecuted(lastExecuted time.Time) error {
	w.updateChan <- lastExecuted

	return nil
}

type clientSubscribedWatcher struct {
	trigger       models.Trigger
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error
	updateChan    chan time.Time
	lastExecuted  time.Time
	logger        log.Logger
}

func (w *clientSubscribedWatcher) Start(ctx context.Context) {
	// TODO
}

func (w *clientSubscribedWatcher) UpdateLastExecuted(time.Time) error {
	// TODO
	return nil
}

type clientUnsubscribedWatcher struct {
	trigger       models.Trigger
	triggeredChan chan<- events.TriggerEvent
	errorChan     chan<- error
	updateChan    chan time.Time
	lastExecuted  time.Time
	logger        log.Logger
}

func (w *clientUnsubscribedWatcher) Start(ctx context.Context) {
	// TODO
}

func (w *clientUnsubscribedWatcher) UpdateLastExecuted(time.Time) error {
	// TODO
	return nil
}
