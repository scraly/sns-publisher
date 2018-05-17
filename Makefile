GO   := GO15VENDOREXPERIMENT=1 go
pkgs  = $(shell $(GO) list ./... | grep -v /vendor/)
IMAGE_NAME ?= "scraly/sns-publisher"
http_proxy ?= ""
https_proxy ?= ""
TEST_OUT_FILE ?= "gotest.out"
TEST_REPORT ?= "report.xml"

all: format build

get:
	# govendor
	go get -u -v github.com/kardianos/govendor
	#deps
	# go get -u github.com/gocql/gocql
	go get -u github.com/gorilla/mux
	# go testing
	go get -u github.com/jstemmer/go-junit-report
	go get github.com/smartystreets/goconvey
	# docs
	# npm install -g bootprint bootprint-openapi html-inline

license:
	@echo ">> checking license from dependencies with wwhrd"
	@$(GO) get -u github.com/frapposelli/wwhrd
	@wwhrd check

test:
	@echo ">> running tests"
	@$(GO) test -short $(pkgs) -cover

test-jenkins-format:
	@echo ">> running tests and transform with go-junit-report"
	@go get -u github.com/jstemmer/go-junit-report
	@$(GO) test -v -cover $(pkgs) 2>&1 | go-junit-report > ${TEST_REPORT}

vendor:
	govendor init
	govendor add -v +external
	govendor update -v +vendor

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

build:
	@echo ">> building binaries"
	@./ci/build.sh

pack: build
	@echo ">> packing all binaries"
	@upx -9 bin/*

docker: pack docker-ci

docker-ci :
	@docker build -t $(IMAGE_NAME):$(shell cat version/VERSION)-$(shell git rev-parse --short HEAD) --build-arg http_proxy=${http_proxy} --build-arg https_proxy=${https_proxy} --no-cache --rm .
	@docker tag $(IMAGE_NAME):$(shell cat version/VERSION)-$(shell git rev-parse --short HEAD) $(IMAGE_NAME):latest

docker-release: docker
	@docker tag $(IMAGE_NAME):$(shell cat version/VERSION)-$(shell git rev-parse --short HEAD) $(IMAGE_NAME):$(shell cat version/VERSION)

clean:
	-rm -f ${TEST_REPORT}
	-rm -rf bin/

.PHONY: all get format build test docker vendor clean
