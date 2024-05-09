
SHELL := bash # the shell used internally by Make

GOCMD ?= $(shell which go)

.PHONY: all build 

all: build

build:
	${GOCMD} build -o build/storeverif ./cmd/storeverif
