package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	c2pb "github.com/teserakt-io/c2/pkg/pb"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

func TestC2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClientFactory := pb.NewMockC2PbClientFactory(mockCtrl)
	c2 := NewC2(mockClientFactory)
	ctx := context.Background()

	t.Run("NewClientKey creates expected request", func(t *testing.T) {
		clientName := "expectedClientName"

		expectedError := errors.New("expected error response")
		expectedRequest := &c2pb.NewClientKeyRequest{Client: &c2pb.Client{Name: clientName}}

		mockClient := pb.NewMockC2PbClient(mockCtrl)
		mockClient.EXPECT().NewClientKey(gomock.Any(), expectedRequest).Return(nil, expectedError)
		mockClient.EXPECT().Close()

		mockClientFactory.EXPECT().Create().Return(mockClient, nil)

		err := c2.NewClientKey(ctx, clientName)
		if err != expectedError {
			t.Errorf("Expected err to be %v, got %v", expectedError, err)
		}
	})

	t.Run("NewTopicKey creates expected request", func(t *testing.T) {
		topicID := "topicID"

		expectedError := errors.New("expected error response")
		expectedRequest := &c2pb.NewTopicRequest{Topic: topicID}

		mockClient := pb.NewMockC2PbClient(mockCtrl)
		mockClient.EXPECT().NewTopic(gomock.Any(), expectedRequest).Return(nil, expectedError)
		mockClient.EXPECT().Close()

		mockClientFactory.EXPECT().Create().Return(mockClient, nil)

		err := c2.NewTopicKey(ctx, topicID)
		if err != expectedError {
			t.Errorf("Expected err to be %v, got %v", expectedError, err)
		}
	})

	t.Run("SubscribeToClientStream creates expected request", func(t *testing.T) {
		mockClient := pb.NewMockC2PbClient(mockCtrl)

		expectedRequest := &c2pb.SubscribeToEventStreamRequest{}
		expectedStream := NewMockC2EventStreamClient(mockCtrl)

		mockClient.EXPECT().SubscribeToEventStream(gomock.Any(), expectedRequest).Return(expectedStream, nil)

		mockClientFactory.EXPECT().Create().Return(mockClient, nil)
		stream, err := c2.SubscribeToEventStream(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if reflect.DeepEqual(stream, expectedStream) == false {
			t.Errorf("Expected stream to be %#v, got %#v", expectedStream, stream)
		}
	})
}
