pkgs = $(shell go list ./...)
version = $(shell git describe --tags --match='v*' | sed 's/^v//')

default: build

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@go fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@go vet $(pkgs)

test:
	@echo ">> testing code"
	@go test $(pkgs)

build:
	@echo ">> building binaries"
	@CGO_ENABLED=0 go build -o terraform-provider-sops

.PHONY: all style format build vet
