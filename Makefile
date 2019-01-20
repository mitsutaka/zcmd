test:
	go get -v -t -d ./...
	golangci-lint run --presets bugs,unused,format,complexity,performance --verbose
	go test -v ./...

setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin
