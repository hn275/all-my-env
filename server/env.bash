export GITHUB_CLIENT_SECRET="asldkfjsadklfjk"
export GITHUB_CLIENT_ID="asldkfjsadklfjk"
export ROW_KEY_SECRET="asdflkjasldkfjasldkjfasldkjfasas"
export JWT_SECRET="asdflkjasldkfjasldkjfasldkjfasas"

function test() {
    if [ -z "${1}" ];then
        echo "Testing all packages"
        # go test $(go list ./... | grep -v lib) -coverprofile cover.out
        go test ./... -coverprofile cover.out
    else
        echo "Testing package ${1}"
        go test $1 -coverprofile cover.out
    fi
    
    [[ -f cover.out ]] && rm cover.out
}
