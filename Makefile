GIT_REF := $(shell git describe --always)
VERSION ?= commit-$(GIT_REF)
GOCMD := GO111MODULE=on go

.PHONY: dep
dep:
	$(GOCMD) mod vendor -v

.PHONY: tidy
tidy:
	$(GOCMD) mod tidy

.PHONY: build
build: dep plugin protogen
	CGO_ENABLED=0 $(GOCMD) build -o bin/server \
        -ldflags "-X main.version=$(VERSION)" \
        github.com/rerorero/prerogel

.PHONY: test
test: dep plugin protogen
	@$(GOCMD) test -v ./...

.PHONY: coverage
coverage: dep protogen
	@$(GOCMD) test -coverpkg=./... -coverprofile=coverage.txt ./...


PROTO_FILES := $(shell find . \( -path "./vendor" \) -prune -o -type f -name '*.proto' -print)
PROTO_GEN_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
PROTOC := protoc --gogoslick_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,plugins=grpc:.

plugin:
	@$(GOCMD) get google.golang.org/grpc
	@$(GOCMD) get github.com/golang/protobuf/protoc-gen-go
	@$(GOCMD) get github.com/gogo/protobuf/proto
	@$(GOCMD) get github.com/gogo/protobuf/protoc-gen-gogoslick
	@$(GOCMD) get github.com/gogo/protobuf/gogoproto

protogen:
	@for d in $(PROTO_GEN_FILES); do cd `dirname $$d`; $(PROTOC) --proto_path=. --proto_path=$(GOPATH)/src ./*.proto; done

