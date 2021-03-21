IMG ?= w6dio/process-rest:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

REF=$(shell git symbolic-ref --quiet HEAD 2> /dev/null)
VERSION=$(shell basename $(REF) )
VCS_REF=$(shell git rev-parse HEAD)
GOVERSION=$(shell go version | awk '{ print $3 }' | sed 's/go//')
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=$(shell uname -s | tr "[:upper:]" "[:lower:]")
GOARCH=$(shell uname -p)

all: process-rest

# Run tests
test: fmt vet
	go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

# Build process-rest binary
process-rest: fmt vet vendor
	VERSION=${VERSION/refs\/heads\//}
	go build -ldflags="-X 'main.Version=${VERSION}' -X 'main.Revision=${VCS_REF}' -X 'main.GoVersion=go${GOVERSION}' -X 'main.Built=${BUILD_DATE}' -X 'main.OsArch=${GOOS}/${GOARCH}'" -mod=vendor -a -o bin/process-rest cmd/process-rest/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run cmd/process-rest/main.go -config config/tests/config.yaml -log-level 2

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

vendor:
	go mod vendor

# Build the docker image
build: test
	docker build  --build-arg=VERSION=${VERSION} --build-arg=VCS_REF=${VCS_REF} --build-arg=BUILD_DATE=${BUILD_DATE}  -t ${IMG} .

# Push the docker image
push:
	docker push ${IMG}