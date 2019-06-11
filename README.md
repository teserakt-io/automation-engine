# c2ae

[![pipeline status](https://gitlab.com/Teserakt/c2ae/badges/master/pipeline.svg)](https://gitlab.com/Teserakt/c2ae/commits/master)
[![coverage report](https://gitlab.com/Teserakt/c2ae/badges/master/coverage.svg)](https://gitlab.com/Teserakt/c2ae/commits/master)

## c2ae-api

The api is exposing the c2ae database over grpc, allowing to query the c2ae rules database.
It also start the c2ae engine, which will monitor existing rule triggers and process them if their execution conditions are met.

### Usage

```bash
cp configs/config.yaml.example configs/config.yaml
./bin/c2ae-api
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
