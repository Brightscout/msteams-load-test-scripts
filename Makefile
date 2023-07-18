GO ?= $(shell command -v go 2> /dev/null)

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -gcflags "all=-N -l" -trimpath -o dist/load-tester;

init_users:
	dist/load-tester init_users

create_channels:
	dist/load-tester create_channels

create_chats:
	dist/load-tester create_chats

clear_store:
	dist/load-tester clear_store

create_posts:
	k6 run k6/createPosts.js

check-style: 
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi; \

	@echo Running golangci-lint
	golangci-lint run --timeout 15m0s ./...
