test:
	go build -v ./...
	go test -race -v ./...
	golangci-lint run --enable-all -D gosec,dupl,funlen --verbose

lint-all:
	golangci-lint run --verbose

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

mod:
	go mod tidy
	git add go.mod
