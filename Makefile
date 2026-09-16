BINARY_NAME=mdread
PREFIX?=$(HOME)/go/bin

.PHONY: all build test clean install demo

all: test build

build:
	mkdir -p bin
	go build -ldflags="-s -w" -o bin/$(BINARY_NAME) ./cmd/mdread

test:
	go test -v ./...

install: build
	mkdir -p $(PREFIX)
	rm -f $(PREFIX)/$(BINARY_NAME)
	cp bin/$(BINARY_NAME) $(PREFIX)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(PREFIX)/$(BINARY_NAME)"

demo: build
	./bin/$(BINARY_NAME) --no-pager examples/sample.md

clean:
	rm -rf bin
