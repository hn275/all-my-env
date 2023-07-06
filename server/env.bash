[[ -z ${GITHUB_CLIENT_SECRET} ]] && export GITHUB_CLIENT_SECRET="asldkfjsadklfjk"
[[ -z ${GITHUB_CLIENT_ID} ]] && export GITHUB_CLIENT_ID="asldkfjsadklfjk"
[[ -z ${ROW_KEY_SECRET} ]] && export ROW_KEY_SECRET="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${JWT_SECRET} ]] && export JWT_SECRET="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${POSTGRES_PASSWORD} ]] && export POSTGRES_PASSWORD="password"
[[ -z ${POSTGRES_USER} ]] && export POSTGRES_USER="username"
[[ -z ${POSTGRES_DB} ]] && export POSTGRES_DB="envhub"
[[ -z ${POSTGRES_PORT} ]] && export POSTGRES_PORT="5432"
[[ -z ${POSTGRES_HOST} ]] && export POSTGRES_HOST="localhost"
[[ -z ${POSTGRES_SSLMODE} ]] && export POSTGRES_SSLMODE="disable"

export POSTGRES_DSN="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}"

MIGRATION_DIR="./db/migrations"
function test() {
    if [ -z "${1}" ];then
        echo "Testing all packages"
        go test ./... -coverprofile cover.out
    else
        echo "Testing package ${1}"
        go test $1 -coverprofile cover.out
    fi
    
    [[ -f cover.out ]] && rm cover.out
}

function dbnew() {
    migrate create -ext sql -dir db/migrations -seq $1
}

function dbup() {
    migrate -database $POSTGRES_DSN -path $MIGRATION_DIR up 1
}

function dbdown() {
    migrate -database $POSTGRES_DSN -path ./db/migrations down 1
}

function dbfix() {
    migrate -database $POSTGRES_DSN -path ./db/migrations force $1
}

function dbview() {
    docker exec -it envhub-db psql -U username envhub
}
