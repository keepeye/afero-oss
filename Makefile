.PHONY: test fmt lint

test:
	go test -race -v ./...

fmt:
	golangci-lint fmt ./...

lint:
	golangci-lint run ./...