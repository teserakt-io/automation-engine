package engine

import (
	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2se/internal/engine/watchers"
	"gitlab.com/teserakt/c2se/internal/services"
)

// ScriptEngine interface describe the public methods available on the script engine
type ScriptEngine interface {
	Start() error
	Stop()
}

type scriptEngine struct {
	ruleService        services.RuleService
	ruleWatcherFactory watchers.RuleWatcherFactory
	logger             log.Logger

	ruleWatchers []watchers.RuleWatcher
}

var _ ScriptEngine = &scriptEngine{}

// NewScriptEngine creates a new script engine
func NewScriptEngine(
	ruleService services.RuleService,
	ruleWatcherFactory watchers.RuleWatcherFactory,
	logger log.Logger,
) ScriptEngine {
	return &scriptEngine{
		ruleService:        ruleService,
		ruleWatcherFactory: ruleWatcherFactory,
		logger:             logger,
	}
}

func (e *scriptEngine) Start() error {
	rules, err := e.ruleService.All()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleWatcher := e.ruleWatcherFactory.Create(rule)
		e.ruleWatchers = append(e.ruleWatchers, ruleWatcher)

		go ruleWatcher.Start()
	}

	return nil
}

func (e *scriptEngine) Stop() {
	for _, w := range e.ruleWatchers {
		if err := w.Stop(); err != nil {
			e.logger.Log("msg", "error while stopping ruleWatcher", "error", err)
		}
	}

	e.ruleWatchers = []watchers.RuleWatcher{}

	e.logger.Log("msg", "stopped script engine")
}
