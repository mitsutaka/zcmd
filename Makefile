test:
	GO111MODULE=on go build -mod=vendor -v ./...
	GO111MODULE=on go test -mod=vendor -race -v ./...
	golangci-lint run --enable-all -D gosec,gochecknoglobals --verbose

lint-all:
	golangci-lint run --enable-all --verbose

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

mod:
	go mod tidy
	go mod vendor
	git add -f vendor
	git add go.mod
