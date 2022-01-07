BINDIR := $(CURDIR)/bin
BINNAME ?= htmltomd
MAIN := ./cmd/htmltomd
LDFLAGS :=

GIT_TAG = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
ifneq ($(GIT_TAG),)
	LDFLAGS += -X github.com/david-mk-lawrence/html-to-md/internal/version.version=${GIT_TAG}
endif

# Rebuild the binary if any source files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) $(MAIN)

install:
	go install $(MAIN)

test:
	go test -cover -run . ./...
