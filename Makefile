SHELL := /bin/bash -o pipefail

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

.PHONY: help

help:
	@echo "Usage: make <TARGET>"
	@echo ""
	@echo "Available targets are:"
	@echo ""
	@echo "    generate-gql                Generate Go files from *.graphqls"
	@echo "    build                       Build the binary"
	@echo "    install gqlgen              Instal gqlgen in your project directory"
	@echo ""

.PHONY: generate-gql
generate-gql: 
	@scripts/gqlgen.sh -m $(call args,module)

.PHONY: build
build:
	@mkdir -p bin
	@go build -mod vendor -o bin/go-project main.go
	@echo "build done"

.PHONY: install-gqlgen
install-gqlgen:
	@go get github.com/99designs/gqlgen 
	@echo "gqlgen installed!"

.PHONY: run
run:
	@go run main.go serve
	@echo "service started!"