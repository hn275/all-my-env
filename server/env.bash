[[ -z ${GITHUB_CLIENT_SECRET} ]] && export GITHUB_CLIENT_SECRET="asldkfjsadklfjk"
[[ -z ${GITHUB_CLIENT_ID} ]] && export GITHUB_CLIENT_ID="asldkfjsadklfjk"
[[ -z ${ROW_KEY} ]] && export ROW_KEY="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${JWT_SECRET} ]] && export JWT_SECRET="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${POSTGRES_PASSWORD} ]] && export POSTGRES_PASSWORD="password"
[[ -z ${POSTGRES_USER} ]] && export POSTGRES_USER="username"
[[ -z ${POSTGRES_DB} ]] && export POSTGRES_DB="envhub"
[[ -z ${POSTGRES_PORT} ]] && export POSTGRES_PORT="5432"
[[ -z ${POSTGRES_HOST} ]] && export POSTGRES_HOST="localhost"
[[ -z ${POSTGRES_SSLMODE} ]] && export POSTGRES_SSLMODE="disable"
export ENVHUB_PATH="${PWD}"

export POSTGRES_DSN="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}"

MIGRATION_DIR="./db/migrations"
function gotest() {
    if [ -z "${1}" ];then
        shift 1
        echo "Testing all packages"
        go test -v ./... "$@" -coverprofile cover.out 
    else
        pkg=$1
        echo "Testing package ${pkg}"
        shift 1
        go test -v "${pkg}" -coverprofile cover.out "$@"
    fi
    
    [[ -f cover.out ]] && rm cover.out
}

function db() {
    case $1 in
        view)
            docker exec -it envhub-db psql -U username envhub
            ;;
        dbml)
            pg-to-dbml -c=${POSTGRES_DSN}
            ;;
        mock)
            go run $ENVHUB_PATH/scripts/main.go mock
            ;;
        *)
            echo "usage: [view|dbml|seed]"
            ;;
    esac
}
