BUILD_FLAGS=-ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,osusergo,netgo

all: clean build-ui build
with-tests: clean test-ui build-ui test build
only-tests: clean test-ui test

test-ui:
	cd ui && npm install && npm run test -- --watch=false --browsers=ChromeHeadless

build-ui:
	mkdir -p public
	rm -f public/*
	cd ui && npm install && npm run build -- --prod --output-path=../public

test:
	go test ./...

build:
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o gopds -v

linux-amd64:
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o gopds-linux-amd64 -v

linux-armv6:
	CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=6 go build $(BUILD_FLAGS) -o gopds-linux-armv6 -v

clean:
	go clean
	rm -f gopds

.PHONY: clean test-ui build-ui test linux-amd64 linux-armv6 build
