.DEFAULT_GOAL := build

export GO111MODULE=on
export CGO_ENABLED=false
BINARY=smallblog

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build
	go build -o $(BINARY)

.PHONY: css
css: ## Build the CSS and map file from scss
	sass css/style.scss:assets/style.min.css --style compressed

.PHONY: install
install: ## Build and install
	go install

.PHONY: run
run: ## Runs the server
	@go run main.go

.PHONY: lint
lint: ## Runs the linter
	$(GOPATH)/bin/golangci-lint run --exclude-use-default=false

.PHONY: test
test: ## Run the unit test suite
	go test -race -coverprofile="coverage.txt" ./...

.PHONY: ttest
ttest: ## Run the unit test suite and parse it with tparse
	go test -race -coverprofile="coverage.txt" ./... -json | tparse -all

.PHONY: clean
clean: ## Remove the binary
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi
