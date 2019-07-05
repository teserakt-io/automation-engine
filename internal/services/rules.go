package services

//go:generate mockgen -destination=rules_mocks.go -package=services -self_package gitlab.com/teserakt/c2ae/internal/services gitlab.com/teserakt/c2ae/internal/services RuleService

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"go.opencensus.io/trace"

	"gitlab.com/teserakt/c2ae/internal/models"
)

// TriggerReader defines methods to read triggers
type TriggerReader interface {
	TriggerByID(ctx context.Context, triggerID int) (models.Trigger, error)
}

// TriggerWriter defines methods to write triggers
type TriggerWriter interface {
	DeleteTriggers(ctx context.Context, triggers ...models.Trigger) error
}

// TargetReader defines methods to read targets
type TargetReader interface {
	TargetByID(ctx context.Context, targetID int) (models.Target, error)
}

// TargetWriter defines methods to write Targets
type TargetWriter interface {
	DeleteTargets(ctx context.Context, targets ...models.Target) error
}

// RuleReader defines methods availble to read rules from database
type RuleReader interface {
	All(ctx context.Context) ([]models.Rule, error)
	ByID(ctx context.Context, ruleID int) (models.Rule, error)
}

// RuleWriter defines methods available to write rules
type RuleWriter interface {
	Save(ctx context.Context, rule *models.Rule) error
	Delete(ctx context.Context, rule models.Rule) error
}

// RuleService defines methods to interact with rules models and database
type RuleService interface {
	RuleReader
	RuleWriter

	TargetReader
	TargetWriter

	TriggerReader
	TriggerWriter
}

type ruleService struct {
	db        models.Database
	validator models.Validator
}

var _ RuleService = &ruleService{}

// NewRuleService creates a new RuleService
func NewRuleService(db models.Database, validator models.Validator) RuleService {
	return &ruleService{
		db:        db,
		validator: validator,
	}
}

// All retrieves all rules from database
func (s *ruleService) All(ctx context.Context) ([]models.Rule, error) {
	ctx, span := trace.StartSpan(ctx, "RuleService.All")
	defer span.End()

	rules := []models.Rule{}
	if result := s.gorm().Find(&rules); result.Error != nil {
		return nil, result.Error
	}

	return rules, nil
}

// Save either creates or updates given rule in database
func (s *ruleService) Save(ctx context.Context, rule *models.Rule) error {
	ctx, span := trace.StartSpan(ctx, "RuleService.Save")
	defer span.End()

	if err := s.validator.ValidateRule(*rule); err != nil {
		return fmt.Errorf("rule validation failed: %v", err)
	}

	if result := s.gorm().Save(rule); result.Error != nil {
		return result.Error
	}

	return nil
}

// ByID retrieves a rule by its ID
func (s *ruleService) ByID(ctx context.Context, ruleID int) (models.Rule, error) {
	ctx, span := trace.StartSpan(ctx, "RuleService.ByID")
	defer span.End()

	r := models.Rule{}
	if result := s.gorm().First(&r, ruleID); result.Error != nil {
		return r, result.Error
	}

	return r, nil
}

// TriggerByID retrieves a trigger by its ID
func (s *ruleService) TriggerByID(ctx context.Context, triggerID int) (models.Trigger, error) {
	ctx, span := trace.StartSpan(ctx, "RuleService.TriggerByID")
	defer span.End()

	t := models.Trigger{}
	if result := s.gorm().First(&t, triggerID); result.Error != nil {
		return t, result.Error
	}

	return t, nil
}

// TargetByID retrieves a target by its ID
func (s *ruleService) TargetByID(ctx context.Context, targetID int) (models.Target, error) {
	ctx, span := trace.StartSpan(ctx, "RuleService.TargetByID")
	defer span.End()

	t := models.Target{}
	if result := s.gorm().First(&t, targetID); result.Error != nil {
		return t, result.Error
	}

	return t, nil
}

// Delete removes given rule and associated triggers / targets from database
func (s *ruleService) Delete(ctx context.Context, rule models.Rule) error {
	ctx, span := trace.StartSpan(ctx, "RuleService.Delete")
	defer span.End()

	if result := s.gorm().Delete(rule); result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteTriggers will delete all given triggers in a single batch
func (s *ruleService) DeleteTriggers(ctx context.Context, triggers ...models.Trigger) error {
	ctx, span := trace.StartSpan(ctx, "RuleService.DeleteTriggers")
	defer span.End()

	var triggerIds []int
	for _, trigger := range triggers {
		triggerIds = append(triggerIds, trigger.ID)
	}

	if len(triggerIds) > 0 {
		if result := s.gorm().Delete(models.Trigger{}, "id IN (?)", triggerIds); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// DeleteTargets will delete all given targets in a single batch
func (s *ruleService) DeleteTargets(ctx context.Context, targets ...models.Target) error {
	ctx, span := trace.StartSpan(ctx, "RuleService.DeleteTargets")
	defer span.End()

	var targetIds []int
	for _, target := range targets {
		targetIds = append(targetIds, target.ID)
	}

	if len(targetIds) > 0 {
		if result := s.gorm().Delete(models.Target{}, "id IN (?)", targetIds); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (s *ruleService) gorm() *gorm.DB {
	return s.db.Connection().Set("gorm:auto_preload", true)
}
