package engine

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/teserakt-io/automation-engine/internal/engine/watchers"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// AutomationEngine interface describe the public methods available on the automation engine
type AutomationEngine interface {
	Start(context.Context) error
}

type automationEngine struct {
	ruleService        services.RuleService
	ruleWatcherFactory watchers.RuleWatcherFactory
	logger             log.FieldLogger
}

var _ AutomationEngine = &automationEngine{}

// NewAutomationEngine creates a new automation engine
func NewAutomationEngine(
	ruleService services.RuleService,
	ruleWatcherFactory watchers.RuleWatcherFactory,
	logger log.FieldLogger,
) AutomationEngine {
	return &automationEngine{
		ruleService:        ruleService,
		ruleWatcherFactory: ruleWatcherFactory,
		logger:             logger,
	}
}

func (e *automationEngine) Start(ctx context.Context) error {
	rules, err := e.ruleService.All(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleWatcher := e.ruleWatcherFactory.Create(rule)
		e.logger.WithField("rule", rule.ID).Info("started ruleWatcher")
		go ruleWatcher.Start(ctx)
	}

	return nil
}
