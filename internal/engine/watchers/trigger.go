// Copyright 2020 Teserakt AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package watchers

//go:generate mockgen -copyright_file ../../../doc/COPYRIGHT_TEMPLATE.txt -destination=trigger_mocks.go -package watchers -self_package github.com/teserakt-io/automation-engine/internal/engine/watchers github.com/teserakt-io/automation-engine/internal/engine/watchers TriggerWatcher

import (
	"context"
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
	c2pb "github.com/teserakt-io/c2/pkg/pb"

	"github.com/teserakt-io/automation-engine/internal/events"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
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
	logger        log.FieldLogger

	updateChan   chan time.Time
	lastExecuted time.Time
}

func (w *schedulerWatcher) Start(ctx context.Context) {

	logger := w.logger.WithFields(log.Fields{
		"trigger": w.trigger.ID,
		"rule":    w.trigger.RuleID,
	})

	logger.Info("started trigger schedulerWatcher")
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
			logger.WithError(ctx.Err()).Warn("stopping trigger schedulerWatcher")
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
	logger        log.FieldLogger

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

	logger := w.logger.WithFields(log.Fields{
		"trigger": w.trigger.ID,
		"rule":    w.trigger.RuleID,
		"event":   settings.EventType,
	})

	logger.Info("started trigger eventWatcher")

	for {
		select {
		case <-ctx.Done():
			logger.WithError(ctx.Err()).Warn("stopping trigger eventWatcher")
			return

		case evt := <-lis.C():
			origCounter := state.Counter

			if w.matchTargets(evt) {
				// Increment trigger counter in state
				state.Counter++
			}

			if state.Counter >= settings.MaxOccurrence {
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
				logger.WithField("state", state).Info("saving trigger state")
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
			w.logger.WithField("type", target.Type).Warn("unknown target type")
		}
	}

	return false
}
