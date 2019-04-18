package engine

import (
	"errors"
	"testing"

	"gitlab.com/teserakt/c2se/internal/events"

	"github.com/golang/mock/gomock"
	"gitlab.com/teserakt/c2se/internal/mocks"
	"gitlab.com/teserakt/c2se/internal/models"
)

func setupMocks(t *testing.T) {
	var t1, t2, t3, t4 models.Trigger
	var l1, l2, l3, l4 events.Listener
	var e1, e2, e3, e4 events.Type

	allRules := []models.Rule{
		models.Rule{
			Triggers: []models.Trigger{t1, t2},
		},
		models.Rule{
			Triggers: []models.Trigger{t3},
		},
		models.Rule{
			Triggers: []models.Trigger{t4},
		},
	}

	l2err := errors.New("failed to create l2 listener")

	gomock.InOrder(
		mockRuleService.EXPECT().All().Times(1).Return(allRules, nil),

		mockDispatcher.EXPECT().ClearListeners().Times(1),

		mockTriggerListenerFactory.EXPECT().Create(t1).Times(1).Return(e1, l1, nil),
		mockDispatcher.EXPECT().Register(e1, l1).Times(1),
		mockTriggerListenerFactory.EXPECT().Create(t2).Times(1).Return(e2, l2, l2err),
		mockDispatcher.EXPECT().Register(e2, l2).Times(0),
		mockTriggerListenerFactory.EXPECT().Create(t3).Times(1).Return(e3, l3, nil),
		mockDispatcher.EXPECT().Register(e3, l3).Times(1),
		mockTriggerListenerFactory.EXPECT().Create(t4).Times(1).Return(e4, l4, nil),
		mockDispatcher.EXPECT().Register(e4, l4).Times(1),

		mockDispatcher.EXPECT().Register(events.RulesModifiedType, watcher).Times(1),
	)

}

var watcher RuleWatcher
var mockRuleService *mocks.MockRuleService
var mockDispatcher *mocks.MockDispatcher
var mockTriggerListenerFactory *mocks.MockTriggerListenerFactory

func TestRuleWatcher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRuleService = mocks.NewMockRuleService(mockCtrl)
	mockDispatcher = mocks.NewMockDispatcher(mockCtrl)
	mockTriggerListenerFactory = mocks.NewMockTriggerListenerFactory(mockCtrl)

	watcher = NewRuleWatcher(mockRuleService, mockDispatcher, mockTriggerListenerFactory)

	t.Run("Reload properly reload what's registered on the dispatcher", func(t *testing.T) {
		setupMocks(t)

		err := watcher.Reload()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

	})

	t.Run("Reload still clear and register the watcher to RuleModified event when no rules exists", func(t *testing.T) {
		gomock.InOrder(
			mockRuleService.EXPECT().All().Times(1).Return(nil, nil),
			mockDispatcher.EXPECT().ClearListeners().Times(1),
			mockDispatcher.EXPECT().Register(events.RulesModifiedType, watcher).Times(1),
		)
		err := watcher.Reload()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
	})

	t.Run("Reload returns error when it fail to fetch the rules", func(t *testing.T) {
		expectedErr := errors.New("db error")

		mockRuleService.EXPECT().All().Times(1).Return(nil, expectedErr)

		err := watcher.Reload()
		if err != expectedErr {
			t.Errorf("Expected err to be %s, got %s", expectedErr, err)
		}
	})

	t.Run("OnEvent calls Reload", func(t *testing.T) {
		setupMocks(t)

		err := watcher.OnEvent(nil)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
	})

	t.Run("OnEvent properly forward reload errors", func(t *testing.T) {
		mockRuleService.EXPECT().All().Times(1).Return(nil, errors.New("db error"))

		err := watcher.OnEvent(nil)
		if err == nil {
			t.Error("Expected err to be not nil")
		}
	})
}
