package watchers

import (
	"fmt"
	"time"

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
	triggeredChan         chan events.TriggerEvent
	errorChan             chan<- error
	logger                log.Logger
}

var _ RuleWatcherFactory = &ruleWatcherFactory{}

// NewRuleWatcherFactory creates a new RuleWatcherFactory
func NewRuleWatcherFactory(
	ruleWriter services.RuleWriter,
	triggerWatcherFactory TriggerWatcherFactory,
	actionFactory actions.ActionFactory,
	triggeredChan chan events.TriggerEvent,
	errorChan chan<- error,
	logger log.Logger,
) RuleWatcherFactory {
	return &ruleWatcherFactory{
		ruleWriter:            ruleWriter,
		triggerWatcherFactory: triggerWatcherFactory,
		actionFactory:         actionFactory,
		triggeredChan:         triggeredChan,
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
		triggeredChan:         f.triggeredChan,
		errorChan:             f.errorChan,
		logger:                f.logger,
		stopChan:              make(chan bool),
	}
}

// RuleWatcher defines methods to implement a rule watcher.
// It is responsible of monitoring the rule trigger(s), and execute the rule action
// when the trigger conditions are met.
type RuleWatcher interface {
	Start()
	Stop() error
}

type ruleWatcher struct {
	rule                  models.Rule
	triggerWatcherFactory TriggerWatcherFactory
	actionFactory         actions.ActionFactory
	ruleWriter            services.RuleWriter
	errorChan             chan<- error
	triggeredChan         chan events.TriggerEvent
	logger                log.Logger

	stopChan chan bool
}

func (w *ruleWatcher) Start() {
	w.logger.Log("msg", "started rule watcher", "rule", w.rule.ID)

	var triggerWatchers []TriggerWatcher

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

		go triggerWatcher.Start()
	}

	for {
		select {
		case triggerEvt := <-w.triggeredChan:
			w.logger.Log("msg", "rule triggered", "rule", w.rule.ID, "trigger", triggerEvt.Trigger.ID)
			w.rule.LastExecuted = triggerEvt.Time
			w.ruleWriter.Save(&w.rule)

			for _, triggerWatcher := range triggerWatchers {
				go triggerWatcher.UpdateLastExecuted(triggerEvt.Time)
			}

			action, err := w.actionFactory.Create(w.rule)
			if err != nil {
				w.errorChan <- err

				continue
			}

			action.Execute()

		case <-w.stopChan:
			for _, triggerWatcher := range triggerWatchers {
				if err := triggerWatcher.Stop(); err != nil {
					w.errorChan <- err
				}
			}

			return
		}
	}
}

func (w *ruleWatcher) Stop() error {
	select {
	case w.stopChan <- true:
		w.logger.Log("msg", "stopped ruleWatcher", "rule", w.rule.ID)
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("Couldn't stop ruleWatcher for rule %d, maybe it's already stopped ?", w.rule.ID)
	}

	return nil
}
