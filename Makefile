PKG=github.com/mailru/easyjson
GOPATH=$(PWD):$(PWD)/vendor
BIN_PATH=$(PWD)/bin:$(PWD)/vendor/bin:$(PATH)

all: generate bin/gb
	PATH=$(BIN_PATH) gb build app

gb:
	GOPATH=$(PWD)/vendor go get github.com/constabulary/gb/...

vendor: gb bin/gb
	PATH=$(BIN_PATH) gb vendor restore

bin/gb:
	GOPATH=$(GOPATH) go build -o=./bin/gb ./vendor/src/github.com/constabulary/gb/cmd/gb

bin/include: bin/gb
	PATH=$(BIN_PATH) gb build github.com/gobwas/include

bin/easyjson: bin/gb
	PATH=$(BIN_PATH) gb build github.com/mailru/easyjson/easyjson

generate: bin/gb bin/include bin/easyjson
	PATH=$(BIN_PATH) gb generate

.PHONY: bin/gb bin/include bin/easyjson


