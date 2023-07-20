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

db seed # seed data
```

### Testing

```sh
gotest ./path/to/package
```
