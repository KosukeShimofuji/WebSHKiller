# メタ情報
NAME = webshkiller
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := X'main.version=$(VERSION)' \
    -X'main.revesion=$(REVESION)'

# 必要なツール群をセットアップする
## Setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help
	go get github.com/rackspace/gophercloud
	go get golang.org/x/crypto/ssh
	go get github.com/fatih/color
	go get github.com/kr/pretty
	go get github.com/mitchellh/cli

# テストを実行する
## Run tests
test: deps
	go test $$(glide novendor)

# glideを使って依存パッケージをインストールする
## install dependencies
deps: setup
	glide install

## update dependencies
update: setup
	glide update

## Lint
lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do\
	    golint -set_exit_status $$pkg || exit $$?;\
	done

## Format sorce codes
fmt: setup
	goimports -w $$(glide nv -x)

## build binaries ex. make bin/myproj
bin/%: cmd/%/main.go deps
	go build -ldflags "$(LDFLAGS") -o $@ $<

## show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help

