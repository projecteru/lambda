.PHONY: golang python deps build test

REPO_PATH := gitlab.ricebook.net/platform/lambda
REVISION := $(shell git rev-parse HEAD || unknown)
BUILTAT := $(shell date +%Y-%m-%dT%H:%M:%S)
VERSION := $(shell cat VERSION)
GO_LDFLAGS ?= -s -X $(REPO_PATH)/versioninfo.REVISION=$(REVISION) \
			  -X $(REPO_PATH)/versioninfo.BUILTAT=$(BUILTAT) \
			  -X $(REPO_PATH)/versioninfo.VERSION=$(VERSION)

deps:
	glide i

build: deps
	go build -ldflags "$(GO_LDFLAGS)" -a -tags netgo -installsuffix netgo -o lambda

test: deps
	go vet `go list ./... | grep -v '/vendor/'`
	go test -v `glide nv`
rpm:
	ROOT="`pwd`/build"
	BIN="$ROOT/usr/bin"
	CONF="$ROOT/etc/eru"
	mkdir -p $BIN
	mkdir -p $CONF
	mv lambda $BIN
	mv lambda.yaml.example $CONF
	VERSION=$(cat VERSION)
	echo $VERSION rpm build begin
	fpm -f -s dir -t rpm -n eru-lambda --epoch 0 -v $VERSION --iteration 1.el7 -C $ROOT -p $PWD --verbose --rpm-auto-add-directories --category 'Development/App' --description 'docker eru lambda executor' --url 'http://gitlab.ricebook.net/platform/lambda/' --license 'BSD'  --no-rpm-sign usr etc
	rm -rf $ROOT