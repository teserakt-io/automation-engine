package watchers

import (
	reflect "reflect"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	gomock "github.com/golang/mock/gomock"

	"github.com/teserakt-io/automation-engine/internal/engine/actions"
	"github.com/teserakt-io/automation-engine/internal/events"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
)

func TestRuleWatcherFactory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRuleWriter := services.NewMockRuleService(mockCtrl)
	mockTriggerWatcherFactory := NewMockTriggerWatcherFactory(mockCtrl)
	mockActionFactory := actions.NewMockActionFactory(mockCtrl)

	errorChan := make(chan<- error)

	factory := NewRuleWatcherFactory(
		mockRuleWriter,
		mockTriggerWatcherFactory,
		mockActionFactory,
		errorChan,
		log.NewNopLogger(),
	)

	t.Run("Creates returns a properly initialized RuleWatcher", func(t *testing.T) {
		rule := models.Rule{ID: 1}

		watcher := factory.Create(rule)

		typedWatcher, ok := watcher.(*ruleWatcher)
		if !ok {
			t.Errorf("Expected watcher type to be *ruleWatcher, got %T", watcher)
		}

		if reflect.DeepEqual(typedWatcher.rule, rule) == false {
			t.Errorf("Expected rule to be %#v, got %#v", rule, typedWatcher.rule)
		}

		if reflect.DeepEqual(typedWatcher.ruleWriter, mockRuleWriter) == false {
			t.Errorf("Expected ruleWriter to be %p, got %p", mockRuleWriter, typedWatcher.ruleWriter)
		}

		if reflect.DeepEqual(typedWatcher.triggerWatcherFactory, mockTriggerWatcherFactory) == false {
			t.Errorf(
				"Expected triggerWatcherFactory to be %p, got %p",
				mockTriggerWatcherFactory,
				typedWatcher.triggerWatcherFactory,
			)
		}

		if reflect.DeepEqual(typedWatcher.errorChan, errorChan) == false {
			t.Errorf("Expected errorChan to be %p, got %p", errorChan, typedWatcher.errorChan)
		}

		if reflect.DeepEqual(typedWatcher.actionFactory, mockActionFactory) == false {
			t.Errorf("Expected action factory to be %p, got %p", mockActionFactory, typedWatcher.actionFactory)
		}
	})
}

func TestTriggerWatcherFactory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStreamListenerFactory := events.NewMockStreamListenerFactory(mockCtrl)
	mockTriggerStateService := services.NewMockTriggerStateService(mockCtrl)
	mockValidator := models.NewMockValidator(mockCtrl)

	factory := NewTriggerWatcherFactory(mockStreamListenerFactory, mockTriggerStateService, mockValidator, log.NewNopLogger())

	expectedLastExecuted := time.Now()

	triggeredChan := make(chan<- TriggerEvent)
	errorChan := make(chan<- error)

	t.Run("Factory creates schedulerWatcher", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_TIME_INTERVAL,
		}

		watcher, err := factory.Create(trigger, nil, expectedLastExecuted, triggeredChan, errorChan)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		typedWatcher, ok := watcher.(*schedulerWatcher)
		if !ok {
			t.Fatalf("Expected watcher to be a *schedulerWatcher, got a %T", watcher)
		}

		if reflect.DeepEqual(typedWatcher.trigger, trigger) == false {
			t.Errorf("Expected watcher trigger to be %#v, got %#v", trigger, typedWatcher.trigger)
		}

		if reflect.DeepEqual(typedWatcher.lastExecuted, expectedLastExecuted) == false {
			t.Errorf(
				"Expected default last executed to be %#v, got %#v",
				expectedLastExecuted,
				typedWatcher.lastExecuted,
			)
		}

		if reflect.DeepEqual(typedWatcher.triggeredChan, triggeredChan) == false {
			t.Errorf("Expected watcher triggeredChan to be %#v, got %#v", triggeredChan, typedWatcher.triggeredChan)
		}

		if reflect.DeepEqual(typedWatcher.errorChan, errorChan) == false {
			t.Errorf("Expected watcher errorChan to be %#v, got %#v", errorChan, typedWatcher.errorChan)
		}

		if typedWatcher.updateChan == nil {
			t.Errorf("Expected watcher updateChan to be not nil")
		}
	})

	t.Run("Factory creates eventWatcher", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_EVENT,
		}

		targets := []models.Target{
			models.Target{ID: 1},
			models.Target{ID: 2},
			models.Target{ID: 3},
		}
		watcher, err := factory.Create(trigger, targets, expectedLastExecuted, triggeredChan, errorChan)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		typedWatcher, ok := watcher.(*eventWatcher)
		if !ok {
			t.Fatalf("Expected watcher to be a *schedulerWatcher, got a %T", watcher)
		}

		if reflect.DeepEqual(typedWatcher.trigger, trigger) == false {
			t.Errorf("Expected watcher trigger to be %#v, got %#v", trigger, typedWatcher.trigger)
		}

		if reflect.DeepEqual(typedWatcher.lastExecuted, expectedLastExecuted) == false {
			t.Errorf(
				"Expected default last executed to be %#v, got %#v",
				expectedLastExecuted,
				typedWatcher.lastExecuted,
			)
		}

		if reflect.DeepEqual(typedWatcher.triggeredChan, triggeredChan) == false {
			t.Errorf("Expected watcher triggeredChan to be %#v, got %#v", triggeredChan, typedWatcher.triggeredChan)
		}

		if reflect.DeepEqual(typedWatcher.errorChan, errorChan) == false {
			t.Errorf("Expected watcher errorChan to be %#v, got %#v", errorChan, typedWatcher.errorChan)
		}

		if reflect.DeepEqual(typedWatcher.targets, targets) == false {
			t.Errorf("Expected watcher targets to be %#v, got %#v", targets, typedWatcher.targets)
		}

		if typedWatcher.updateChan == nil {
			t.Errorf("Expected watcher updateChan to be not nil")
		}
	})

	t.Run("Factory returns error on unknow trigger type", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_UNDEFINED_TRIGGER,
		}

		_, err := factory.Create(trigger, nil, expectedLastExecuted, triggeredChan, errorChan)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
