#!/usr/bin/env bash
set -euo pipefail

go install github.com/ashupednekar/slides/render@c2c1437d15dcce98607a8cf3a2660b8e9bb47541
render
wrangler pages deploy dist/ 
