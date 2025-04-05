.PHONY: all build format lint

all: format lint build

build:
	go build .

format:
	go install mvdan.cc/gofumpt@latest
	gofumpt -w .

lint:
	go vet -v ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks all ./...
