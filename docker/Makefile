include Makefile.config

NS ?= michak
VERSION ?= latest

IMAGE_NAME ?= sparks-builder
CONTAINER_NAME ?= $(IMAGE_NAME)
CONTAINER_INSTANCE ?= default

all: build run

build: Dockerfile
	$(PREBUILD)
	docker build -t $(NS)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .
	$(POSTBUILD)
push:
	docker push $(NS)/$(IMAGE_NAME):$(VERSION)

shell:
	docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) -i -t $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION) /bin/bash

run:
	$(PRERUN)
	docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION) $(RUNARGS)
	$(POSTRUN)

rm:
	docker rm $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

release: build
	make push -e VERSION=$(VERSION)

stop:
	docker stop $(NS)/$(IMAGE_NAME):$(VERSION)

clean:
	docker rmi $(NS)/$(IMAGE_NAME):$(VERSION)

reset:
	-docker stop $(shell docker ps -a -q)
	-docker rm $(shell docker ps -a -q)
	-docker rmi $(shell docker images -q)