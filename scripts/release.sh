#!/usr/bin/env bash

set -euo pipefail

main() {
    go install github.com/goreleaser/goreleaser@latest
    GOROOT=~/go1.18 PATH=~/go1.18/bin:$PATH goreleaser build
}

main "$@"
