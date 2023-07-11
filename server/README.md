# Server

## Requirements

- Go 1.20
- Docker/Docker compose
- [Golang migrate](https://github.com/golang-migrate/migrate)

## ENV

- `env.bash` has a collection of functions needed for development

```sh
source env.bash
```

### Generating New DBML file

```sh
dbml
```

### Testing

```sh
# TODO: ignore tmp dir for testing
test ./path/to/package
```
