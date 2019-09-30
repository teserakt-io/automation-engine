package watchers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"

	"github.com/teserakt-io/automation-engine/internal/engine/actions"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/services"
)

func TestRuleWatcher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer func() {
		// Give some time to the goroutine to switch to running state
		// before letting mockCtrl check its expectations.
		time.Sleep(100 * time.Millisecond)
		mockCtrl.Finish()
	}()

	trigger1 := models.Trigger{ID: 1}
	trigger2 := models.Trigger{ID: 2}

	target1 := models.Target{ID: 1}
	target2 := models.Target{ID: 2}

	rule := models.Rule{
		LastExecuted: time.Now(),
		Triggers:     []models.Trigger{trigger1, trigger2},
		Targets:      []models.Target{target1, target2},
	}

	logger := log.NewNopLogger()

	mockRuleWriter := services.NewMockRuleService(mockCtrl)
	mockTriggerWatcherFactory := NewMockTriggerWatcherFactory(mockCtrl)
	mockTriggerWatcher1 := NewMockTriggerWatcher(mockCtrl)
	mockTriggerWatcher2 := NewMockTriggerWatcher(mockCtrl)
	mockActionFactory := actions.NewMockActionFactory(mockCtrl)
	mockAction := actions.NewMockAction(mockCtrl)

	triggeredChan := make(chan TriggerEvent)
	errorChan := make(chan error)

	watcher := &ruleWatcher{
		rule:                  rule,
		ruleWriter:            mockRuleWriter,
		triggerWatcherFactory: mockTriggerWatcherFactory,
		actionFactory:         mockActionFactory,
		triggeredChan:         triggeredChan,
		errorChan:             errorChan,
		logger:                logger,
	}

	t.Run("Start start a triggerWatcher for each trigger", func(t *testing.T) {
		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.Targets, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.Targets, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockTriggerWatcher1.EXPECT().Start(ctx).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})
		mockTriggerWatcher2.EXPECT().Start(ctx).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})

		go watcher.Start(ctx)

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("Error when creating trigger watchers are forwarded to error chan", func(t *testing.T) {
		expectedError := errors.New("triggerWatcherCreate failed")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.Targets, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedError)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.Targets, rule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher2.EXPECT().Start(gomock.Eq(ctx)).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})

		go watcher.Start(ctx)

		select {
		case err := <-errorChan:
			if err != expectedError {
				t.Errorf("Expected err to be %s, got %s", expectedError, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected an error on errorChan")
		}
	})

	t.Run("All triggerWatchers get updated when one of them trigger and action get executed", func(t *testing.T) {
		expectedTime := time.Now()
		modifiedRule := rule
		modifiedRule.LastExecuted = expectedTime

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		newRuleWatcher := &ruleWatcher{
			rule:                  modifiedRule,
			ruleWriter:            mockRuleWriter,
			triggerWatcherFactory: mockTriggerWatcherFactory,
			actionFactory:         mockActionFactory,
			triggeredChan:         triggeredChan,
			errorChan:             errorChan,
			logger:                logger,
		}

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.Targets, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger2, rule.Targets, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher2, nil)

		mockTriggerWatcher1.EXPECT().Start(ctx).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})
		mockTriggerWatcher2.EXPECT().Start(ctx).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})

		mockRuleWriter.EXPECT().Save(gomock.Any(), &modifiedRule).Times(1)

		mockTriggerWatcher1.EXPECT().UpdateLastExecuted(expectedTime).Times(1)
		mockTriggerWatcher2.EXPECT().UpdateLastExecuted(expectedTime).Times(1)

		mockActionFactory.EXPECT().Create(modifiedRule).Times(1).Return(mockAction, nil)
		mockAction.EXPECT().Execute(gomock.Any()).Times(1)

		go newRuleWatcher.Start(ctx)

		triggeredChan <- TriggerEvent{Trigger: modifiedRule.Triggers[1], Time: expectedTime}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("Error is sent on errorChan when the action fail to execute", func(t *testing.T) {
		modifiedRule := models.Rule{
			LastExecuted: time.Now(),
			Triggers:     []models.Trigger{trigger1},
			Targets:      []models.Target{target1, target2},
		}

		mockTriggerWatcherFactory.EXPECT().
			Create(trigger1, rule.Targets, modifiedRule.LastExecuted, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockTriggerWatcher1, nil)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockTriggerWatcher1.EXPECT().Start(ctx).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})

		mockRuleWriter.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1)

		mockTriggerWatcher1.EXPECT().UpdateLastExecuted(gomock.Any()).Times(1)

		expectedError := errors.New("action factory failed to create action")
		mockActionFactory.EXPECT().Create(gomock.Any()).Times(1).Return(nil, expectedError)

		newRuleWatcher := &ruleWatcher{
			rule:                  modifiedRule,
			ruleWriter:            mockRuleWriter,
			triggerWatcherFactory: mockTriggerWatcherFactory,
			actionFactory:         mockActionFactory,
			triggeredChan:         triggeredChan,
			errorChan:             errorChan,
			logger:                logger,
		}

		go newRuleWatcher.Start(ctx)

		select {
		case triggeredChan <- TriggerEvent{Trigger: modifiedRule.Triggers[0], Time: time.Now()}:
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected ruleWatcher triggeredChan to receive messages, but its blocking")
		}

		select {
		case err := <-errorChan:
			if err != expectedError {
				t.Errorf("Expected error to be %s, got %s", expectedError, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected an error when actionFactory failed to create action")
		}
	})
}
