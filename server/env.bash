[[ -z ${GITHUB_CLIENT_ID} ]] && export GITHUB_CLIENT_ID="asldkfjsadklfjk"
[[ -z ${GITHUB_CLIENT_SECRET} ]] && export GITHUB_CLIENT_SECRET="asldkfjsadklfjk"
[[ -z ${VARIABLE_KEY} ]] && export VARIABLE_KEY="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${USER_TOKEN_KEY} ]] && export USER_TOKEN_KEY="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${USER_ID_KEY} ]] && export USER_ID_KEY="asdflkjasldkfjasldkjfasldkjfasas"
[[ -z ${JWT_SECRET} ]] && export JWT_SECRET="asdflkjasldkfjasldkjfasldkjfasas"
export ENVHUB_PATH="${PWD}"


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

# db related functions
export MYSQL="envhubuser:envhubpassword@tcp(127.0.0.1:3306)/envhub?tls=true&interpolateParams=true"
MIGRATION_DIR="${ENVHUB_PATH}/database/migrations"

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
        migrate)
            case $2 in
                 force)
                     migrate -verbose -path ${MIGRATION_DIR} -database ${MYSQL} force $3
                     ;;
                 new)
                     migrate -verbose create -ext sql -dir ${MIGRATION_DIR} -seq $3
                     ;;
                 up)
                     migrate -verbose -database ${MYSQL} -path ${MIGRATION_DIR} up 1
                     ;;
                 down)
                     migrate -verbose -database ${MYSQL} -path ${MIGRATION_DIR} down 1
                     ;;
                 *)
                     echo "usage: [force|new|up|down]"
                     ;;
             esac
             ;;
        *)
            echo "usage: [view|dbml|seed]"
            ;;
    esac
}


