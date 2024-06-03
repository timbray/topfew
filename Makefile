.PHONY: test

all:	test topfew

test:	main.go internal/*.go
	go test ./... && go vet ./...

# local version you can run
topfew:
	go build -o bin/topfew

release:	test
	GOOS=darwin  GOARCH=arm64 go build -o topfew && gzip < topfew>topfew-macos-arm.gz
	GOOS=darwin  GOARCH=amd64 go build -o topfew && gzip < topfew>topfew-macos-x86.gz
	GOOS=linux   GOARCH=amd64 go build -o topfew && gzip < topfew>topfew-linux-x86.gz
	GOOS=linux   GOARCH=arm64 go build -o topfew && gzip < topfew>topfew-linux-arm.gz
	GOOS=windows GOARCH=amd64 go build -o topfew && zip -mq topfew-windows-x86.exe.zip topfew
	GOOS=windows GOARCH=arm64 go build -o topfew && zip -mq topfew-windows-arm.exe.zip topfew
