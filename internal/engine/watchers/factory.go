package watchers

//go:generate mockgen -copyright_file ../../../doc/COPYRIGHT_TEMPLATE.txt -destination=factory_mocks.go -package watchers -self_package github.com/teserakt-io/automation-engine/internal/engine/watchers github.com/teserakt-io/automation-engine/internal/engine/watchers TriggerWatcherFactory,RuleWatcherFactory

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/teserakt-io/automation-engine/internal/engine/actions"
	"github.com/teserakt-io/automation-engine/internal/events"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// RuleWatcherFactory allows to create RuleWatchers
type RuleWatcherFactory interface {
	Create(models.Rule) RuleWatcher
}

type ruleWatcherFactory struct {
	ruleWriter            services.RuleWriter
	triggerWatcherFactory TriggerWatcherFactory
	actionFactory         actions.ActionFactory
	errorChan             chan<- error
	logger                log.FieldLogger
}

var _ RuleWatcherFactory = &ruleWatcherFactory{}

// NewRuleWatcherFactory creates a new RuleWatcherFactory
func NewRuleWatcherFactory(
	ruleWriter services.RuleWriter,
	triggerWatcherFactory TriggerWatcherFactory,
	actionFactory actions.ActionFactory,
	errorChan chan<- error,
	logger log.FieldLogger,
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
		triggeredChan:         make(chan TriggerEvent),
		errorChan:             f.errorChan,
		logger:                f.logger,
	}
}

// TriggerWatcherFactory allows to create trigger watchers from a given trigger
// independently of the trigger type
type TriggerWatcherFactory interface {
	Create(
		trigger models.Trigger,
		targets []models.Target,
		lastExecuted time.Time,
		triggeredChan chan<- TriggerEvent,
		errorChan chan<- error,
	) (TriggerWatcher, error)
}

type triggerWatcherFactory struct {
	logger                log.FieldLogger
	streamListenerFactory events.StreamListenerFactory
	triggerStateService   services.TriggerStateService
	validator             models.TriggerValidator
}

var _ TriggerWatcherFactory = (*triggerWatcherFactory)(nil)
var _ TriggerWatcher = (*schedulerWatcher)(nil)
var _ TriggerWatcher = (*eventWatcher)(nil)

// NewTriggerWatcherFactory creates a new watcher factory for given trigger
func NewTriggerWatcherFactory(
	streamListenerFactory events.StreamListenerFactory,
	triggerStateService services.TriggerStateService,
	validator models.TriggerValidator,
	logger log.FieldLogger,
) TriggerWatcherFactory {
	return &triggerWatcherFactory{
		logger:                logger,
		streamListenerFactory: streamListenerFactory,
		triggerStateService:   triggerStateService,
		validator:             validator,
	}
}

func (f *triggerWatcherFactory) Create(
	trigger models.Trigger,
	targets []models.Target,
	lastExecuted time.Time,
	triggeredChan chan<- TriggerEvent,
	errorChan chan<- error,
) (TriggerWatcher, error) {
	var watcher TriggerWatcher

	switch trigger.TriggerType {
	case pb.TriggerType_TIME_INTERVAL:
		watcher = &schedulerWatcher{
			validator:     f.validator,
			trigger:       trigger,
			triggeredChan: triggeredChan,
			errorChan:     errorChan,
			updateChan:    make(chan time.Time, 1),
			lastExecuted:  lastExecuted,
			logger:        f.logger,
		}

	case pb.TriggerType_EVENT:
		watcher = &eventWatcher{
			triggerStateService:   f.triggerStateService,
			validator:             f.validator,
			trigger:               trigger,
			targets:               targets,
			triggeredChan:         triggeredChan,
			errorChan:             errorChan,
			updateChan:            make(chan time.Time),
			lastExecuted:          lastExecuted,
			logger:                f.logger,
			streamListenerFactory: f.streamListenerFactory,
		}

	default:
		return nil, fmt.Errorf("TriggerWatcherFactory don't know how to handle trigger type %s", trigger.TriggerType)
	}

	return watcher, nil
}
