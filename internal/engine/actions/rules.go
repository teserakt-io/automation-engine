package actions

//go:generate mockgen -destination=rules_mocks.go -package actions -self_package gitlab.com/teserakt/c2se/internal/engine/actions gitlab.com/teserakt/c2se/internal/engine/actions ActionFactory,Action

import (
	"fmt"
	"log"

	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
	"gitlab.com/teserakt/c2se/internal/services"
	e4 "gitlab.com/teserakt/e4common"
)

// ActionFactory is responsible of Aciton creation
type ActionFactory interface {
	Create(models.Rule) (Action, error)
}

// Action describe rule's Action methods
type Action interface {
	Execute()
}

type actionFactory struct {
	c2Client  services.C2
	errorChan chan<- error
}

var _ ActionFactory = &actionFactory{}

// NewActionFactory creates a new ActionFactory
func NewActionFactory(c2Client services.C2, errorChan chan<- error) ActionFactory {
	return &actionFactory{
		c2Client:  c2Client,
		errorChan: errorChan,
	}
}

func (f *actionFactory) Create(rule models.Rule) (Action, error) {
	var action Action

	switch rule.ActionType {
	case pb.ActionType_KEY_ROTATION:
		action = &keyRotationAction{
			targets:   rule.Targets,
			c2Client:  f.c2Client,
			errorChan: f.errorChan,
		}
	default:
		return nil, fmt.Errorf("unknow action type %d", rule.ActionType)
	}

	return action, nil
}

// UnsupportedTargetType is an error returned when trying to execute
// an action which doesn't support the given target type.
type UnsupportedTargetType struct {
	Action         Action
	TargetTypeName string
}

func (e UnsupportedTargetType) Error() string {
	return fmt.Sprintf(
		"ERROR: unsupported action %T for target type %s",
		e.Action,
		e.TargetTypeName,
	)
}

type keyRotationAction struct {
	targets  []models.Target
	c2Client services.C2

	errorChan chan<- error
}

var _ Action = &keyRotationAction{}

func (a *keyRotationAction) Execute() {
	for _, target := range a.targets {
		log.Printf("Executing %T for target: %s", a, target.Expr)
		switch target.Type {
		case pb.TargetType_CLIENT:
			hashedID := e4.HashIDAlias(target.Expr)
			err := a.c2Client.NewClientKey(hashedID)
			if err != nil {
				a.errorChan <- err

				continue
			}
		case pb.TargetType_TOPIC:
			err := a.c2Client.NewTopicKey(target.Expr)
			if err != nil {
				a.errorChan <- err

				continue
			}
		default:
			a.errorChan <- UnsupportedTargetType{Action: a, TargetTypeName: pb.TargetType_name[int32(target.Type)]}

			continue
		}
	}
}
