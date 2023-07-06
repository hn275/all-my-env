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

### Database Migration

```sh
dbnew some_migration_name
dbup # deploy migration
dbdown # rollback migration
dbfix broken_version # if you made a syntax, this is needed before deploying again
```

### Testing

```sh
# TODO: ignore tmp dir for testing
test ./path/to/package
```
