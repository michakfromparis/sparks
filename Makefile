OUTPUT_NAME=$(shell basename $(CURDIR))
OUTPUT_PATH=$(GOPATH)/bin/$(OUTPUT_NAME)

all: test build

deps:
	@go get -d

build:
	@go build -v -o $(OUTPUT_PATH)
	
test:
	@go test -v ./...

clean:
	@go clean
	@rm -f $(OUTPUT_PATH)

run:	build
	@$(OUTPUT_PATH)

# Cross compilation
# build-linux:
#         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v
# docker-build:
#         docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v    
