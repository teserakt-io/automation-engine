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

	log "github.com/sirupsen/logrus"

	"github.com/teserakt-io/automation-engine/internal/engine/watchers"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// AutomationEngine interface describe the public methods available on the automation engine
type AutomationEngine interface {
	Start(context.Context) error
}

type automationEngine struct {
	ruleService        services.RuleService
	ruleWatcherFactory watchers.RuleWatcherFactory
	logger             log.FieldLogger
}

var _ AutomationEngine = &automationEngine{}

// NewAutomationEngine creates a new automation engine
func NewAutomationEngine(
	ruleService services.RuleService,
	ruleWatcherFactory watchers.RuleWatcherFactory,
	logger log.FieldLogger,
) AutomationEngine {
	return &automationEngine{
		ruleService:        ruleService,
		ruleWatcherFactory: ruleWatcherFactory,
		logger:             logger,
	}
}

func (e *automationEngine) Start(ctx context.Context) error {
	rules, err := e.ruleService.All(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleWatcher := e.ruleWatcherFactory.Create(rule)
		e.logger.WithField("rule", rule.ID).Info("started ruleWatcher")
		go ruleWatcher.Start(ctx)
	}

	return nil
}
