#!/usr/bin/env bash
#
# release.sh — cut a new envdir release.
#
# Reads the canonical version from main.go, validates that a matching git tag
# does not already exist, runs the test suite, builds cross-platform binaries
# via goreleaser, tags the commit, pushes the tag, and creates a GitHub
# release.
#
# Usage:
#   scripts/release.sh              # full release
#   scripts/release.sh --dry-run    # print every action, do nothing
#
# Requirements:
#   - go, goreleaser, gh (GitHub CLI, authenticated), git
#   - Clean working tree, on main, synced with origin
#

set -euo pipefail

# ---- argument parsing --------------------------------------------------------

DRY_RUN=0

for arg in "$@"; do
    case "$arg" in
        --dry-run)
            DRY_RUN=1
            ;;
        -h|--help)
            sed -n '2,/^$/p' "$0" | sed 's/^# \{0,1\}//'
            exit 0
            ;;
        *)
            echo "release.sh: unknown argument '$arg'" >&2
            echo "usage: $0 [--dry-run]" >&2
            exit 2
            ;;
    esac
done

# ---- helpers -----------------------------------------------------------------

log() { printf '==> %s\n' "$*"; }
err() { printf '!!  %s\n' "$*" >&2; }

run() {
    printf '    $ %s\n' "$*"
    if [ "$DRY_RUN" -eq 0 ]; then
        "$@"
    fi
}

require_tool() {
    if ! command -v "$1" >/dev/null 2>&1; then
        err "required tool not found on PATH: $1"
        exit 1
    fi
}

# ---- locate repo root --------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

# ---- preflight: tooling ------------------------------------------------------

log "checking required tools"
require_tool go
require_tool goreleaser
require_tool gh
require_tool git

# ---- preflight: Go version is a stable release, not tip ---------------------

go_version=$(go version)
if ! printf '%s' "$go_version" | grep -qE 'go[0-9]+\.[0-9]+(\.[0-9]+)? '; then
    err "go appears to be a development build, not a stable release:"
    err "  $go_version"
    err "install a released Go version before releasing."
    exit 1
fi
log "go version: $go_version"

# ---- preflight: GitHub auth --------------------------------------------------

log "checking GitHub authentication"
if ! gh auth status >/dev/null 2>&1; then
    err "gh is not authenticated. Run 'gh auth login' first."
    exit 1
fi

# ---- preflight: git state ----------------------------------------------------

log "verifying git state"

if ! git diff-index --quiet HEAD -- ; then
    err "working tree has uncommitted changes. Commit or stash first."
    git status --short >&2
    exit 1
fi

current_branch="$(git rev-parse --abbrev-ref HEAD)"
if [ "$current_branch" != "main" ]; then
    err "must release from 'main' branch, currently on '$current_branch'"
    exit 1
fi

run git fetch origin main
local_head="$(git rev-parse HEAD)"
origin_head="$(git rev-parse origin/main)"
if [ "$local_head" != "$origin_head" ]; then
    err "local main ($local_head) is not in sync with origin/main ($origin_head)"
    err "pull or push to align, then retry."
    exit 1
fi

# ---- version: read from main.go ---------------------------------------------

require_tool current_version

log "reading version from main.go"
VERSION=$(current_version main.go)
if [ -z "$VERSION" ]; then
    err "could not extract Version from main.go"
    exit 1
fi

# current_version returns semver; validate it has all three parts.
if ! printf '%s' "$VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
    err "version '$VERSION' is not valid semver (expected MAJOR.MINOR.PATCH)"
    exit 1
fi

TAG="v$VERSION"
log "version: $VERSION (tag: $TAG)"

# ---- preflight: tag must not already exist -----------------------------------

if git rev-parse --verify --quiet "refs/tags/$TAG" >/dev/null; then
    err "tag $TAG already exists locally. Bump the version in main.go or delete the tag."
    exit 1
fi
if git ls-remote --exit-code --tags origin "refs/tags/$TAG" >/dev/null 2>&1; then
    err "tag $TAG already exists on origin. Bump the version in main.go and retry."
    exit 1
fi

# ---- test --------------------------------------------------------------------

log "running tests"
run go test -trimpath -race ./...

# ---- tag + push --------------------------------------------------------------

log "tagging $TAG"
run git tag -a "$TAG" -m "envdir $TAG"
run git push origin "$TAG"

# ---- goreleaser --------------------------------------------------------------

log "building and releasing with goreleaser"
run goreleaser release --clean

log "release $TAG complete"
if [ "$DRY_RUN" -eq 1 ]; then
    log "(dry run: nothing was actually changed)"
fi
