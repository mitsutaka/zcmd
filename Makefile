test:
	go build -mod=vendor -v ./...
	go test -mod=vendor -v ./...
	golangci-lint run --presets bugs,unused,format,complexity,performance -D unparam --verbose

setup:
	GO111MODULE=off go get -v github.com/client9/misspell/cmd/misspell
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin
