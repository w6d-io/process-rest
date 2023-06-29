VERSION    ?= $(shell basename /$(shell git symbolic-ref --quiet HEAD 2> /dev/null ) )
VCS_REF    = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
export ACK_GINKGO_DEPRECATIONS = 1.16.5
export GO111MODULE  := on
export NEXT_TAG     ?=
export CGO_ENABLED   = 1
export LOG_LEVEL     = 2

export PATH         := $(shell pwd)/bin:${PATH}

IMG ?= w6dio/process-rest:latest

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

ifeq (,$(shell go env GOOS))
GOOS       = $(shell echo $OS)
else
GOOS       = $(shell go env GOOS)
endif

ifeq (,$(shell go env GOARCH))
GOARCH     = $(shell echo uname -p)
else
GOARCH     = $(shell go env GOARCH)
endif

ifeq (gsed not found,$(shell which gsed))
SEDBIN=sed
else
SEDBIN=$(shell which gsed)
endif

ifeq (darwin,$(GOOS))
GOTAGS = "-tags=dynamic"
else
GOTAGS =
endif

ifndef (,$(NEXT_TAG))
CHGLOG_FLAG = "--next-tag=$(NEXT_TAG)"
else
CHGLOG_FLAG =
endif

.PHONY: all
all: build

##@ Development

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet $(GOTAGS) ./...

.PHONY: test
test: fmt vet
	go test $(GOTAGS) -v -coverpkg=./... -coverprofile=cover.out ./...
#	@go tool cover -func cover.out | grep total

# Formats the code
.PHONY: format
format: goimports
	$(GOIMPORTS) -w -local github.com/w6d-io,gitlab.w6d.io cmd internal pkg

# Changelog
.PHONY: changelog
changelog: chglog
	$(GITCHGLOG) -o docs/CHANGELOG.md $(CHGLOG_FLAG)

##@ Build

.PHONY: build

##@ Build

.PHONY: build
build: fmt vet
	go build $(GOTAGS) \
       -ldflags="-X 'github.com/w6d-io/process-rest/internal/config.Version=${VERSION}' -X 'github.com/w6d-io/process-rest/internal/config.Revision=${VCS_REF}' -X 'github.com/w6d-io/process-rest/internal/config.Built=${BUILD_DATE}'" \
       -a -o bin/process-rest main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run $(GOTAGS) main.go -config config/tests/config.yaml -log-level 2

.PHONY: docker-build
docker-build: Dockerfile
	docker build --build-arg=VERSION=${VERSION} --build-arg=VCS_REF=${VCS_REF} --build-arg=BUILD_DATE=${BUILD_DATE} -t ${IMG} .

.PHONY: docker-push
docker-push:
	docker push ${IMG}

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOIMPORTS               ?= $(LOCALBIN)/goimports
GITCHGLOG               ?= $(LOCALBIN)/git-chglog

.PHONY: chglog
chglog: $(GITCHGLOG) ## Download git-chglog locally if necessary
$(GITCHGLOG): $(LOCALBIN)
	@test -s $(LOCALBIN)/git-chglog || GOBIN=$(LOCALBIN) go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: goimports
goimports: $(GOIMPORTS) ## Download goimports locally if necessary
$(GOIMPORTS): $(LOCALBIN)
	@test -s $(LOCALBIN)/goimports || GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@latest
