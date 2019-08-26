.PHONY: all clean

build:
	mkdir -p dist
	go build -v -x -o dist/ghostlist cmd/ghostlist/main.go

clean:
	rm dist/ghostlist
