package events

import (
	"fmt"
	"log"

	"github.com/gorhill/cronexpr"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
	"gitlab.com/teserakt/c2se/internal/services"
)

// TriggerListenerFactory defines a way to create dynamic Listeners
type TriggerListenerFactory interface {
	Create(trigger models.Trigger) (Type, Listener, error)
}

type triggerListenerFactory struct {
	ruleService services.RuleService
}

var _ TriggerListenerFactory = &triggerListenerFactory{}

// NewTriggerListenerFactory creates a new listener factory for given trigger
func NewTriggerListenerFactory(ruleService services.RuleService) TriggerListenerFactory {
	return &triggerListenerFactory{
		ruleService: ruleService,
	}
}

func (l *triggerListenerFactory) Create(trigger models.Trigger) (Type, Listener, error) {
	switch trigger.TriggerType {
	case pb.TriggerType_TIME_INTERVAL:
		listener := &SchedulerListener{
			trigger:     &trigger,
			ruleService: l.ruleService,
		}

		return SchedulerTickType, listener, nil
	case pb.TriggerType_CLIENT_SUBSCRIBED:
		listener := &ClientSubscribedListener{
			trigger:     &trigger,
			ruleService: l.ruleService,
		}

		return ClientSubscribedType, listener, nil
	case pb.TriggerType_CLIENT_UNSUBSCRIBED:
		listener := &ClientUnsubscribedListener{
			trigger:     &trigger,
			ruleService: l.ruleService,
		}

		return ClientUnsubscribedType, listener, nil
	default:
		return UnknowEventType, nil, fmt.Errorf("no event callback exists for trigger type %s", trigger.TriggerType)
	}
}

// SchedulerListener implements a listener for Scheduler events
type SchedulerListener struct {
	trigger     *models.Trigger
	ruleService services.RuleService
}

var _ Listener = &SchedulerListener{}

// OnEvent implements Listener interface
func (l *SchedulerListener) OnEvent(evt Event) error {
	settings := &pb.TriggerSettingsTimeInterval{}
	if err := settings.Decode(l.trigger.Settings); err != nil {
		return err
	}

	eventValue, ok := evt.Value().(SchedulerEventValue)
	if !ok {
		return fmt.Errorf(
			"Wrong event data received on SchedulerListener: expected SchedulerEventValue, got %T",
			evt.Value(),
		)
	}

	expr, err := cronexpr.Parse(settings.Expr)
	if err != nil {
		return err
	}

	rule := l.trigger.Rule
	nextTime := expr.Next(rule.LastExecuted)
	if nextTime.Before(eventValue.Time) {
		log.Printf("Trigger %d-%d triggered !", rule.ID, l.trigger.ID)
		rule.LastExecuted = eventValue.Time
		l.ruleService.Save(l.trigger.Rule)

		// TODO: do rule.Action on rule.Targets :)
	}

	return nil
}

// ClientSubscribedListener implements a listener for ClientSubscribed events
type ClientSubscribedListener struct {
	trigger     *models.Trigger
	ruleService services.RuleService
}

var _ Listener = &ClientSubscribedListener{}

// OnEvent implements Listener interface
func (l *ClientSubscribedListener) OnEvent(evt Event) error {
	log.Printf("onClientSubscribed with event %#v\n", evt)

	return nil
}

// ClientUnsubscribedListener implements a listener for ClientUnsubscribed events
type ClientUnsubscribedListener struct {
	trigger     *models.Trigger
	ruleService services.RuleService
}

// OnEvent implements Listener interface
func (l *ClientUnsubscribedListener) OnEvent(evt Event) error {
	log.Printf("onClientUnsubscribed with event %#v\n", evt)

	return nil
}
