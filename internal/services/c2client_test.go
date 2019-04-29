package services

import (
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"gitlab.com/teserakt/c2se/internal/pb"
	e4 "gitlab.com/teserakt/e4common"
)

func TestC2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockC2Requester := NewMockC2Requester(mockCtrl)

	c2 := NewC2(mockC2Requester)

	t.Run("NewClientKey creates expected request", func(t *testing.T) {
		clientID := []byte("expectedClientID")

		expectedRequest := e4.C2Request{
			Command: e4.C2Request_NEW_CLIENT_KEY,
			Id:      clientID,
		}
		expectedError := errors.New("expected error response")

		mockC2Requester.EXPECT().C2Request(expectedRequest).Return(e4.C2Response{}, expectedError)

		err := c2.NewClientKey(clientID)
		if err != expectedError {
			t.Errorf("Expected err to be %s, got %s", expectedError, err)
		}
	})

	t.Run("NewTopicKey creates expected request", func(t *testing.T) {
		topicID := "topicID"

		expectedRequest := e4.C2Request{
			Command: e4.C2Request_NEW_TOPIC,
			Topic:   topicID,
		}
		expectedError := errors.New("expected error response")

		mockC2Requester.EXPECT().C2Request(expectedRequest).Return(e4.C2Response{}, expectedError)

		err := c2.NewTopicKey(topicID)
		if err != expectedError {
			t.Errorf("Expected err to be %s, got %s", expectedError, err)
		}
	})

}

func TestC2Requester(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockC2PbClientFactory := pb.NewMockC2PbClientFactory(mockCtrl)
	mockC2PbClient := pb.NewMockC2PbClient(mockCtrl)

	requester := NewC2Requester(mockC2PbClientFactory)

	t.Run("C2Request properly calls the C2PbClient", func(t *testing.T) {
		expectedRequest := e4.C2Request{
			Command: e4.C2Request_GET_TOPICS,
			Id:      []byte("something"),
		}

		expectedResponse := e4.C2Response{
			Success: true,
			Count:   42,
		}

		mockC2PbClientFactory.EXPECT().Create().Return(mockC2PbClient, nil)
		mockC2PbClient.EXPECT().C2Command(gomock.Any(), &expectedRequest).Return(&expectedResponse, nil)
		mockC2PbClient.EXPECT().Close()

		resp, err := requester.C2Request(expectedRequest)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if reflect.DeepEqual(resp, expectedResponse) == false {
			t.Errorf("Expected response to be %#v, got %#v", expectedResponse, resp)
		}
	})

	t.Run("C2Request handle c2PbClientFactory errors", func(t *testing.T) {
		factoryError := errors.New("factory error")
		mockC2PbClientFactory.EXPECT().Create().Return(nil, factoryError)

		_, err := requester.C2Request(e4.C2Request{})
		if err != factoryError {
			t.Errorf("Expected err to be %s, got %s", factoryError, err)
		}
	})

	t.Run("C2Request handle c2PbClient errors", func(t *testing.T) {
		clientError := errors.New("client error")
		mockC2PbClientFactory.EXPECT().Create().Return(mockC2PbClient, nil)
		mockC2PbClient.EXPECT().C2Command(gomock.Any(), gomock.Any()).Return(nil, clientError)
		mockC2PbClient.EXPECT().Close()

		_, err := requester.C2Request(e4.C2Request{})
		if err != clientError {
			t.Errorf("Expected err to be %s, got %s", clientError, err)
		}
	})

	t.Run("C2Request convert unsuccessful responses to errors", func(t *testing.T) {
		expectedErrorMsg := "response error"
		mockC2PbClientFactory.EXPECT().Create().Return(mockC2PbClient, nil)
		mockC2PbClient.EXPECT().C2Command(gomock.Any(), gomock.Any()).Return(&e4.C2Response{
			Success: false,
			Err:     expectedErrorMsg,
		}, nil)
		mockC2PbClient.EXPECT().Close()

		_, err := requester.C2Request(e4.C2Request{})
		if err == nil || err.Error() != expectedErrorMsg {
			t.Errorf("Expected err to be %s, got %s", expectedErrorMsg, err)
		}
	})
}
