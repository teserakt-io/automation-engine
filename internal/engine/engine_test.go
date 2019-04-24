package engine

import (
	"errors"
	"testing"

	"gitlab.com/teserakt/c2se/internal/engine/watchers"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/services"

	"github.com/golang/mock/gomock"
)

func TestScriptEngine(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockRuleService := services.NewMockRuleService(mockCtrl)
	mockRuleWatcherFactory := watchers.NewMockRuleWatcherFactory(mockCtrl)

	engine := NewScriptEngine(mockRuleService, mockRuleWatcherFactory)

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

		mockRuleWatcher1.EXPECT().Start().Times(1)
		mockRuleWatcher2.EXPECT().Start().Times(1)
		mockRuleWatcher3.EXPECT().Start().Times(1)

		mockRuleWatcher1.EXPECT().Stop().Times(1).Return(errors.New("failed to stop"))
		mockRuleWatcher2.EXPECT().Stop().Times(1)
		mockRuleWatcher3.EXPECT().Stop().Times(1)

		err := engine.Start()
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		engine.Stop()
	})

	t.Run("Start returns error when it fail to fetch the rules", func(t *testing.T) {
		expectedError := errors.New("ruleService All() failed")
		mockRuleService.EXPECT().All().Times(1).Return(nil, expectedError)

		err := engine.Start()
		if err != expectedError {
			t.Errorf("Expected error to be %s, got %s", expectedError, err)
		}
	})
}
