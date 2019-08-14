
# Rules

## Fields

- **ID**: an unique identifier for the rule, auto generated on creation.
- **Description**: short text explaining the role of this rule.
- **ActionType**: identifier of what will get done when the rule get executed. See below for available values.
- **LastExecuted**: hold the timestamp when the rule action was last executed. When the rule is created, it is set to the default value `0001-01-01 00:00:00 +0000 UTC`
- **Triggers**: a set of triggers attached to this rule
- **Targets**: a set of targets attached to this rule

## Available action types

| **Action type** | **Description** |
| --- | --- |
| KEY_ROTATION | Send a key renewal request for every targets to the C2 server |
