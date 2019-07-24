package models

//go:generate mockgen -destination validators_mocks.go -package=models -self_package gitlab.com/teserakt/c2ae/internal/models gitlab.com/teserakt/c2ae/internal/models Validator

import (
	"errors"
	"fmt"
	"regexp"

	"gitlab.com/teserakt/c2ae/internal/pb"
)

// Validation errors
var (
	ErrUndefinedAction        = errors.New("rule action is undefined")
	ErrUnknowActionType       = errors.New("rule action type is unknown")
	ErrUndefinedTriggerType   = errors.New("trigger type is undefined")
	ErrUnsupportedTriggerType = errors.New("trigger type is not supported")
	ErrTargetExprRequired     = errors.New("target expr is required")
)

// Validator defines an interface for models validation
type Validator interface {
	ValidateRule(rule Rule) error
	ValidateTrigger(trigger Trigger) error
	ValidateTarget(target Target) error
}

type validator struct {
}

var _ Validator = &validator{}

// NewValidator creates a new models Validator
func NewValidator() Validator {
	return &validator{}
}

// ValidateRule will check if given rule is valid, and returns an error when not.
func (v *validator) ValidateRule(rule Rule) error {
	if rule.ActionType == pb.ActionType_UNDEFINED_ACTION {
		return ErrUndefinedAction
	}

	if _, ok := pb.ActionType_name[int32(rule.ActionType)]; !ok {
		return ErrUnknowActionType
	}

	for _, trigger := range rule.Triggers {
		if err := v.ValidateTrigger(trigger); err != nil {
			return fmt.Errorf("trigger validation failed: %v", err)
		}
	}

	for _, target := range rule.Targets {
		if err := v.ValidateTarget(target); err != nil {
			return fmt.Errorf("target validation failed: %v", err)
		}
	}

	return nil
}

// ValidateTrigger will check if given trigger is valid, and returns an error when not.
func (v *validator) ValidateTrigger(trigger Trigger) error {
	if trigger.TriggerType == pb.TriggerType_UNDEFINED_TRIGGER {
		return ErrUndefinedTriggerType
	}

	settings, err := pb.Decode(trigger.TriggerType, trigger.Settings)
	if err != nil {
		return err
	}

	if err := settings.Validate(); err != nil {
		return err
	}

	return nil
}

// ValidateTarget will check if given target is valid, and returns an error when not.
func (v *validator) ValidateTarget(target Target) error {
	if len(target.Expr) == 0 {
		return ErrTargetExprRequired
	}

	if _, err := regexp.Compile(target.Expr); err != nil {
		return fmt.Errorf("target expr regexp is invalid: %v", err)
	}

	return nil
}