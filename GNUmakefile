version = $(shell git describe --tags --match='v*' | sed 's/^v//')

default: build

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

vet:
	@echo ">> vetting code"
	@go vet ./...

test:
	@echo ">> testing code"
	@go test -v ./...

build:
	@echo ">> building binaries"
	@CGO_ENABLED=0 go build -o terraform-provider-sops

.PHONY: all style vet test build
