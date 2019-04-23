package watchers

import (
	reflect "reflect"
	"testing"
	"time"

	"gitlab.com/teserakt/c2se/internal/events"

	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
)

func TestTriggerWatcherFactory(t *testing.T) {
	factory := NewTriggerWatcherFactory()

	expectedLastExecuted := time.Now()

	triggeredChan := make(chan<- events.TriggerEvent)
	errorChan := make(chan<- error)

	t.Run("Factory creates schedulerWatcher", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_TIME_INTERVAL,
		}

		watcher, err := factory.Create(trigger, expectedLastExecuted, triggeredChan, errorChan)
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

		if typedWatcher.stopChan == nil {
			t.Errorf("Expected watcher stopChan to be not nil")
		}
	})

	t.Run("Factory creates clientSubscribedWatcher", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_CLIENT_SUBSCRIBED,
		}

		watcher, err := factory.Create(trigger, expectedLastExecuted, triggeredChan, errorChan)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		typedWatcher, ok := watcher.(*clientSubscribedWatcher)
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

		if typedWatcher.stopChan == nil {
			t.Errorf("Expected watcher stopChan to be not nil")
		}
	})

	t.Run("Factory creates clientUnsubscribedWatcher", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_CLIENT_UNSUBSCRIBED,
		}

		watcher, err := factory.Create(trigger, expectedLastExecuted, triggeredChan, errorChan)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		typedWatcher, ok := watcher.(*clientUnsubscribedWatcher)
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

		if typedWatcher.stopChan == nil {
			t.Errorf("Expected watcher stopChan to be not nil")
		}
	})

	t.Run("Factory returns error on unknow trigger type", func(t *testing.T) {
		trigger := models.Trigger{
			TriggerType: pb.TriggerType_UNDEFINED_TRIGGER,
		}

		_, err := factory.Create(trigger, expectedLastExecuted, triggeredChan, errorChan)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
