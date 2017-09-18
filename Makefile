.PHONY: golang python deps build test

REPO_PATH := github.com/projecteru2/lambda
REVISION := $(shell git rev-parse HEAD || unknown)
BUILTAT := $(shell date +%Y-%m-%dT%H:%M:%S)
VERSION := $(shell cat VERSION)
GO_LDFLAGS ?= -s -X $(REPO_PATH)/versioninfo.REVISION=$(REVISION) \
			  -X $(REPO_PATH)/versioninfo.BUILTAT=$(BUILTAT) \
			  -X $(REPO_PATH)/versioninfo.VERSION=$(VERSION)

deps:
	glide i

build: deps
	go build -ldflags "$(GO_LDFLAGS)" -a -tags netgo -installsuffix netgo -o eru-lambda

test: deps
	go vet `go list ./... | grep -v '/vendor/'`
	go test -v `glide nv`
