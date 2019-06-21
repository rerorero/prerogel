GIT_REF := $(shell git describe --always)
VERSION ?= commit-$(GIT_REF)
GOCMD := GO111MODULE=on go

.PHONY: dep
dep:
	$(GOCMD) mod vendor -v

.PHONY: tidy
tidy:
	$(GOCMD) mod tidy

.PHONY: run
run:
	$(GOCMD) run main.go plugin.go

