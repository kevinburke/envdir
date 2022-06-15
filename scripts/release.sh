#!/usr/bin/env bash

set -euo pipefail

main() {
    go install github.com/goreleaser/goreleaser@latest
    goreleaser --version
}

main "$@"
