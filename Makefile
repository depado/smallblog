.DEFAULT_GOAL := build

export GO111MODULE=on
export CGO_ENABLED=0
export VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
export BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
export BUILDDATE=$(shell LANG=en_us_88591 date)
BINARY=smallblog
LDFLAGS=-ldflags "-X 'github.com/Depado/smallblog/cmd.Version=$(VERSION)' \
		-X 'github.com/Depado/smallblog/cmd.Build=$(BUILD)' \
		-X 'github.com/Depado/smallblog/cmd.Time=$(BUILDDATE)' -s -w"
PACKEDFLAGS=-ldflags "-X 'github.com/Depado/smallblog/cmd.Version=$(VERSION)' \
		-X 'github.com/Depado/smallblog/cmd.Build=$(BUILD)' \
		-X 'github.com/Depado/smallblog/cmd.Time=$(BUILDDATE)' \
		-X 'github.com/Depado/smallblog/cmd.Packer=upx --best --lzma' -s -w"

.PHONY: help
help: ## Display help text for makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build
	go build $(LDFLAGS) -o $(BINARY)

.PHONY: tmp
tmp: ## Build and output the binary in /tmp
	go build $(LDFLAGS) -o /tmp/$(BINARY)

.PHONY: packed
packed: ## Build a packed version of the binary
	go build $(PACKEDFLAGS) -o $(BINARY)
	upx --best --lzma $(BINARY)

.PHONY: docker
docker: ## Build the docker image
	docker build -t $(BINARY):latest -t $(BINARY):$(BUILD) -f Dockerfile .

.PHONY: release
release: ## Create a new release on Github
	goreleaser

.PHONY: snapshot
snapshot: ## Create a new snapshot release
	goreleaser --snapshot --rm-dist

.PHONY: lint
lint: ## Runs the linter
	$(GOPATH)/bin/golangci-lint run --exclude-use-default=false

.PHONY: test
test: ## Run the test suite
	CGO_ENABLED=1 go test -race -coverprofile="coverage.txt" ./...

.PHONY: clean
clean: ## Remove the binary
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi
