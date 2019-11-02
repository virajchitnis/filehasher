# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=filehasher
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_LINUX_ARM=$(BINARY_LINUX)_arm

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_LINUX_ARM)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
#deps:
#        $(GOGET) github.com/markbates/goth
#        $(GOGET) github.com/markbates/pop


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v

build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_LINUX_ARM) -v
