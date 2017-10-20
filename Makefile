SHELL := /bin/bash
VERSION = $(shell git describe --tags --always --dirty)

.PHONY: help
help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@echo
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s : | sed -e 's/^/  /')"

.PHONY: test
test: ## Run all tests
	go test -v ./...

.PHONY: version
version: ## Update version string in version.go
	# echo ${VERSION}
	sed -i '' -E 's/Version = "[A-Z0-9a-z\.\-]+"/Version = "${VERSION}"/g' version.go
