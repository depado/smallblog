.PHONY: all clean test docker latest release

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

BINARY=smallblog
VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all:
	go build -o $(BINARY) $(LDFLAGS)

docker:
	docker build \
		-t $(BINARY):latest \
		-t $(BINARY):$(VERSION) \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile --no-cache .

latest:
	docker build \
		-t $(BINARY):latest \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile --no-cache .

test:
	go test ./...

clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi