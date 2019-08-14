
# Triggers

## Fields

- **ID**: an unique identifier for the trigger, auto generated on creation.
- **RuleID**: the identifier of the rule this target is belonging to.
- **Type**: the type of the trigger. See below for available values.
- **Settings**: holds the trigger configuration data. See below for possible values.
- **State**: the trigger internal state. See below for details.

## Available trigger types

| **Trigger type** | **Description** |
| --- | --- |
| TIME_INTERVAL | Makes this trigger watching a cron expression, and compare it to the lastExecuted field of the rule. When the cron expression is due, the rule action is executed |
| EVENT | Makes the trigger listen for C2 events. It executes the rule action when a configured amount of matching events from the C2 server has been received. A event is *matching* if its type correspond to the configured EventType, and at least one of the rule targets is matching the event source or target fields. *CLIENT* targets  will be checked against event Source field, *TOPIC* targets against event Target field, and *ANY* on both. |

## TIME_INTERVAL Trigger

### Settings

| **Field** | **Type** | **Description** | **Example** |
| --- | --- | --- | --- |
| Expr | string | A valid cron expression as defined by [Wikipedia](https://en.wikipedia.org/wiki/Cron#CRON_expression) defining the interval of expected execution | */5 * * * 1-5 *# every 5 minutes, from monday to friday* |

### State

This trigger type doesn't persist any state.

### EVENT Trigger

### Settings

| **Field** | **Type** | **Description** | **Example** |
| --- | --- | --- | --- |
| EventType | string | A C2 event type (one of CLIENT_SUBSCRIBED or CLIENT_UNSUBSCRIBED, see C2 api.proto `EventType` definition for complete list) | CLIENT_SUBSCRIBED |
| MaxOccurence | int | A positive number of matching events to be received before the rule action get executed. Those event must match both the EventType and at least one of the rule defined targets, | 5 |

### State

This trigger will old a counter in its *State* field, which get incremented upon receiving events matching its settings. This internal counter is persisted in database every time it changes, and is compared with the MaxOccurence setting on each events. When it match or exceed the MaxOccurence value, the rule action get triggered, and the counter reset to 0.

> A rule modification does not reset the internal counter. So if the actual counter hold, let's say, the value 5, and the MaxOccurence setting is modified from 10 to 3, the rule will trigger as soon as the next matching event is received as `Counter(6) >= MaxOccurence(3)`. The counter is then reset to 0, and it will need 3 more matching events to trigger again.
