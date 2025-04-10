.PHONY: test fmt lint

test:
	go test -race -v ./...

fmt:
	gofumpt -w ./
	goimports -w  -local github.com/messikiller/afero-oss ./
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...