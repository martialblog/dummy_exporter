.PHONY: test coverage lint vet

COMMIT := $(shell git rev-parse HEAD)
VERSION?=latest

build:
	go build
release:
	go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"
lint:
	go fmt $(go list ./... | grep -v /vendor/)
vet:
	go vet $(go list ./... | grep -v /vendor/)
test:
	go test -v -race ./...
coverage:
	go test -v -cover -coverprofile=coverage.out ./... &&\
	go tool cover -html=coverage.out -o coverage.html
clean:
	rm -f dummy_exporter dummy_exporter.test *.benchmark *.profile
