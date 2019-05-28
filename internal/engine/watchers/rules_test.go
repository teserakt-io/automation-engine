package watchers

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"

	"gitlab.com/teserakt/c2ae/internal/engine/actions"
	"gitlab.com/teserakt/c2ae/internal/events"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/services"
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
	mockActionFactory := actions.NewMockActionFactory(mockCtrl)
	mockAction := actions.NewMockAction(mockCtrl)

	triggeredChan := make(chan events.TriggerEvent, 10)
	errorChan := make(chan error)

	watcher := &ruleWatcher{
		rule:                  rule,
		ruleWriter:            mockRuleWriter,
		triggerWatcherFactory: mockTriggerWatcherFactory,
		actionFactory:         mockActionFactory,
		triggeredChan:         triggeredChan,
		errorChan:             errorChan,
		logger:                log.NewNopLogger(),
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

		ctx, cancel := context.WithCancel(context.Background())

		mockTriggerWatcher1.EXPECT().Start(gomock.Any()).Times(1)
		mockTriggerWatcher2.EXPECT().Start(gomock.Any()).Times(1)

		go watcher.Start(ctx)

		cancel()

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}
	})

	t.Run("Error when creating trigger watchers are forwarded to error chan", func(t *testing.T) {
		expectedError := errors.New("triggerWatcherCreate failed")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedError)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start(gomock.Eq(ctx)).Times(0)
		mockTriggerWatcher2.EXPECT().Start(gomock.Eq(ctx)).Times(1)

		go watcher.Start(ctx)

		select {
		case err := <-errorChan:
			if err != expectedError {
				t.Errorf("Expected err to be %s, got %s", expectedError, err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an error on errorChan")
		}

		cancel()
	})

	t.Run("All triggerWatchers get updated when one of them trigger and action get executed", func(t *testing.T) {
		expectedTime := time.Now()
		modifiedRule := rule
		modifiedRule.LastExecuted = expectedTime

		ctx, cancel := context.WithCancel(context.Background())

		newRuleWatcher := &ruleWatcher{
			rule:                  modifiedRule,
			ruleWriter:            mockRuleWriter,
			triggerWatcherFactory: mockTriggerWatcherFactory,
			actionFactory:         mockActionFactory,
			triggeredChan:         triggeredChan,
			errorChan:             errorChan,
			logger:                log.NewNopLogger(),
		}

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start(gomock.Any()).Times(1)
		mockTriggerWatcher2.EXPECT().Start(gomock.Any()).Times(1)

		mockRuleWriter.EXPECT().Save(&modifiedRule).Times(1)

		mockTriggerWatcher1.EXPECT().UpdateLastExecuted(expectedTime).Times(1)
		mockTriggerWatcher2.EXPECT().UpdateLastExecuted(expectedTime).Times(1)

		mockActionFactory.EXPECT().Create(modifiedRule).Times(1).Return(mockAction, nil)
		mockAction.EXPECT().Execute().Times(1)

		go newRuleWatcher.Start(ctx)

		triggeredChan <- events.TriggerEvent{Trigger: rule.Triggers[1], Time: expectedTime}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(100 * time.Millisecond):
		}

		cancel()
	})

	t.Run("Error is sent on errorChan when the action fail to execute", func(t *testing.T) {
		modifiedRule := models.Rule{
			LastExecuted: time.Now(),
			Triggers:     []models.Trigger{trigger1},
		}

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		ctx, cancel := context.WithCancel(context.Background())

		mockTriggerWatcher1.EXPECT().Start(ctx).Times(1)

		mockRuleWriter.EXPECT().Save(gomock.Any()).Times(1)

		mockTriggerWatcher1.EXPECT().UpdateLastExecuted(gomock.Any()).Times(1)

		expectedError := errors.New("action factory failed to create action")
		mockActionFactory.EXPECT().Create(gomock.Any()).Times(1).Return(nil, expectedError)
		mockAction.EXPECT().Execute().Times(0)

		newRuleWatcher := &ruleWatcher{
			rule:                  modifiedRule,
			ruleWriter:            mockRuleWriter,
			triggerWatcherFactory: mockTriggerWatcherFactory,
			actionFactory:         mockActionFactory,
			triggeredChan:         triggeredChan,
			errorChan:             errorChan,
			logger:                log.NewNopLogger(),
		}

		go newRuleWatcher.Start(ctx)

		triggeredChan <- events.TriggerEvent{Trigger: rule.Triggers[1], Time: time.Now()}

		select {
		case err := <-errorChan:
			if err != expectedError {
				t.Errorf("Expected error to be %s, got %s", expectedError, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected an error when actionFactory failed to create action")
		}

		cancel()
	})
}

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
