## Shortu

Shortu is a simple url shortener written in Go :)

### Local Development Prerequisite

- docker, docker-compose

### Usage

1. Clone this repo

2. Create your own secrets file in conf, name it as "<env>.secrets.env", such as "dev.secrets.env". You can refer to `.sample`.

3. Run deploy script
```bash
./scripts/deploy.sh <env>
```

### Run tests (Including E2E test)

1. (Optional, required if you need to run e2e) Make sure a fresh instance of service is running (Follow [Usage](#usage) to run it)

2. Run tests with (ignore `[env=e2e]` if not to run e2e)

```bash
$ [env=e2e] go test ./...
```

### DB Migration

This project uses tern as database schema migration tool.


### Mocks

This project uses `mockgen` to generate some mocks for UT.

### TODOS

- better logging
- support k8s deployment
