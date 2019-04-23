package watchers

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gitlab.com/teserakt/c2se/internal/events"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/services"
)

func TestRuleWatcher(t *testing.T) {
	var wg sync.WaitGroup

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

	ruleWatcher := NewRuleWatcher(
		rule,
		mockRuleWriter,
		mockTriggerWatcherFactory,
		triggeredChan,
		errorChan,
	)

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

		go func() {
			ruleWatcher.Start()
		}()
		wg.Add(1)

		go func() {
			ruleWatcher.Stop()
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}

		wg.Wait()
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

		go func() {
			ruleWatcher.Start()
		}()
		wg.Add(1)

		select {
		case err := <-errorChan:
			if err != expectedError {
				t.Errorf("Expected err to be %s, got %s", expectedError, err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an error on errorChan")
		}

		go func() {
			ruleWatcher.Stop()
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("TriggerWatchers get updated when one trigger", func(t *testing.T) {

		expectedTime := time.Now()

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

		modifiedRule := rule
		modifiedRule.LastExecuted = expectedTime
		mockRuleWriter.EXPECT().Save(&modifiedRule).Times(1)

		mockTriggerWatcher1.EXPECT().UpdateLastExecuted(expectedTime).Times(1)
		mockTriggerWatcher2.EXPECT().UpdateLastExecuted(expectedTime).Times(1)

		mockTriggerWatcher1.EXPECT().Stop().Times(1)
		mockTriggerWatcher2.EXPECT().Stop().Times(1)

		go func() {
			ruleWatcher.Start()
		}()
		wg.Add(1)

		triggeredChan <- events.TriggerEvent{Trigger: rule.Triggers[1], Time: expectedTime}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error on errorChan, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}

		go func() {
			ruleWatcher.Stop()
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()

		wg.Wait()
	})
}
