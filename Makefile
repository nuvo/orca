GIT_TAG := $(shell git describe --tags --always)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := "-s -w -X main.GitTag=${GIT_TAG} -X main.GitCommit=${GIT_COMMIT}"

export GO111MODULE:=on

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
	go mod download
