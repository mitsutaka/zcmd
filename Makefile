test:
	go get -v -t -d ./...
	golangci-lint run --presets bugs,unused,format,complexity,performance --verbose
	go test -v ./...

setup:
	GO111MODULE=off go get -v github.com/client9/misspell/cmd/misspell
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin
