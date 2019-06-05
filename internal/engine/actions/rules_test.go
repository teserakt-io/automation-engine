package actions

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"

	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	"gitlab.com/teserakt/c2ae/internal/services"
)

func TestKeyRotationAction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockC2Client := services.NewMockC2(mockCtrl)

	errorChan := make(chan error)

	t.Run("Execute calls the expected C2 client method", func(t *testing.T) {
		targets := []models.Target{
			models.Target{Type: pb.TargetType_CLIENT, Expr: "client1"},
			models.Target{Type: pb.TargetType_TOPIC, Expr: "topic1"},
			models.Target{Type: pb.TargetType_ANY, Expr: "n/a"},
			models.Target{Type: pb.TargetType_CLIENT, Expr: "client2"},
			models.Target{Type: pb.TargetType_TOPIC, Expr: "topic2"},
		}

		gomock.InOrder(
			mockC2Client.EXPECT().NewClientKey(gomock.Any(), "client1"),
			mockC2Client.EXPECT().NewTopicKey(gomock.Any(), "topic1"),
			mockC2Client.EXPECT().NewClientKey(gomock.Any(), "client2"),
			mockC2Client.EXPECT().NewTopicKey(gomock.Any(), "topic2"),
		)

		action := &keyRotationAction{
			targets:   targets,
			c2Client:  mockC2Client,
			errorChan: errorChan,
			logger:    log.NewNopLogger(),
		}

		go action.Execute(context.Background())

		select {
		case err := <-errorChan:
			terr, ok := err.(UnsupportedTargetType)
			if !ok {
				t.Errorf("Expected an UnsupportedTargetType error, got %T", err)
			}

			if reflect.DeepEqual(terr.Action, action) == false {
				t.Errorf("Expected error action to be %#v, got %#v", action, terr.Action)
			}

			if terr.TargetTypeName != pb.TargetType_name[int32(pb.TargetType_ANY)] {
				t.Errorf(
					"Expected error target type name to be %s, got %s",
					pb.TargetType_name[int32(pb.TargetType_ANY)],
					terr.TargetTypeName,
				)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected an UnsupportedTargetType error")
		}

		select {
		case err := <-errorChan:
			t.Errorf("Expected only 1 error, got more: %s", err)
		case <-time.After(10 * time.Millisecond):
		}
	})

	t.Run("Execute forward the C2 client errors to the errorChan", func(t *testing.T) {
		targets := []models.Target{
			models.Target{Type: pb.TargetType_CLIENT, Expr: "client1"},
			models.Target{Type: pb.TargetType_TOPIC, Expr: "topic1"},
		}

		action := &keyRotationAction{
			targets:   targets,
			c2Client:  mockC2Client,
			errorChan: errorChan,
			logger:    log.NewNopLogger(),
		}

		client1Err := errors.New("client1 error")
		topic1Err := errors.New("topic1 error")

		gomock.InOrder(
			mockC2Client.EXPECT().NewClientKey(gomock.Any(), "client1").Return(client1Err),
			mockC2Client.EXPECT().NewTopicKey(gomock.Any(), "topic1").Return(topic1Err),
		)

		go action.Execute(context.Background())

		select {
		case err := <-errorChan:
			if err != client1Err {
				t.Errorf("Expected error to be %s, got %s", client1Err, err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected %s error", client1Err)
		}

		select {
		case err := <-errorChan:
			if err != topic1Err {
				t.Errorf("Expected error to be %s, got %s", topic1Err, err)
			}
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Expected %s error", topic1Err)
		}

		select {
		case err := <-errorChan:
			t.Errorf("Expected no error, got %s", err)
		case <-time.After(10 * time.Millisecond):
		}
	})
}

func TestActionFactory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockC2Client := services.NewMockC2(mockCtrl)

	errorChan := make(chan error)

	factory := NewActionFactory(mockC2Client, errorChan, log.NewNopLogger())
	t.Run("Create keyRotationAction returns expected struct", func(t *testing.T) {
		rule := models.Rule{
			ActionType: pb.ActionType_KEY_ROTATION,
			Targets: []models.Target{
				models.Target{ID: 1},
				models.Target{ID: 2},
			},
		}

		action, err := factory.Create(rule)
		if err != nil {
			t.Errorf("Expected create to not return error, got %s", err)
		}

		typedAction, ok := action.(*keyRotationAction)
		if !ok {
			t.Errorf("Expected the action to be a *keyRotationAction, got %T", action)
		}

		if reflect.DeepEqual(typedAction.targets, rule.Targets) == false {
			t.Errorf("Expected action targets to be %#v, got %#v", rule.Targets, typedAction.targets)
		}

		if reflect.DeepEqual(typedAction.c2Client, mockC2Client) == false {
			t.Errorf("Expected C2 client to be %p, got %p", mockC2Client, typedAction.c2Client)
		}

		if typedAction.errorChan != errorChan {
			t.Errorf("Expected action errorChan to be %p, got %p", errorChan, typedAction.errorChan)
		}
	})

	t.Run("Create returns error on unsupported action type", func(t *testing.T) {
		rule := models.Rule{
			ActionType: pb.ActionType_UNDEFINED_ACTION,
		}

		action, err := factory.Create(rule)
		if err == nil {
			t.Errorf("Expected an error when creating an unsupported type of action")
		}

		if action != nil {
			t.Errorf("Expected action to be nil, got %#v", action)
		}
	})
}
