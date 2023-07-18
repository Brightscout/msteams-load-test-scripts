GO ?= $(shell command -v go 2> /dev/null)

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -gcflags "all=-N -l" -trimpath -o dist/load-tester;

init_users:
	dist/load-tester init_users

create_channels:
	dist/load-tester create_channels

clear_store:
	dist/load-tester clear_store
