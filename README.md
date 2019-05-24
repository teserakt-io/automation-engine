# c2ae

[![pipeline status](https://gitlab.com/Teserakt/c2ae/badges/master/pipeline.svg)](https://gitlab.com/Teserakt/c2ae/commits/master)
[![coverage report](https://gitlab.com/Teserakt/c2ae/badges/master/coverage.svg)](https://gitlab.com/Teserakt/c2ae/commits/master)

## c2ae-api

The api is exposing the c2ae database over grpc, allowing to query the c2ae rules database.
It also start the c2ae engine, which will monitor existing rule triggers and process them if their execution conditions are met.

### Usage

```bash
./bin/c2ae-api -db /tmp/c2ae.db -addr 127.0.0.1:5556 -c2cert /path/to/c2/cert.pem
```

### c2ae engine

The c2ae engine is responsible of monitoring every existing rules, and trigger their actions when one of their trigger's condition is met.
It is started on the background of the c2ae-api and spawns a goroutine for each rules, and another one for each rule's trigger.

## c2ae-cli

The cli client allow to define new rules and list currently defined ones by interacting with the api.

### Usage

```bash
./bin/c2ae-cli --help
```

### Auto completion

Auto completion helper script can be sourced in current session or added to .bashrc with:

```bash
. <(./bin/c2ae-cli completion)
# Or for zsh (probably incomplete until https://github.com/spf13/cobra/pull/646 get merged)
. <(./bin/c2ae-cli completion --zsh)
```
It will provide auto completion for the various enums available

### Examples

#### Setting up key rotation every 2 minutes for some clients

```
### First create a new rule:
c2ae-cli create --action=KEY_ROTATION --description "Rotate client1 & client2 keys every 2 minutes"
# Rule #1 created!

### Now add targets:
c2ae-cli add-target --rule=1 --type=CLIENT --expr="client1"
# New target successfully added on rule #1
c2ae-cli add-target --rule=1 --type=CLIENT --expr="client2"
# New target successfully added on rule #1

### And finally set the trigger:
c2ae-cli add-trigger --rule=1 --type=TIME_INTERVAL --setting expr="*/2 * * * *"
# New trigger successfully added on rule #1

# And done ! Now the API will have auto loaded the newly created trigger and
# started a goroutine to make it execute at specified time interval.
```

## Development

Start api with:
```
go run cmd/api/c2ae-api.go -db /tmp/c2ae.db -c2cert /path/to/c2/cert.pem
```

Run cli with:
```
go run cmd/cli/c2ae-cli.go --help
```

A Makefile is also provided with various targets, like build, running tests, getting coverage, generating the mocks / protobuf...
Run ```make``` for the full list of targets and descriptions.
