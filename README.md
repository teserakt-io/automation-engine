# c2se

[![pipeline status](https://gitlab.com/Teserakt/c2se/badges/master/pipeline.svg)](https://gitlab.com/Teserakt/c2se/commits/master)
[![coverage report](https://gitlab.com/Teserakt/c2se/badges/master/coverage.svg)](https://gitlab.com/Teserakt/c2se/commits/master)

## c2se-api

The api is exposing the c2se database over grpc, allowing to query the c2se rules database.
It also start the c2se engine, which will monitor existing rule triggers and process them if their execution conditions are met.

### Usage

```bash
./bin/c2se-api -db /tmp/c2se.db -addr 127.0.0.1:5556
```

### c2se engine

The c2se engine is responsible of monitoring existing triggers on every rules, and trigger their actions when their condition is met.
Its core component is an event dispatcher, where every trigger get registered and wait for system events.

An internal scheduler is also started, dispatching _tick_ events every seconds.

Each defined triggers get registered on the dispatcher from their type, which will make them receive every system events. In case the trigger isn't able to process the event, it will be discarded to not lock the whole system.

Also, the dispatcher get notified on every rules modification through the API. It will then stop and reload every triggers from the database and restart their routines.

#### Handling new events

Here are the steps required to support new events in the c2se engine:

- Create a new events.Type (and update events.EventStrings)
- Register some listeners on the dispatcher for this new events.Type
- Finally, call dispatcher.Dispatch with this new events.Type, and all registered listeners will get notified.

## c2se-cli

The cli client allow to define new rules and list currently defined ones by interacting with the api.

### Usage

```bash
./bin/c2se-cli --help
```

### Auto completion

Auto completion helper script can be sourced in current session or added to .bashrc with:

```bash
. <(./bin/c2se-cli completion)
# Or for zsh (probably incomplete until https://github.com/spf13/cobra/pull/646 get merged)
. <(./bin/c2se-cli completion --zsh)
```
It will provide auto completion for the various enums available
