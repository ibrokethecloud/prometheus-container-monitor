SHELL := /bin/bash
WORKING_DIR := $(shell pwd)

IMAGE_NAME = prometheus-container-monitor
IMAGE_VERSION = latest
IMAGE_TAG = gmehta3/$(IMAGE_NAME):$(IMAGE_VERSION)

GOOS?=linux darwin
GOARCH?=amd64
GOLDFLAGS?=-ldflags "-s" -a -installsuffix cgo -o ./bin/$(IMAGE_NAME)

.PHONY: build push docker-build docker-clean docker-prepare dep dep-update

all:: build

release:: docker-build push

push::
	@docker push $(IMAGE_TAG)

dep::
	@go get -u github.com/golang/dep/...

dep-update:: dep
	@dep init
	@dep ensure -update

build::	
	@echo ">> building"
	@for arch in ${GOARCH}; do \
		for os in ${GOOS}; do \
			echo ">>>> $${os}/$${arch}"; \
			env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} go build ${GOLDFLAGS}; \
		done; \
	done

docker-clean::
	@echo ">> cleaning (docker)"
	@docker rmi $(IMAGE_NAME)-build &>/dev/null || true

docker-prepare:: docker-clean
	@echo ">> preparing (docker)"
	@docker build --no-cache -t $(IMAGE_NAME)-build -f Dockerfile_build .

docker-build:: docker-prepare
	@echo ">> building (docker)"
	@test -f $@.cid && { docker rm -f $$(cat $@.cid) && rm $@.cid; } || true;
	@docker run -t --cidfile="$@.cid" \
		-v "$$PWD":"/go/src/$(IMAGE_NAME)" $(IMAGE_NAME)-build
	@docker stop $$(cat $@.cid)
	@docker rm $$(cat $@.cid)
	@docker build -t $(IMAGE_TAG) $(WORKING_DIR)
