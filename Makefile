all: test build

fmt:
	go fmt ./pkg/... ./cmd/...

vet:
	go vet ./pkg/... ./cmd/...

# Run tests
test: fmt vet
	go test ./test/... ./pkg/... -coverprofile cover.out

# Build orca binary
build: fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/orca cmd/orca.go