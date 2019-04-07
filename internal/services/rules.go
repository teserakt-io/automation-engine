package services

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/teserakt/c2se/internal/models"
)

// RuleService defines methods to interact with rules models and database
type RuleService interface {
	All() ([]models.Rule, error)
	ByID(ruleID int) (*models.Rule, error)
	Save(rule *models.Rule) error
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

func (s *ruleService) All() ([]models.Rule, error) {
	var rules []models.Rule
	if result := s.gorm().Find(&rules); result.Error != nil {
		return nil, result.Error
	}

	return rules, nil
}

func (s *ruleService) Save(rule *models.Rule) error {
	if s.gorm().NewRecord(*rule) {
		if result := s.gorm().Create(rule); result.Error != nil {
			return result.Error
		}
	} else {

		if result := s.gorm().Save(rule); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (s *ruleService) ByID(ruleID int) (*models.Rule, error) {

	r := &models.Rule{}

	if result := s.gorm().First(&r, ruleID); result.Error != nil {
		return nil, result.Error
	}

	return r, nil
}

func (s *ruleService) gorm() *gorm.DB {
	return s.db.Connection()
}
