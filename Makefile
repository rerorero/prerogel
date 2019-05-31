GIT_REF := $(shell git describe --always)
VERSION ?= commit-$(GIT_REF)

.PHONY: dep
dep:
	GO111MODULE=on go mod vendor -v

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

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


PROTO_FILES := $(shell find . -type f -name '*.proto' -print)
PROTO_GEN_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
PROTOC := protoc --gogoslick_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,plugins=grpc:.

protogen: $(PROTO_GEN_FILES)
%.pb.go: %.proto
	cd $(dir $<); $(PROTOC) --proto_path=. --proto_path=$(GOPATH)/src ./*.proto
