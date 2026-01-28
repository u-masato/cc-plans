.PHONY: build install clean test

BINARY=cc-plans
INSTALL_PATH=$(HOME)/.local/bin

build:
	go build -o $(BINARY) ./cmd/cc-plans

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY) $(INSTALL_PATH)/

clean:
	rm -f $(BINARY)

test:
	go test ./...
