.PHONY: help build/% run/% generate

# Which target to run by default (when no target is passed to make)
.DEFAULT_GOAL := help

help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

generate: ## Generates openapi code
	go generate ./...

build/%: ## Build an executable
	CGO_ENABLED=0 \
	go build \
		-mod=vendor \
		-o ./bin/$* \
		./cmd/$*

run/%: build/% ## Run the executable
	./bin/$*
