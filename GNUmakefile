export CGO_ENABLED = 0
VERSION = $(shell git describe --tags --match='v*' --always)

CROSSBUILD_OS   = linux windows darwin
CROSSBUILD_ARCH = 386 amd64
OSARCH_COMBOS   = $(foreach os,$(CROSSBUILD_OS),$(addprefix $(os)_,$(CROSSBUILD_ARCH)))

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
	@go build -o terraform-provider-sops

crossbuild: $(GOPATH)/bin/gox
	@echo ">> cross-building"
	@gox -arch="$(CROSSBUILD_ARCH)" -os="$(CROSSBUILD_OS)" -output="binaries/{{.OS}}_{{.Arch}}/terraform-provider-sops_$(VERSION)"

$(GOPATH)/bin/gox:
	# Need to disable modules for this to not pollute go.mod
	@GO111MODULE=off go get -u github.com/mitchellh/gox

release: crossbuild bin/hub
	@echo ">> uploading release ${VERSION}"
	@mkdir -p releases
	@set -e; for OSARCH in $(OSARCH_COMBOS); do \
		zip -j releases/terraform-provider-sops_$(VERSION)_$$OSARCH.zip binaries/$$OSARCH/terraform-provider-sops_* > /dev/null; \
		./bin/hub release edit -a "releases/terraform-provider-sops_$(VERSION)_$$OSARCH.zip#terraform-provider-sops_$(VERSION)_$$OSARCH.zip" ${VERSION}; \
	done

bin/hub:
	@mkdir -p bin
	@curl -sL 'https://github.com/github/hub/releases/download/v2.14.1/hub-linux-amd64-2.14.1.tgz' | \
		tar -xzf - --strip-components 2 -C bin --wildcards '*/bin/hub'

.PHONY: all style vet test build crossbuild release
