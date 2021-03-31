.PHONY: test

all:	test build linux

test:	*/*.go
	go test -v ./... && go vet ./...

build:	*/*.go
	go build -o bin/tf .

linux:	*/*.go
	GOOS=linux go build -o bin/ltf .

