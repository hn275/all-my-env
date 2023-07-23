# Server

## Requirements

- Go 1.20
- Docker

## ENV

- `env.bash` has a collection of functions needed for development

```sh
source env.bash
```

## Get started

### Database functions

```sh
db view # access psql shell in docker

db dbml # generate a new schemas file, this need pg-to-dbml installed
        # https://github.com/papandreou/pg-to-dbml

db mock # mock rows
```

### Testing

```sh
gotest ./path/to/package # make sure to have docker container running and mock data in db
# or test all packages with
gotest
```
