#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$repo_root"

(cd render && go run .)

wrangler pages deploy dist/ "$@"
