// Copyright 2020 Teserakt AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	log "github.com/sirupsen/logrus"

	c2pb "github.com/teserakt-io/c2/pkg/pb"

	"github.com/teserakt-io/automation-engine/internal/services"
)

func TestStreamer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer func() {
		// Give some time to the goroutine to switch to running state
		// before letting mockCtrl check its expectations.
		time.Sleep(100 * time.Millisecond)

		mockCtrl.Finish()
	}()

	c2ClientMock := services.NewMockC2(mockCtrl)
	logger := log.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("Add / Remove listeners properly update the streamer", func(t *testing.T) {
		streamer := NewStreamer(c2ClientMock, logger)

		if reflect.DeepEqual(streamer.Listeners(), []StreamListener{}) == false {
			t.Errorf("Expected no listeners, got %#v", streamer.Listeners())
		}

		lis1 := NewMockStreamListener(mockCtrl)
		streamer.AddListener(lis1)

		if reflect.DeepEqual(streamer.Listeners(), []StreamListener{lis1}) == false {
			t.Errorf("Expected no listeners, got %#v", streamer.Listeners())
		}

		lis2 := NewMockStreamListener(mockCtrl)
		streamer.AddListener(lis2)

		if reflect.DeepEqual(streamer.Listeners(), []StreamListener{lis1, lis2}) == false {
			t.Errorf("Expected no listeners, got %#v", streamer.Listeners())
		}

		err := streamer.RemoveListener(lis1)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if reflect.DeepEqual(streamer.Listeners(), []StreamListener{lis2}) == false {
			t.Errorf("Expected no listeners, got %#v", streamer.Listeners())
		}

		err = streamer.RemoveListener(lis1)
		if err != ErrListenerNotFound {
			t.Errorf("Expected error to be %v, got %v", ErrListenerNotFound, err)
		}
	})

	t.Run("StartStream start streaming events from c2Client to all listeners", func(t *testing.T) {
		streamer := NewStreamer(c2ClientMock, logger)

		ctx, cancel := context.WithCancel(context.Background())

		evt := c2pb.Event{Type: c2pb.EventType_CLIENT_SUBSCRIBED, Source: "src1", Target: "target1"}

		streamMock := services.NewMockC2EventStreamClient(mockCtrl)
		streamMock.EXPECT().Recv().Return(&evt, nil).MinTimes(1)

		c2ClientMock.EXPECT().SubscribeToEventStream(ctx).Return(streamMock, nil)

		lis1 := NewMockStreamListener(mockCtrl)
		lis1.EXPECT().onEvent(evt).MinTimes(1)
		lis2 := NewMockStreamListener(mockCtrl)
		lis2.EXPECT().onEvent(evt).MinTimes(1)

		streamer.AddListener(lis1)
		streamer.AddListener(lis2)

		go streamer.StartStream(ctx)

		<-time.After(1 * time.Millisecond)

		cancel()
	})
}
