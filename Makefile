# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=filehasher
BINARY_DIR=bin
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_LINUX_ARM=$(BINARY_LINUX)_arm

all: build
build: bin
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v
# test:
# 	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf bin
run: bin
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v ./...
	./$(BINARY_DIR)/$(BINARY_NAME)
deps:
	$(GOGET) github.com/mattn/go-sqlite3

# Pre-requisites
bin:
	mkdir $@


# Cross compilation
build-linux: bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_LINUX) -v

build-linux-arm: bin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_LINUX_ARM) -v
