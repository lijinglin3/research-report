all: fmt vet golangci-lint test

check: fmt vet golangci-lint

fmt:
	test `go fmt ./... | wc -l` -eq 0

vet:
	go vet ./...

golangci-lint:
	golangci-lint run ./...

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
