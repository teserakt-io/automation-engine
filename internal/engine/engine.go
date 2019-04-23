package engine

import (
	"log"

	"gitlab.com/teserakt/c2se/internal/events"

	"gitlab.com/teserakt/c2se/internal/engine/watchers"
	"gitlab.com/teserakt/c2se/internal/services"
)

// ScriptEngine interface describe the public methods available on the script engine
type ScriptEngine interface {
	Start() error
	Stop()
}

type scriptEngine struct {
	ruleService           services.RuleService
	triggerWatcherFactory watchers.TriggerWatcherFactory
	triggeredChan         chan events.TriggerEvent
	errorChan             chan<- error

	ruleWatchers []watchers.RuleWatcher
}

var _ ScriptEngine = &scriptEngine{}

// NewScriptEngine creates a new script engine
func NewScriptEngine(
	ruleService services.RuleService,
	triggerWatcherFactory watchers.TriggerWatcherFactory,
	triggeredChan chan events.TriggerEvent,
	errorChan chan<- error,
) ScriptEngine {
	return &scriptEngine{
		ruleService:           ruleService,
		triggerWatcherFactory: triggerWatcherFactory,
		triggeredChan:         triggeredChan,
		errorChan:             errorChan,
	}
}

func (e *scriptEngine) Start() error {
	rules, err := e.ruleService.All()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleWatcher := watchers.NewRuleWatcher(
			rule,
			e.ruleService,
			e.triggerWatcherFactory,
			e.triggeredChan,
			e.errorChan,
		)

		e.ruleWatchers = append(e.ruleWatchers, ruleWatcher)

		go ruleWatcher.Start()
	}

	return nil
}

func (e *scriptEngine) Stop() {
	for _, w := range e.ruleWatchers {
		if err := w.Stop(); err != nil {
			e.errorChan <- err
		}
	}

	e.ruleWatchers = []watchers.RuleWatcher{}
	log.Println("Stopped script engine")
}
