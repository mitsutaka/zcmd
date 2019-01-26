test:
	GO111MODULE=on go build -mod=vendor -v ./...
	GO111MODULE=on go test -mod=vendor -race -v ./...
	golangci-lint run --presets bugs,unused,format,complexity,performance -D unparam --verbose
	misspell -error $(shell go list -mod=vendor ./... | grep -v /vendor/)

setup:
	GO111MODULE=off go get -v -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -v -u github.com/golangci/golangci-lint/cmd/golangci-lint

module:
	go mod tidy
	go mod vendor
	git add -f vendor
