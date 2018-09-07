SHELL=/bin/bash

GOVENDOR := $(shell command -v govendor)

GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "-CHANGES" || true)
GIT_DESCRIBE := $(shell git describe --tags --always)

VERSION := $(patsubst v%,%,$(GIT_DESCRIBE)$(GIT_DIRTY))

LDFLAGS := "-X main.Version=$(VERSION)"

DEBFILE := segment-rpc-legacy_$(VERSION)_amd64.deb

bin/rpc: dep
	mkdir -p bin
	go build -o bin/rpc ./cmd/rpc

bin/rpc-linux-amd64: dep
	mkdir -p bin
	env GOOS=linux GOARCH=amd64 go build -o bin/rpc-linux-amd64 ./cmd/rpc

$(DEBFILE): bin/rpc-linux-amd64
	fpm \
		-s dir \
		-t deb \
		-n segment-rpc-legacy \
		-v $(VERSION) \
		-m sre-team@segment.com \
		--vendor "Segment.io, Inc." \
		./bin/rpc-linux-amd64=/usr/bin/rpc-legacy

deb: $(DEBFILE)

upload-deb: $(DEBFILE)
	package_cloud push segment/infra/ubuntu/xenial $(DEBFILE)

dep:
ifndef GOVENDOR
	go get -u github.com/kardianos/govendor
endif
	govendor fetch +outside
	govendor sync

clean:
	rm -f bin/* *.deb

.PHONY: deb upload-deb clean
