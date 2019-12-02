package actions

//go:generate mockgen -destination=rules_mocks.go -package actions -self_package github.com/teserakt-io/automation-engine/internal/engine/actions github.com/teserakt-io/automation-engine/internal/engine/actions ActionFactory,Action

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"

	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// ActionFactory is responsible of Aciton creation
type ActionFactory interface {
	Create(models.Rule) (Action, error)
}

// Action describe rule's Action methods
type Action interface {
	Execute(context.Context)
}

type actionFactory struct {
	c2Client  services.C2
	errorChan chan<- error
	logger    log.FieldLogger
}

var _ ActionFactory = &actionFactory{}

// NewActionFactory creates a new ActionFactory
func NewActionFactory(c2Client services.C2, errorChan chan<- error, logger log.FieldLogger) ActionFactory {
	return &actionFactory{
		c2Client:  c2Client,
		errorChan: errorChan,
		logger:    logger,
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
			logger:    f.logger,
		}
	default:
		return nil, fmt.Errorf("unknown action type %d", rule.ActionType)
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
	logger   log.FieldLogger

	errorChan chan<- error
}

var _ Action = &keyRotationAction{}

func (a *keyRotationAction) Execute(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "KeyRotationAction.Execute")
	defer span.End()

	for _, target := range a.targets {
		a.logger.WithFields(log.Fields{
			"action":     "keyRotation",
			"target":     target.Expr,
			"targetType": pb.TargetType_name[int32(target.Type)],
		}).Info("executing action")

		switch target.Type {
		case pb.TargetType_CLIENT:
			// TODO: for now we expect target to be defined with exact names of client.
			// But we may want later to allow some wildcards to target multiple clients at once
			// like weather-station-*, which should match weather-station-east, weather-station-west...
			// But it might not be possible to fetch all existing client names (huge number)
			// Maybe we could just send the wildcarded expression to the C2 and let it deal with it and
			// match the clients directly from a DB query.
			err := a.c2Client.NewClientKey(ctx, target.Expr)
			if err != nil {
				a.errorChan <- err

				continue
			}
		case pb.TargetType_TOPIC:
			err := a.c2Client.NewTopicKey(ctx, target.Expr)
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
