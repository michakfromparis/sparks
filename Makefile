# Constants
OUTPUT_NAME=$(shell basename $(CURDIR))
IMPORT_PATH=$(subst $(GOPATH)/src/,,$(CURDIR))
OUTPUT_DIRECTORY=$(CURDIR)/build
# OUTPUT_DIRECTORY=$(HOME)/.cache/go-docker-build/$(OUTPUT_NAME)

# Host OS Detection
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


all: test build

check:
	@if [ "$(HOST_OS)" = 'unknown' ]; then echo "FATAL: Could not detect HOST_OS"; exit 1; fi

deps: check
	go get -u golang.org/x/tools/cmd/goimports

deps-dev: check deps
	go get -u github.com/spf13/cobra/cobra

style: check
	goimports -l -w .

build: check style
	go build -v -o "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"
	@du -h "$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"
	
build-docker: check
	mkdir -p "$(OUTPUT_DIRECTORY)/linux"
	docker run --rm -it                         \
		-v "$(GOPATH)/src":/go/src              \
		-v "$(OUTPUT_DIRECTORY)/linux":/build   \
		-w "/go/src/$(IMPORT_PATH)"             \
		golang:latest                           \
		make build-docker-linux
	@du -h "$(OUTPUT_DIRECTORY)/linux/$(OUTPUT_NAME)"

build-docker-linux: check deps
	go build -v -o "/build/$(OUTPUT_NAME)"

test: check
	go test -v ./...

clean: check
	go clean
	rm -rf "$(OUTPUT_DIRECTORY)"

run: check build
	"$(OUTPUT_DIRECTORY)/$(HOST_OS)/$(OUTPUT_NAME)"

