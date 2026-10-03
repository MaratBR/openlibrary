#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 || ! $1 =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    echo "Usage: $0 <lowercase-hyphenated-feature-name>" >&2
    exit 2
fi

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
spec_dir="$repo_root/specs/$1"
template_dir="$repo_root/docs/development/spec-templates"

if [[ -e $spec_dir ]]; then
    echo "Spec already exists: $spec_dir" >&2
    exit 1
fi

mkdir "$spec_dir"
cp "$template_dir"/{spec,plan,tasks}.md "$spec_dir/"
echo "Created $spec_dir"
