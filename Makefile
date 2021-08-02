# Copyright (c) 2021. Quirino Gervacio
# MIT License. All Rights Reserved

.PHONY: help test

help:
	@echo "Usage: make [target]"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "\033[36m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

tag ?= latest

t: test
test: ## (t) Test
	@go clean -testcache
	@go test ./... -v -coverprofile cp.out
	@go tool cover -html=cp.out

r: run
run: ## (r) Run main.go
	@go run main.go \
	    --specFile test/spec.yaml \
        --logLevel debug

u: up
up: ## (u) Run a stack for testing
	@TAG=$(tag) docker-compose up

d: down
down: ## (d) Stop the running stack
	@TAG=$(tag) docker-compose down

b: build
build: ## (b) Build image (param: tag)
	@TAG=$(tag) docker-compose build

p: push
push: ## (p) Push image (param: tag)
	@TAG=$(tag) docker-compose push
