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
    # if [ -z "${1}" ];then
    #     echo "Testing all packages"
    #     go test ./... -coverprofile cover.out
    # else
    #     echo "Testing package ${1}"
    #     go test $1 -coverprofile cover.out
    # fi
    
    echo "Testing package ${1}"
    go test $1 -coverprofile cover.out
    [[ -f cover.out ]] && rm cover.out
}

function dbview() {
    docker exec -it envhub-db psql -U username envhub
}

function dbml() {
    pg-to-dbml -c=${POSTGRES_DSN}
}
