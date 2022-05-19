.PHONY: build clean tool lint help

all: build

build:
	cp config.toml bin/.
	go build -o bin/rexpromagent main.go

buildLinux:
	cp config.toml bin/.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/rexpromagent_linux_amd64 main.go

tool:
	go tool vet . |&amp; grep -v vendor; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf bin/*
	go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
