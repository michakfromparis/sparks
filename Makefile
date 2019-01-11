OUTPUT_NAME=$(shell basename $(CURDIR))

all: test build

deps:
	go get -d

build: 
	@go build -v -o $(OUTPUT_NAME)
test: 
	@go test -v ./...
clean: 
	@go clean
	@rm -f $(OUTPUT_NAME)
run:	build
	@./$(OUTPUT_NAME)

# Cross compilation
# build-linux:
#         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v
# docker-build:
#         docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v    
