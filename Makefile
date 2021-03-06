PKG=github.com/mailru/easyjson
GOPATH=$(PWD):$(PWD)/vendor
BIN_PATH=$(PWD)/bin:$(PATH)

all: generate bin/gb
	PATH=$(BIN_PATH) gb build app

vendor:
	GOPATH=$(PWD)/vendor go get github.com/constabulary/gb/...
	PATH=$(PWD)/vendor/bin:$(PATH) gb vendor restore

bin/gb:
	GOPATH=$(GOPATH) go build -o=./bin/gb ./vendor/src/github.com/constabulary/gb/cmd/gb

bin/include: bin/gb
	PATH=$(BIN_PATH) gb build github.com/gobwas/include

bin/easyjson: bin/gb
	PATH=$(BIN_PATH) gb build github.com/mailru/easyjson/easyjson

generate: bin/gb bin/include bin/easyjson
	PATH=$(BIN_PATH) gb generate

install:
	install -d /usr/local/kursobot
	install -d /usr/local/kursobot/bin
	install -m0755 ./bin/app /usr/local/kursobot/bin/app
	install -m0644 ./etc/example.toml /usr/local/kursobot/kursobot.toml.example
	install -m0644 ./etc/supervisord.ini /usr/local/kursobot/kursobotd.ini.example

.PHONY: all install vendor bin/gb bin/include bin/easyjson