package engine

//go:generate mockgen -destination=../mocks/engine_watchers.go -package=mocks gitlab.com/teserakt/c2se/internal/engine RuleWatcher

import (
	"fmt"
	"log"

	"gitlab.com/teserakt/c2se/internal/events"
	"gitlab.com/teserakt/c2se/internal/services"
)

// RuleWatcher defines methods to implement a rule storage
type RuleWatcher interface {
	events.Listener
	Reload() error
}

type ruleWatcher struct {
	reader                 services.RuleReader
	dispatcher             events.Dispatcher
	triggerListenerFactory events.TriggerListenerFactory
}

// NewRuleWatcher creates a watcher for changes on rules, and registering them on the dispatcher
func NewRuleWatcher(
	reader services.RuleReader,
	dispatcher events.Dispatcher,
	triggerListenerFactory events.TriggerListenerFactory,
) RuleWatcher {
	return &ruleWatcher{
		reader:                 reader,
		dispatcher:             dispatcher,
		triggerListenerFactory: triggerListenerFactory,
	}
}

// Reload refresh the rule list from the rule reader
func (s *ruleWatcher) Reload() error {
	log.Println("RuleWatcher started reloading rules")

	rules, err := s.reader.All()
	if err != nil {
		return err
	}

	s.dispatcher.ClearListeners()

	for _, rule := range rules {
		for _, trigger := range rule.Triggers {
			eventType, listener, err := s.triggerListenerFactory.Create(trigger)
			if err != nil {
				log.Printf("ruleWatcher error while creating trigger %d listener: %s", trigger.ID, err)

				continue
			}

			s.dispatcher.Register(eventType, listener)
			log.Printf(
				"Registered listener %p for event %s, rule %d and trigger %d",
				listener,
				events.EventStrings[eventType],
				trigger.RuleID,
				trigger.ID,
			)
		}
	}

	s.dispatcher.Register(events.RulesModifiedType, s)

	log.Println("RuleWatcher finished reloading rules")

	return nil
}

// On implements events.Listener
func (s *ruleWatcher) OnEvent(evt events.Event) error {
	if err := s.Reload(); err != nil {
		return fmt.Errorf("error while reloading ruleWatcher: %s", err)
	}

	return nil
}
