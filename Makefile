HAS_DEP := $(shell command -v dep;)
DEP_VERSION := v0.5.0
TEST_FILES := $(shell find ./test -type f -name "*.go")
GIT_TAG := $(shell git describe --tags --always)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.GitTag=${GIT_TAG} -X main.GitCommit=${GIT_COMMIT}"
DIST := $(CURDIR)/dist

all: bootstrap test build

fmt:
	go fmt ./pkg/... ./cmd/...

vet:
	go vet ./pkg/... ./cmd/...

# Run tests
test: fmt vet
	for f in $(TEST_FILES); do go test -v $$f; done 

# Build orca binary
build: fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o bin/orca cmd/orca.go

bootstrap:
ifndef HAS_DEP
	wget -q -O $(GOPATH)/bin/dep https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-linux-amd64
	chmod +x $(GOPATH)/bin/dep
endif
	dep ensure

dist:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o orca cmd/orca.go
	tar -zcvf $(DIST)/orca-linux-$(GIT_TAG).tgz orca
	rm orca
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o orca cmd/orca.go
	tar -zcvf $(DIST)/orca-macos-$(GIT_TAG).tgz orca
	rm orca
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o orca.exe cmd/orca.go
	tar -zcvf $(DIST)/orca-windows-$(GIT_TAG).tgz orca.exe
	rm orca.exe
