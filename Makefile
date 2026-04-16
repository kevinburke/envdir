.PHONY: vendor

DIFFER := $(GOPATH)/bin/differ

$(DIFFER):
	go install github.com/kevinburke/differ@latest

lint:
	go vet -trimpath ./...
	staticcheck ./...

test:
	go test -race ./...

ci-diffs: $(DIFFER)
	$(DIFFER) go mod tidy
	$(DIFFER) goimports -w .
	$(DIFFER) go fmt ./...

vendor:
	go mod vendor

release:
	bash scripts/release.sh
