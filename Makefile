.PHONY: vendor

test:
	go test -race ./...

vendor:
	go mod vendor

release:
	go install -trimpath -v ./...
	envdir envs/release bash scripts/release.sh
