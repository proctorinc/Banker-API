# Banker API
Banker is a banking API that does stuff with money. Very cool

## Setup
Prereqs
- Brew (mac only) https://mac.install.guide/homebrew/3
- Go https://go.dev/doc/install
- Air https://github.com/air-verse/air

## How to run
To start the server, run Air (May depend upon where you have air configured)
```sh
~/.air
```

## Development
To generate sqlc types after making changes to the db/schema.sql or schema.graphql, run:
```sh
./scripts/generate.sh
```
This command runs ```sqlc generate``` and ```gqlgen generate``` to compile sql and graphql schema files into go
