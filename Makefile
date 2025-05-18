SHELL := /bin/bash
.PHONY: help config build-agent run-agent test alpha beta minor patch release

help:
	@echo "Makefile Commands:"
	@echo "  config               - Set up the environment."
	@echo "  build-agent          - Build kresa-agent"	
	@echo "  run-agent            - Build and run kresa-agent in-place"
	@echo "  test                 - Run tests"
	@echo "  alpha                - Generate changelog and create an alpha tag."
	@echo "  beta                 - Generate changelog and create an beta tag."
	@echo "  minor                - Generate changelog and create a minor tag."
	@echo "  patch                - Generate changelog and create a patch tag."
	@echo "  release              - Generate changelog and create a release tag."

all: config build run

config:
	@echo "Installing required tools"
	go mod tidy
	go mod vendor

build-agent:
	@echo "Building go-rssagg"
	go build -C cmd/kresa-agent

run-agent:
	@echo "Running go-rssagg"
	go build -C cmd/kresa-agent && cmd/kresa-agent/kresa-agent --env .env
	
test:
	@echo "Running tests"
#	./util/test.sh

alpha:
	@echo "Generating changelog and tag"
	commit-and-tag-version --prerelease alpha

beta:
	@echo "Generating changelog and tag"
	commit-and-tag-version --prerelease beta

minor:
	@echo "Generating changelog and tag"
	commit-and-tag-version --release-as minor

patch:
	@echo "Generating changelog and tag"
	commit-and-tag-version --release-as patch

release:
	@echo "Generating changelog and tag"
	commit-and-tag-version
