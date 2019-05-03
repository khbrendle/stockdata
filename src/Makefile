# Go Parameters

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mlManager

# all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
# run:
# 	$(GOBUILD) -o $(BINARY_NAME) -v
# 	./$(BINARY_NAME)
test:
	$(GOTEST) -v ./... #-bench=. -benchmem
# deps:
# 	$(GOGET) github.com/someone/something

# # Cross compilation
# build-linux:
#   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
# docker-build:
#   docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
