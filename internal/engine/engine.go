package engine

import (
	"gitlab.com/teserakt/c2se/internal/events"
)

// ScriptEngine interface describe the public methods available on the script engine
type ScriptEngine interface {
	Run() error
}

type scriptEngine struct {
	scheduler   Scheduler
	dispatcher  events.Dispatcher
	ruleWatcher RuleWatcher
}

var _ ScriptEngine = &scriptEngine{}

// NewScriptEngine creates a new script engine
func NewScriptEngine(scheduler Scheduler, dispatcher events.Dispatcher, ruleWatcher RuleWatcher) ScriptEngine {
	return &scriptEngine{
		scheduler:   scheduler,
		dispatcher:  dispatcher,
		ruleWatcher: ruleWatcher,
	}
}

func (e *scriptEngine) Run() error {
	go e.scheduler.Start()
	go e.dispatcher.Start()

	if err := e.ruleWatcher.Reload(); err != nil {
		return err
	}

	return nil
}
