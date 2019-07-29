package watchers

import (
	"context"
	"errors"
	reflect "reflect"
	"testing"
	"time"

	c2pb "gitlab.com/teserakt/c2/pkg/pb"
	"gitlab.com/teserakt/c2ae/internal/services"

	"gitlab.com/teserakt/c2ae/internal/events"

	"github.com/golang/mock/gomock"

	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
)

func TestSchedulerTriggerWatcher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockValidator := models.NewMockTriggerValidator(mockCtrl)

	triggerSettings := pb.TriggerSettingsTimeInterval{
		Expr: "* * * * *",
	}
	encodedSettings, err := triggerSettings.Encode()
	if err != nil {
		t.Fatalf("Failed to encode settings: %s", err)
	}

	trigger := models.Trigger{Settings: encodedSettings}

	triggeredChan := make(chan TriggerEvent)
	updateChan := make(chan time.Time)
	errorChan := make(chan error)

	watcher := &schedulerWatcher{
		validator:     mockValidator,
		trigger:       trigger,
		triggeredChan: triggeredChan,
		logger:        log.NewNopLogger(),

		updateChan: updateChan,
		errorChan:  errorChan,
	}

	t.Run("Start launch the schedulerWatcher and properly trigger", func(t *testing.T) {
		initialLastExecuted := time.Now().Add(-2 * time.Minute)
		ctx, cancel := context.WithCancel(context.Background())

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

		go watcher.Start(ctx)

		var triggerEvt TriggerEvent

		select {
		case triggerEvt = <-triggeredChan:
			time.Sleep(100 * time.Millisecond)
		case err := <-errorChan:
			t.Errorf("Expected no errors, got %s", err)
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected to receive a triggerEvent")
		}

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
			validator: mockValidator,

			trigger:       invalidTrigger,
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,
		}

		ctx, cancel := context.WithCancel(context.Background())

		mockValidator.EXPECT().ValidateTrigger(invalidTrigger).Return(errors.New("bad trigger"))

		go invalidWatcher.Start(ctx)

		select {
		case err := <-errorChan:
			if _, ok := err.(InvalidTrigger); !ok {
				t.Errorf("Expected error to be of type InvalidTrigger, got %T", err)
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
			validator: mockValidator,

			trigger:       invalidTrigger,
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,
		}

		ctx, cancel := context.WithCancel(context.Background())

		mockValidator.EXPECT().ValidateTrigger(invalidTrigger).Return(errors.New("bad trigger"))

		go invalidWatcher.Start(ctx)

		select {
		case err := <-errorChan:
			if _, ok := err.(InvalidTrigger); !ok {
				t.Errorf("Expected error to be of type InvalidTrigger, got %T", err)
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

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

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
		initialLastExecuted := time.Now().Add(2 * time.Minute)
		watcher.lastExecuted = initialLastExecuted

		ctx, cancel := context.WithCancel(context.Background())

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

		go watcher.Start(ctx)

		var triggerEvt TriggerEvent

		select {
		case triggerEvt = <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", triggerEvt)
		case err := <-errorChan:
			t.Errorf("Expected no errors, got %s", err)
		case <-time.After(100 * time.Millisecond):
		}

		if watcher.lastExecuted != initialLastExecuted {
			t.Errorf("Expected watcher lastExecuted to be %#v, got %#v", initialLastExecuted, watcher.lastExecuted)
		}

		cancel()
	})

}

func TestEventTriggerWatcher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer func() {
		// Give some room for the mockCtrl to finish recording calls made by goroutines...
		time.Sleep(500 * time.Millisecond)
		mockCtrl.Finish()
	}()

	mockValidator := models.NewMockTriggerValidator(mockCtrl)
	mockStreamListenerFactory := events.NewMockStreamListenerFactory(mockCtrl)
	mockTriggerStateService := services.NewMockTriggerStateService(mockCtrl)

	t.Run("Start properly return errors with invalid trigger", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		triggerSettings := pb.TriggerSettingsEvent{}
		triggeredChan := make(chan TriggerEvent)
		updateChan := make(chan time.Time)
		errorChan := make(chan error)

		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Fatalf("failed to encode trigger settings: %v", err)
		}

		trigger := models.Trigger{
			TriggerType: pb.TriggerType_EVENT,
			Settings:    encodedSettings,
		}

		watcher := &eventWatcher{
			validator:             mockValidator,
			streamListenerFactory: mockStreamListenerFactory,
			triggerStateService:   mockTriggerStateService,

			trigger: trigger,
			targets: []models.Target{
				models.Target{Type: pb.TargetType_TOPIC, Expr: "testTopic"},
			},
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,
		}

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(errors.New("bad trigger"))

		go watcher.Start(ctx)

		select {
		case err := <-errorChan:
			if _, ok := err.(InvalidTrigger); !ok {
				t.Errorf("Expected error to be of type InvalidTrigger, got %T", err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an error when starting eventWatcher with invalid trigger")
		}

	})

	t.Run("Start with a valid trigger listen for events and skip or trigger", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		triggerSettings := pb.TriggerSettingsEvent{
			EventType:    pb.EventTypeClientSubscribed,
			MaxOccurence: 2,
		}

		triggeredChan := make(chan TriggerEvent)
		updateChan := make(chan time.Time)
		errorChan := make(chan error)

		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Fatalf("failed to encode trigger settings: %v", err)
		}

		triggerState := models.TriggerState{
			TriggerID: 1,
		}

		trigger := models.Trigger{
			ID:          1,
			TriggerType: pb.TriggerType_EVENT,
			Settings:    encodedSettings,
		}

		initialLastExecuted := time.Now().Add(-2 * time.Minute)
		targetTopic := "testTopic1"

		watcher := &eventWatcher{
			validator:             mockValidator,
			streamListenerFactory: mockStreamListenerFactory,
			triggerStateService:   mockTriggerStateService,
			trigger:               trigger,
			targets: []models.Target{
				models.Target{Type: pb.TargetType_TOPIC, Expr: targetTopic},
			},
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,

			lastExecuted: initialLastExecuted,
		}

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

		mockStreamListener := events.NewMockStreamListener(mockCtrl)
		mockStreamListener.EXPECT().Close()

		eventChan := make(chan c2pb.Event, 1)

		mockStreamListener.EXPECT().C().Return(eventChan).AnyTimes()
		mockStreamListenerFactory.EXPECT().Create(events.DefaultListenerBufSize, pb.EventTypeClientSubscribed).Return(mockStreamListener)

		mockTriggerStateService.EXPECT().ByTriggerID(gomock.Any(), trigger.ID).Return(triggerState, nil)

		go watcher.Start(ctx)

		//  Unknow target
		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "client1", Target: "unknow"}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}

		// Valid target
		mockTriggerStateService.EXPECT().Save(gomock.Any(), gomock.Any())

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "client1", Target: targetTopic}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}

		// Valid target again, expecting trigger
		mockTriggerStateService.EXPECT().Save(gomock.Any(), gomock.Any())

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "client1", Target: targetTopic}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			if reflect.DeepEqual(trigger, evt.Trigger) == false {
				t.Errorf("Expected event trigger to be %#v, got %#v", trigger, evt.Trigger)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected a trigger event, got timeout")
		}

		// 3rd valid target event, counter must have reset and not trigger again
		mockTriggerStateService.EXPECT().Save(gomock.Any(), gomock.Any())

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "client1", Target: targetTopic}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}

		cancel() // do not defer, or mockCtrl will miss some calls...
	})

	t.Run("watcher properly filter events for TargetType_ANY target type", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		triggerSettings := pb.TriggerSettingsEvent{
			EventType:    pb.EventTypeClientSubscribed,
			MaxOccurence: 1,
		}

		triggeredChan := make(chan TriggerEvent)
		updateChan := make(chan time.Time)
		errorChan := make(chan error)

		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Fatalf("failed to encode trigger settings: %v", err)
		}

		triggerState := models.TriggerState{
			TriggerID: 1,
		}

		trigger := models.Trigger{
			ID:          1,
			TriggerType: pb.TriggerType_EVENT,
			Settings:    encodedSettings,
		}

		initialLastExecuted := time.Now().Add(-2 * time.Minute)
		target := "TargetType_ANY"

		watcher := &eventWatcher{
			validator:             mockValidator,
			streamListenerFactory: mockStreamListenerFactory,
			triggerStateService:   mockTriggerStateService,
			trigger:               trigger,
			targets: []models.Target{
				models.Target{Type: pb.TargetType_ANY, Expr: target},
			},
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,

			lastExecuted: initialLastExecuted,
		}

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

		mockStreamListener := events.NewMockStreamListener(mockCtrl)
		mockStreamListener.EXPECT().Close()

		eventChan := make(chan c2pb.Event, 1)

		mockStreamListener.EXPECT().C().Return(eventChan).AnyTimes()
		mockStreamListenerFactory.EXPECT().Create(events.DefaultListenerBufSize, pb.EventTypeClientSubscribed).Return(mockStreamListener)

		mockTriggerStateService.EXPECT().ByTriggerID(gomock.Any(), trigger.ID).Return(triggerState, nil)

		go watcher.Start(ctx)

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: target, Target: "unknow"}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case <-triggeredChan:
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected watcher to trigger, got timeout")
		}

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "", Target: target}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case <-triggeredChan:
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected watcher to trigger, got timeout")
		}

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "something", Target: "something else"}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}

		cancel()
	})

	t.Run("watcher properly filter events for TargetType_CLIENT target type", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		triggerSettings := pb.TriggerSettingsEvent{
			EventType:    pb.EventTypeClientSubscribed,
			MaxOccurence: 1,
		}

		triggeredChan := make(chan TriggerEvent)
		updateChan := make(chan time.Time)
		errorChan := make(chan error)

		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Fatalf("failed to encode trigger settings: %v", err)
		}

		triggerState := models.TriggerState{
			TriggerID: 1,
		}

		trigger := models.Trigger{
			ID:          1,
			TriggerType: pb.TriggerType_EVENT,
			Settings:    encodedSettings,
		}

		initialLastExecuted := time.Now().Add(-2 * time.Minute)
		target := "client1"

		watcher := &eventWatcher{
			validator:             mockValidator,
			streamListenerFactory: mockStreamListenerFactory,
			triggerStateService:   mockTriggerStateService,
			trigger:               trigger,
			targets: []models.Target{
				models.Target{Type: pb.TargetType_CLIENT, Expr: target},
			},
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,

			lastExecuted: initialLastExecuted,
		}

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

		mockStreamListener := events.NewMockStreamListener(mockCtrl)
		mockStreamListener.EXPECT().Close()

		eventChan := make(chan c2pb.Event, 1)

		mockStreamListener.EXPECT().C().Return(eventChan).AnyTimes()
		mockStreamListenerFactory.EXPECT().Create(events.DefaultListenerBufSize, pb.EventTypeClientSubscribed).Return(mockStreamListener)

		mockTriggerStateService.EXPECT().ByTriggerID(gomock.Any(), trigger.ID).Return(triggerState, nil)

		go watcher.Start(ctx)

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: target, Target: "unknow"}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case <-triggeredChan:
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected watcher to trigger, got timeout")
		}

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "", Target: target}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}

		eventChan <- c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "something", Target: "something else"}
		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %v", err)
		case evt := <-triggeredChan:
			t.Errorf("Expected watcher to not trigger, got trigger event: %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}

		cancel()
	})

	t.Run("UpdateLastExecuted properly update the watcher lastExecuted", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		triggerSettings := pb.TriggerSettingsEvent{
			EventType:    pb.EventTypeClientSubscribed,
			MaxOccurence: 1,
		}

		triggeredChan := make(chan TriggerEvent)
		updateChan := make(chan time.Time)
		errorChan := make(chan error)

		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Fatalf("failed to encode trigger settings: %v", err)
		}

		triggerState := models.TriggerState{
			TriggerID: 1,
		}

		trigger := models.Trigger{
			ID:          1,
			TriggerType: pb.TriggerType_EVENT,
			Settings:    encodedSettings,
		}

		initialLastExecuted := time.Now().Add(-2 * time.Minute)
		updatedLastExecuted := time.Now().Add(1 * time.Minute)
		target := "client1"

		watcher := &eventWatcher{
			validator:             mockValidator,
			streamListenerFactory: mockStreamListenerFactory,
			triggerStateService:   mockTriggerStateService,
			trigger:               trigger,
			targets: []models.Target{
				models.Target{Type: pb.TargetType_CLIENT, Expr: target},
			},
			triggeredChan: triggeredChan,
			logger:        log.NewNopLogger(),

			updateChan: updateChan,
			errorChan:  errorChan,

			lastExecuted: initialLastExecuted,
		}

		mockValidator.EXPECT().ValidateTrigger(trigger).Return(nil)

		mockStreamListener := events.NewMockStreamListener(mockCtrl)
		mockStreamListener.EXPECT().Close()

		eventChan := make(chan c2pb.Event, 1)

		mockStreamListener.EXPECT().C().Return(eventChan).AnyTimes()
		mockStreamListenerFactory.EXPECT().Create(events.DefaultListenerBufSize, pb.EventTypeClientSubscribed).Return(mockStreamListener)

		mockTriggerStateService.EXPECT().ByTriggerID(gomock.Any(), trigger.ID).Return(triggerState, nil)

		go watcher.Start(ctx)

		if err := watcher.UpdateLastExecuted(updatedLastExecuted); err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error while waiting for lastExecuted to be updated, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}

		cancel()

		if watcher.lastExecuted != updatedLastExecuted {
			t.Errorf("Expected lastExecuted to be %s, got %s", updatedLastExecuted, watcher.lastExecuted)
		}
	})

}
