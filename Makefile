.PHONY: all clean

build:
	mkdir -p dist
	go build -v -o dist/ghostlist cmd/ghostlist/main.go

clean:
	rm dist/ghostlist

test:
	cd pkg/ghostlist && go test || (echo "Tests failed"; exit 1)

all: test build
