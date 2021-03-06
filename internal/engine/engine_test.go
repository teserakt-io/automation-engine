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

package engine

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"github.com/teserakt-io/automation-engine/internal/engine/watchers"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/services"
)

func TestAutomationEngine(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer func() {
		// Give some time to the goroutine to switch to running state
		// before letting mockCtrl check its expectations.
		time.Sleep(100 * time.Millisecond)

		mockCtrl.Finish()
	}()

	mockRuleService := services.NewMockRuleService(mockCtrl)
	mockRuleWatcherFactory := watchers.NewMockRuleWatcherFactory(mockCtrl)

	logger := log.New()
	logger.SetOutput(ioutil.Discard)

	engine := NewAutomationEngine(mockRuleService, mockRuleWatcherFactory, logger)

	rules := []models.Rule{
		models.Rule{ID: 1},
		models.Rule{ID: 2},
		models.Rule{ID: 3},
	}

	mockRuleWatcher1 := watchers.NewMockRuleWatcher(mockCtrl)
	mockRuleWatcher2 := watchers.NewMockRuleWatcher(mockCtrl)
	mockRuleWatcher3 := watchers.NewMockRuleWatcher(mockCtrl)

	t.Run("Start properly start a rule watcher for every rule", func(t *testing.T) {
		mockRuleService.EXPECT().All(gomock.Any()).Times(1).Return(rules, nil)

		mockRuleWatcherFactory.EXPECT().Create(rules[0]).Times(1).Return(mockRuleWatcher1)
		mockRuleWatcherFactory.EXPECT().Create(rules[1]).Times(1).Return(mockRuleWatcher2)
		mockRuleWatcherFactory.EXPECT().Create(rules[2]).Times(1).Return(mockRuleWatcher3)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockRuleWatcher1.EXPECT().Start(gomock.Any()).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})
		mockRuleWatcher2.EXPECT().Start(gomock.Any()).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})
		mockRuleWatcher3.EXPECT().Start(gomock.Any()).Times(1).DoAndReturn(func(ctx context.Context) {
			<-ctx.Done()
		})

		err := engine.Start(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Start returns error when it fail to fetch the rules", func(t *testing.T) {
		expectedError := errors.New("ruleService All() failed")
		mockRuleService.EXPECT().All(gomock.Any()).Times(1).Return(nil, expectedError)

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
		}()

		err := engine.Start(ctx)
		if err != expectedError {
			t.Errorf("Expected error to be %v, got %v", expectedError, err)
		}
	})
}
