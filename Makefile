.PHONY: vendor

test:
	go test -race ./...

vendor:
	dep ensure
