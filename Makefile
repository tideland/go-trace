SHELL=/bin/bash
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCOVER=$(GOCMD) tool cover
GOLINT=golint


GO111MODULE=on


.PHONY: setup
setup: ## Install all the build and lint dependencies
	$(GOGET) -u golang.org/x/tools/cmd/cover
	$(GOGET) -u golang.org/x/lint/golint


.PHONY: build
build: ## Build the library
	$(GOBUILD) -v  ./...


.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && $(GOTEST) -v -race -covermode=atomic -coverprofile=coverage.txt -timeout=30s ./...


.PHONY: clean
clean: ## Clean 
	$(GOCLEAN)


.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	$(GOCOVER) -html=coverage.txt


.PHONY: lint
lint: ## Run the linter
	$(GOLINT) ./...


.PHONY: ci
ci: lint test ## Run all the tests and code checks


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
