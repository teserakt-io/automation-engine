# C2 Scripting Engine

## Script format and database

A script is a sequence of lines of one the two following format:

```
C <client id> N
T <topic> N
```

where

* `C` lines define a rule to `client id`'s key every `N` hours
* `T` lines define a rule to update `topic`'s key every `N` hours

The topic is an UTF-8 string, the client is a valid client identifier.

`N` must be a positive integer. If `N` is set to zero, then any rule in
the database with the given topic or client id is removed from the
database.

The database includes a table with the following schema:

```
ID | type (C or T) | topic or client id | frequency (hours) | last update (in Unix seconds)
```

Initially the database is empty.


## Script reader: c2ser

`c2ser` reads an e4s scripts, verifies its validity (discarding the file
if any non-empty line is not a valid line), and updates the
database accordingly (adding all rules to the database, removing rules
for which zero is given as an update frequency).


## Service: c2se

`c2se` is the service that sends requests to `c2backend` corresponding
to the rules in the database.
Every hour, for each database entry, `c2se` does the following:
If `current time - last update >= frequency`, then:

* For `C` rules: send a `SetClientKey` request
* For `T` rules: send a `SetTopicKey` request to all clients registered
  to the given topic



