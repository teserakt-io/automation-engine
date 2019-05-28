package watchers

import (
	"context"
	reflect "reflect"
	"testing"
	"time"

	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/events"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
)

func TestTriggerWatcherFactory(t *testing.T) {
	factory := NewTriggerWatcherFactory(log.NewNopLogger())

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

func TestSchedulerTriggerWatcher(t *testing.T) {

	triggerSettings := pb.TriggerSettingsTimeInterval{
		Expr: "* * * * *",
	}
	encodedSettings, err := triggerSettings.Encode()
	if err != nil {
		t.Fatalf("Failed to encode settings: %s", err)
	}

	trigger := models.Trigger{Settings: encodedSettings}

	triggeredChan := make(chan events.TriggerEvent)
	updateChan := make(chan time.Time)
	errorChan := make(chan error)

	watcher := &schedulerWatcher{
		trigger:       trigger,
		triggeredChan: triggeredChan,
		logger:        log.NewNopLogger(),

		updateChan: updateChan,
		errorChan:  errorChan,
	}

	t.Run("Start launch the schedulerWatcher and properly trigger", func(t *testing.T) {

		initialLastExecuted := time.Now().Add(-2 * time.Minute)

		ctx, cancel := context.WithCancel(context.Background())

		go watcher.Start(ctx)

		var triggerEvt events.TriggerEvent

		select {
		case triggerEvt = <-triggeredChan:
		case err := <-errorChan:
			t.Errorf("Expected no errors, got %s", err)
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected to receive a triggerEvent")
		}

		time.Sleep(10 * time.Millisecond)
		cancel()

		if reflect.DeepEqual(triggerEvt.Trigger, trigger) == false {
			t.Errorf("Expected triggerEvent to contains trigger %#v, got %#v", trigger, triggerEvt.Trigger)
		}

		if !triggerEvt.Time.After(initialLastExecuted) {
			t.Errorf("Expected event time to be after initial lastExecuted")
		}

		if watcher.lastExecuted != triggerEvt.Time {
			t.Errorf("Expected watcher lastExecuted to be %#v, got %#v", triggerEvt.Time, watcher.lastExecuted)
		}
	})

	t.Run("Start handles triggers with invalid settings", func(t *testing.T) {
		invalidTrigger := models.Trigger{
			Settings: nil,
		}

		invalidWatcher := &schedulerWatcher{
			trigger:       invalidTrigger,
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,
		}

		ctx, cancel := context.WithCancel(context.Background())

		go invalidWatcher.Start(ctx)

		select {
		case err := <-errorChan:
			if _, ok := err.(InvalidTriggerSettings); !ok {
				t.Errorf("Expected error to be of type InvalidTriggerSettings, got %T", err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected to get an error")
		}

		cancel()
	})

	t.Run("Start handles triggers with invalid cron expressions", func(t *testing.T) {

		invalidExprSettings := pb.TriggerSettingsTimeInterval{
			Expr: "invalid",
		}

		encodedSettings, err := invalidExprSettings.Encode()
		if err != nil {
			t.Fatalf("Expected no error while encoding settings, got %s", err)
		}

		invalidTrigger := models.Trigger{
			Settings: encodedSettings,
		}

		invalidWatcher := &schedulerWatcher{
			trigger:       invalidTrigger,
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,
		}

		ctx, cancel := context.WithCancel(context.Background())

		go invalidWatcher.Start(ctx)

		select {
		case err := <-errorChan:
			if _, ok := err.(InvalidCronExpr); !ok {
				t.Errorf("Expected error to be of type InvalidCronExpr, got %T", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected to get an error")
		}

		cancel()
	})

	t.Run("UpdateLastExecuted properly update the watcher lastExecuted", func(t *testing.T) {
		watcher.lastExecuted = time.Now()
		updatedTime := time.Now().Add(1 * time.Second)

		ctx, cancel := context.WithCancel(context.Background())

		go watcher.Start(ctx)

		err := watcher.UpdateLastExecuted(updatedTime)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error while waiting for lastExecuted to be updated, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}

		cancel()

		if watcher.lastExecuted != updatedTime {
			t.Errorf("Expected lastExecuted to be %s, got %s", updatedTime, watcher.lastExecuted)
		}
	})

	t.Run("Watcher doesn't trigger if lastExecuted is too recent", func(t *testing.T) {
		// TODO
	})

}
