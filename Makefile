# constants

# default source directory
SOURCE_DIRECTORY=$(GOPATH)/src
# default build output directory
OUTPUT_DIRECTORY=$(CURDIR)/build
# default to current directory name
OUTPUT_NAME=$(shell basename $(CURDIR))
# substracting GOPATH from source directory to deduct import path. i.e. github.com/user/project
IMPORT_PATH=$(subst $(SOURCE_DIRECTORY)/,,$(CURDIR))
# sparks sdk root directory
SPARKS_SDK_ROOT=$(HOME)/Sources/Sparks
# GOPATH must be set
ifndef GOPATH
$(error GOPATH is not set)
endif

# host os detection
HOST_OS=unknown
ifeq ($(OS),Windows_NT)
	HOST_OS=windows
else
	UNAME=$(shell uname -s)
	ifeq ($(UNAME),Linux)
		HOST_OS=linux
	endif
	ifeq ($(UNAME),Darwin)
		HOST_OS=osx
	endif
endif

# default to test and build
all: test build

# check build environment consistency
check:
	@if [ "$(HOST_OS)" = 'unknown' ]; then echo "FATAL: Could not detect HOST_OS"; exit 1; fi

# install build dependencies
deps: check
	go get -d
	go get golang.org/x/tools/cmd/goimports
	go get github.com/kisielk/errcheck
	go get github.com/maruel/panicparse
	go get golang.org/x/lint/golint

# install development build dependencies
deps-dev: check deps
	go get golang.org/x/tools/cmd/gorename
	# go get github.com/spf13/cobra/cobra

# format go code
format: check
	-goimports -l -w .

# run golint on all source tree
lint: check
	golint ./...

# check for untested errors in code
errcheck: check
	-errcheck ./...

# build
build: check format lint errcheck
	go build -v -o "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"
	@du -h "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"

# install into "$(GOPATH)/bin"
install: build
	@cp "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)" "$(GOPATH)/bin/"

# build linux binary inside a docker container
build-docker: check
	mkdir -p "$(OUTPUT_DIRECTORY)/linux"
	docker run --rm -it                         \
		-v "$(GOPATH)/src":/go/src              \
		-v "$(OUTPUT_DIRECTORY)/linux":/build   \
		-w "/go/src/$(IMPORT_PATH)"             \
		golang:latest                           \
		make build-docker-linux
	@du -h "$(OUTPUT_DIRECTORY)/linux/$(OUTPUT_NAME)"

# build linux binary inside a docker container
run-docker: check
	docker run --rm -it                         \
		-v "$(OUTPUT_DIRECTORY)/linux":/build   \
		-v "$(SPARKS_SDK_ROOT)":/sparks 		\
		golang:latest                           \
		/build/$(OUTPUT_NAME) $(ARGS)

# linux build inside a docker container, see build-docker rule above
build-docker-linux: check deps
	go build -v -o "/build/$(OUTPUT_NAME)"

# run tests
test: check
	go test -v ./...

# clean build artefacts
clean: check
	go clean
	rm -rf "$(OUTPUT_DIRECTORY)"

# run built binary
run: check build
	"$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)" $(ARGS) 2>&1 | panicparse

# build and run in docker
rund: build-docker run-docker
