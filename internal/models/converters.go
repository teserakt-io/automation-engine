package models

import (
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/teserakt/c2se/internal/pb"
)

// RuleToPb converts a models.Rule to a pb.Rule
func RuleToPb(rule Rule) (*pb.Rule, error) {
	lastExecuted, err := ptypes.TimestampProto(rule.LastExecuted)
	if err != nil {
		return nil, err
	}

	targets, err := TargetsToPb(rule.Targets)
	if err != nil {
		return nil, err
	}

	triggers, err := TriggersToPb(rule.Triggers)
	if err != nil {
		return nil, err
	}

	return &pb.Rule{
		Id:           int32(rule.ID),
		Action:       rule.ActionType,
		Description:  rule.Description,
		Targets:      targets,
		Triggers:     triggers,
		LastExectued: lastExecuted,
	}, nil
}

// RulesToPb converts a []models.Rule to a []pb.Rule
func RulesToPb(rules []Rule) ([]*pb.Rule, error) {
	var out []*pb.Rule
	for _, r := range rules {
		cr, err := RuleToPb(r)
		if err != nil {
			return nil, err
		}

		out = append(out, cr)
	}

	return out, nil
}

// TargetToPb converts a models.Target to a pb.Target
func TargetToPb(target Target) (*pb.Target, error) {
	return &pb.Target{
		Id:   int32(target.ID),
		Expr: target.Expr,
		Type: target.Type,
	}, nil
}

// TargetsToPb converts a []models.Target to a []pb.Target
func TargetsToPb(targets []Target) ([]*pb.Target, error) {
	var out []*pb.Target
	for _, t := range targets {
		tc, err := TargetToPb(t)
		if err != nil {
			return nil, err
		}
		out = append(out, tc)
	}

	return out, nil
}

// TriggerToPb converts a models.Trigger to a pb.Trigger
func TriggerToPb(trigger Trigger) (*pb.Trigger, error) {
	return &pb.Trigger{
		Id:       int32(trigger.ID),
		Type:     trigger.TriggerType,
		Settings: trigger.Settings,
		State:    trigger.State,
	}, nil
}

// TriggersToPb converts a []models.Trigger to a []pb.Trigger
func TriggersToPb(triggers []Trigger) ([]*pb.Trigger, error) {
	var out []*pb.Trigger
	for _, t := range triggers {
		tc, err := TriggerToPb(t)
		if err != nil {
			return nil, err
		}
		out = append(out, tc)
	}

	return out, nil
}

// PbToRule converts a pb.Rule to a models.Rule
func PbToRule(rule *pb.Rule) (Rule, error) {

	lastExecuted, err := ptypes.Timestamp(rule.LastExectued)
	if err != nil {
		return Rule{}, err
	}

	targets, err := PbToTargets(rule.Targets)
	if err != nil {
		return Rule{}, err
	}

	triggers, err := PbToTriggers(rule.Triggers)
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

// PbToTrigger converts a pb.Trigger to a models.Trigger
func PbToTrigger(trigger *pb.Trigger) (Trigger, error) {
	return Trigger{
		ID:          int(trigger.Id),
		TriggerType: trigger.Type,
		Settings:    trigger.Settings,
		State:       trigger.State,
	}, nil
}

// PbToTriggers convers a []pb.Trigger to a []models.Trigger
func PbToTriggers(triggers []*pb.Trigger) ([]Trigger, error) {
	var out []Trigger
	for _, t := range triggers {
		tc, err := PbToTrigger(t)
		if err != nil {
			return nil, err
		}

		out = append(out, tc)
	}

	return out, nil
}

// PbToTarget converts a pb.Target to a models.Target
func PbToTarget(target *pb.Target) (Target, error) {
	return Target{
		ID:   int(target.Id),
		Type: target.Type,
		Expr: target.Expr,
	}, nil
}

// PbToTargets converts a []pb.Target to a []models.Target
func PbToTargets(targets []*pb.Target) ([]Target, error) {
	var out []Target
	for _, t := range targets {
		tc, err := PbToTarget(t)
		if err != nil {
			return nil, err
		}

		out = append(out, tc)
	}

	return out, nil
}
