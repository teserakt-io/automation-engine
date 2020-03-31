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

package actions

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
)

func TestKeyRotationAction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockC2Client := services.NewMockC2(mockCtrl)

	errorChan := make(chan error)

	logger := log.New()
	logger.SetOutput(ioutil.Discard)

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
			logger:    logger,
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
}

func TestActionFactory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockC2Client := services.NewMockC2(mockCtrl)

	errorChan := make(chan error)

	logger := log.New()
	logger.SetOutput(ioutil.Discard)

	factory := NewActionFactory(mockC2Client, errorChan, logger)
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
