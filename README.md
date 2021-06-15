## Shortu

Shortu is a simple url shortener written in Go :)

## Local Development Prerequisite

- docker, docker-compose

## Usage

1. Create your own secrets file in conf, name it as "<env>.secrets.env", such as "dev.secrets.env". You can refer to `.sample`.


2. Run deploy script (if you want to run e2e, provide environment variable `RUN_E2E` with any non empty value)
```bash
$RUN_E2E=1 ./scripts/deploy.sh <env>
```


## DB Migration

This project uses tern as database schema migration tool.


## Mocks

This project uses `mockgen` to generate some mocks for UT.

## TODO

- Use connection tool for cache
- logging