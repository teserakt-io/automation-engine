# C2 Automation Engine WebUI

## Goal

Providing a web interface to the automation engine, allowing to visualize and update its configuration.
This interface should integrate on the existing C2 web UI (as an additionnal `Automation Engine` tab in the menu)
Ideally, this new tab visibility should be toggleable from a configuration option so we could expose a 'Automation engine enabled' web UI, as well as others with just the C2 features.

## C2AE Overview

The C2AE goal is to allow definition of *rules*, defining an action to be triggered under given *conditions* on given *targets*

Thus we've defined 3 entities:
 - `Rule`: defines an action to be performed (based on the rule `type`)
 - `Trigger`: defines the execution conditions (based on the trigger `type` and `settings`)
 - `Target`: defines who/what will get the action performed on.

Let's see with the following rule as an example:
```
Perform a KEY_ROTATION action, every hour, for a TOPIC identified by "/test/topic"
```

The `Rule` action is KEY_ROTATION, with a single `Trigger` of type TIME_INTERVAL and settings defining an hourly expression (more on that later)
And a single `Target` of type TOPIC, where the topic name is "/test/topic"

Now a `Rule` can hold several `Triggers` and `Targets` as well:

```
Perform a KEY_ROTATION action, every 30 minutes, or when a client subscribe to the topic "/test/topic", for a TOPIC identified by "/test/topic" and any CLIENTS matching the regexp "weather-station-.*"
```

## WebUI user stories

As a webUI operator, I want to:
- List the existing rules as a table, with their triggers / targets count, and their lastExectued
- Inspect a single rule
    - See all rule targets type / expression
    - See all rule triggers type / settings
- Create a new rule
    - Define new target(s) at the same time
    - Define new trigger(s) at the same time
- Update a rule
    - Modify rule action / description
    - Add a new trigger
    - Add a new target
    - Modify a previously created trigger
    - Modify a previously created target
    - Remove a previously created trigger
    - Remove a previously created target
- Delete a rule
    - Confirm the deletion

## Api specifications

See [api.swagger.json](./api.swagger.json)

### Notes on the swagger file

- From the enums, the "UNDEFINED_*" values shouldn't appear on the UI.
- Taking apart the pbAction's UNDEFINED_ACTION, there is only one KEY_ROTATION action left, but consider more to be added in the futur
- When creating or updating a rule, pushing new pbTrigger or pbTarget objects in its triggers / targets properties will create them linked to the current pbRule. The id field of the new pbTrigger / pbTarget can be ommited in this case (or set to 0, but better to not set it at all).
- The pbTrigger settings field is expected to be a json object, transmitted as a base64 string to the API. Its content is fixed by the pbTriggerType, see `pbTriggerType settings` section below for details
- The target expression field is expecting a valid regular expression.

### pbTriggerType settings

Each trigger type may defines it's own settings format.
Settings are represented as a json object, and must be transmitted as a base64 encoded string to the API.

#### pbTriggerType.TIME_INTERVAL

For the TIME_INTERVAL trigger type, settings just hold an `expr` string, representing a valid cron expression (as defined here: https://en.wikipedia.org/wiki/Cron#CRON_expression)

```
{"expr": "* * * * *"}
```

#### pbTriggerType.CLIENT_SUBSCRIBED

Still WIP and to be defined. It can be presented as an empty textarea / empty json object for now.

#### pbTriggerType.CLIENT_UNSUBSCRIBED

Still WIP and to be defined. It can be presented as an empty textarea / empty json object for now.
