# constants
OUTPUT_NAME=$(shell basename $(CURDIR))        # default to current directory name
IMPORT_PATH=$(subst $(GOPATH)/src/,,$(CURDIR)) # substracting GOPATH from current directory
OUTPUT_DIRECTORY=$(CURDIR)/build               # build output directory

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


# install development build dependencies
deps-dev: check deps
	go get github.com/spf13/cobra/cobra

# format go code
format: check
	goimports -l -w .

# build
build: check format install
	go build -v -o "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"
	@du -h "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"

# install into "$(GOPATH)/bin"
install:
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

# linux build called inside the docker container
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
	"$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"
