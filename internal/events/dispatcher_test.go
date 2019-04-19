package events

import (
	"reflect"
	"testing"
	"time"
)

type mockListener struct {
	CallCount   int
	Events      []Event
	OnEventFunc func(evt Event) error
}

var _ Listener = &mockListener{}

func (l *mockListener) OnEvent(evt Event) error {
	l.Events = append(l.Events, evt)
	l.CallCount++

	if l.OnEventFunc != nil {
		return l.OnEventFunc(evt)
	}

	return nil
}

func assertEventContains(t *testing.T, evt Event, eventType Type, source, value interface{}) {
	if evt.Type() != eventType {
		t.Errorf("Expected event type to be %d, got %d", eventType, evt.Type())
	}

	if evt.Source() != source {
		t.Errorf("Expected event source to be %s, got %s", source, evt.Source())
	}

	if evt.Value() != value {
		t.Errorf("Expected event value to be %s, got %s", value, evt.Value())
	}
}

func TestDispatcher(t *testing.T) {
	t.Run("Dispatch properly notify registered listeners", func(t *testing.T) {
		dispatcher := NewDispatcher(100, 100)
		go dispatcher.Start()
		defer dispatcher.ClearListeners()

		event1 := &event{eventType: UnknowEventType, source: "source1", value: "value1"}
		event2 := &event{eventType: UnknowEventType, source: "source2", value: "value2"}
		event3 := &event{eventType: SchedulerTickType, source: "source3", value: "value3"}

		listener1 := &mockListener{}
		listener2 := &mockListener{}
		listener3 := &mockListener{}

		dispatcher.Register(UnknowEventType, listener1)
		dispatcher.Register(UnknowEventType, listener2)
		dispatcher.Register(SchedulerTickType, listener3)

		dispatcher.Dispatch(event1.Type(), event1.Source(), event1.Value())
		dispatcher.Dispatch(event2.Type(), event2.Source(), event2.Value())
		dispatcher.Dispatch(event3.Type(), event3.Source(), event3.Value())

		// Wait a bit for events to process
		time.Sleep(10 * time.Millisecond)

		if listener1.CallCount != 2 {
			t.Errorf("Expected listener 1 to have been called twice, got %d", listener1.CallCount)
		}
		if listener2.CallCount != 2 {
			t.Errorf("Expected listener 2 to have been called twice, got %d", listener2.CallCount)
		}
		if listener3.CallCount != 1 {
			t.Errorf("Expected listener 3 to have been called once, got %d", listener3.CallCount)
		}

		if reflect.DeepEqual(listener1.Events, []Event{event1, event2}) == false {
			t.Errorf(
				"Expected listener 1 to have received events %#v, got %#v",
				[]Event{event1, event2},
				listener1.Events,
			)
		}

		if reflect.DeepEqual(listener2.Events, []Event{event1, event2}) == false {
			t.Errorf(
				"Expected listener 2 to have received events %#v, got %#v",
				[]Event{event1, event2},
				listener2.Events,
			)
		}

		if reflect.DeepEqual(listener3.Events, []Event{event3}) == false {
			t.Errorf(
				"Expected listener 3 to have received events %#v, got %#v",
				[]Event{event3},
				listener3.Events,
			)
		}
	})

	t.Run("Events get discarded if the listener is too slow to process them", func(t *testing.T) {
		dispatcher := NewDispatcher(100, 1)
		go dispatcher.Start()
		defer dispatcher.ClearListeners()

		listener := &mockListener{
			OnEventFunc: func(evt Event) error {
				time.Sleep(10 * time.Millisecond)

				return nil
			},
		}

		dispatcher.Register(UnknowEventType, listener)

		for i := 0; i < 10; i++ {
			dispatcher.Dispatch(UnknowEventType, nil, nil)
		}

		if listener.CallCount >= 10 {
			t.Error("Expected some events to get discared")
		}
	})
}
