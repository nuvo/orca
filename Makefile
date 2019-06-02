HAS_DEP := $(shell command -v dep;)
DEP_VERSION := v0.5.0
GIT_TAG := $(shell git describe --tags --always)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := "-s -w -X main.GitTag=${GIT_TAG} -X main.GitCommit=${GIT_COMMIT}"

all: bootstrap test build

fmt:
	go fmt ./...

vet:
	go vet ./...

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Build orca binary
build: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o bin/orca cmd/orca.go

# Build orca docker image
docker: build
	cp bin/orca orca
	docker build -t nuvo/orca:latest .
	rm orca

bootstrap:
ifndef HAS_DEP
	wget -q -O $(GOPATH)/bin/dep https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-linux-amd64
	chmod +x $(GOPATH)/bin/dep
endif
	dep ensure
