if ! command -v sqlc >/dev/null 2>&1
then
    echo "sqlc not found: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    exit 1
fi

if ! command -v gow >/dev/null 2>&1
then
    echo "sqlc not found: go install github.com/mitranim/gow@latest"
    exit 1
fi

if ! command -v templ >/dev/null 2>&1
then
    echo "sqlc not found: go install github.com/a-h/templ/cmd/templ@latest"
    exit 1
fi