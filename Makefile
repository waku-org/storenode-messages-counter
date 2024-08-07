
.PHONY: all build lint-install lint

all: build

build:
	go build -tags=gowaku_no_rln -o build/storemsgcounter ./cmd/storemsgcounter

lint-install:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		bash -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint:
	@echo "lint"
	@golangci-lint run ./...
