.PHONY: run build test

run:
	go run main.go

build:
	go build -o assisko main.go

test:
	go test ./...
