package engine

import (
	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/engine/watchers"
	"gitlab.com/teserakt/c2ae/internal/services"
)

// AutomationEngine interface describe the public methods available on the automation engine
type AutomationEngine interface {
	Start() error
	Stop()
}

type automationEngine struct {
	ruleService        services.RuleService
	ruleWatcherFactory watchers.RuleWatcherFactory
	logger             log.Logger

	ruleWatchers []watchers.RuleWatcher
}

var _ AutomationEngine = &automationEngine{}

// NewAutomationEngine creates a new automation engine
func NewAutomationEngine(
	ruleService services.RuleService,
	ruleWatcherFactory watchers.RuleWatcherFactory,
	logger log.Logger,
) AutomationEngine {
	return &automationEngine{
		ruleService:        ruleService,
		ruleWatcherFactory: ruleWatcherFactory,
		logger:             logger,
	}
}

func (e *automationEngine) Start() error {
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

func (e *automationEngine) Stop() {
	for _, w := range e.ruleWatchers {
		if err := w.Stop(); err != nil {
			e.logger.Log("msg", "error while stopping ruleWatcher", "error", err)
		}
	}

	e.ruleWatchers = []watchers.RuleWatcher{}

	e.logger.Log("msg", "stopped automation engine")
}
