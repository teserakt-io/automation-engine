package watchers

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gitlab.com/teserakt/c2se/internal/events"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/services"
)

func TestRuleWatcher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	trigger1 := models.Trigger{ID: 1}
	trigger2 := models.Trigger{ID: 2}

	rule := models.Rule{
		LastExecuted: time.Now(),
		Triggers:     []models.Trigger{trigger1, trigger2},
	}

	mockRuleWriter := services.NewMockRuleService(mockCtrl)
	mockTriggerWatcherFactory := NewMockTriggerWatcherFactory(mockCtrl)
	mockTriggerWatcher1 := NewMockTriggerWatcher(mockCtrl)
	mockTriggerWatcher2 := NewMockTriggerWatcher(mockCtrl)

	triggeredChan := make(chan events.TriggerEvent, 10)
	errorChan := make(chan error)

	watcher := &ruleWatcher{
		rule:                  rule,
		ruleWriter:            mockRuleWriter,
		triggerWatcherFactory: mockTriggerWatcherFactory,
		triggeredChan:         triggeredChan,
		errorChan:             errorChan,
		stopChan:              make(chan bool),
	}

	t.Run("Start start a triggerWatcher for each triggers", func(t *testing.T) {
		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start().Times(1)
		mockTriggerWatcher2.EXPECT().Start().Times(1)

		mockTriggerWatcher1.EXPECT().Stop().Times(1)
		mockTriggerWatcher2.EXPECT().Stop().Times(1)

		go watcher.Start()

		watcher.Stop()

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}
	})

	t.Run("Error when creating trigger watchers are forwarded to error chan", func(t *testing.T) {
		expectedError := errors.New("triggerWatcherCreate failed")

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedError)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start().Times(0)
		mockTriggerWatcher2.EXPECT().Start().Times(1)

		mockTriggerWatcher2.EXPECT().Stop().Times(1)

		go watcher.Start()

		select {
		case err := <-errorChan:
			if err != expectedError {
				t.Errorf("Expected err to be %s, got %s", expectedError, err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an error on errorChan")
		}

		watcher.Stop()
	})

	t.Run("All triggerWatchers get updated when one of them trigger", func(t *testing.T) {
		expectedTime := time.Now()
		modifiedRule := rule
		modifiedRule.LastExecuted = expectedTime

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start().Times(1)
		mockTriggerWatcher2.EXPECT().Start().Times(1)

		mockRuleWriter.EXPECT().Save(&modifiedRule).Times(1)

		mockTriggerWatcher1.EXPECT().UpdateLastExecuted(expectedTime).Times(1)
		mockTriggerWatcher2.EXPECT().UpdateLastExecuted(expectedTime).Times(1)

		mockTriggerWatcher1.EXPECT().Stop().Times(1)
		mockTriggerWatcher2.EXPECT().Stop().Times(1)

		newRuleWatcher := &ruleWatcher{
			rule:                  modifiedRule,
			ruleWriter:            mockRuleWriter,
			triggerWatcherFactory: mockTriggerWatcherFactory,
			triggeredChan:         triggeredChan,
			errorChan:             errorChan,
			stopChan:              make(chan bool),
		}

		go newRuleWatcher.Start()

		triggeredChan <- events.TriggerEvent{Trigger: rule.Triggers[1], Time: expectedTime}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}

		newRuleWatcher.Stop()
	})

	t.Run("Stopping a non running RuleWatcher returns an error", func(t *testing.T) {
		err := watcher.Stop()
		if err == nil {
			t.Errorf("Expected an error")
		}
	})

	t.Run("Stop forward errors when a triggerwatcher fail to stop", func(t *testing.T) {
		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start().Times(1)
		mockTriggerWatcher2.EXPECT().Start().Times(1)

		expectedErr := errors.New("triggerWatcher stop failure")
		mockTriggerWatcher1.EXPECT().Stop().Times(1).Return(expectedErr)
		mockTriggerWatcher2.EXPECT().Stop().Times(1)

		go watcher.Start()

		watcher.Stop()

		select {
		case err := <-errorChan:
			if err != expectedErr {
				t.Errorf("Expected error to be %s, got %s", expectedErr, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected an error")
		}
	})
}

func TestRuleWatcherFactory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRuleWriter := services.NewMockRuleService(mockCtrl)
	mockTriggerWatcherFactory := NewMockTriggerWatcherFactory(mockCtrl)

	triggeredChan := make(chan events.TriggerEvent)
	errorChan := make(chan<- error)

	factory := NewRuleWatcherFactory(
		mockRuleWriter,
		mockTriggerWatcherFactory,
		triggeredChan,
		errorChan,
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

		if reflect.DeepEqual(typedWatcher.triggeredChan, triggeredChan) == false {
			t.Errorf("Expected errorChan to be %p, got %p", triggeredChan, typedWatcher.triggeredChan)
		}

		if reflect.DeepEqual(typedWatcher.errorChan, errorChan) == false {
			t.Errorf("Expected errorChan to be %p, got %p", errorChan, typedWatcher.errorChan)
		}

		if typedWatcher.stopChan == nil {
			t.Error("Expected stopChan to be not nil")
		}

	})
}
