.PHONY: test

all:	test build linux

test:	internal/*.go
	cd internal && go test -v && go vet

build:	*/*.go
	go build -o bin/tf main/main.go

linux:	*/*.go
	GOOS=linux go build -o bin/ltf main/main.go

