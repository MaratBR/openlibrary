#!/bin/sh

set -eu

if [ "$#" -eq 0 ]; then
    echo "No tools specified for verification" >&2
    exit 1
fi

for tool in "$@"; do
    if command -v "$tool" >/dev/null 2>&1; then
        continue
    fi

    case "$tool" in
        sqlc) hint="go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest" ;;
        gow) hint="go install github.com/mitranim/gow@latest" ;;
        templ) hint="go install github.com/a-h/templ/cmd/templ@latest" ;;
        *) hint="install $tool and ensure it is available on PATH" ;;
    esac

    echo "$tool not found: $hint" >&2
    exit 1
done
