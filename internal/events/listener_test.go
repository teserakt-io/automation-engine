package events

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	c2pb "gitlab.com/teserakt/c2/pkg/pb"
	pb "gitlab.com/teserakt/c2ae/internal/pb"
)

func TestStreamListenerFactory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStreamer := NewMockStreamer(mockCtrl)
	f := NewStreamListenerFactory(mockStreamer)

	t.Run("Create returns a new streamListener", func(t *testing.T) {
		expectedChanBufSize := 10
		expectedWhitelist := []pb.EventType{pb.EventTypeClientSubscribed, pb.EventTypeClientUnsubscribed}

		mockStreamer.EXPECT().AddListener(gomock.AssignableToTypeOf(&streamListener{}))

		lis := f.Create(expectedChanBufSize, expectedWhitelist...)
		typedLis := lis.(*streamListener)
		if cap(typedLis.eventChan) != expectedChanBufSize {
			t.Errorf("Expected listener eventChan capacity to be %d, got %d", expectedChanBufSize, cap(typedLis.eventChan))
		}

		if reflect.DeepEqual(typedLis.eventTypeWhitelist, expectedWhitelist) == false {
			t.Errorf("Expected whitelist to be %#v, got %#v", expectedWhitelist, typedLis.eventTypeWhitelist)
		}

		if typedLis.streamer != mockStreamer {
			t.Errorf("Expected streamer to be %#v, got %#v", mockStreamer, typedLis.streamer)
		}
	})
}

func TestStreamListener(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStreamer := NewMockStreamer(mockCtrl)

	t.Run("listener channel contains only whitelisted events", func(t *testing.T) {
		lis := &streamListener{
			eventChan:          make(chan c2pb.Event, 5),
			eventTypeWhitelist: []pb.EventType{pb.EventTypeClientSubscribed},
			streamer:           mockStreamer,
		}

		evt1 := c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED}
		evt2 := c2pb.Event{Type: c2pb.EventType_CLIENT_UNSUBSCRIBED}
		evt3 := c2pb.Event{Type: c2pb.EventType_UNDEFINED}
		evt4 := c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED}

		lis.onEvent(evt1)
		lis.onEvent(evt2)
		lis.onEvent(evt3)
		lis.onEvent(evt4)

		select {
		case evt := <-lis.C():
			if reflect.DeepEqual(evt, evt1) == false {
				t.Errorf("Expected 1st event to be %#v, got %#v", evt1, evt)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an event, got timeout")
		}

		select {
		case evt := <-lis.C():
			if reflect.DeepEqual(evt, evt4) == false {
				t.Errorf("Expected 2nd event to be %#v, got %#v", evt4, evt)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an event, got timeout")
		}

		select {
		case evt := <-lis.C():
			t.Errorf("Expected no more events, got %#v", evt)
		case <-time.After(10 * time.Millisecond):
		}
	})

	t.Run("listener drop oldest message when its channel is full", func(t *testing.T) {
		lis := &streamListener{
			eventChan:          make(chan c2pb.Event, 2),
			eventTypeWhitelist: []pb.EventType{pb.EventTypeClientSubscribed},
			streamer:           mockStreamer,
		}

		evt1 := c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "src1"}
		evt2 := c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "src2"}
		evt3 := c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "src3"}

		lis.onEvent(evt1)
		lis.onEvent(evt2)
		lis.onEvent(evt3)

		select {
		case evt := <-lis.C():
			if reflect.DeepEqual(evt, evt2) == false {
				t.Errorf("Expected event to be %#v, got %#v", evt2, evt)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an event, got timeout")
		}
	})

	t.Run("Closing listeners remove it from dispatcher", func(t *testing.T) {
		lis := &streamListener{
			streamer: mockStreamer,
		}

		mockStreamer.EXPECT().RemoveListener(lis)

		if err := lis.Close(); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
