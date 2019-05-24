package events

import (
	"time"

	"gitlab.com/teserakt/c2ae/internal/models"
)

// SchedulerEventValue holds values transmitted on SchedulerTickType events
type SchedulerEventValue struct {
	Time time.Time
}

// TriggerEvent holds values transmitted when a trigger trigger
type TriggerEvent struct {
	Trigger models.Trigger
	Time    time.Time
}
