.PHONY: vendor

test:
	go test -race ./...

vendor:
	go mod vendor
