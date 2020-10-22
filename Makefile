.PHONY: build

# Build targets
TARGETS ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 windows/amd64
TARGET_OBJS ?= darwin-amd64.tar.gz darwin-amd64.tar.gz.sha256 linux-amd64.tar.gz linux-amd64.tar.gz.sha256 linux-386.tar.gz linux-386.tar.gz.sha256 linux-arm.tar.gz linux-arm.tar.gz.sha256 linux-arm64.tar.gz linux-arm64.tar.gz.sha256 windows-amd64.zip windows-amd64.zip.sha256

GO        	?= go
BINDIR    	:= $(CURDIR)/bin
BINARIES  	:= ferry
TAG			:= $(shell git describe --tags --abbrev=0)
GOFLAGS   	:=
LDFLAGS   	:= -w -s -X github.com/gojira/ferry/cmd.Version=$(TAG)
PKG       	:= $(shell go mod vendor)

# Go Package required
PKG_GOX := github.com/mitchellh/gox@v1.0.1

build: test
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(BINARIES)'

build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -parallel=3 -output="dist/$(BINARIES)-{{.OS}}-{{.Arch}}" -osarch='$(TARGETS)' $(GOFLAGS) $(if $(TAGS),-tags '$(TAGS)',) -ldflags '$(LDFLAGS)'

# fmt will fix the golang source style in place.
fmt:
	$(GO) fmt ./...

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	rm -Rf bin/*
	go clean