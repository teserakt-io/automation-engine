package models

//go:generate mockgen -destination converters_mocks.go -package=models -self_package github.com/teserakt-io/automation-engine/internal/models github.com/teserakt-io/automation-engine/internal/models Converter

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

// Converter interface defines methods to switch between protobuf and models types.
type Converter interface {
	RuleToPb(Rule) (*pb.Rule, error)
	RulesToPb([]Rule) ([]*pb.Rule, error)

	TargetToPb(Target) (*pb.Target, error)
	TargetsToPb([]Target) ([]*pb.Target, error)

	TriggerToPb(Trigger) (*pb.Trigger, error)
	TriggersToPb([]Trigger) ([]*pb.Trigger, error)

	PbToRule(*pb.Rule) (Rule, error)
	PbToRules([]*pb.Rule) ([]Rule, error)

	PbToTarget(*pb.Target) (Target, error)
	PbToTargets([]*pb.Target) ([]Target, error)

	PbToTrigger(*pb.Trigger) (Trigger, error)
	PbToTriggers([]*pb.Trigger) ([]Trigger, error)
}

type converter struct{}

var _ Converter = &converter{}

// NewConverter creates a new Converter
func NewConverter() Converter {
	return &converter{}
}

// RuleToPb converts a models.Rule to a pb.Rule
func (c *converter) RuleToPb(rule Rule) (*pb.Rule, error) {
	lastExecuted, err := ptypes.TimestampProto(rule.LastExecuted)
	if err != nil {
		return nil, err
	}

	targets, err := c.TargetsToPb(rule.Targets)
	if err != nil {
		return nil, err
	}

	triggers, err := c.TriggersToPb(rule.Triggers)
	if err != nil {
		return nil, err
	}

	return &pb.Rule{
		Id:           int32(rule.ID),
		Action:       rule.ActionType,
		Description:  rule.Description,
		Targets:      targets,
		Triggers:     triggers,
		LastExecuted: lastExecuted,
	}, nil
}

// RulesToPb converts a []models.Rule to a []pb.Rule
func (c *converter) RulesToPb(rules []Rule) ([]*pb.Rule, error) {
	var out []*pb.Rule
	for _, r := range rules {
		cr, err := c.RuleToPb(r)
		if err != nil {
			return nil, err
		}

		out = append(out, cr)
	}

	return out, nil
}

// TargetToPb converts a models.Target to a pb.Target
func (c *converter) TargetToPb(target Target) (*pb.Target, error) {
	return &pb.Target{
		Id:   int32(target.ID),
		Expr: target.Expr,
		Type: target.Type,
	}, nil
}

// TargetsToPb converts a []models.Target to a []pb.Target
func (c *converter) TargetsToPb(targets []Target) ([]*pb.Target, error) {
	var out []*pb.Target
	for _, t := range targets {
		tc, err := c.TargetToPb(t)
		if err != nil {
			return nil, err
		}
		out = append(out, tc)
	}

	return out, nil
}

// TriggerToPb converts a models.Trigger to a pb.Trigger
func (c *converter) TriggerToPb(trigger Trigger) (*pb.Trigger, error) {
	return &pb.Trigger{
		Id:       int32(trigger.ID),
		Type:     trigger.TriggerType,
		Settings: trigger.Settings,
	}, nil
}

// TriggersToPb converts a []models.Trigger to a []pb.Trigger
func (c *converter) TriggersToPb(triggers []Trigger) ([]*pb.Trigger, error) {
	var out []*pb.Trigger
	for _, t := range triggers {
		tc, err := c.TriggerToPb(t)
		if err != nil {
			return nil, err
		}
		out = append(out, tc)
	}

	return out, nil
}

// PbToRule converts a pb.Rule to a models.Rule
func (c *converter) PbToRule(rule *pb.Rule) (Rule, error) {
	lastExecuted, err := ptypes.Timestamp(rule.LastExecuted)
	if err != nil {
		return Rule{}, err
	}

	targets, err := c.PbToTargets(rule.Targets)
	if err != nil {
		return Rule{}, err
	}

	triggers, err := c.PbToTriggers(rule.Triggers)
	if err != nil {
		return Rule{}, err
	}

	return Rule{
		ID:           int(rule.Id),
		ActionType:   rule.Action,
		Description:  rule.Description,
		LastExecuted: lastExecuted,
		Targets:      targets,
		Triggers:     triggers,
	}, nil
}

func (c *converter) PbToRules(rules []*pb.Rule) ([]Rule, error) {
	var out []Rule
	for _, t := range rules {
		tc, err := c.PbToRule(t)
		if err != nil {
			return nil, err
		}

		out = append(out, tc)
	}

	return out, nil
}

// PbToTrigger converts a pb.Trigger to a models.Trigger
func (c *converter) PbToTrigger(trigger *pb.Trigger) (Trigger, error) {
	return Trigger{
		ID:          int(trigger.Id),
		TriggerType: trigger.Type,
		Settings:    trigger.Settings,
	}, nil
}

// PbToTriggers convers a []pb.Trigger to a []models.Trigger
func (c *converter) PbToTriggers(triggers []*pb.Trigger) ([]Trigger, error) {
	var out []Trigger
	for _, t := range triggers {
		tc, err := c.PbToTrigger(t)
		if err != nil {
			return nil, err
		}

		out = append(out, tc)
	}

	return out, nil
}

// PbToTarget converts a pb.Target to a models.Target
func (c *converter) PbToTarget(target *pb.Target) (Target, error) {
	return Target{
		ID:   int(target.Id),
		Type: target.Type,
		Expr: target.Expr,
	}, nil
}

// PbToTargets converts a []pb.Target to a []models.Target
func (c *converter) PbToTargets(targets []*pb.Target) ([]Target, error) {
	var out []Target
	for _, t := range targets {
		tc, err := c.PbToTarget(t)
		if err != nil {
			return nil, err
		}

		out = append(out, tc)
	}

	return out, nil
}
