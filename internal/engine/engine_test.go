package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"

	"gitlab.com/teserakt/c2ae/internal/engine/watchers"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/services"
)

func TestAutomationEngine(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRuleService := services.NewMockRuleService(mockCtrl)
	mockRuleWatcherFactory := watchers.NewMockRuleWatcherFactory(mockCtrl)

	engine := NewAutomationEngine(mockRuleService, mockRuleWatcherFactory, log.NewNopLogger())

	rules := []models.Rule{
		models.Rule{ID: 1},
		models.Rule{ID: 2},
		models.Rule{ID: 3},
	}

	mockRuleWatcher1 := watchers.NewMockRuleWatcher(mockCtrl)
	mockRuleWatcher2 := watchers.NewMockRuleWatcher(mockCtrl)
	mockRuleWatcher3 := watchers.NewMockRuleWatcher(mockCtrl)

	t.Run("Start properly start a rule watcher for every rules", func(t *testing.T) {
		mockRuleService.EXPECT().All().Times(1).Return(rules, nil)

		mockRuleWatcherFactory.EXPECT().Create(rules[0]).Times(1).Return(mockRuleWatcher1)
		mockRuleWatcherFactory.EXPECT().Create(rules[1]).Times(1).Return(mockRuleWatcher2)
		mockRuleWatcherFactory.EXPECT().Create(rules[2]).Times(1).Return(mockRuleWatcher3)

		ctx, cancel := context.WithCancel(context.Background())

		mockRuleWatcher1.EXPECT().Start(ctx).Times(1)
		mockRuleWatcher2.EXPECT().Start(ctx).Times(1)
		mockRuleWatcher3.EXPECT().Start(ctx).Times(1)

		err := engine.Start(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		cancel()
	})

	t.Run("Start returns error when it fail to fetch the rules", func(t *testing.T) {
		expectedError := errors.New("ruleService All() failed")
		mockRuleService.EXPECT().All().Times(1).Return(nil, expectedError)

		ctx, cancel := context.WithCancel(context.Background())

		err := engine.Start(ctx)
		if err != expectedError {
			t.Errorf("Expected error to be %v, got %v", expectedError, err)
		}

		cancel()
	})
}
