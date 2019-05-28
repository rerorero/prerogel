GIT_REF := $(shell git describe --always)
VERSION ?= commit-$(GIT_REF)

GOPATH := $(shell go env GOPATH)

$(GOPATH)/bin/dep:
	@go get github.com/golang/dep/cmd/dep

.PHONY: dep
dep: $(GOPATH)/bin/dep
	@dep ensure -v

.PHONY: dep-vendor-only
dep-vendor-only: $(GOPATH)/bin/dep
	@dep ensure -v -vendor-only

.PHONY: build
build: dep-vendor-only
	CGO_ENABLED=0 go build -o bin/server \
        -ldflags "-X main.version=$(VERSION)" \
        github.com/rerorero/graph/cmd/server

.PHONY: test
test:
	@go test -v -race ./...

.PHONY: coverage
coverage:
	@go test -race -coverpkg=./... -coverprofile=coverage.txt ./...

