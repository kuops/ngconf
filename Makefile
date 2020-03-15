# ngconf make files
BIN := ngconf
VERSION := v0.1
GOARCH := amd64
DARWIN := darwin-$(GOARCH)
LINUX := linux-$(GOARCH)

.PHONY: build-linux build-darwin all

build-linux: 
	@GOOS=linux go build -o $(BIN)-$(VERSION)-$(LINUX)

build-darwin:
	@GOOS=darwin go build -o $(BIN)-$(VERSION)-$(DARWIN)

all: build-linux build-darwin

clean:
	@find . -type f -name "$(BIN)-$(VERSION)-*"|xargs rm -f