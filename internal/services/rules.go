package services

//go:generate mockgen -destination=../mocks/services_rule.go -package=mocks gitlab.com/teserakt/c2se/internal/services RuleService

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/teserakt/c2se/internal/models"
)

// RuleService defines methods to interact with rules models and database
type RuleService interface {
	All() ([]models.Rule, error)
	ByID(ruleID int) (models.Rule, error)
	Save(rule *models.Rule) error
	Delete(rule models.Rule) error

	TriggerByID(triggerID int) (models.Trigger, error)
	TargetByID(targetID int) (models.Target, error)
}

type ruleService struct {
	db models.Database
}

var _ RuleService = &ruleService{}

// NewRuleService creates a new RuleService
func NewRuleService(db models.Database) RuleService {
	return &ruleService{
		db: db,
	}
}

// All retrieves all rules from database
func (s *ruleService) All() ([]models.Rule, error) {
	var rules []models.Rule
	if result := s.gorm().Find(&rules); result.Error != nil {
		return nil, result.Error
	}

	return rules, nil
}

// Save either creates or updates given rule in database
func (s *ruleService) Save(rule *models.Rule) error {
	if result := s.gorm().Save(rule); result.Error != nil {
		return result.Error
	}

	return nil
}

// ByID retrieves a rule by its ID
func (s *ruleService) ByID(ruleID int) (models.Rule, error) {
	r := models.Rule{}

	if result := s.gorm().First(&r, ruleID); result.Error != nil {
		return r, result.Error
	}

	return r, nil
}

// TriggerByID retrieves a trigger by its ID
func (s *ruleService) TriggerByID(triggerID int) (models.Trigger, error) {
	t := models.Trigger{}

	if result := s.gorm().First(&t, triggerID); result.Error != nil {
		return t, result.Error
	}

	return t, nil
}

// TargetByID retrieves a target by its ID
func (s *ruleService) TargetByID(targetID int) (models.Target, error) {
	t := models.Target{}

	if result := s.gorm().First(&t, targetID); result.Error != nil {
		return t, result.Error
	}

	return t, nil
}

// Delete removes given rule and associated triggers / targets from database
func (s *ruleService) Delete(rule models.Rule) error {
	if result := s.gorm().Delete(rule); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *ruleService) gorm() *gorm.DB {
	return s.db.Connection().Set("gorm:auto_preload", true)
}
