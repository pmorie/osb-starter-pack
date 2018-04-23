SOURCES       := $(shell find . -name '*.go' -not -path "*/vendor/*")
SOURCE_DIRS    = middleware
.DEFAULT_GOAL := check


fmtcheck: ## Check go formatting
	@gofmt -l $(SOURCES) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

vet: ## Run go vet
	@go tool vet ./middleware

test: ## Run unit tests
	@go test -cover ./...

lint: ## Run golint
	@golint -set_exit_status $(addsuffix /... , $(SOURCE_DIRS))

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: fmtcheck vet lint test lint help
