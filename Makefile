BINDIR      := $(CURDIR)/bin
VERSION ?= dev

build-cli:
	go build -ldflags="-w -s -X 'github.com/porter-dev/porter/cli.Version=${VERSION}'" -a -tags cli -o $(BINDIR)/switchboard ./cli

build-cli-dev:
	go build -tags cli -o $(BINDIR)/switchboard ./cli