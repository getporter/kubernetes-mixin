MIXIN = kubernetes
PKG = get.porter.sh/mixin/$(MIXIN)
SHELL = bash

GO = GO111MODULE=on go

PORTER_HOME ?= $(HOME)/.porter

COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags 2> /dev/null || echo v0)
PERMALINK ?= $(shell git describe --tags --exact-match &> /dev/null && echo latest || echo canary)

LDFLAGS = -w -X $(PKG)/pkg.Version=$(VERSION) -X $(PKG)/pkg.Commit=$(COMMIT)
XBUILD = CGO_ENABLED=0 $(GO) build -a -tags netgo -ldflags '$(LDFLAGS)'
BINDIR = bin/mixins/$(MIXIN)

CLIENT_PLATFORM ?= $(shell go env GOOS)
CLIENT_ARCH ?= $(shell go env GOARCH)
RUNTIME_PLATFORM ?= linux
RUNTIME_ARCH ?= amd64
SUPPORTED_PLATFORMS = linux darwin windows
SUPPORTED_ARCHES = amd64

ifeq ($(CLIENT_PLATFORM),windows)
FILE_EXT=.exe
else ifeq ($(RUNTIME_PLATFORM),windows)
FILE_EXT=.exe
else
FILE_EXT=
endif

REGISTRY ?= $(USER)

.PHONY: build
build: build-client build-runtime

build-runtime:
	mkdir -p $(BINDIR)
	GOARCH=$(RUNTIME_ARCH) GOOS=$(RUNTIME_PLATFORM) $(GO) build -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(MIXIN)-runtime$(FILE_EXT) ./cmd/$(MIXIN)

build-client:
	mkdir -p $(BINDIR)
	$(GO) build -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(MIXIN)$(FILE_EXT) ./cmd/$(MIXIN)

xbuild-all:
	$(foreach OS, $(SUPPORTED_PLATFORMS), \
		$(foreach ARCH, $(SUPPORTED_ARCHES), \
				$(MAKE) $(MAKE_OPTS) CLIENT_PLATFORM=$(OS) CLIENT_ARCH=$(ARCH) MIXIN=$(MIXIN) xbuild; \
		))

xbuild: $(BINDIR)/$(VERSION)/$(MIXIN)-$(CLIENT_PLATFORM)-$(CLIENT_ARCH)$(FILE_EXT)
$(BINDIR)/$(VERSION)/$(MIXIN)-$(CLIENT_PLATFORM)-$(CLIENT_ARCH)$(FILE_EXT):
	mkdir -p $(dir $@)
	GOOS=$(CLIENT_PLATFORM) GOARCH=$(CLIENT_ARCH) $(XBUILD) -o $@ ./cmd/$(MIXIN)

test: test-unit test-integration
	$(BINDIR)/$(MIXIN)$(FILE_EXT) version

test-unit: build
	$(GO) test ./...

test-integration: xbuild
	# Test against the cross-built client binary that we will publish
	cp $(BINDIR)/$(VERSION)/$(MIXIN)-$(CLIENT_PLATFORM)-$(CLIENT_ARCH)$(FILE_EXT) $(BINDIR)/$(MIXIN)$(FILE_EXT)
	$(GO) test -tags=integration ./tests/...

publish: bin/porter$(FILE_EXT)
	go run mage.go Publish $(MIXIN)

bin/porter$(FILE_EXT):
	curl -fsSLo bin/porter$(FILE_EXT) https://cdn.porter.sh/canary/porter-$(CLIENT_PLATFORM)-$(CLIENT_ARCH)$(FILE_EXT)
	chmod +x bin/porter$(FILE_EXT)

install:
	mkdir -p $(PORTER_HOME)/mixins/$(MIXIN)/runtimes
	install $(BINDIR)/$(MIXIN)$(FILE_EXT) $(PORTER_HOME)/mixins/$(MIXIN)/$(MIXIN)$(FILE_EXT)
	install $(BINDIR)/$(MIXIN)-runtime$(FILE_EXT) $(PORTER_HOME)/mixins/$(MIXIN)/runtimes/$(MIXIN)-runtime$(FILE_EXT)

clean:
	-rm -fr bin/
