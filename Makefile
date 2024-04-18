.PHONY: test

all:	test tf

test:	main.go internal/*.go
	go test ./... && go vet ./...

# local version you can run
tf:
	go build -o bin/tf

release:	test
	GOOS=darwin  GOARCH=arm64 go build -o tf && gzip < tf > tf-macos-arm.gz
	GOOS=darwin  GOARCH=amd64 go build -o tf && gzip < tf > tf-macos-x86.gz
	GOOS=linux   GOARCH=amd64 go build -o tf && gzip < tf > tf-linux-x86.gz
	GOOS=linux   GOARCH=arm64 go build -o tf && gzip < tf > tf-linux-arm.gz
	GOOS=windows GOARCH=amd64 go build -o tf && zip -mq tf-windows-x86.exe.zip tf
	GOOS=windows GOARCH=arm64 go build -o tf && zip -mq tf-windows-arm.exe.zip tf
