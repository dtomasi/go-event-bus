SHELL := /bin/bash

setup: install-pre-commit-hooks
	go mod vendor

install-pre-commit-hooks:
	pre-commit install --install-hooks

run-pre-commit:
	pre-commit run -a

test:
	go test -v -race ./...

coverage:
	go test -v -race -cover -covermode=atomic ./...
