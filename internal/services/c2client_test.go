package services

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"

	c2pb "gitlab.com/teserakt/c2/pkg/pb"

	"gitlab.com/teserakt/c2ae/internal/pb"
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

		mockClientFactory.EXPECT().Create().Return(mockClient, nil)

		err := c2.NewClientKey(ctx, clientName)
		if err != expectedError {
			t.Errorf("Expected err to be %s, got %s", expectedError, err)
		}
	})

	t.Run("NewTopicKey creates expected request", func(t *testing.T) {
		topicID := "topicID"

		expectedError := errors.New("expected error response")
		expectedRequest := &c2pb.NewTopicRequest{Topic: topicID}

		mockClient := pb.NewMockC2PbClient(mockCtrl)
		mockClient.EXPECT().NewTopic(gomock.Any(), expectedRequest).Return(nil, expectedError)

		mockClientFactory.EXPECT().Create().Return(mockClient, nil)

		err := c2.NewTopicKey(ctx, topicID)
		if err != expectedError {
			t.Errorf("Expected err to be %s, got %s", expectedError, err)
		}
	})

}
