#!/usr/bin/env bash

set -euo pipefail

main() {
    export GO111MODULE=on
    go install github.com/goreleaser/goreleaser@latest
    GOROOT=~/go1.20 PATH=~/go1.20/bin:$PATH goreleaser release
}

main "$@"
