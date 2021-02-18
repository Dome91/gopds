BUILD_FLAGS=-ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,osusergo,netgo

all: clean install-dependencies build-ui build
with-tests: clean install-dependencies test-ui build-ui test build
only-tests: clean install-dependencies test-ui test

install-dependencies:
	go mod download
	cd ui && npm install

test-ui:
	cd ui && npm run test -- --watch=false --browsers=ChromeHeadless

build-ui:
	rm -rf public/assets
	rm -f public/*.txt
	rm -f public/*.ico
	rm -f public/*.html
	rm -f public/*.js
	rm -f public/*.css
	cd ui && npm run build -- --prod  --delete-output-path=false --output-path=../public

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
