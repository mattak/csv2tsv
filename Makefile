GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_DIR=bin
BINARY_NAME=csv2tsv
TARGET_FILE=./cmd/csv2tsv/main.go

all: build
install: build system_install

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) $(TARGET_FILE)

.PHONY: run
run:
	$(GORUN) $(TARGET_FILE)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -r $(BINARY_DIR)

.PHONY: system_install
system_install:
	cd cmd/csv2tsv && go install

