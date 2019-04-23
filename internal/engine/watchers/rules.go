package watchers

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/teserakt/c2se/internal/events"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/services"
)

// RuleWatcher defines methods to implement a rule storage
type RuleWatcher interface {
	Start()
	Stop() error
}

type ruleWatcher struct {
	rule                  models.Rule
	triggerWatcherFactory TriggerWatcherFactory
	ruleWriter            services.RuleWriter
	errorChan             chan<- error
	triggeredChan         chan events.TriggerEvent

	stopChan chan bool
}

// NewRuleWatcher creates a watcher for changes on rules, and registering them on the dispatcher
func NewRuleWatcher(
	rule models.Rule,
	ruleWriter services.RuleWriter,
	triggerWatcherFactory TriggerWatcherFactory,
	triggeredChan chan events.TriggerEvent,
	errorChan chan<- error,
) RuleWatcher {
	return &ruleWatcher{
		rule:                  rule,
		ruleWriter:            ruleWriter,
		triggerWatcherFactory: triggerWatcherFactory,
		triggeredChan:         triggeredChan,
		errorChan:             errorChan,
		stopChan:              make(chan bool),
	}
}

func (w *ruleWatcher) Start() {
	log.Printf("Started rule watcher for rule %d", w.rule.ID)

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
			log.Printf("Rule %d triggered from trigger %d", w.rule.ID, triggerEvt.Trigger.ID)
			w.rule.LastExecuted = triggerEvt.Time
			w.ruleWriter.Save(&w.rule)

			for _, triggerWatcher := range triggerWatchers {
				go triggerWatcher.UpdateLastExecuted(triggerEvt.Time)
			}

			// TODO perform the rule.Action !
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
		log.Printf("Stopped ruleWatcher for rule %d", w.rule.ID)
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("Couldn't stop ruleWatcher for rule %d, maybe it's already stopped ?", w.rule.ID)
	}

	return nil
}
