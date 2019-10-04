GOFLAGS = -mod=vendor
export GOFLAGS

test:
	GO111MODULE=on go build -v ./...
	GO111MODULE=on go test -race -v ./...
	golangci-lint run --enable-all -D gosec,dupl,funlen --verbose

lint-all:
	golangci-lint run --enable-all --verbose

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

mod:
	go mod tidy
	go mod vendor
	git add -f vendor
	git add go.mod
