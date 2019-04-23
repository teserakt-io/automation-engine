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

The c2se engine is responsible of monitoring every existing rules, and trigger their actions when one of their trigger's condition is met.
Its started on the background of the c2se-api and spawns a goroutine for each rules, and another one for each rule's trigger.

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
