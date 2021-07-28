test:
	go build -v ./...
	go test -race -v ./...
	golangci-lint run --verbose

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

mod:
	go mod tidy
	git add go.mod
