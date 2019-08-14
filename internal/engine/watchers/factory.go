package watchers

//go:generate mockgen -destination=factory_mocks.go -package watchers -self_package gitlab.com/teserakt/c2ae/internal/engine/watchers gitlab.com/teserakt/c2ae/internal/engine/watchers TriggerWatcherFactory,RuleWatcherFactory

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"gitlab.com/teserakt/c2ae/internal/engine/actions"
	"gitlab.com/teserakt/c2ae/internal/events"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	"gitlab.com/teserakt/c2ae/internal/services"
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
	logger                log.Logger
}

var _ RuleWatcherFactory = &ruleWatcherFactory{}

// NewRuleWatcherFactory creates a new RuleWatcherFactory
func NewRuleWatcherFactory(
	ruleWriter services.RuleWriter,
	triggerWatcherFactory TriggerWatcherFactory,
	actionFactory actions.ActionFactory,
	errorChan chan<- error,
	logger log.Logger,
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
	logger                log.Logger
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
	logger log.Logger,
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
