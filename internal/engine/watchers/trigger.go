package watchers

//go:generate mockgen -destination=trigger_mocks.go -package watchers -self_package gitlab.com/teserakt/c2ae/internal/engine/watchers gitlab.com/teserakt/c2ae/internal/engine/watchers TriggerWatcher

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/teserakt/c2ae/internal/services"

	"github.com/go-kit/kit/log"
	"github.com/gorhill/cronexpr"

	c2pb "gitlab.com/teserakt/c2/pkg/pb"

	"gitlab.com/teserakt/c2ae/internal/events"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
)

// TriggerEvent holds values transmitted when a trigger trigger
type TriggerEvent struct {
	Trigger models.Trigger
	Time    time.Time
}

// TriggerWatcher defines an interface for types watching on a trigger
type TriggerWatcher interface {
	Start(context.Context)
	UpdateLastExecuted(time.Time) error
}

// InvalidTrigger describe an error returned when the trigger is invalid
type InvalidTrigger struct {
	Err error
}

func (e InvalidTrigger) Error() string {
	return e.Err.Error()
}

type schedulerWatcher struct {
	validator models.TriggerValidator

	trigger       models.Trigger
	triggeredChan chan<- TriggerEvent
	errorChan     chan<- error
	logger        log.Logger

	updateChan   chan time.Time
	lastExecuted time.Time
}

func (w *schedulerWatcher) Start(ctx context.Context) {
	w.logger.Log("msg", "started trigger schedulerWatcher", "trigger", w.trigger.ID, "rule", w.trigger.RuleID)
	// Validate trigger
	if err := w.validator.ValidateTrigger(w.trigger); err != nil {
		w.errorChan <- InvalidTrigger{fmt.Errorf("failed to validate trigger: %v", err)}
		return
	}

	// Decode trigger settings
	settings := &pb.TriggerSettingsTimeInterval{}
	if err := settings.Decode(w.trigger.Settings); err != nil {
		w.errorChan <- InvalidTrigger{fmt.Errorf("failed to decode trigger settings: %v", err)}
		return
	}

	expr, err := cronexpr.Parse(settings.Expr)
	if err != nil {
		w.errorChan <- InvalidTrigger{fmt.Errorf("failed to parse cron expression: %v", err)}
		return
	}

	for {
		var delay time.Duration

		nextTime := expr.Next(w.lastExecuted)

		if now := time.Now(); nextTime.After(now) {
			delay = nextTime.Sub(now)
		}

		trigger := time.After(delay)
		select {
		case <-ctx.Done():
			w.logger.Log("msg", "stopping trigger schedulerWatcher", "trigger", w.trigger.ID, "rule", w.trigger.RuleID, "reason", ctx.Err())
			return

		case <-trigger:
			now := time.Now()

			w.triggeredChan <- TriggerEvent{
				Trigger: w.trigger,
				Time:    now,
			}
			w.lastExecuted = now

		case w.lastExecuted = <-w.updateChan:
		}
	}
}

func (w *schedulerWatcher) UpdateLastExecuted(lastExecuted time.Time) error {
	w.updateChan <- lastExecuted

	return nil
}

type eventWatcher struct {
	streamListenerFactory events.StreamListenerFactory
	triggerStateService   services.TriggerStateService
	validator             models.TriggerValidator

	trigger       models.Trigger
	targets       []models.Target
	triggeredChan chan<- TriggerEvent
	errorChan     chan<- error
	logger        log.Logger

	updateChan   chan time.Time
	lastExecuted time.Time
}

func (w *eventWatcher) Start(ctx context.Context) {
	// Validate trigger
	if err := w.validator.ValidateTrigger(w.trigger); err != nil {
		w.errorChan <- InvalidTrigger{fmt.Errorf("failed to validate trigger: %v", err)}
		return
	}

	// Decode settings
	settings := &pb.TriggerSettingsEvent{}
	if err := settings.Decode(w.trigger.Settings); err != nil {
		w.errorChan <- InvalidTrigger{fmt.Errorf("failed to decode trigger settings: %v", err)}
		return
	}

	lis := w.streamListenerFactory.Create(events.DefaultListenerBufSize, settings.EventType)
	defer lis.Close()

	state, err := w.triggerStateService.ByTriggerID(ctx, w.trigger.ID)
	if err != nil {
		w.errorChan <- fmt.Errorf("failed to fetch trigger state: %v", err)

		return
	}

	w.logger.Log("msg", "started trigger eventWatcher", "trigger", w.trigger.ID, "rule", w.trigger.RuleID, "event", settings.EventType)

	for {
		select {
		case <-ctx.Done():
			w.logger.Log("msg", "stopping trigger eventWatcher", "trigger", w.trigger.ID, "rule", w.trigger.RuleID, "reason", ctx.Err())
			return

		case evt := <-lis.C():
			origCounter := state.Counter

			if w.matchTargets(evt) {
				// Increment trigger counter in state
				state.Counter++
			}

			if state.Counter >= settings.MaxOccurence {
				//Trigger the rule action and reset the counter
				now := time.Now()
				w.triggeredChan <- TriggerEvent{
					Trigger: w.trigger,
					Time:    now,
				}
				w.lastExecuted = now
				state.Counter = 0
			}

			// Save state when counter has been modified
			if state.Counter != origCounter {
				w.logger.Log("msg", "saving trigger state", "state", state, "trigger", w.trigger.ID)
				if err := w.triggerStateService.Save(ctx, &state); err != nil {
					w.errorChan <- fmt.Errorf("failed to save trigger state: %v", err)
				}
			}

		case w.lastExecuted = <-w.updateChan:
		}
	}
}

func (w *eventWatcher) UpdateLastExecuted(lastExecuted time.Time) error {
	w.updateChan <- lastExecuted

	return nil
}

func (w *eventWatcher) matchTargets(evt c2pb.Event) bool {
	for _, target := range w.targets {
		switch target.Type {
		case pb.TargetType_CLIENT:
			if target.Expr == evt.Source {
				return true
			}
		case pb.TargetType_TOPIC:
			if target.Expr == evt.Target {
				return true
			}
		case pb.TargetType_ANY:
			if target.Expr == evt.Source || target.Expr == evt.Target {
				return true
			}
		default:
			w.logger.Log("msg", "unknow target type", "type", target.Type)
		}
	}

	return false
}
