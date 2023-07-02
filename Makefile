.PHONY: build lint

GOOS    := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH  := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
GOENV   := GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO      := $(GOENV) go

COMMIT    := $(shell git describe --no-match --always --dirty)
BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE := $(shell date -Iseconds)

REPO := github.com/lobshunter/dorisctl

LDFLAGS := -w -s
LDFLAGS += -X "$(REPO)/pkg/version.gitHash=$(COMMIT)"
LDFLAGS += -X "$(REPO)/pkg/version.gitBranch=$(BRANCH)"
LDFLAGS += -X "$(REPO)/pkg/version.buildDate=$(BUILD_DATE)"

default: lint build

build:
	$(GO) build -ldflags '$(LDFLAGS)' -o bin/dorisctl-${GOOS}-${GOARCH} ./cmd/dorisctl

lint:
	golangci-lint run