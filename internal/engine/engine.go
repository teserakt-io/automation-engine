package engine

import (
	"context"

	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/engine/watchers"
	"gitlab.com/teserakt/c2ae/internal/services"
)

// AutomationEngine interface describe the public methods available on the automation engine
type AutomationEngine interface {
	Start(context.Context) error
}

type automationEngine struct {
	ruleService        services.RuleService
	ruleWatcherFactory watchers.RuleWatcherFactory
	logger             log.Logger
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

func (e *automationEngine) Start(ctx context.Context) error {
	rules, err := e.ruleService.All()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleWatcher := e.ruleWatcherFactory.Create(rule)
		e.logger.Log("msg", "started ruleWatcher", "rule", rule.ID)
		go ruleWatcher.Start(ctx)
	}

	return nil
}
