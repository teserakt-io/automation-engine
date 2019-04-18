package engine

import (
	"testing"
	"time"

	"gitlab.com/teserakt/c2se/internal/events"

	"github.com/golang/mock/gomock"
	"gitlab.com/teserakt/c2se/internal/mocks"
)

func TestScheduler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDispatcher := mocks.NewMockDispatcher(mockCtrl)

	s := NewScheduler(5*time.Millisecond, mockDispatcher)

	t.Run("Scheduler calls dispatcher at proper interval", func(t *testing.T) {
		mockDispatcher.EXPECT().Dispatch(events.SchedulerTickType, s, gomock.Any()).Times(10)

		go s.Start()
		time.Sleep(50 * time.Millisecond)
		err := s.Stop()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
	})

	t.Run("Scheduler stop when not started returns error", func(t *testing.T) {
		err := s.Stop()
		if err != ErrNotStarted {
			t.Errorf("Expected err to be %s, got %s", ErrNotStarted, err)
		}
	})
}
