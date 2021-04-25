.PHONY: all clean plugins

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))

DOCKER_IMAGE := logark/logark

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CMD:=logark

BIN_DIR:=dist
PLUGIN_DIR:=$(BIN_DIR)/plugins
CMD_DIR:=./cmd
    
all: plugins binary

plugins:
	GO111MODULE=on $(GOBUILD) -buildmode=plugin -o $(PLUGIN_DIR)/mutate.so ./plugins/filters/mutate/...
	GO111MODULE=on $(GOBUILD) -buildmode=plugin -o $(PLUGIN_DIR)/test.so ./plugins/filters/test/...
	GO111MODULE=on $(GOBUILD) -buildmode=plugin -o $(PLUGIN_DIR)/prune.so ./plugins/filters/prune/...

$(BIN_DIR):
	mkdir -p $@

binary: $(BIN_DIR)
	GO111MODULE=on $(GOBUILD) -o ./$(BIN_DIR)/$(CMD) $(CMD_DIR)/$(CMD) 

docker: binary
	docker build -t $(DOCKER_IMAGE) .

clean: 
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -f $(CMD)

run-dev:
	GO111MODULE=on go build $(CMD_DIR)/$(CMD)
	./$(CMD)