include make_env
# constants
DOCKER_IMAGE_TAG = sparks-builder
# default build output directory
ROOT_DIRECTORY = $(CURDIR)/..
# default build output directory
BUILD_DIRECTORY = $(ROOT_DIRECTORY)/build
# sparks sdk root directory
SPARKS_SDK_ROOT ?= $(HOME)/Sources/Sparks

# default to build
all: build

NS ?= michak
VERSION ?= latest

IMAGE_NAME ?= sparks-builder
CONTAINER_NAME ?= sparks-builder
CONTAINER_INSTANCE ?= default

.PHONY: build push shell run rm release

build: Dockerfile
    docker build -t $(NS)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .

push:
    docker push $(NS)/$(IMAGE_NAME):$(VERSION)

shell:
    docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) -i -t $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION) /bin/bash

run:
    docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION)

start:
    docker run -d --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION)

stop:
    docker stop $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

rm:
    docker rm $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

release: build
    make push -e VERSION=$(VERSION)

# build
build: check
	cp $(BUILD_DIRECTORY)/linux/$(OUTPUT_NAME) docker
	docker build -t $(DOCKER_IMAGE_TAG) docker
	rm docker/sparks
	@du -h "$(BUILD_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"

# clean build artefacts
clean: check

# install into "$(GOPATH)/bin"
install: build

# build linux binary inside a docker container
run: check
	mkdir -p build
	docker run --rm -it                         \
		-v "$(SPARKS_SDK_ROOT)":/sparks 					\
		-v build:/build   												\
		$(DOCKER_IMAGE_TAG)                      	\
		bash -c "sparks build --v --linux /sparks"

