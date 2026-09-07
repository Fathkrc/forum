.PHONY: run test vet fmt check-format build check

run:
	go run ./cmd/forum

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

check-format:
	test -z "$$(gofmt -l .)"

build:
	mkdir -p bin
	go build -o bin/movies-api ./cmd/forum

check: check-format vet test build
