
# Targets

## Fields

- **ID**: an unique identifier for the target, auto generated on creation.
- **RuleID**: the identifier of the rule this target is belonging to.
- **Type**: the type of the target. See below for available values.
- **Expr**: hold an identifier of the target, usually its name.

## Available target types

| **Target type** | **Description** |
| --- | --- |
| TOPIC | Identify target as a C2 topic, making the target Expr field match a topic name on the C2 server |
| CLIENT | Identify target as a C2 client, making the target Expr field match a client name on the C2 server |
| ANY | Wildcard target identifier, making the target Expr field match either a topic or a client name on the C2 server |
