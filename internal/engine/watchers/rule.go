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

//go:generate mockgen -copyright_file ../../../doc/COPYRIGHT_TEMPLATE.txt -destination=rule_mocks.go -package watchers -self_package github.com/teserakt-io/automation-engine/internal/engine/watchers github.com/teserakt-io/automation-engine/internal/engine/watchers RuleWatcher

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"

	"github.com/teserakt-io/automation-engine/internal/engine/actions"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// RuleWatcher defines methods to implement a rule watcher.
// It is responsible of monitoring the rule trigger(s), and execute the rule action
// when the trigger conditions are met.
type RuleWatcher interface {
	Start(context.Context)
}

type ruleWatcher struct {
	rule                  models.Rule
	triggerWatcherFactory TriggerWatcherFactory
	actionFactory         actions.ActionFactory
	ruleWriter            services.RuleWriter
	errorChan             chan<- error
	triggeredChan         chan TriggerEvent
	logger                log.FieldLogger
}

func (w *ruleWatcher) Start(ctx context.Context) {
	var triggerWatchers []TriggerWatcher

	if len(w.rule.Triggers) == 0 {
		w.logger.WithField("rule", w.rule.ID).Warn("rule has no triggers")
		return
	}

	for _, trigger := range w.rule.Triggers {
		triggerWatcher, err := w.triggerWatcherFactory.Create(
			trigger,
			w.rule.Targets,
			w.rule.LastExecuted,
			w.triggeredChan,
			w.errorChan,
		)

		if err != nil {
			w.errorChan <- err

			continue
		}

		triggerWatchers = append(triggerWatchers, triggerWatcher)

		go triggerWatcher.Start(ctx)
	}

	for {
		select {
		case triggerEvt := <-w.triggeredChan:
			ctx, span := trace.StartSpan(ctx, "RuleWatcher.RuleTriggered")
			span.Annotate([]trace.Attribute{
				trace.Int64Attribute("ruleID", int64(w.rule.ID)),
				trace.Int64Attribute("triggerID", int64(triggerEvt.Trigger.ID)),
			}, "Rule triggered")

			w.logger.WithFields(log.Fields{
				"rule":    w.rule.ID,
				"trigger": triggerEvt.Trigger.ID,
			}).Info("rule triggered")

			w.rule.LastExecuted = triggerEvt.Time
			w.ruleWriter.Save(ctx, &w.rule)

			for _, triggerWatcher := range triggerWatchers {
				if err := triggerWatcher.UpdateLastExecuted(triggerEvt.Time); err != nil {
					w.errorChan <- err

					continue
				}
			}

			action, err := w.actionFactory.Create(w.rule)
			if err != nil {
				w.errorChan <- err

				continue
			}

			action.Execute(ctx)
			span.End()
		case <-ctx.Done():
			w.logger.WithError(ctx.Err()).WithField("rule", w.rule.ID).Warn("stopping ruleWatcher")

			return
		}
	}
}
