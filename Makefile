GIT_REF := $(shell git describe --always)
VERSION ?= commit-$(GIT_REF)

.PHONY: dep
dep:
	GO111MODULE=on go mod vendor -v

.PHONY: build
build:
	CGO_ENABLED=0 go build -o bin/server \
        -ldflags "-X main.version=$(VERSION)" \
        github.com/rerorero/graph

.PHONY: test
test:
	GO111MODULE=on go test -v -race ./...

.PHONY: coverage
coverage:
	GO111MODULE=on go test -race -coverpkg=./... -coverprofile=coverage.txt ./...

