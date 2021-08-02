# Copyright (c) 2021. Quirino Gervacio
# MIT License. All Rights Reserved

.PHONY: help test

help:
	@echo "Usage: make [target]"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "\033[36m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

tag ?= latest
export TAG=$(tag)
export BIAPI_AK=`cat test/biapiAk.test`
export BIAPI_SK=`cat test/biapiSk.test`
export TAAPI_SK=`cat test/taapiSk.test`
export EM_USER=`cat test/emUser.test`
export EM_PASS=`cat test/emPass.test`
export NOT_EM_LIST=`cat test/emUser.test`

t: test
test: ## (t) Test
	@go clean -testcache
	@go test ./... -v -coverprofile cp.out
	@go tool cover -html=cp.out

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
	@TAG=$(tag) \
        BIAPI_AK=$(BIAPI_AK) \
        BIAPI_SK=$(BIAPI_SK) \
        TAAPI_SK=$(TAAPI_SK) \
        EM_USER=$(EM_USER) \
        EM_PASS=$(EM_PASS) \
        NOTI_EM_LIST=$(NOTI_EM_LIST) docker-compose push
