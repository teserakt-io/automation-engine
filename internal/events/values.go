package events

import "time"

// SchedulerEventValue holds values transmitted on SchedulerTickType events
type SchedulerEventValue struct {
	Time time.Time
}
