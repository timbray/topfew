
all:	test build

test:	internal/*.go
	cd internal && go test -v

build:	*/*.go
	go build -o bin/tf main/main.go


