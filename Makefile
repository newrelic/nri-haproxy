WORKDIR      := $(shell pwd)
TARGET       := target
TARGET_DIR    = $(WORKDIR)/$(TARGET)
NATIVEOS	 := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH	 := $(shell go version | awk -F '[ /]' '{print $$5}')
INTEGRATION  := haproxy
BINARY_NAME   = nri-$(INTEGRATION)
GO_PKGS      := $(shell go list ./... )
GO_FILES     := ./src/
GOTOOLS       = github.com/axw/gocov/gocov \
		github.com/AlekSi/gocov-xml \

.DEFAULT_GOAL := docker-make

all: check-version clean validate test compile

build: docker-make

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

tools: check-version
	@echo "=== $(INTEGRATION) === [ tools ]: Installing tools required by the project..."
	@go get $(GOTOOLS)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.27.0

tools-update: check-version
	@echo "=== $(INTEGRATION) === [ tools-update ]: Updating tools required by the project..."
	@go get -u $(GOTOOLS)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.27.0

deps: tools

validate: deps
	@echo "=== $(INTEGRATION) === [ validate ]: Running lints..."
	@go mod download
	@golangci-lint run


compile: deps
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

compile-only:
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

test: deps
	@echo "=== $(INTEGRATION) === [ test ]: Running unit tests..."
	@gocov test -race $(GO_FILES)... | gocov-xml > coverage.xml

# Include thematic Makefiles
include Makefile-*.mk

check-version:
ifdef GOOS
ifneq "$(GOOS)" "$(NATIVEOS)"
	$(error GOOS is not $(NATIVEOS). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
endif
ifdef GOARCH
ifneq "$(GOARCH)" "$(NATIVEARCH)"
	$(error GOARCH variable is not $(NATIVEARCH). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
endif

.PHONY: all build clean tools tools-update deps validate compile test check-version
