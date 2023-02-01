.PHONY: build

VERSION := $(shell git rev-parse HEAD)

build:
	go build -ldflags "-X main.build=$(VERSION)"
tidy:
	go mod tidy
fmt:
	go fmt *.go
test:
	go test -v
