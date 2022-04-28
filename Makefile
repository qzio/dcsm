NAME:=$(shell basename $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST)))))

.PHONY: build install build_linux_amd64 build_linux_arm64 build_darwin_arm64

.DEFAULT: build

build:
	CGO_ENABLED=0 go build

build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(NAME)_linux_amd64

build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(NAME)_linux_arm64

build_darwin_arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(NAME)_darwin_arm64

install:
	go install
