package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.com/teserakt/c2se/internal/mocks"
)

func TestEngine(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockScheduler := mocks.NewMockScheduler(mockCtrl)
	mockDispatcher := mocks.NewMockDispatcher(mockCtrl)
	mockRuleWatcher := mocks.NewMockRuleWatcher(mockCtrl)

	e := NewScriptEngine(mockScheduler, mockDispatcher, mockRuleWatcher)

	t.Run("Engine Run start expected services", func(t *testing.T) {
		mockScheduler.EXPECT().Start().Times(1)
		mockDispatcher.EXPECT().Start().Times(1)
		mockRuleWatcher.EXPECT().Reload().Times(1)

		err := e.Run()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		// Wait a bit as the engine spawn goroutines, which may still be starting when Run() returns
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Engine handle Reload errors", func(t *testing.T) {
		expectedErr := errors.New("expectedErr")

		mockScheduler.EXPECT().Start().Times(1)
		mockDispatcher.EXPECT().Start().Times(1)
		mockRuleWatcher.EXPECT().Reload().Times(1).Return(expectedErr)

		err := e.Run()
		if err != expectedErr {
			t.Errorf("Expected err to be %s, got %s", expectedErr, err)
		}

		// Wait a bit as the engine spawn goroutines, which may still be starting when Run() returns
		time.Sleep(10 * time.Millisecond)
	})
}
