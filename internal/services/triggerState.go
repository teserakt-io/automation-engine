package services

//go:generate mockgen -destination=triggerState_mocks.go -package=services -self_package github.com/teserakt-io/automation-engine/internal/services github.com/teserakt-io/automation-engine/internal/services TriggerStateService

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
