package watchers

//go:generate mockgen -destination=rule_mocks.go -package watchers -self_package gitlab.com/teserakt/c2ae/internal/engine/watchers gitlab.com/teserakt/c2ae/internal/engine/watchers RuleWatcher

import (
	"context"

	"github.com/go-kit/kit/log"
	"go.opencensus.io/trace"

	"gitlab.com/teserakt/c2ae/internal/engine/actions"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/services"
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
	logger                log.Logger
}

func (w *ruleWatcher) Start(ctx context.Context) {
	var triggerWatchers []TriggerWatcher

	if len(w.rule.Triggers) == 0 {
		w.logger.Log("msg", "rule has no triggers", "rule", w.rule.ID)
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

			w.logger.Log("msg", "rule triggered", "rule", w.rule.ID, "trigger", triggerEvt.Trigger.ID)
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
			w.logger.Log("msg", "stopping ruleWatcher", "rule", w.rule.ID, "reason", ctx.Err())

			return
		}
	}
}