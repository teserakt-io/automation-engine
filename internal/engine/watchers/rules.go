package watchers

import (
	"context"

	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/engine/actions"
	"gitlab.com/teserakt/c2ae/internal/events"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/services"
)

//go:generate mockgen -destination=rules_mocks.go -package watchers -self_package gitlab.com/teserakt/c2ae/internal/engine/watchers gitlab.com/teserakt/c2ae/internal/engine/watchers RuleWatcherFactory,RuleWatcher

// RuleWatcherFactory allows to create RuleWatchers
type RuleWatcherFactory interface {
	Create(models.Rule) RuleWatcher
}

type ruleWatcherFactory struct {
	ruleWriter            services.RuleWriter
	triggerWatcherFactory TriggerWatcherFactory
	actionFactory         actions.ActionFactory
	errorChan             chan<- error
	logger                log.Logger
}

var _ RuleWatcherFactory = &ruleWatcherFactory{}

// NewRuleWatcherFactory creates a new RuleWatcherFactory
func NewRuleWatcherFactory(
	ruleWriter services.RuleWriter,
	triggerWatcherFactory TriggerWatcherFactory,
	actionFactory actions.ActionFactory,
	errorChan chan<- error,
	logger log.Logger,
) RuleWatcherFactory {
	return &ruleWatcherFactory{
		ruleWriter:            ruleWriter,
		triggerWatcherFactory: triggerWatcherFactory,
		actionFactory:         actionFactory,
		errorChan:             errorChan,
		logger:                logger,
	}
}

func (f *ruleWatcherFactory) Create(rule models.Rule) RuleWatcher {
	return &ruleWatcher{
		rule:                  rule,
		ruleWriter:            f.ruleWriter,
		triggerWatcherFactory: f.triggerWatcherFactory,
		actionFactory:         f.actionFactory,
		triggeredChan:         make(chan events.TriggerEvent),
		errorChan:             f.errorChan,
		logger:                f.logger,
	}
}

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
	triggeredChan         chan events.TriggerEvent
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

			w.logger.Log("msg", "rule triggered", "rule", w.rule.ID, "trigger", triggerEvt.Trigger.ID)
			w.rule.LastExecuted = triggerEvt.Time
			w.ruleWriter.Save(&w.rule)

			for _, triggerWatcher := range triggerWatchers {
				triggerWatcher.UpdateLastExecuted(triggerEvt.Time)
			}

			action, err := w.actionFactory.Create(w.rule)
			if err != nil {
				w.errorChan <- err

				continue
			}

			action.Execute()

		case <-ctx.Done():
			w.logger.Log("msg", "stopping ruleWatcher", "rule", w.rule.ID, "reason", ctx.Err())

			return
		}
	}
}
