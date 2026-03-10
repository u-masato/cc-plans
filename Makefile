.PHONY: build install clean test lint fmt setup

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

lint:
	go tool golangci-lint run

fmt:
	goimports -w .

setup:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/evilmartians/lefthook/v2@latest
	lefthook install --force
