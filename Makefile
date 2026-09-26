BINARY  := sysmon
PREFIX  ?= /usr
DESTDIR ?=
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build test lint cover install clean deb docker

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test -v ./...

lint:
	golangci-lint run

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "打开 coverage.html 查看覆盖率"

install: build
	install -D -m 0755 $(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)
	install -D -m 0644 configs/sysmon.yaml $(DESTDIR)/etc/sysmon/sysmon.yaml
	install -D -m 0644 deploy/sysmon.service $(DESTDIR)/lib/systemd/system/sysmon.service

clean:
	rm -f $(BINARY) coverage.out coverage.html *.deb

deb: build
	bash scripts/package-deb.sh

docker:
	docker build -t sysmon:$(VERSION) .
