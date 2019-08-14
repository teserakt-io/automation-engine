# c2ae

[![pipeline status](https://gitlab.com/Teserakt/c2ae/badges/master/pipeline.svg)](https://gitlab.com/Teserakt/c2ae/commits/master)
[![coverage report](https://gitlab.com/Teserakt/c2ae/badges/master/coverage.svg)](https://gitlab.com/Teserakt/c2ae/commits/master)


## Introduction

The automation engine aim to ease and automate the key management of E4 by providing a way to define policies (or rules) for key renewal, or any other operations, which need to be performed under certain events or conditions, to keep the system communications secure.

It defines 3 main components:

- rules

A rule is a simple container entity. It holds an action type, a description, a list of targets, a list of triggers and a timestamp. It defines what has to be done when it get executed (from its action type), and when it was last executed (from its timestamp).
The list of available actions is defined in the [proto file](./api.proto) (see `ActionType`). A current available action is, for example, a key rotation (`ActionType.KEY_ROTATION`).

For more details, see [the rules documentation](./doc/rules.md)

- targets

A target define who/what the rule action will be executed for. It has a type (see available types in [proto file](./api.proto) > `TargetType`) and an expression (the identifier of the target). When a rule is triggered, it will execute its action for each of its targets.
For example, we can define a rule with the action `KEY_ROTATION`, and several targets, a `TOPIC` target type, with expression `/devices/groupA`, and another `CLIENT` target type, with expression `secure-thing-XYZ`. This means every time the rule get executed, the topic identified by `/devices/groupA` and the client identified by `secure-thing-XYZ` will have their key renewed.

A generic target can also be defined, to allow matching only by it's identifier, using the `ANY` type.

For more details, see [the targets documentation](./doc/targets.md)

- triggers

A trigger defines the condition to decide if the rule action must be executed. It holds a type, a settings map (content being type dependant), and an internal state map.
The list of available trigger types is defined in the [proto file](./api.proto) (see `TriggerType`) and their respective settings definition is availabe [here](./internal/pb/triggerSettings.go).
For example, a trigger can be of type `TIME_INTERVAL`, meaning it require an `Expr` setting to be defined to a cron expression. This trigger will then monitor the rule *last executed* timestamp against the cron expression, and notify the rule to execute when its due to.

For more details, see [the triggers documentation](./doc/triggers.md)

## Automation Engine API

The c2ae-api is exposing http and grpc endpoints, allowing to create, read, update or delete rules.
It also start the c2ae internal engine, which will monitor the existing triggers and launch their rule action if conditions are met.

### Usage

Generate a certificate if needed, and start the binary:

```bash
# Init config
cp configs/config.yaml.example configs/config.yaml
# Generate TLS certificate
openssl req -nodes -newkey rsa:2048 -keyout configs/c2ae-key.pem -x509 -sha256 -days 365 -out configs/c2ae-cert.pem -subj "/CN=localhost" -addext "subjectAltName = 'IP:127.0.0.1'"

# Retrieve c2 certificate
cp /path/to/c2/configs/c2-cert.pem configs/c2-cert.pem

# Run api server
./bin/c2ae-api
```

### Automation Engine

The c2ae engine is responsible of monitoring every existing rules, and trigger their actions when one of the rule's trigger condition is met.
It is started on the background of the API server, and spawns a goroutine for each rules, and another one for each rule's trigger.

On startup, the engine will also subscribe to an event stream over GRPC on the C2 server (`SubscribeToEventStream`). This connection will be kept open at all time to allow reception of C2 events. If the connection is lost, the engine will automatically retry to reconnect every seconds and will log an error until it succeed.

## Automation Engine CLI

The cli client allow to define new rules and list currently defined ones by interacting with the api.

### Usage

```bash
./bin/c2ae-cli --help
```

It require a C2AE-API running and can be specified where to connect to using the `--endpoint` and `--cert` global flags.

example (those are also default values):
```
./bin/c2ae-cli --endpoint 127.0.0.1:5556 --cert configs/c2ae-cert.pem list
```

### CLI client auto completion

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

### Run from docker image

The CI automatically push docker images of C2AE API and CLI after each successfull builds and for each branches.

List of available C2 images: https://gitlab.com/Teserakt/c2ae/container_registry

#### API

The c2ae-api server can be started like so:
```
# Replace <BRANCH_NAME> with the actual branch you want to pull the image from, like master, or devel, or tag...
docker run -it --name c2ae-api --rm -v $(pwd)/configs:/opt/e4/configs -e C2AE_LISTEN_ADDR=0.0.0.0:5556 -p 5556:5556 registry.gitlab.com/teserakt/c2ae/api:<BRANCH_NAME>
```

It just require a volume to the configs folder (Depending on your configuration, you may also need to get another volumes for the certificate and keys if they're not in the configs folder) and the ports for the GRPC api (which can be removed if not used)

See `internal/config/config.go` `ViperCfgFields()` for the full list of available environment variables.

#### CLI

```
# Replace <BRANCH_NAME> with the actual branch you want to pull the image from, like master, or devel, or tag...
# Replace <COMMAND> with the actual command to execute
docker run -it  --rm --link c2ae-api -e C2AE_API_ENDPOINT="c2ae-api:5556" registry.gitlab.com/teserakt/c2ae/cli:<BRANCH_NAME> <COMMAND>
```

## Development

A Makefile is provided with various targets, like build, running tests, getting coverage, generating the mocks / protobuf...
Run ```make``` for the full list of targets and descriptions.
