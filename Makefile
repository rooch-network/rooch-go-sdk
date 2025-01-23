BUILDDIR ?= $(CURDIR)/build
TOOLS_DIR := tools

GO_BIN := ${GOPATH}/bin

DOCKER := $(shell which docker)
CUR_DIR := $(shell pwd)

ldflags := $(LDFLAGS)
build_tags := $(BUILD_TAGS)
build_args := $(BUILD_ARGS)

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static" -v
endif

ifeq ($(VERBOSE),true)
	build_args += -v
endif

BUILD_TARGETS := build install
BUILD_FLAGS := --tags "$(build_tags)" --ldflags '$(ldflags)'

all: build install

build: BUILD_ARGS := $(build_args) -o $(BUILDDIR)

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	CGO_CFLAGS="-O -D__BLST_PORTABLE__" go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-docker:
	$(DOCKER) build --tag rooch-network/rooch-go-sdk -f Dockerfile \
		$(shell git rev-parse --show-toplevel)

.PHONY: build build-docker

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...
