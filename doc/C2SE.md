# C2 scripting

C2 scripts allow to define key management policies that are automatically enforced. 
Currently, you can define key rotation policies per client or per topic, to ensure that a key is updated every `N` hour, where `N` is the period defined in the script.

Future rules will include the capability to:

* Add or remove right of a client to a certain topic at a future time
* Register new clients to C2 automatically at a given time
* Register new topics to C2 automatically at a given time

## Script format and database

A script is a sequence of lines of one the two following format:

```
C <client> N
T <topic> N
```

where

* `C` lines define a rule to `client`'s key every `N` hours
* `T` lines define a rule to update `topic`'s key every `N` hours

The topic is an UTF-8 string, the client is a client identifier alias (that is, a string). 
The period `N` is a positive integer, at most 9000.
If `N` is set to zero, then any rule in the database with the given topic or client id is deleted from the database.
Comments in a script file should start with a `#`. Blank lines are ignored.


The database includes two table with the following schemas:

For client keys:
```
ID | client id alias (unique) | key period (hours) | last update (in Unix seconds)
```
For topics keys:
```
ID | topic (unique) | key period (hours) | last update (in Unix seconds)
```
Here the `ID` fields serve as primary key, because string and byte arrays cannot be used as primary keys in some databases.

The client id is an alias, from which the actual id is later computed in the scripting engine.


## Script reader: c2sr

`c2ser` is the command-line utility that takes one or more e4s scripts as arguments, and for each script does the following:

1. Verifies the script its validity, discarding the file entirely if any non-empty line is not a valid line
1. Updates the database according to the rules in the script, such that:
    - If a rule has a period 0, then any database entry with the rule's client id or topic is deleted
    - If the database already includes a rule for the given client id or topic already, then this rule is *not* overwritten (and the new proposed rule is ignored)


## Script engine: c2se

`c2se` is the service that sends requests to `c2backend` corresponding to the rules in the database.
Every hour, for each database entry, `c2se` does the following:
If `current time - last update >= frequency`, then:

* For `C` rules: send a `SetClientKey` request 
* For `T` rules: send a `SetTopicKey` request to all clients registered to the given topic



