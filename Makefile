VERSION=$(shell git describe --always --long --dirty)
.PHONY: all clean

all: test build

build:
	mkdir -p dist
	go build -v -ldflags="-X 'github.com/stevec7/ghostlist/cmd/ghostlist/version.version=${VERSION}'" -o dist/ghostlist cmd/ghostlist/main.go

clean:
	rm dist/ghostlist

test:
	cd pkg/ghostlist && go test  || (echo "Tests failed"; exit 1)

