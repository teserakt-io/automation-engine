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

package services

//go:generate mockgen -copyright_file ../../doc/COPYRIGHT_TEMPLATE.txt -destination=triggerState_mocks.go -package=services -self_package github.com/teserakt-io/automation-engine/internal/services github.com/teserakt-io/automation-engine/internal/services TriggerStateService

import (
	"context"

	"github.com/jinzhu/gorm"
	"go.opencensus.io/trace"

	"github.com/teserakt-io/automation-engine/internal/models"
)

// TriggerStateService defines a service for managing triggerState models.
type TriggerStateService interface {
	Save(context.Context, *models.TriggerState) error
	ByTriggerID(context.Context, int) (models.TriggerState, error)
}

type triggerStateService struct {
	db models.Database
}

var _ TriggerStateService = (*triggerStateService)(nil)

// NewTriggerStateService creates a new service for handling triggerState models
func NewTriggerStateService(db models.Database) TriggerStateService {
	return &triggerStateService{
		db: db,
	}
}

func (s *triggerStateService) Save(ctx context.Context, state *models.TriggerState) error {
	_, span := trace.StartSpan(ctx, "TriggerStateService.Save")
	defer span.End()

	if result := s.gorm().Save(state); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *triggerStateService) ByTriggerID(ctx context.Context, triggerID int) (models.TriggerState, error) {
	_, span := trace.StartSpan(ctx, "TriggerStateService.ByTriggerID")
	defer span.End()

	t := models.TriggerState{
		TriggerID: triggerID,
	}
	if result := s.gorm().First(&t, "trigger_id = ?", triggerID); result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound { // Ignore record not found errors, as we returns a fresh triggerState in any cases.
			return t, result.Error
		}
	}

	return t, nil
}

func (s *triggerStateService) gorm() *gorm.DB {
	return s.db.Connection().Set("gorm:auto_preload", true)
}
