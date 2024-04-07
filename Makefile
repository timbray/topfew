.PHONY: test

all:	test build 

test:	*/*.go
	go test -v ./... && go vet ./...

build:	bin/macos-arm/tf bin/macos-x86/tf bin/linux-x86/tf bin/linux-arm/tf

bin/macos-arm/tf:	*/*.go
	GOOS=darwin GOARCH=arm64 go build -o bin/macos-arm/tf

bin/macos-x86/tf:	*/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/macos-x86/tf

bin/linux-x86/tf:	*/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux-x86/tf

bin/linux-arm/tf:	*/*.go
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm/tf


